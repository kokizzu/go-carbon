package cache

import (
	"testing"

	"github.com/go-graphite/go-carbon/points"
)

func TestInFlight(t *testing.T) {
	var data []points.Point

	c := New()

	c.Add(points.OnePoint("hello.world", 42, 10))

	m1 := c.WriteoutQueue().Get(nil)
	p1, _ := c.PopNotConfirmed(m1)
	if !p1.Eq(points.OnePoint("hello.world", 42, 10)) {
		t.FailNow()
	}

	data = c.Get("hello.world")
	if len(data) != 1 || data[0].Value != 42 {
		t.FailNow()
	}

	c.Add(points.OnePoint("hello.world", 43, 10))

	// 42 in flight, 43 in cache
	data = c.Get("hello.world")
	if len(data) != 2 || data[0].Value != 42 || data[1].Value != 43 {
		t.FailNow()
	}

	m2 := c.WriteoutQueue().Get(nil)
	p2, _ := c.PopNotConfirmed(m2)
	if !p2.Eq(points.OnePoint("hello.world", 43, 10)) {
		t.FailNow()
	}

	// 42, 43 in flight
	data = c.Get("hello.world")
	if len(data) != 2 || data[0].Value != 42 || data[1].Value != 43 {
		t.FailNow()
	}

	c.Confirm(p1)

	c.Add(points.OnePoint("hello.world", 44, 10))
	m3 := c.WriteoutQueue().Get(nil)
	p3, _ := c.PopNotConfirmed(m3)
	if !p3.Eq(points.OnePoint("hello.world", 44, 10)) {
		t.FailNow()
	}

	// 43, 44 in flight
	data = c.Get("hello.world")
	if len(data) != 2 || data[0].Value != 43 || data[1].Value != 44 {
		t.FailNow()
	}
}

func TestRequeueFailedWrite(t *testing.T) {
	c := New()
	if !c.IsEmpty() {
		t.Fatal("new cache is not empty")
	}

	c.Add(points.OnePoint("hello.world", 42, 10))
	failed, ok := c.PopNotConfirmed("hello.world")
	if !ok {
		t.Fatal("failed to pop point")
	}
	c.Add(points.OnePoint("hello.world", 43, 11))
	c.SetThrottle(func(*points.Points, bool) bool { return true })

	c.Requeue(failed)

	data := c.Get("hello.world")
	if len(data) != 2 || data[0].Value != 42 || data[1].Value != 43 {
		t.Fatalf("requeued data = %+v; want failed point before newer point", data)
	}
	if got := c.Size(); got != 2 {
		t.Fatalf("cache size = %d; want 2", got)
	}
	if got := c.NotConfirmedLength(); got != 0 {
		t.Fatalf("not-confirmed batches = %d; want 0", got)
	}

	combined, ok := c.PopNotConfirmed("hello.world")
	if !ok || c.IsEmpty() {
		t.Fatal("in-flight batch must keep the cache non-empty")
	}
	c.Confirm(combined)
	if !c.IsEmpty() {
		t.Fatal("cache is not empty after confirmation")
	}
}

func BenchmarkPopNotConfirmed(b *testing.B) {
	c := New()
	p1 := points.OnePoint("hello.world", 42, 10)
	var p2 *points.Points

	for n := 0; n < b.N; n++ {
		c.Add(p1)
		p2, _ = c.PopNotConfirmed("hello.world")
		c.Confirm(p2)
	}

	if !p1.Eq(p2) {
		b.FailNow()
	}
}

func BenchmarkPopNotConfirmed100(b *testing.B) {
	c := New()

	for i := 0; i < 100; i++ {
		c.Add(points.OnePoint("hello.world", 42, 10))
		c.PopNotConfirmed("hello.world")
	}

	p1 := points.OnePoint("hello.world", 42, 10)
	var p2 *points.Points

	for n := 0; n < b.N; n++ {
		c.Add(p1)
		p2, _ = c.PopNotConfirmed("hello.world")
		c.Confirm(p2)
	}

	if !p1.Eq(p2) {
		b.FailNow()
	}
}

func BenchmarkPop(b *testing.B) {
	c := New()
	p1 := points.OnePoint("hello.world", 42, 10)
	var p2 *points.Points

	for n := 0; n < b.N; n++ {
		c.Add(p1)
		p2, _ = c.Pop("hello.world")
	}

	if !p1.Eq(p2) {
		b.FailNow()
	}
}
func BenchmarkGet(b *testing.B) {
	c := New()
	c.Add(points.OnePoint("hello.world", 42, 10))

	var d []points.Point
	for n := 0; n < b.N; n++ {
		d = c.Get("hello.world")
	}

	if len(d) != 1 {
		b.FailNow()
	}
}

func BenchmarkGetNotConfirmed1(b *testing.B) {
	c := New()

	c.Add(points.OnePoint("hello.world", 42, 10))
	c.PopNotConfirmed("hello.world")

	var d []points.Point
	for n := 0; n < b.N; n++ {
		d = c.Get("hello.world")
	}

	if len(d) != 1 {
		b.FailNow()
	}
}

func BenchmarkGetNotConfirmed100(b *testing.B) {
	c := New()

	for i := 0; i < 100; i++ {
		c.Add(points.OnePoint("hello.world", 42, 10))
		c.PopNotConfirmed("hello.world")
	}

	var d []points.Point
	for n := 0; n < b.N; n++ {
		d = c.Get("hello.world")
	}

	if len(d) != 100 {
		b.FailNow()
	}
}

func BenchmarkGetNotConfirmed100Miss(b *testing.B) {
	c := New()

	for i := 0; i < 100; i++ {
		c.Add(points.OnePoint("hello.world", 42, 10))
		c.PopNotConfirmed("hello.world")
	}

	var d []points.Point
	for n := 0; n < b.N; n++ {
		d = c.Get("metric.name")
	}

	if d != nil {
		b.FailNow()
	}
}
