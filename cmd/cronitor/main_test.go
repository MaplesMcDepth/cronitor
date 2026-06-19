package main

import (
	"strings"
	"testing"
	"time"
)

func TestValidateField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		field string
		min   int
		max   int
		want  bool
	}{
		{name: "wildcard", field: "*", min: 0, max: 59, want: true},
		{name: "step", field: "*/15", min: 0, max: 59, want: true},
		{name: "range", field: "1-5", min: 1, max: 31, want: true},
		{name: "list", field: "1,15,30", min: 0, max: 59, want: true},
		{name: "single value", field: "9", min: 0, max: 23, want: true},
		{name: "out of range", field: "61", min: 0, max: 59, want: false},
		{name: "invalid range end", field: "1-99", min: 1, max: 31, want: false},
		{name: "invalid list item", field: "1,a", min: 0, max: 59, want: false},
		{name: "invalid step", field: "*/x", min: 0, max: 59, want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := validateField(tt.field, tt.min, tt.max); got != tt.want {
				t.Fatalf("validateField(%q, %d, %d) = %v, want %v", tt.field, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestExplainField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		field    string
		singular string
		plural   string
		want     string
	}{
		{name: "wildcard", field: "*", singular: "minute", plural: "minutes", want: ""},
		{name: "step", field: "*/5", singular: "minute", plural: "minutes", want: "every 5 minutes"},
		{name: "list", field: "1,2,3", singular: "minute", plural: "minutes", want: "at 1,2,3 minutes"},
		{name: "range", field: "1-5", singular: "day", plural: "days", want: "from 1-5"},
		{name: "single value", field: "9", singular: "hour", plural: "hours", want: "9 hour"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := explainField(tt.field, tt.singular, tt.plural); got != tt.want {
				t.Fatalf("explainField(%q, %q, %q) = %q, want %q", tt.field, tt.singular, tt.plural, got, tt.want)
			}
		})
	}
}

func TestMatchField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		field string
		value int
		want  bool
	}{
		{name: "wildcard", field: "*", value: 17, want: true},
		{name: "step match", field: "*/15", value: 30, want: true},
		{name: "step miss", field: "*/15", value: 31, want: false},
		{name: "range match", field: "9-17", value: 12, want: true},
		{name: "range miss", field: "9-17", value: 20, want: false},
		{name: "list match", field: "1,5,10", value: 5, want: true},
		{name: "list miss", field: "1,5,10", value: 2, want: false},
		{name: "single value", field: "7", value: 7, want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := matchField(tt.field, tt.value, 0, 59); got != tt.want {
				t.Fatalf("matchField(%q, %d, 0, 59) = %v, want %v", tt.field, tt.value, got, tt.want)
			}
		})
	}
}

func TestFindNext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		expr   string
		after  time.Time
		want   time.Time
	}{
		{
			name:  "next weekday interval",
			expr:  "*/15 9-17 * * 1-5",
			after: time.Date(2026, time.June, 19, 8, 59, 30, 0, time.UTC),
			want:  time.Date(2026, time.June, 19, 9, 0, 0, 0, time.UTC),
		},
		{
			name:  "roll to next day",
			expr:  "0 9 * * *",
			after: time.Date(2026, time.June, 19, 9, 1, 0, 0, time.UTC),
			want:  time.Date(2026, time.June, 20, 9, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := strings.Fields(tt.expr)
			if got := findNext(fields, tt.after); !got.Equal(tt.want) {
				t.Fatalf("findNext(%q, %s) = %s, want %s", tt.expr, tt.after, got, tt.want)
			}
		})
	}
}
