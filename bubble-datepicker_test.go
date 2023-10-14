package datepicker

import (
	"testing"
	"time"
)

var halloween = time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC)
var thanksgiving = time.Date(2023, time.November, 23, 0, 0, 0, 0, time.UTC)
var xmas = time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)

func TestNew(t *testing.T) {
	h := halloween
	if m := New(h); m.Time != h {
		t.Errorf("expected `New` method to return a model with time instance")
	}
}

func TestSetFocus(t *testing.T) {
	tests := []struct {
		input Focus
		want  Focus
	}{
		{input: FocusNone, want: FocusNone},
		{input: FocusCalendar, want: FocusCalendar},
		{input: FocusHeaderMonth, want: FocusHeaderMonth},
		{input: FocusHeaderYear, want: FocusHeaderYear},
	}

	model := New(halloween)
	for i, test := range tests {
		model.SetFocus(test.input)
		if got := model.Focused; test.want != got {
			t.Errorf("TestSetFocus failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestBlur(t *testing.T) {
	tests := []struct {
		input Focus
		want  Focus
	}{
		{input: FocusNone, want: FocusNone},
		{input: FocusCalendar, want: FocusNone},
		{input: FocusHeaderMonth, want: FocusNone},
		{input: FocusHeaderYear, want: FocusNone},
	}

	model := New(halloween)
	for i, test := range tests {
		model.SetFocus(test.input)
		model.Blur()
		if got := model.Focused; test.want != got {
			t.Errorf("TestBlur failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestSetTime(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: halloween},
		{input: thanksgiving, want: thanksgiving},
		{input: xmas, want: xmas},
	}
	model := New(time.Now())
	for i, test := range tests {
		model.SetTime(test.input)
		if got := model.Time; test.want != got {
			t.Errorf("TestSetTime failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestLastWeek(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.October, 24, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2023, time.November, 16, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2023, time.December, 18, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.LastWeek()
		if got := model.Time; test.want != got {
			t.Errorf("TestLastWeek failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestNextWeek(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.November, 7, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2023, time.November, 30, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.NextWeek()
		if got := model.Time; test.want != got {
			t.Errorf("TestNextWeek failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestYesterday(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.October, 30, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2023, time.November, 22, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2023, time.December, 24, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.Yesterday()
		if got := model.Time; test.want != got {
			t.Errorf("TestYesterday failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestTomorrow(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2023, time.November, 24, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2023, time.December, 26, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.Tomorrow()
		if got := model.Time; test.want != got {
			t.Errorf("TestTomorrow failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestLastMonth(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.September, 31, 0, 0, 0, 0, time.UTC)}, // normalizes
		{input: thanksgiving, want: time.Date(2023, time.October, 23, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2023, time.November, 25, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.LastMonth()
		if got := model.Time; test.want != got {
			t.Errorf("TestLastMonth failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestNextMonth(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2023, time.November, 31, 0, 0, 0, 0, time.UTC)}, // normalizes
		{input: thanksgiving, want: time.Date(2023, time.December, 23, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.NextMonth()
		if got := model.Time; test.want != got {
			t.Errorf("TestNextMonth failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestLastYear(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2022, time.October, 31, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2022, time.November, 23, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2022, time.December, 25, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.LastYear()
		if got := model.Time; test.want != got {
			t.Errorf("TestLastYear failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}

func TestNextYear(t *testing.T) {
	tests := []struct {
		input time.Time
		want  time.Time
	}{
		{input: halloween, want: time.Date(2024, time.October, 31, 0, 0, 0, 0, time.UTC)},
		{input: thanksgiving, want: time.Date(2024, time.November, 23, 0, 0, 0, 0, time.UTC)},
		{input: xmas, want: time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)},
	}
	for i, test := range tests {
		model := New(test.input)
		model.NextYear()
		if got := model.Time; test.want != got {
			t.Errorf("TestNextYear failure - index: %d - want: '%s' got: '%s'", i, test.want, got)
		}
	}
}
