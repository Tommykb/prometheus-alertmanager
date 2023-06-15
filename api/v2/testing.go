// Copyright 2022 Prometheus Team
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

package v2

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/prometheus/alertmanager/provider"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	open_api_models "github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/alertmanager/silence/silencepb"
)

func createSilence(t *testing.T, ID, creator string, start, ends time.Time) (open_api_models.PostableSilence, []byte) {
	t.Helper()

	comment := "test"
	matcherName := "a"
	matcherValue := "b"
	isRegex := false
	startsAt := strfmt.DateTime(start)
	endsAt := strfmt.DateTime(ends)

	sil := open_api_models.PostableSilence{
		ID: ID,
		Silence: open_api_models.Silence{
			Matchers:  open_api_models.Matchers{&open_api_models.Matcher{Name: &matcherName, Value: &matcherValue, IsRegex: &isRegex}},
			StartsAt:  &startsAt,
			EndsAt:    &endsAt,
			CreatedBy: &creator,
			Comment:   &comment,
		},
	}
	b, err := json.Marshal(&sil)
	require.NoError(t, err)

	return sil, b
}

func createSilenceMatcher(t *testing.T, name, pattern string, matcherType silencepb.Matcher_Type) *silencepb.Matcher {
	t.Helper()

	return &silencepb.Matcher{
		Name:    name,
		Pattern: pattern,
		Type:    matcherType,
	}
}

func createLabelMatcher(t *testing.T, name, value string, matchType labels.MatchType) *labels.Matcher {
	t.Helper()

	matcher, _ := labels.NewMatcher(matchType, name, value)
	return matcher
}

// fakeAlerts is a struct implementing the provider.Alerts interface for tests.
type fakeAlerts struct {
	fps    map[model.Fingerprint]int
	alerts []*types.Alert
	err    error
}

func newFakeAlerts(alerts []*types.Alert) *fakeAlerts {
	fps := make(map[model.Fingerprint]int)
	for i, a := range alerts {
		fps[a.Fingerprint()] = i
	}
	f := &fakeAlerts{
		alerts: alerts,
		fps:    fps,
	}
	return f
}

func (f *fakeAlerts) Subscribe() provider.AlertIterator           { return nil }
func (f *fakeAlerts) Get(model.Fingerprint) (*types.Alert, error) { return nil, nil }
func (f *fakeAlerts) Put(alerts ...*types.Alert) error {
	return f.err
}

func (f *fakeAlerts) GetPending() provider.AlertIterator {
	ch := make(chan *types.Alert)
	done := make(chan struct{})
	go func() {
		defer close(ch)
		for _, a := range f.alerts {
			ch <- a
		}
	}()
	return provider.NewAlertIterator(ch, done, f.err)
}
