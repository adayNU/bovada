package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/adayNU/bovada"
)

const (
	mlb = "mlb"
	nba = "nba"
	nhl = "nhl"

	today    = "today"
	tomorrow = "tomorrow"
)

var (
	league   = flag.String("l", "mlb", `Specify which league to query events for ("mlb", "nba"", "nhl").`)
	num      = flag.Int("n", 1, "Number of events to chose a bet for.")
	date     = flag.String("d", "all", `Limits events to a certain day ("today", "tomorrow", "all"')`)
	upcoming = flag.Bool("u", false, "Upcoming events only")
)

var leagueMap = map[string]string{
	mlb: bovada.MLBPath,
	nba: bovada.NBAPath,
	nhl: bovada.NHLPath,
}

func main() {
	flag.Parse()

	if _, ok := leagueMap[*league]; !ok {
		fmt.Println("Unknown league: ", *league)
		os.Exit(1)
	}

	var opts = bovada.NewQueryOpts()

	switch *date {
	case today:
		opts.TodayOnly()
	case tomorrow:
		opts.TodayOnly()
	default:
		// No-Op.
	}

	opts.UpcomingOnly(*upcoming)

	var c = bovada.NewClient(http.DefaultClient)
	var r, err = c.GetEvents(leagueMap[*league], opts)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		os.Exit(1)
	}

	if len(r.Events) < *num {
		fmt.Printf("Only (%d) events available, wanted (%d)\n", len(r.Events), *num)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r.Events), func(i, j int) { r.Events[i], r.Events[j] = r.Events[j], r.Events[i] })

	for i, event := range r.Events {
		if i == *num {
			break
		}

		rand.Seed(time.Now().UnixNano())
		fmt.Println(event.Competitors[rand.Intn(2)].Name)
	}
}
