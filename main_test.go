package main

import (
	"fmt"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

func TestCompareSummary(t *testing.T) {
	want := 1

	c1 := ics.NewCalendar()
	c1.SetMethod(ics.MethodRequest)
	NewEvent(c1, "Summary", time.Now(), time.Now())
	NewEvent(c1, "1234567", time.Now(), time.Now())

	c2 := ics.NewCalendar()
	c2.SetMethod(ics.MethodRequest)
	NewEvent(c2, "Summary", time.Now(), time.Now())
	NewEvent(c2, "abcdefg", time.Now(), time.Now())

	diff, err := CompareEvents(c1, c2)
	if len(diff) != want || err != nil {
		t.Fatalf(`diff = %d, %v, want %d, error`, len(diff), err, want)
	}
}

func NewEvent(calendar *ics.Calendar, summary string, start time.Time, end time.Time) {
	e := calendar.AddEvent(fmt.Sprintf("%s@domain.com", uuid.New().String()[:26]))
	e.SetCreatedTime(time.Now())
	e.SetDtStampTime(time.Now())
	e.SetModifiedAt(time.Now())
	e.SetStartAt(start)
	e.SetEndAt(end)
	e.SetSummary(summary)
}

func TestCompareStart(t *testing.T) {
	want := 1

	c1 := ics.NewCalendar()
	c1.SetMethod(ics.MethodRequest)
	NewEvent(c1, "Summary", time.Now(), time.Now())

	c2 := ics.NewCalendar()
	c2.SetMethod(ics.MethodRequest)
	NewEvent(c2, "Summary", time.Now().Add(time.Hour), time.Now())

	diff, err := CompareEvents(c1, c2)
	if len(diff) != want || err != nil {
		t.Fatalf(`diff = %d, %v, want %d, error`, len(diff), err, want)
	}
}

func TestCompareEnd(t *testing.T) {
	want := 1

	c1 := ics.NewCalendar()
	c1.SetMethod(ics.MethodRequest)
	NewEvent(c1, "Summary", time.Now(), time.Now())

	c2 := ics.NewCalendar()
	c2.SetMethod(ics.MethodRequest)
	NewEvent(c2, "Summary", time.Now(), time.Now().Add(time.Hour))

	diff, err := CompareEvents(c1, c2)
	if len(diff) != want || err != nil {
		t.Fatalf(`diff = %d, %v, want %d, error`, len(diff), err, want)
	}
}
