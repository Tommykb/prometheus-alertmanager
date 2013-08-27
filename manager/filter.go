// Copyright 2013 Prometheus Team
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

package manager

import (
	"fmt"
	"hash/fnv"
	"regexp"
)

type Filters []*Filter

type Filter struct {
	Name         *regexp.Regexp
	Value        *regexp.Regexp
	NamePattern  string
	ValuePattern string

	fingerprint uint64
}

func NewFilter(namePattern string, valuePattern string) *Filter {
	summer := fnv.New64a()
	fmt.Fprintf(summer, namePattern, valuePattern)

	return &Filter{
		Name:         regexp.MustCompile("^" + namePattern + "$"),
		Value:        regexp.MustCompile("^" + valuePattern + "$"),
		NamePattern:  namePattern,
		ValuePattern: valuePattern,
		fingerprint:  summer.Sum64(),
	}
}

func (f *Filter) Handles(l AlertLabelSet) bool {
	for k, v := range l {
		if f.Name.MatchString(k) && f.Value.MatchString(v) {
			return true
		}
	}

	return false
}

func (f Filters) Handles(l AlertLabelSet) bool {
	fCount := len(f)
	fMatch := 0

	for _, filter := range f {
		if filter.Handles(l) {
			fMatch++
		}
	}

	return fCount == fMatch
}

func (f Filters) Filter(l AlertLabelSets) AlertLabelSets {
	out := AlertLabelSets{}
	for _, labels := range l {
		if f.Handles(labels) {
			out = append(out, labels)
		}
	}
	return out
}

func (f Filters) fingerprint() uint64 {
	summer := fnv.New64a()

	for i, f := range f {
		fmt.Fprintln(summer, i, f.fingerprint)
	}

	return summer.Sum64()
}
