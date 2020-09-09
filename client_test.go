package bovada

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-check/check"
)

type clientSuite struct{}

func (cs *clientSuite) TestNewClient(c *check.C) {
	var tcs = []struct {
		name string
		in   *http.Client
		exp  *http.Client
	}{
		{
			name: "nil client passed, returns default client",
			in:   nil,
			exp:  http.DefaultClient,
		},
		{
			name: "non-nil client passed, sets input",
			in:   &http.Client{Timeout: 100},
			exp:  &http.Client{Timeout: 100},
		},
	}

	for _, tc := range tcs {
		c.Log(tc.name)

		c.Check(NewClient(tc.in).client, check.DeepEquals, tc.exp)
	}
}

func (cs *clientSuite) TestQueryOpts(c *check.C) {
	// Freeze time for test.
	var d = time.Now()
	now = func() time.Time {
		return d
	}

	// Compute expected query param values.
	var sEOD = strconv.Itoa(minutesLeftInDay(d))
	var sEOT = strconv.Itoa(minutesLeftInDay(d) + minutesInDay)

	var tcs = []struct {
		name string
		opts *queryOpts
		exp  url.Values
	}{
		{
			name: "Default opts.",
			opts: NewQueryOpts(),
			exp: url.Values{
				langKey: []string{"en"},
			},
		},
		{
			name: "Default opts + today.",
			opts: NewQueryOpts().TodayOnly(),
			exp: url.Values{
				langKey:      []string{"en"},
				startTimeKey: []string{sEOD},
			},
		},
		{
			name: "Default opts + tomorrow.",
			opts: NewQueryOpts().TomorrowOnly(),
			exp: url.Values{
				langKey:            []string{"en"},
				startTimeKey:       []string{sEOT},
				startTimeOffsetKey: []string{sEOD},
			},
		},
		{
			name: "Default opts + tomorrow + today (today should overwrite).",
			opts: NewQueryOpts().TomorrowOnly().TodayOnly(),
			exp: url.Values{
				langKey:      []string{"en"},
				startTimeKey: []string{sEOD},
			},
		},
		{
			name: "Default opts + today + tomorrow (tomorrow should overwrite).",
			opts: NewQueryOpts().TodayOnly().TomorrowOnly(),
			exp: url.Values{
				langKey:            []string{"en"},
				startTimeKey:       []string{sEOT},
				startTimeOffsetKey: []string{sEOD},
			},
		},
		{
			name: "Default opts + upcoming true.",
			opts: NewQueryOpts().UpcomingOnly(true),
			exp: url.Values{
				langKey:         []string{"en"},
				upcomingOnlyKey: []string{"true"},
			},
		},
		{
			name: "Default opts + upcoming false.",
			opts: NewQueryOpts().UpcomingOnly(false),
			exp: url.Values{
				langKey:         []string{"en"},
				upcomingOnlyKey: []string{"false"},
			},
		},
		{
			name: "Default opts + upcoming true + upcoming false (false should overwrite).",
			opts: NewQueryOpts().UpcomingOnly(true).UpcomingOnly(false),
			exp: url.Values{
				langKey:         []string{"en"},
				upcomingOnlyKey: []string{"false"},
			},
		},
		{
			name: "Default opts + upcoming true + today + upcoming false + tomorrow (tomorrow + false should overwrite).",
			opts: NewQueryOpts().UpcomingOnly(true).TodayOnly().UpcomingOnly(false).TomorrowOnly(),
			exp: url.Values{
				langKey:            []string{"en"},
				upcomingOnlyKey:    []string{"false"},
				startTimeKey:       []string{sEOT},
				startTimeOffsetKey: []string{sEOD},
			},
		},
	}

	for _, tc := range tcs {
		c.Log(tc.name)

		c.Check(tc.opts.query, check.DeepEquals, tc.exp)
	}
}

// For now just very naively hit the Bovada endpoint
// and expect 200 response and nil errors.
func (cs *clientSuite) TestGetEvents(c *check.C) {
	var tcs = []struct {
		path string
		opts *queryOpts
	}{
		{
			path: MLBPath,
			opts: NewQueryOpts(),
		},
		{
			path: NBAPath,
			opts: NewQueryOpts().TodayOnly(),
		},
		{
			path: NHLPath,
			opts: NewQueryOpts().UpcomingOnly(true),
		},
	}

	for _, tc := range tcs {
		var cl = NewClient(http.DefaultClient)

		var resp, err = cl.GetEvents(tc.path, tc.opts)
		c.Check(resp, check.NotNil)
		c.Check(err, check.IsNil)
	}
}

func (cs *clientSuite) TestMinutesLeftInDay(c *check.C) {
	var tcs = []struct {
		name    string
		t       time.Time
		expDay  int
		expWeek int
	}{
		{
			name:    "Hours no minutes, Wednesday.",
			t:       time.Date(2020, time.August, 12, 1, 0, 0, 0, time.UTC),
			expDay:  22*60 + 60,
			expWeek: (22*60 + 60) + (4 * 24 * 60),
		},
		{
			name:    "Hours and minutes, Thursday.",
			t:       time.Date(2020, time.August, 13, 1, 1, 0, 0, time.UTC),
			expDay:  22*60 + 59,
			expWeek: (22*60 + 59) + (3 * 24 * 60),
		},
		{
			name:    "Hours, minutes, seconds, nanoseconds, Sunday.",
			t:       time.Date(2020, time.August, 16, 1, 1, 1, 1, time.UTC),
			expDay:  22*60 + 59,
			expWeek: 22*60 + 59,
		},
		{
			name:    "Hours, minutes, seconds, nanoseconds (afternoon), Monday.",
			t:       time.Date(2020, time.August, 17, 13, 1, 1, 1, time.UTC),
			expDay:  10*60 + 59,
			expWeek: 10*60 + 59 + (6 * 24 * 60),
		},
	}

	for _, tc := range tcs {
		c.Log(tc.name)

		c.Check(minutesLeftInDay(tc.t), check.Equals, tc.expDay)
		c.Check(minutesLeftInWeek(tc.t), check.Equals, tc.expWeek)
	}
}

var _ = check.Suite(&clientSuite{})

// Hook up to test runner.
func Test(t *testing.T) { check.TestingT(t) }
