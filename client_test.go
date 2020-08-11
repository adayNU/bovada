package bovada

import (
	"net/http"
	"net/url"
	"testing"

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
				startTimeKey: []string{startTimeToday},
			},
		},
		{
			name: "Default opts + tomorrow.",
			opts: NewQueryOpts().TomorrowOnly(),
			exp: url.Values{
				langKey:      []string{"en"},
				startTimeKey: []string{startTimeTomorrow},
			},
		},
		{
			name: "Default opts + tomorrow + today (today should overwrite).",
			opts: NewQueryOpts().TomorrowOnly().TodayOnly(),
			exp: url.Values{
				langKey:      []string{"en"},
				startTimeKey: []string{startTimeToday},
			},
		},
		{
			name: "Default opts + today + tomorrow (tomorrow should overwrite).",
			opts: NewQueryOpts().TodayOnly().TomorrowOnly(),
			exp: url.Values{
				langKey:      []string{"en"},
				startTimeKey: []string{startTimeTomorrow},
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
				langKey:         []string{"en"},
				upcomingOnlyKey: []string{"false"},
				startTimeKey:    []string{startTimeTomorrow},
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

var _ = check.Suite(&clientSuite{})

// Hook up to test runner.
func Test(t *testing.T) { check.TestingT(t) }
