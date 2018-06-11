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

package dispatch

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"golang.org/x/net/context"

	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/provider"
	"github.com/prometheus/alertmanager/types"
)

// Dispatcher sorts incoming alerts into aggregation groups and
// assigns the correct notifiers to each.
type Dispatcher struct {
	route  *Route
	alerts provider.Alerts
	stage  notify.Stage

	marker  types.Marker
	timeout func(time.Duration) time.Duration

	aggrGroups map[*Route]map[model.Fingerprint]*aggrGroup
	mtx        sync.RWMutex

	done   chan struct{}
	ctx    context.Context
	cancel func()

	logger log.Logger
}

// NewDispatcher returns a new Dispatcher.
func NewDispatcher(
	ap provider.Alerts,
	r *Route,
	s notify.Stage,
	mk types.Marker,
	to func(time.Duration) time.Duration,
	l log.Logger,
) *Dispatcher {
	disp := &Dispatcher{
		alerts:  ap,
		stage:   s,
		route:   r,
		marker:  mk,
		timeout: to,
		logger:  log.With(l, "component", "dispatcher"),
	}
	return disp
}

// Run starts dispatching alerts incoming via the updates channel.
func (d *Dispatcher) Run() {
	d.done = make(chan struct{})

	d.mtx.Lock()
	d.aggrGroups = map[*Route]map[model.Fingerprint]*aggrGroup{}
	d.mtx.Unlock()

	d.ctx, d.cancel = context.WithCancel(context.Background())

	d.run(d.alerts.Subscribe())
	close(d.done)
}

// AlertBlock contains a list of alerts associated with a set of
// routing options.
type AlertBlock struct {
	RouteOpts *RouteOpts  `json:"routeOpts"`
	Alerts    []*APIAlert `json:"alerts"`
}

// APIAlert is the API representation of an alert, which is a regular alert
// annotated with silencing and inhibition info.
type APIAlert struct {
	*model.Alert
	Status      types.AlertStatus `json:"status"`
	Receivers   []string          `json:"receivers"`
	Fingerprint string            `json:"fingerprint"`
}

// AlertGroup is a list of alert blocks grouped by the same label set.
type AlertGroup struct {
	Labels   model.LabelSet `json:"labels"`
	GroupKey string         `json:"groupKey"`
	Blocks   []*AlertBlock  `json:"blocks"`
}

// AlertOverview is a representation of all active alerts in the system.
type AlertOverview []*AlertGroup

func (ao AlertOverview) Swap(i, j int)      { ao[i], ao[j] = ao[j], ao[i] }
func (ao AlertOverview) Less(i, j int) bool { return ao[i].Labels.Before(ao[j].Labels) }
func (ao AlertOverview) Len() int           { return len(ao) }

func matchesFilterLabels(a *APIAlert, matchers []*labels.Matcher) bool {
	for _, m := range matchers {
		if v, prs := a.Labels[model.LabelName(m.Name)]; !prs || !m.Matches(string(v)) {
			return false
		}
	}

	return true
}

// Groups populates an AlertOverview from the dispatcher's internal state.
func (d *Dispatcher) Groups(matchers []*labels.Matcher) AlertOverview {
	overview := AlertOverview{}

	d.mtx.RLock()
	defer d.mtx.RUnlock()

	seen := map[model.Fingerprint]*AlertGroup{}

	for route, ags := range d.aggrGroups {
		for _, ag := range ags {
			alertGroup, ok := seen[ag.fingerprint()]
			if !ok {
				alertGroup = &AlertGroup{Labels: ag.labels}
				alertGroup.GroupKey = ag.GroupKey()

				seen[ag.fingerprint()] = alertGroup
			}

			now := time.Now()

			var apiAlerts []*APIAlert
			for _, a := range types.Alerts(ag.alertSlice()...) {
				if !a.EndsAt.IsZero() && a.EndsAt.Before(now) {
					continue
				}
				status := d.marker.Status(a.Fingerprint())
				aa := &APIAlert{
					Alert:       a,
					Status:      status,
					Fingerprint: a.Fingerprint().String(),
				}

				if !matchesFilterLabels(aa, matchers) {
					continue
				}

				apiAlerts = append(apiAlerts, aa)
			}
			if len(apiAlerts) == 0 {
				continue
			}

			alertGroup.Blocks = append(alertGroup.Blocks, &AlertBlock{
				RouteOpts: &route.RouteOpts,
				Alerts:    apiAlerts,
			})

			overview = append(overview, alertGroup)
		}
	}

	sort.Sort(overview)

	return overview
}

func (d *Dispatcher) run(it provider.AlertIterator) {
	cleanup := time.NewTicker(30 * time.Second)
	defer cleanup.Stop()

	defer it.Close()

	for {
		select {
		case alert, ok := <-it.Next():
			if !ok {
				// Iterator exhausted for some reason.
				if err := it.Err(); err != nil {
					level.Error(d.logger).Log("msg", "Error on alert update", "err", err)
				}
				return
			}

			level.Debug(d.logger).Log("msg", "Received alert", "alert", alert)

			// Log errors but keep trying.
			if err := it.Err(); err != nil {
				level.Error(d.logger).Log("msg", "Error on alert update", "err", err)
				continue
			}

			for _, r := range d.route.Match(alert.Labels) {
				d.processAlert(alert, r)
			}

		case <-cleanup.C:
			d.mtx.Lock()

			for _, groups := range d.aggrGroups {
				for _, ag := range groups {
					if ag.empty() {
						ag.stop()
						delete(groups, ag.fingerprint())
					}
				}
			}

			d.mtx.Unlock()

		case <-d.ctx.Done():
			return
		}
	}
}

// Stop the dispatcher.
func (d *Dispatcher) Stop() {
	if d == nil || d.cancel == nil {
		return
	}
	d.cancel()
	d.cancel = nil

	<-d.done
}

// notifyFunc is a function that performs notifcation for the alert
// with the given fingerprint. It aborts on context cancelation.
// Returns false iff notifying failed.
type notifyFunc func(context.Context, ...*types.Alert) bool

// processAlert determines in which aggregation group the alert falls
// and inserts it.
func (d *Dispatcher) processAlert(alert *types.Alert, route *Route) {
	groupLabels := model.LabelSet{}

	for ln, lv := range alert.Labels {
		if _, ok := route.RouteOpts.GroupBy[ln]; ok {
			groupLabels[ln] = lv
		}
	}

	fp := groupLabels.Fingerprint()

	d.mtx.Lock()
	group, ok := d.aggrGroups[route]
	if !ok {
		group = map[model.Fingerprint]*aggrGroup{}
		d.aggrGroups[route] = group
	}
	d.mtx.Unlock()

	// If the group does not exist, create it.
	ag, ok := group[fp]
	if !ok {
		ag = newAggrGroup(d.ctx, groupLabels, route, d.timeout, d.logger)
		group[fp] = ag

		go ag.run(func(ctx context.Context, alerts ...*types.Alert) bool {
			_, _, err := d.stage.Exec(ctx, d.logger, alerts...)
			if err != nil {
				level.Error(d.logger).Log("msg", "Notify for alerts failed", "num_alerts", len(alerts), "err", err)
			}
			return err == nil
		})
	}

	ag.insert(alert)
}

// aggrGroup aggregates alert fingerprints into groups to which a
// common set of routing options applies.
// It emits notifications in the specified intervals.
type aggrGroup struct {
	labels   model.LabelSet
	opts     *RouteOpts
	logger   log.Logger
	routeKey string

	ctx     context.Context
	cancel  func()
	done    chan struct{}
	next    *time.Timer
	timeout func(time.Duration) time.Duration

	mtx        sync.RWMutex
	alerts     map[model.Fingerprint]*types.Alert
	hasFlushed bool
}

// newAggrGroup returns a new aggregation group.
func newAggrGroup(ctx context.Context, labels model.LabelSet, r *Route, to func(time.Duration) time.Duration, logger log.Logger) *aggrGroup {
	if to == nil {
		to = func(d time.Duration) time.Duration { return d }
	}
	ag := &aggrGroup{
		labels:   labels,
		routeKey: r.Key(),
		opts:     &r.RouteOpts,
		timeout:  to,
		alerts:   map[model.Fingerprint]*types.Alert{},
	}
	ag.ctx, ag.cancel = context.WithCancel(ctx)

	ag.logger = log.With(logger, "aggrGroup", ag)

	// Set an initial one-time wait before flushing
	// the first batch of notifications.
	ag.next = time.NewTimer(ag.opts.GroupWait)

	return ag
}

func (ag *aggrGroup) fingerprint() model.Fingerprint {
	return ag.labels.Fingerprint()
}

func (ag *aggrGroup) GroupKey() string {
	return fmt.Sprintf("%s:%s", ag.routeKey, ag.labels)
}

func (ag *aggrGroup) String() string {
	return ag.GroupKey()
}

func (ag *aggrGroup) alertSlice() []*types.Alert {
	ag.mtx.RLock()
	defer ag.mtx.RUnlock()

	var alerts []*types.Alert
	for _, a := range ag.alerts {
		alerts = append(alerts, a)
	}
	return alerts
}

func (ag *aggrGroup) run(nf notifyFunc) {
	ag.done = make(chan struct{})

	defer close(ag.done)
	defer ag.next.Stop()

	for {
		select {
		case now := <-ag.next.C:
			// Give the notifications time until the next flush to
			// finish before terminating them.
			ctx, cancel := context.WithTimeout(ag.ctx, ag.timeout(ag.opts.GroupInterval))

			// The now time we retrieve from the ticker is the only reliable
			// point of time reference for the subsequent notification pipeline.
			// Calculating the current time directly is prone to flaky behavior,
			// which usually only becomes apparent in tests.
			ctx = notify.WithNow(ctx, now)

			// Populate context with information needed along the pipeline.
			ctx = notify.WithGroupKey(ctx, ag.GroupKey())
			ctx = notify.WithGroupLabels(ctx, ag.labels)
			ctx = notify.WithReceiverName(ctx, ag.opts.Receiver)
			ctx = notify.WithRepeatInterval(ctx, ag.opts.RepeatInterval)

			// Wait the configured interval before calling flush again.
			ag.mtx.Lock()
			ag.next.Reset(ag.opts.GroupInterval)
			ag.hasFlushed = true
			ag.mtx.Unlock()

			ag.flush(func(alerts ...*types.Alert) bool {
				return nf(ctx, alerts...)
			})

			cancel()

		case <-ag.ctx.Done():
			return
		}
	}
}

func (ag *aggrGroup) stop() {
	// Calling cancel will terminate all in-process notifications
	// and the run() loop.
	ag.cancel()
	<-ag.done
}

// insert inserts the alert into the aggregation group.
func (ag *aggrGroup) insert(alert *types.Alert) {
	ag.mtx.Lock()
	defer ag.mtx.Unlock()

	ag.alerts[alert.Fingerprint()] = alert

	// Immediately trigger a flush if the wait duration for this
	// alert is already over.
	if !ag.hasFlushed && alert.StartsAt.Add(ag.opts.GroupWait).Before(time.Now()) {
		ag.next.Reset(0)
	}
}

func (ag *aggrGroup) empty() bool {
	ag.mtx.RLock()
	defer ag.mtx.RUnlock()

	return len(ag.alerts) == 0
}

// flush sends notifications for all new alerts.
func (ag *aggrGroup) flush(notify func(...*types.Alert) bool) {
	if ag.empty() {
		return
	}
	ag.mtx.Lock()

	var (
		alerts      = make(map[model.Fingerprint]*types.Alert, len(ag.alerts))
		alertsSlice = make([]*types.Alert, 0, len(ag.alerts))
	)
	for fp, alert := range ag.alerts {
		alerts[fp] = alert
		alertsSlice = append(alertsSlice, alert)
	}

	sort.SliceStable(alertsSlice, func(i, j int) bool {
		// Look at labels.job, then labels.instance.
		for _, override_key := range [...]model.LabelName{"job", "instance"} {
			key_i, ok_i := alertsSlice[i].Labels[override_key]
			if !ok_i {
				return false
			}
			key_j, ok_j := alertsSlice[j].Labels[override_key]
			if !ok_j {
				return true
			}

			if key_i != key_j {
				return key_i < key_j
			}
		}

		return alertsSlice[i].Labels.Before(alertsSlice[j].Labels)
	})

	ag.mtx.Unlock()

	level.Debug(ag.logger).Log("msg", "Flushing", "alerts", fmt.Sprintf("%v", alertsSlice))

	if notify(alertsSlice...) {
		ag.mtx.Lock()
		for fp, a := range alerts {
			// Only delete if the fingerprint has not been inserted
			// again since we notified about it.
			if a.Resolved() && ag.alerts[fp] == a {
				delete(ag.alerts, fp)
			}
		}
		ag.mtx.Unlock()
	}
}
