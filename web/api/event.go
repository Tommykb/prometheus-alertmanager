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

package api

import (
	"log"
	"net/http"

	"github.com/prometheus/alert_manager/manager"
)

func (s AlertManagerService) AddEvents(es manager.Events) {
	for i, ev := range es {
		if ev.Summary == "" || ev.Description == "" {
			log.Printf("Missing field in event %d: %s", i, ev)
			rb := s.ResponseBuilder()
			rb.SetResponseCode(http.StatusBadRequest)
			return
		}
		if _, ok := ev.Labels[manager.EventNameLabel]; !ok {
			log.Printf("Missing alert name label in event %d: %s", i, ev)
			rb := s.ResponseBuilder()
			rb.SetResponseCode(http.StatusBadRequest)
			return
		}
	}

	err := s.Aggregator.Receive(es)
	if err != nil {
		log.Println("Error during aggregation:", err)
		rb := s.ResponseBuilder()
		rb.SetResponseCode(http.StatusServiceUnavailable)
	}
}
