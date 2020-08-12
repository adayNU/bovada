package bovada

const (
	// Host is the hostname for Bovada.
	Host = "www.bovada.lv"
	// MLBPath is the path for the MLB events endpoint.
	MLBPath = "services/sports/event/coupon/events/A/description/baseball/mlb"
	// NBAPath is the path for the NBA events endpoint.
	NBAPath = "services/sports/event/coupon/events/A/description/basketball/nba"
	// NHLPath is the path for the NHL events endpoint.
	NHLPath = "services/sports/event/coupon/events/A/description/hockey/nhl"

	upcomingOnlyKey = "preMatchOnly"
	langKey         = "lang"

	// Bovada uses minutes as their units for their time limiting parameters.

	// startTimeKey represents the time (in minutes from now)
	// which the event must start prior to.
	startTimeKey = "startTimeLimit"
	// startTimeOffsetKey represents the time (in minutes from now)
	// which the event must start after.
	startTimeOffsetKey = "startTimeOffset"
)

// path includes metadata about the returned events.
type path struct {
	ID          string `json:"id"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Order       int    `json:"order"`
	Leaf        bool   `json:"leaf"`
	Current     bool   `json:"current"`
}

// competitor describes information about the competitors
// in an event.
type competitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Home bool   `json:"home"`
}

// event is the canonical event type.
type event struct {
	ID            string       `json:"id"`
	Description   string       `json:"description"`
	Type          string       `json:"type"`
	Link          string       `json:"link"`
	Status        string       `json:"status"`
	Sport         string       `json:"sport"`
	StartTime     int64        `json:"startTime"`
	Live          bool         `json:"live"`
	AwayTeamFirst bool         `json:"awayTeamFirst"`
	DenySameGame  string       `json:"denySameGame"`
	TeaserAllowed bool         `json:"teaserAllowed"`
	CompetitionID string       `json:"competitionId"`
	Notes         string       `json:"notes"`
	NumMarkets    int          `json:"numMarkets"`
	LastModified  int64        `json:"lastModified"`
	Competitors   []competitor `json:"competitors"`
}

// EventResponse is the data returned from an API call.
type EventResponse struct {
	Paths  []path  `json:"path"`
	Events []event `json:"events"`
}
