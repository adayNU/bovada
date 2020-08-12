// Package bovada provides a means of querying the Bovada API
// to determine upcoming events for US professional sports leagues.
package bovada

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var now = time.Now

// Client implements methods for interacting with the Bovada API.
type Client struct {
	client *http.Client
}

// NewClient returns a new |*Client|. If a nil
// *http.Client is passed, it uses the http.DefaultClient.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}

	return &Client{
		client: c,
	}
}

// queryOpts hold the query parameters which reflect the
// options specified by the methods called upon it.
type queryOpts struct {
	query url.Values
}

// NewQueryOpts returns a |*queryOpts| with the language
// set to English.
func NewQueryOpts() *queryOpts {
	return &queryOpts{
		query: url.Values{
			langKey: []string{"en"},
		},
	}
}

// Upcoming will limit the query to upcoming events only
// based on the value of |b|. If this was previously set,
// this call will overwrite the old value.
func (q *queryOpts) UpcomingOnly(b bool) *queryOpts {
	q.query.Set(upcomingOnlyKey, strconv.FormatBool(b))
	return q
}

// TodayOnly will limit the query to only events starting today.
// If the start date parameter has already been set, it will overwrite it.
// It will also clear the start time offset, as it is by definition today
// so no filter is necessary there.
func (q *queryOpts) TodayOnly() *queryOpts {
	var t = minutesLeftInDay(now())
	q.query.Del(startTimeOffsetKey)
	q.query.Set(startTimeKey, strconv.Itoa(t))

	return q
}

// TomorrowOnly will limit the query to only events starting tomorrow.
// If the start date parameter has already been set, it will overwrite it.
func (q *queryOpts) TomorrowOnly() *queryOpts {
	var t = minutesLeftInDay(now())
	q.query.Set(startTimeOffsetKey, strconv.Itoa(t))
	q.query.Set(startTimeKey, strconv.Itoa(t+minutesInDay))

	return q
}

// GetEvents queries the API for the given path and with the specified options.
func (c *Client) GetEvents(path string, opts *queryOpts) (*EventResponse, error) {
	if opts == nil {
		opts = NewQueryOpts()
	}
	return c.getEvents(path, opts.query)
}

// getEvents queries the API with the specified path and query parameters.
func (c *Client) getEvents(path string, query url.Values) (*EventResponse, error) {
	var url = url.URL{
		Scheme:   "https",
		Host:     Host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	var b = &bytes.Buffer{}

	var r, err = http.NewRequest("GET", url.String(), b)
	if err != nil {
		return nil, errors.New("error building request: " + err.Error())
	}

	var resp *http.Response
	resp, err = c.client.Do(r)
	if err != nil {
		return nil, errors.New("error issuing request: " + err.Error())
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body: " + err.Error())
	}

	var e = make([]*EventResponse, 1)
	err = json.Unmarshal(body, &e)
	if err != nil {
		return nil, errors.New("error parsing response: " + err.Error())
	}

	return e[0], nil
}

// minutesInDay is the number of minutes in a day.
const minutesInDay = 60 * 24 // 1,440.

// minutesLeftInDay returns the number of minutes remaining
// in the day (rounded down) given time |t|.
func minutesLeftInDay(t time.Time) int {
	return (23-t.Hour())*60 + (60 - t.Minute())
}
