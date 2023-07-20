package main

import (
	"testing"
	"time"

	"snippetbox.lets-go/internal/assert"
)

func TestHumanDate(t *testing.T) {

	// slice of anonymous structs containing the test case name
	// input to our test func humanDate(), and the expected value
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		// test empty time or "zero" valued time
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		// test central europe time (CET) as the proper fixed time zone
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	// loop over our tests
	for _, tt := range tests {
		// run a sub-test which is basically running the func + an assertion statement
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}

	// initialize a new time.Time object and pass it to the humanDate function
	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
	hd := humanDate(tm)

	// check that the output from hd func is in the format we expected
	// if it isn't as expected, use the t.Errorf() func to indicate that the test
	// failed
	if hd != "17 Mar 2022 at 10:15" {
		t.Errorf("got %q; wanted %q", hd, "17 Mar 2022 at 10:15")
	}

}
