// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inhibit

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/alertmanager/provider/mem"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/client_golang/prometheus"
)

// BenchmarkMutes benchmarks the Mutes method for the Muter interface
// for different numbers of inhibition rules.
func BenchmarkMutes(b *testing.B) {
	b.Run("1 inhibition rule, 1 inhibiting alert", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitionRules(b, 1))
	})
	b.Run("10 inhibition rules, 1 inhibiting alert", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitionRules(b, 10))
	})
	b.Run("100 inhibition rules, 1 inhibiting alert", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitionRules(b, 100))
	})
	b.Run("1000 inhibition rules, 1 inhibiting alert", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitionRules(b, 1000))
	})
	b.Run("10000 inhibition rules, 1 inhibiting alert", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitionRules(b, 10000))
	})
	b.Run("1 inhibition rule, 10 inhibiting alerts", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitingAlerts(b, 10))
	})
	b.Run("1 inhibition rule, 100 inhibiting alerts", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitingAlerts(b, 100))
	})
	b.Run("1 inhibition rule, 1000 inhibiting alerts", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitingAlerts(b, 1000))
	})
	b.Run("1 inhibition rule, 10000 inhibiting alerts", func(b *testing.B) {
		benchmarkMutes(b, benchmarkNumInhibitingAlerts(b, 10000))
	})
}

type benchmarkOptions struct {
	// n is the total number of inhibition rules to benchmark.
	n int
	// newRuleFunc creates the next inhibition rule in the benchmark.
	// It is called n times.
	newRuleFunc func(idx int) config.InhibitRule
	// newAlertsFunc creates teh inhibiting alerts for each inhibition
	// rule in the benchmark. It is called n times.
	newAlertsFunc func(idx int, r config.InhibitRule) []types.Alert
	// benchFunc runs the benchmark.
	benchFunc func(mutesFunc func(model.LabelSet) bool) error
}

// benchmarkNumInhibitionRules creates N inhibition rules with different source
// matchers, but the same target matchers. The source matchers are suffixed with
// the position of the inhibition rule in the list. For example, foo=bar1, foo=bar2,
// etc.
func benchmarkNumInhibitionRules(b *testing.B, n int) benchmarkOptions {
	return benchmarkOptions{
		n: n,
		newRuleFunc: func(idx int) config.InhibitRule {
			return config.InhibitRule{
				SourceMatchers: config.Matchers{
					mustNewMatcher(b, labels.MatchEqual, "foo", "bar"+strconv.Itoa(idx)),
				},
				TargetMatchers: config.Matchers{
					mustNewMatcher(b, labels.MatchEqual, "bar", "baz"),
				},
			}
		},
		newAlertsFunc: func(idx int, _ config.InhibitRule) []types.Alert {
			return []types.Alert{{
				Alert: model.Alert{
					Labels: model.LabelSet{
						"foo": model.LabelValue("bar" + strconv.Itoa(idx)),
					},
				},
			}}
		}, benchFunc: func(mutesFunc func(set model.LabelSet) bool) error {
			if ok := mutesFunc(model.LabelSet{"bar": "baz"}); !ok {
				return errors.New("expected bar=baz to be muted")
			}
			return nil
		},
	}
}

// benchmarkNumInhibitingAlerts creates 1 inhibition rule, but with N matching alerts.
func benchmarkNumInhibitingAlerts(b *testing.B, n int) benchmarkOptions {
	return benchmarkOptions{
		n: 1,
		newRuleFunc: func(_ int) config.InhibitRule {
			return config.InhibitRule{
				SourceMatchers: config.Matchers{
					mustNewMatcher(b, labels.MatchEqual, "foo", "bar"),
				},
				TargetMatchers: config.Matchers{
					mustNewMatcher(b, labels.MatchEqual, "bar", "baz"),
				},
			}
		},
		newAlertsFunc: func(_ int, _ config.InhibitRule) []types.Alert {
			var alerts []types.Alert
			for i := 0; i < n; i++ {
				alerts = append(alerts, types.Alert{
					Alert: model.Alert{
						Labels: model.LabelSet{
							"foo": "bar",
							"idx": model.LabelValue(strconv.Itoa(i)),
						},
					},
				})
			}
			return alerts
		},
		benchFunc: func(mutesFunc func(set model.LabelSet) bool) error {
			if ok := mutesFunc(model.LabelSet{"bar": "baz"}); !ok {
				return errors.New("expected bar=baz to be muted")
			}
			return nil
		},
	}
}

func benchmarkMutes(b *testing.B, opts benchmarkOptions) {
	r := prometheus.NewRegistry()
	m := types.NewMarker(r)
	s, err := mem.NewAlerts(context.TODO(), m, time.Minute, nil, log.NewNopLogger(), r)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	alerts, rules := benchmarkFromOptions(opts)
	for _, a := range alerts {
		tmp := a
		if err = s.Put(&tmp); err != nil {
			b.Fatal(err)
		}
	}

	ih := NewInhibitor(s, rules, m, log.NewNopLogger())
	defer ih.Stop()
	go ih.Run()

	// Wait some time for the inhibitor to seed its cache.
	waitDuration := time.Millisecond * time.Duration(len(alerts))
	if waitDuration > time.Second {
		waitDuration = time.Second
	}
	<-time.After(waitDuration)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		require.NoError(b, opts.benchFunc(ih.Mutes))
	}
}

func benchmarkFromOptions(opts benchmarkOptions) ([]types.Alert, []config.InhibitRule) {
	var (
		alerts = make([]types.Alert, 0, opts.n)
		rules  = make([]config.InhibitRule, 0, opts.n)
	)
	for i := 0; i < opts.n; i++ {
		r := opts.newRuleFunc(i)
		alerts = append(alerts, opts.newAlertsFunc(i, r)...)
		rules = append(rules, r)
	}
	return alerts, rules
}

func mustNewMatcher(b *testing.B, op labels.MatchType, name, value string) *labels.Matcher {
	m, err := labels.NewMatcher(op, name, value)
	require.NoError(b, err)
	return m
}
