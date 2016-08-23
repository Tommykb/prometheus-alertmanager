package mesh

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
	"github.com/satori/go.uuid"
	"github.com/weaveworks/mesh"
)

func TestReplaceFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "replace_file")
	if err != nil {
		t.Fatal(err)
	}
	origFilename := filepath.Join(dir, "testfile")

	of, err := os.Create(origFilename)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := openReplace(filepath.Join(dir, "testfile"))
	if err != nil {
		t.Fatalf("Creating test file failed: %s", err)
	}
	if _, err := nf.Write([]byte("test")); err != nil {
		t.Fatalf("Writing replace file failed: %s", err)
	}

	if nf.Name() == of.Name() {
		t.Fatalf("Replacement file must not have same name as original")
	}
	if err := nf.Close(); err != nil {
		t.Fatalf("Closing replace file failed: %s", err)
	}
	of.Close()

	ofr, err := os.Open(origFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer ofr.Close()

	res, err := ioutil.ReadAll(ofr)
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != "test" {
		t.Fatalf("File contents do not match; got %q, expected %q", string(res), "test")
	}
}

func TestNotificationInfosOnGossip(t *testing.T) {
	var (
		now = utcNow()
	)
	cases := []struct {
		initial map[notificationKey]notificationEntry
		msg     map[notificationKey]notificationEntry
		delta   map[notificationKey]notificationEntry
		final   map[notificationKey]notificationEntry
	}{
		{
			initial: map[notificationKey]notificationEntry{},
			msg: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now, time.Time{}},
			},
			delta: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now, time.Time{}},
			},
			final: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now, time.Time{}},
			},
		}, {
			initial: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now, time.Time{}},
			},
			msg: map[notificationKey]notificationEntry{
				{"recv1", 123}: {false, now.Add(time.Minute), time.Time{}},
			},
			delta: map[notificationKey]notificationEntry{
				{"recv1", 123}: {false, now.Add(time.Minute), time.Time{}},
			},
			final: map[notificationKey]notificationEntry{
				{"recv1", 123}: {false, now.Add(time.Minute), time.Time{}},
			},
		}, {
			initial: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now.Add(time.Minute), time.Time{}},
			},
			msg: map[notificationKey]notificationEntry{
				{"recv1", 123}: {false, now, time.Time{}},
			},
			delta: map[notificationKey]notificationEntry{},
			final: map[notificationKey]notificationEntry{
				{"recv1", 123}: {true, now.Add(time.Minute), time.Time{}},
			},
		},
	}

	for _, c := range cases {
		ni, err := NewNotificationInfos(log.Base(), time.Hour, "")
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range c.initial {
			ni.st.set[k] = v
		}

		b, err := encodeNotificationSet(c.msg)
		if err != nil {
			t.Fatal(err)
		}
		// OnGossip expects the delta but an empty set to be replaced with nil.
		d, err := ni.OnGossip(b)
		if err != nil {
			t.Errorf("%v OnGossip %v: %s", c.initial, c.msg, err)
			continue
		}
		want := c.final
		if have := ni.st.set; !reflect.DeepEqual(want, have) {
			t.Errorf("%v OnGossip %v: want %v, have %v", c.initial, c.msg, want, have)
		}

		want = c.delta
		if len(c.delta) == 0 {
			want = nil
		}
		if d != nil {
			if have := d.(*notificationState).set; !reflect.DeepEqual(want, have) {
				t.Errorf("%v OnGossip %v: want %v, have %v", c.initial, c.msg, want, have)
			}
		} else if want != nil {
			t.Errorf("%v OnGossip %v: want nil", c.initial, c.msg)
		}
	}

	for _, c := range cases {
		ni, err := NewNotificationInfos(log.Base(), time.Hour, "")
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range c.initial {
			ni.st.set[k] = v
		}

		b, err := encodeNotificationSet(c.msg)
		if err != nil {
			t.Fatal(err)
		}
		// OnGossipBroadcast expects the provided delta as is.
		d, err := ni.OnGossipBroadcast(mesh.UnknownPeerName, b)
		if err != nil {
			t.Errorf("%v OnGossipBroadcast %v: %s", c.initial, c.msg, err)
			continue
		}
		want := c.final
		if have := ni.st.set; !reflect.DeepEqual(want, have) {
			t.Errorf("%v OnGossip %v: want %v, have %v", c.initial, c.msg, want, have)
		}

		want = c.delta
		if have := d.(*notificationState).set; !reflect.DeepEqual(want, have) {
			t.Errorf("%v OnGossipBroadcast %v: want %v, have %v", c.initial, c.msg, want, have)
		}
	}

	for _, c := range cases {
		ni, err := NewNotificationInfos(log.Base(), time.Hour, "")
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range c.initial {
			ni.st.set[k] = v
		}

		b, err := encodeNotificationSet(c.msg)
		if err != nil {
			t.Fatal(err)
		}
		// OnGossipUnicast always expects the full state back.
		err = ni.OnGossipUnicast(mesh.UnknownPeerName, b)
		if err != nil {
			t.Errorf("%v OnGossip %v: %s", c.initial, c.msg, err)
			continue
		}

		want := c.final
		if have := ni.st.set; !reflect.DeepEqual(want, have) {
			t.Errorf("%v OnGossip %v: want %v, have %v", c.initial, c.msg, want, have)
		}
	}
}

func TestNotificationInfosSet(t *testing.T) {
	var (
		now       = utcNow()
		retention = time.Hour
	)
	cases := []struct {
		initial map[notificationKey]notificationEntry
		input   []*types.NotificationInfo
		update  map[notificationKey]notificationEntry
		final   map[notificationKey]notificationEntry
	}{
		{
			initial: map[notificationKey]notificationEntry{},
			input: []*types.NotificationInfo{
				{
					Alert:     0x10,
					Receiver:  "recv1",
					Resolved:  false,
					Timestamp: now,
				},
			},
			update: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {false, now, now.Add(retention)},
			},
			final: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {false, now, now.Add(retention)},
			},
		},
		{
			// In this testcase we the second input update is already state
			// respective to the current state. We currently do not prune it
			// from the update as it's not a common occurrence.
			// The update is okay to propagate but the final state must correctly
			// drop it.
			initial: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {false, now, now.Add(retention)},
				{"recv2", 0x10}: {false, now.Add(10 * time.Minute), now.Add(retention).Add(10 * time.Minute)},
			},
			input: []*types.NotificationInfo{
				{
					Alert:     0x10,
					Receiver:  "recv1",
					Resolved:  true,
					Timestamp: now.Add(10 * time.Minute),
				},
				{
					Alert:     0x10,
					Receiver:  "recv2",
					Resolved:  true,
					Timestamp: now,
				},
				{
					Alert:     0x20,
					Receiver:  "recv2",
					Resolved:  false,
					Timestamp: now,
				},
			},
			update: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {true, now.Add(10 * time.Minute), now.Add(retention).Add(10 * time.Minute)},
				{"recv2", 0x10}: {true, now, now.Add(retention)},
				{"recv2", 0x20}: {false, now, now.Add(retention)},
			},
			final: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {true, now.Add(10 * time.Minute), now.Add(retention).Add(10 * time.Minute)},
				{"recv2", 0x10}: {false, now.Add(10 * time.Minute), now.Add(retention).Add(10 * time.Minute)},
				{"recv2", 0x20}: {false, now, now.Add(retention)},
			},
		},
	}

	for _, c := range cases {
		ni, err := NewNotificationInfos(log.Base(), retention, "")
		if err != nil {
			t.Fatal(err)
		}
		tg := &testGossip{}
		ni.Register(tg)
		ni.st = &notificationState{set: c.initial}

		if err := ni.Set(c.input...); err != nil {
			t.Errorf("Insert failed: %s", err)
			continue
		}
		// Verify the correct state afterwards.
		if have := ni.st.set; !reflect.DeepEqual(have, c.final) {
			t.Errorf("Wrong final state %v, expected %v", have, c.final)
			continue
		}

		// Verify that we gossiped the correct update.
		if have := tg.updates[0].(*notificationState).set; !reflect.DeepEqual(have, c.update) {
			t.Errorf("Wrong gossip update %v, expected %v", have, c.update)
			continue
		}
	}
}

func TestNotificationInfosGet(t *testing.T) {
	var (
		now       = utcNow()
		retention = time.Hour
	)
	type query struct {
		recv string
		fps  []model.Fingerprint
		want []*types.NotificationInfo
	}
	cases := []struct {
		state   map[notificationKey]notificationEntry
		queries []query
	}{
		{
			state: map[notificationKey]notificationEntry{
				{"recv1", 0x10}: {true, now.Add(time.Minute), now.Add(retention)},
				{"recv1", 0x30}: {true, now.Add(time.Minute), now.Add(retention)},
				{"recv2", 0x10}: {false, now.Add(time.Minute), now.Add(retention)},
				{"recv2", 0x20}: {false, now, now.Add(retention)},
				// Expired results must be filtered.
				{"recv1", 0x30}: {false, now.Add(2 * retention), now.Add(-retention)},
			},
			queries: []query{
				{
					recv: "recv1",
					fps:  []model.Fingerprint{0x1000, 0x10, 0x20, 0x30},
					want: []*types.NotificationInfo{
						nil,
						{
							Alert:     0x10,
							Receiver:  "recv1",
							Resolved:  true,
							Timestamp: now.Add(time.Minute),
						},
						nil,
						nil,
					},
				},
				{
					recv: "unknown",
					fps:  []model.Fingerprint{0x10, 0x1000},
					want: []*types.NotificationInfo{nil, nil},
				},
			},
		},
	}
	for _, c := range cases {
		ni, err := NewNotificationInfos(log.Base(), retention, "")
		if err != nil {
			t.Fatal(err)
		}
		ni.st = &notificationState{
			set: c.state,
			now: func() time.Time { return now },
		}
		for _, q := range c.queries {
			have, err := ni.Get(q.recv, q.fps...)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(have, q.want) {
				t.Errorf("%v %v expected result %v, got %v", q.recv, q.fps, q.want, have)
			}
		}
	}
}

func TestSilencesSet(t *testing.T) {
	var (
		now      = utcNow()
		id1      = uuid.NewV4()
		matchers = types.NewMatchers(types.NewMatcher("a", "b"))
	)
	cases := []struct {
		input  *types.Silence
		update map[uuid.UUID]*types.Silence
		fail   bool
	}{
		{
			// Set an invalid silence.
			input: &types.Silence{},
			fail:  true,
		},
		{
			// Set a silence including ID.
			input: &types.Silence{
				ID:        id1,
				Matchers:  matchers,
				StartsAt:  now.Add(time.Minute),
				EndsAt:    now.Add(time.Hour),
				CreatedBy: "x",
				Comment:   "x",
			},
			update: map[uuid.UUID]*types.Silence{
				id1: &types.Silence{
					ID:        id1,
					Matchers:  matchers,
					StartsAt:  now.Add(time.Minute),
					EndsAt:    now.Add(time.Hour),
					UpdatedAt: now,
					CreatedBy: "x",
					Comment:   "x",
				},
			},
		},
	}
	for i, c := range cases {
		t.Logf("Test case %d", i)

		s, err := NewSilences(nil, log.Base(), time.Hour, "")
		if err != nil {
			t.Fatal(err)
		}
		tg := &testGossip{}
		s.Register(tg)
		s.st.now = func() time.Time { return now }

		beforeID := c.input.ID

		uid, err := s.Set(c.input)
		if err != nil {
			if c.fail {
				continue
			}
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if c.fail {
			t.Errorf("Expected error but got none")
			continue
		}

		if beforeID != uuid.Nil && uid != beforeID {
			t.Errorf("Silence ID unexpectedly changed: before %q, after %q", beforeID, uid)
			continue
		}

		// Verify the update propagated.
		if have := tg.updates[0].(*silenceState).m; !reflect.DeepEqual(have, c.update) {
			t.Errorf("Update did not match")
			t.Errorf("%s", pretty.Compare(have, c.update))
		}
	}
}

func TestSilencesQuery(t *testing.T) {
	var (
		t0       = time.Now()
		silences = testNewSilences(t, t0)
	)

	n := 500
	insert := make([]*types.Silence, n)
	for i := 0; i < n; i++ {
		insert[i] = createNewSilence(t, silences, t0, i)
	}

	pairs := []queryPair{
		queryPair{
			n:      10,
			offset: 0,
		},
		queryPair{
			n:      10,
			offset: 2,
		},
		queryPair{
			n:      100,
			offset: 4,
		},
	}

	for _, p := range pairs {
		res, err := silences.Query(p.n, p.offset)
		if err != nil {
			t.Fatalf("Retrieval failed: %s", err)
		}

		s := res.Silences

		start := p.offset * p.n
		end := start + p.n
		if end > n {
			t.Fatalf("your test data doesn't include the range you're requesting: insert[%d:%d] (max index %d)", start, end, n)
		}
		expected := append([]*types.Silence{}, insert[start:end]...)

		if len(s) != p.n {
			t.Fatalf("incorrect number of silences returned: wanted %d, got %d", p.n, len(s))
		}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("incorrect silences returned\n")
			t.Fatalf(pretty.Compare(s, expected))
		}
	}
}

func TestSilencesQueryExceedsAvailable(t *testing.T) {
	var (
		t0       = time.Now()
		silences = testNewSilences(t, t0)
	)

	n := 50
	insert := make([]*types.Silence, n)
	for i := 0; i < n; i++ {
		insert[i] = createNewSilence(t, silences, t0, i)
	}

	res, err := silences.Query(n*2, 0)
	if err != nil {
		t.Fatalf("Retrieval failed: %s", err)
	}

	if len(res.Silences) != n {
		t.Fatalf("incorrect silences length: wanted %d, got %d", n, len(res.Silences))
	}
}

func TestSilencesQueryOffsetOutOfBounds(t *testing.T) {
	var (
		t0       = time.Now()
		silences = testNewSilences(t, t0)
	)

	n := 50
	insert := make([]*types.Silence, n)
	for i := 0; i < n; i++ {
		insert[i] = createNewSilence(t, silences, t0, i)
	}

	_, err := silences.Query(n*2, 20)
	if err != types.ErrRequestExceedsAvailable {
		t.Fatalf("expected error, got none")
	}
}

func testNewSilences(t *testing.T, t0 time.Time) *Silences {
	silences, err := NewSilences(nil, log.Base(), time.Hour, "")
	if err != nil {
		t.Fatal(err)
	}

	tg := &testGossip{}
	silences.Register(tg)
	silences.st.now = func() time.Time { return t0 }

	return silences
}

type queryPair struct {
	n, offset int
}

func createNewSilence(t *testing.T, s *Silences, t0 time.Time, i int) *types.Silence {
	sil := &types.Silence{
		Matchers:  types.NewMatchers(types.NewMatcher("a", "b")),
		StartsAt:  t0.Add(time.Duration(i) * time.Minute),
		EndsAt:    t0.Add((time.Duration(i) + 1) * time.Minute),
		CreatedBy: "user",
		Comment:   "another test comment",
	}
	uid, err := s.Set(sil)
	if err != nil {
		t.Fatalf("Insert failed: %s", err)
	}
	sil.ID = uid
	return sil
}

// testGossip implements the mesh.Gossip interface. Received broadcast
// updates are appended to a list.
type testGossip struct {
	updates []mesh.GossipData
}

func (g *testGossip) GossipUnicast(dst mesh.PeerName, msg []byte) error {
	panic("not implemented")
}

func (g *testGossip) GossipBroadcast(update mesh.GossipData) {
	g.updates = append(g.updates, update)
}
