// Copyright 2018 Prometheus Team
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

package template

import (
	"encoding/json"
	tmplhtml "html/template"
	"net/url"
	"sync"
	"testing"
	tmpltext "text/template"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/prometheus/alertmanager/types"
)

func TestPairNames(t *testing.T) {
	pairs := Pairs{
		{"name1", "value1"},
		{"name2", "value2"},
		{"name3", "value3"},
	}

	expected := []string{"name1", "name2", "name3"}
	require.EqualValues(t, expected, pairs.Names())
}

func TestPairValues(t *testing.T) {
	pairs := Pairs{
		{"name1", "value1"},
		{"name2", "value2"},
		{"name3", "value3"},
	}

	expected := []string{"value1", "value2", "value3"}
	require.EqualValues(t, expected, pairs.Values())
}

func TestPairsString(t *testing.T) {
	pairs := Pairs{{"name1", "value1"}}
	require.Equal(t, "name1=value1", pairs.String())
	pairs = append(pairs, Pair{"name2", "value2"})
	require.Equal(t, "name1=value1, name2=value2", pairs.String())
}

func TestKVSortedPairs(t *testing.T) {
	kv := KV{"d": "dVal", "b": "bVal", "c": "cVal"}

	expectedPairs := Pairs{
		{"b", "bVal"},
		{"c", "cVal"},
		{"d", "dVal"},
	}

	for i, p := range kv.SortedPairs() {
		require.EqualValues(t, p.Name, expectedPairs[i].Name)
		require.EqualValues(t, p.Value, expectedPairs[i].Value)
	}

	// validates alertname always comes first
	kv = KV{"d": "dVal", "b": "bVal", "c": "cVal", "alertname": "alert", "a": "aVal"}

	expectedPairs = Pairs{
		{"alertname", "alert"},
		{"a", "aVal"},
		{"b", "bVal"},
		{"c", "cVal"},
		{"d", "dVal"},
	}

	for i, p := range kv.SortedPairs() {
		require.EqualValues(t, p.Name, expectedPairs[i].Name)
		require.EqualValues(t, p.Value, expectedPairs[i].Value)
	}
}

func TestKVRemove(t *testing.T) {
	kv := KV{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}

	kv = kv.Remove([]string{"key2", "key4"})

	expected := []string{"key1", "key3"}
	require.EqualValues(t, expected, kv.Names())
}

func TestAlertsFiring(t *testing.T) {
	alerts := Alerts{
		{Status: string(model.AlertFiring)},
		{Status: string(model.AlertResolved)},
		{Status: string(model.AlertFiring)},
		{Status: string(model.AlertResolved)},
		{Status: string(model.AlertResolved)},
	}

	for _, alert := range alerts.Firing() {
		if alert.Status != string(model.AlertFiring) {
			t.Errorf("unexpected status %q", alert.Status)
		}
	}
}

func TestAlertsResolved(t *testing.T) {
	alerts := Alerts{
		{Status: string(model.AlertFiring)},
		{Status: string(model.AlertResolved)},
		{Status: string(model.AlertFiring)},
		{Status: string(model.AlertResolved)},
		{Status: string(model.AlertResolved)},
	}

	for _, alert := range alerts.Resolved() {
		if alert.Status != string(model.AlertResolved) {
			t.Errorf("unexpected status %q", alert.Status)
		}
	}
}

func TestData(t *testing.T) {
	u, err := url.Parse("http://example.com/")
	require.NoError(t, err)
	tmpl := &Template{ExternalURL: u}
	startTime := time.Time{}.Add(1 * time.Second)
	endTime := time.Time{}.Add(2 * time.Second)
	expStartTime := Time{startTime}
	expEndTime := Time{endTime}

	for _, tc := range []struct {
		receiver    string
		groupLabels model.LabelSet
		alerts      []*types.Alert

		exp *Data
	}{
		{
			receiver: "webhook",
			exp: &Data{
				Receiver:          "webhook",
				Status:            "resolved",
				Alerts:            Alerts{},
				GroupLabels:       KV{},
				CommonLabels:      KV{},
				CommonAnnotations: KV{},
				ExternalURL:       u.String(),
			},
		},
		{
			receiver: "webhook",
			groupLabels: model.LabelSet{
				model.LabelName("job"): model.LabelValue("foo"),
			},
			alerts: []*types.Alert{
				{
					Alert: model.Alert{
						StartsAt: startTime,
						Labels: model.LabelSet{
							model.LabelName("severity"): model.LabelValue("warning"),
							model.LabelName("job"):      model.LabelValue("foo"),
						},
						Annotations: model.LabelSet{
							model.LabelName("description"): model.LabelValue("something happened"),
							model.LabelName("runbook"):     model.LabelValue("foo"),
						},
					},
				},
				{
					Alert: model.Alert{
						StartsAt: startTime,
						EndsAt:   endTime,
						Labels: model.LabelSet{
							model.LabelName("severity"): model.LabelValue("critical"),
							model.LabelName("job"):      model.LabelValue("foo"),
						},
						Annotations: model.LabelSet{
							model.LabelName("description"): model.LabelValue("something else happened"),
							model.LabelName("runbook"):     model.LabelValue("foo"),
						},
					},
				},
			},
			exp: &Data{
				Receiver: "webhook",
				Status:   "firing",
				Alerts: Alerts{
					{
						Status:      "firing",
						Labels:      KV{"severity": "warning", "job": "foo"},
						Annotations: KV{"description": "something happened", "runbook": "foo"},
						StartsAt:    expStartTime,
						Fingerprint: "9266ef3da838ad95",
					},
					{
						Status:      "resolved",
						Labels:      KV{"severity": "critical", "job": "foo"},
						Annotations: KV{"description": "something else happened", "runbook": "foo"},
						StartsAt:    expStartTime,
						EndsAt:      expEndTime,
						Fingerprint: "3b15fd163d36582e",
					},
				},
				GroupLabels:       KV{"job": "foo"},
				CommonLabels:      KV{"job": "foo"},
				CommonAnnotations: KV{"runbook": "foo"},
				ExternalURL:       u.String(),
			},
		},
		{
			receiver:    "webhook",
			groupLabels: model.LabelSet{},
			alerts: []*types.Alert{
				{
					Alert: model.Alert{
						StartsAt: startTime,
						Labels: model.LabelSet{
							model.LabelName("severity"): model.LabelValue("warning"),
							model.LabelName("job"):      model.LabelValue("foo"),
						},
						Annotations: model.LabelSet{
							model.LabelName("description"): model.LabelValue("something happened"),
							model.LabelName("runbook"):     model.LabelValue("foo"),
						},
					},
				},
				{
					Alert: model.Alert{
						StartsAt: startTime,
						EndsAt:   endTime,
						Labels: model.LabelSet{
							model.LabelName("severity"): model.LabelValue("critical"),
							model.LabelName("job"):      model.LabelValue("bar"),
						},
						Annotations: model.LabelSet{
							model.LabelName("description"): model.LabelValue("something else happened"),
							model.LabelName("runbook"):     model.LabelValue("bar"),
						},
					},
				},
			},
			exp: &Data{
				Receiver: "webhook",
				Status:   "firing",
				Alerts: Alerts{
					{
						Status:      "firing",
						Labels:      KV{"severity": "warning", "job": "foo"},
						Annotations: KV{"description": "something happened", "runbook": "foo"},
						StartsAt:    expStartTime,
						Fingerprint: "9266ef3da838ad95",
					},
					{
						Status:      "resolved",
						Labels:      KV{"severity": "critical", "job": "bar"},
						Annotations: KV{"description": "something else happened", "runbook": "bar"},
						StartsAt:    expStartTime,
						EndsAt:      expEndTime,
						Fingerprint: "c7e68cb08e3e67f9",
					},
				},
				GroupLabels:       KV{},
				CommonLabels:      KV{},
				CommonAnnotations: KV{},
				ExternalURL:       u.String(),
			},
		},
	} {
		tc := tc
		t.Run("", func(t *testing.T) {
			got := tmpl.Data(tc.receiver, tc.groupLabels, tc.alerts...)
			require.Equal(t, tc.exp, got)
		})
	}
}

func TestTemplateExpansion(t *testing.T) {
	tmpl, err := FromGlobs([]string{})
	require.NoError(t, err)

	for _, tc := range []struct {
		title string
		in    string
		data  interface{}
		html  bool

		exp  string
		fail bool
	}{
		{
			title: "Template without action",
			in:    `abc`,
			exp:   "abc",
		},
		{
			title: "Template with simple action",
			in:    `{{ "abc" }}`,
			exp:   "abc",
		},
		{
			title: "Template with invalid syntax",
			in:    `{{ `,
			fail:  true,
		},
		{
			title: "Template using toUpper",
			in:    `{{ "abc" | toUpper }}`,
			exp:   "ABC",
		},
		{
			title: "Template using toLower",
			in:    `{{ "ABC" | toLower }}`,
			exp:   "abc",
		},
		{
			title: "Template using title",
			in:    `{{ "abc" | title }}`,
			exp:   "Abc",
		},
		{
			title: "Template using TrimSpace",
			in:    `{{ " a b c " | trimSpace }}`,
			exp:   "a b c",
		},
		{
			title: "Template using positive match",
			in:    `{{ if match "^a" "abc"}}abc{{ end }}`,
			exp:   "abc",
		},
		{
			title: "Template using negative match",
			in:    `{{ if match "abcd" "abc" }}abc{{ end }}`,
			exp:   "",
		},
		{
			title: "Template using join",
			in:    `{{ . | join "," }}`,
			data:  []string{"a", "b", "c"},
			exp:   "a,b,c",
		},
		{
			title: "Text template without HTML escaping",
			in:    `{{ "<b>" }}`,
			exp:   "<b>",
		},
		{
			title: "HTML template with escaping",
			in:    `{{ "<b>" }}`,
			html:  true,
			exp:   "&lt;b&gt;",
		},
		{
			title: "HTML template using safeHTML",
			in:    `{{ "<b>" | safeHtml }}`,
			html:  true,
			exp:   "<b>",
		},
		{
			title: "Template using reReplaceAll",
			in:    `{{ reReplaceAll "ab" "AB" "abcdabcda"}}`,
			exp:   "ABcdABcda",
		},
		{
			title: "Template using stringSlice",
			in:    `{{ with .GroupLabels }}{{ with .Remove (stringSlice "key1" "key3") }}{{ .SortedPairs.Values }}{{ end }}{{ end }}`,
			data: Data{
				GroupLabels: KV{
					"key1": "key1",
					"key2": "key2",
					"key3": "key3",
					"key4": "key4",
				},
			},
			exp: "[key2 key4]",
		},
		{
			title: "Template StartsAt and EndsAt time formatting",
			in:    `{{ range .Alerts }}{{ .StartsAt }} - {{ .EndsAt }}{{ end }}`,
			data: Data{
				Alerts: Alerts{
					Alert{
						StartsAt: Time{time.Date(2023, 11, 14, 15, 49, 0, 0, time.UTC)},
						EndsAt:   Time{time.Date(2023, 11, 14, 18, 22, 36, 2000000, time.UTC)},
					},
				},
			},
			exp: "2023-11-14 15:49:00 +0000 UTC - 2023-11-14 18:22:36.002 +0000 UTC",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			f := tmpl.ExecuteTextString
			if tc.html {
				f = tmpl.ExecuteHTMLString
			}
			got, err := f(tc.in, tc.data)
			if tc.fail {
				require.NotNil(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.exp, got)
		})
	}
}

func TestTemplateExpansionWithOptions(t *testing.T) {
	testOptionWithAdditionalFuncs := func(funcs FuncMap) Option {
		return func(text *tmpltext.Template, html *tmplhtml.Template) {
			text.Funcs(tmpltext.FuncMap(funcs))
			html.Funcs(tmplhtml.FuncMap(funcs))
		}
	}
	for _, tc := range []struct {
		options []Option
		title   string
		in      string
		data    interface{}
		html    bool

		exp  string
		fail bool
	}{
		{
			title:   "Test custom function",
			options: []Option{testOptionWithAdditionalFuncs(FuncMap{"printFoo": func() string { return "foo" }})},
			in:      `{{ printFoo }}`,
			exp:     "foo",
		},
		{
			title:   "Test Default function with additional function added",
			options: []Option{testOptionWithAdditionalFuncs(FuncMap{"printFoo": func() string { return "foo" }})},
			in:      `{{ toUpper "test" }}`,
			exp:     "TEST",
		},
		{
			title:   "Test custom function is overridden by the DefaultFuncs",
			options: []Option{testOptionWithAdditionalFuncs(FuncMap{"toUpper": func(s string) string { return "foo" }})},
			in:      `{{ toUpper "test" }}`,
			exp:     "TEST",
		},
		{
			title: "Test later Option overrides the previous",
			options: []Option{
				testOptionWithAdditionalFuncs(FuncMap{"printFoo": func() string { return "foo" }}),
				testOptionWithAdditionalFuncs(FuncMap{"printFoo": func() string { return "bar" }}),
			},
			in:  `{{ printFoo }}`,
			exp: "bar",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			tmpl, err := FromGlobs([]string{}, tc.options...)
			require.NoError(t, err)
			f := tmpl.ExecuteTextString
			if tc.html {
				f = tmpl.ExecuteHTMLString
			}
			got, err := f(tc.in, tc.data)
			if tc.fail {
				require.NotNil(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.exp, got)
		})
	}
}

// This test asserts that template functions are thread-safe.
func TestTemplateFuncs(t *testing.T) {
	tmpl, err := FromGlobs([]string{})
	require.NoError(t, err)

	for _, tc := range []struct {
		title string
		in    string
		data  interface{}
		exp   string
	}{{
		title: "Template using toUpper",
		in:    `{{ "abc" | toUpper }}`,
		exp:   "ABC",
	}, {
		title: "Template using toLower",
		in:    `{{ "ABC" | toLower }}`,
		exp:   "abc",
	}, {
		title: "Template using title",
		in:    `{{ "abc" | title }}`,
		exp:   "Abc",
	}, {
		title: "Template using trimSpace",
		in:    `{{ " abc " | trimSpace }}`,
		exp:   "abc",
	}, {
		title: "Template using join",
		in:    `{{ . | join "," }}`,
		data:  []string{"abc", "def"},
		exp:   "abc,def",
	}, {
		title: "Template using match",
		in:    `{{ match "[a-z]+" "abc" }}`,
		exp:   "true",
	}, {
		title: "Template using reReplaceAll",
		in:    `{{ reReplaceAll "ab" "AB" "abc" }}`,
		exp:   "ABc",
	}} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			wg := sync.WaitGroup{}
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					got, err := tmpl.ExecuteTextString(tc.in, tc.data)
					require.NoError(t, err)
					require.Equal(t, tc.exp, got)
				}()
			}
			wg.Wait()
		})
	}
}

func TestTemplateDataJSONFormat(t *testing.T) {
	data := Data{
		Receiver: "some-receiver",
		Status:   "resolved",
		Alerts: Alerts{
			Alert{
				Status:       "resolved",
				Labels:       KV{"testAlertLabel": "testAlertLabelValue"},
				Annotations:  KV{"testAlertAnnotation": "testAlertAnnotationValue"},
				StartsAt:     Time{time.Date(2023, 11, 14, 15, 49, 0, 0, time.UTC)},
				EndsAt:       Time{time.Date(2023, 11, 14, 18, 22, 36, 2000000, time.UTC)},
				GeneratorURL: "localhost:9090",
				Fingerprint:  "someFingerprint",
			},
		},
		GroupLabels:       KV{"testGroupLabel": "testGroupLabelValue"},
		CommonLabels:      KV{"testCommonLabel": "testCommonLabelValue"},
		CommonAnnotations: KV{"testCommonAnnotation": "testCommonAnnotationValue"},
		ExternalURL:       "localhost:9090",
	}

	jsonBytes, err := json.Marshal(data)

	require.NoError(t, err)
	require.JSONEq(
		t,
		`{
			"receiver":"some-receiver",
			"status":"resolved",
			"alerts":[
				{
					"status":"resolved",
					"labels":{"testAlertLabel":"testAlertLabelValue"},
					"annotations":{"testAlertAnnotation":"testAlertAnnotationValue"},
					"startsAt":"2023-11-14T15:49:00.000Z",
					"endsAt":"2023-11-14T18:22:36.002Z",
					"generatorURL":"localhost:9090",
					"fingerprint":"someFingerprint"
				}
			],
			"groupLabels":{"testGroupLabel":"testGroupLabelValue"},
			"commonLabels":{"testCommonLabel":"testCommonLabelValue"},
			"commonAnnotations":{"testCommonAnnotation":"testCommonAnnotationValue"},
			"externalURL":"localhost:9090"
		}`,
		string(jsonBytes),
	)
}
