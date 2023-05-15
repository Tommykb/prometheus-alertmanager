package matchers

import (
	"reflect"
	"testing"

	"github.com/prometheus/alertmanager/pkg/labels"
)

func TestMatchers(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  labels.Matchers
		err   string
	}{
		{
			input: `{foo="bar"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				return append(ms, m)
			}(),
		},
		{
			input: `{foo=~"bar.*"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchRegexp, "foo", "bar.*")
				return append(ms, m)
			}(),
		},
		{
			input: `{foo!="bar"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchNotEqual, "foo", "bar")
				return append(ms, m)
			}(),
		},
		{
			input: `{foo!~"bar.*"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchNotRegexp, "foo", "bar.*")
				return append(ms, m)
			}(),
		},
		{
			input: `{foo="bar", baz!="quux"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotEqual, "baz", "quux")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `{foo="bar", baz!~"quux.*"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotRegexp, "baz", "quux.*")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `{foo="bar",baz!~".*quux", derp="wat"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotRegexp, "baz", ".*quux")
				m3, _ := labels.NewMatcher(labels.MatchEqual, "derp", "wat")
				return append(ms, m, m2, m3)
			}(),
		},
		{
			input: `{foo="bar", baz!="quux", derp="wat"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotEqual, "baz", "quux")
				m3, _ := labels.NewMatcher(labels.MatchEqual, "derp", "wat")
				return append(ms, m, m2, m3)
			}(),
		},
		{
			input: `{foo="bar", baz!~".*quux.*", derp="wat"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotRegexp, "baz", ".*quux.*")
				m3, _ := labels.NewMatcher(labels.MatchEqual, "derp", "wat")
				return append(ms, m, m2, m3)
			}(),
		},
		{
			input: `{foo="bar", instance=~"some-api.*"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchRegexp, "instance", "some-api.*")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `{foo=""}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "")
				return append(ms, m)
			}(),
		},
		{
			input: `{foo="bar,quux", job="job1"}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar,quux")
				m2, _ := labels.NewMatcher(labels.MatchEqual, "job", "job1")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `{foo = "bar", dings != "bums", }`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotEqual, "dings", "bums")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `foo=bar,dings!=bums`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "foo", "bar")
				m2, _ := labels.NewMatcher(labels.MatchNotEqual, "dings", "bums")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `{quote="She said: \"Hi, ladies! That's gender-neutral…\""}`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "quote", `She said: "Hi, ladies! That's gender-neutral…"`)
				return append(ms, m)
			}(),
		},
		{
			input: `statuscode=~"5.."`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchRegexp, "statuscode", "5..")
				return append(ms, m)
			}(),
		},
		{
			input: `tricky=~~~`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchRegexp, "tricky", "~~")
				return append(ms, m)
			}(),
		},
		{
			input: `trickier==\\=\=\"`,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchEqual, "trickier", `=\=\="`)
				return append(ms, m)
			}(),
		},
		{
			input: `contains_quote != "\"" , contains_comma !~ "foo,bar" , `,
			want: func() labels.Matchers {
				ms := labels.Matchers{}
				m, _ := labels.NewMatcher(labels.MatchNotEqual, "contains_quote", `"`)
				m2, _ := labels.NewMatcher(labels.MatchNotRegexp, "contains_comma", "foo,bar")
				return append(ms, m, m2)
			}(),
		},
		{
			input: `job=`,
			want: func() labels.Matchers {
				m, _ := labels.NewMatcher(labels.MatchEqual, "job", "")
				return labels.Matchers{m}
			}(),
		},
		{
			input: `job="value`,
			err:   `matcher value contains unescaped double quote: "value`,
		},
		{
			input: `job=value"`,
			err:   `matcher value contains unescaped double quote: value"`,
		},
		{
			input: `trickier==\\=\=\""`,
			err:   `matcher value contains unescaped double quote: =\\=\=\""`,
		},
		{
			input: `contains_unescaped_quote = foo"bar`,
			err:   `matcher value contains unescaped double quote: foo"bar`,
		},
		{
			input: `{invalid-name = "valid label"}`,
			err:   `bad matcher format: invalid-name = "valid label"`,
		},
		{
			input: `{foo=~"invalid[regexp"}`,
			err:   "error parsing regexp: missing closing ]: `[regexp)$`",
		},
		// Double escaped strings.
		{
			input: `"{foo=\"bar"}`,
			err:   `bad matcher format: "{foo=\"bar"`,
		},
		{
			input: `"foo=\"bar"`,
			err:   `bad matcher format: "foo=\"bar"`,
		},
		{
			input: `"foo=\"bar\""`,
			err:   `bad matcher format: "foo=\"bar\""`,
		},
		{
			input: `"foo=\"bar\"`,
			err:   `bad matcher format: "foo=\"bar\"`,
		},
		{
			input: `"{foo=\"bar\"}"`,
			err:   `bad matcher format: "{foo=\"bar\"}"`,
		},
		{
			input: `"foo="bar""`,
			err:   `bad matcher format: "foo="bar""`,
		},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil && tc.err == "" {
				t.Fatalf("got error where none expected: %v", err)
			}
			if err == nil && tc.err != "" {
				t.Fatalf("expected error but got none: %v", tc.err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("labels not equal:\ngot  %v\nwant %v", got, tc.want)
			}
		})
	}
}
