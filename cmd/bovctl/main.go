package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/adayNU/bovada"
)

const (
	mlb = "mlb"
	nba = "nba"
	nhl = "nhl"
	nfl = "nfl"
	cfb = "cfb"

	today    = "today"
	tomorrow = "tomorrow"
	week = "week"
)

var (
	league   = flag.String("l", "mlb", `Specify which league to query events for ("mlb", "nba", "nhl", "nfl", "cfb").`)
	num      = flag.Int("n", 1, "Number of events to chose a bet for.")
	date     = flag.String("d", "all", `Limits events to a certain day ("today", "tomorrow", "all"')`)
	upcoming = flag.Bool("u", false, "Upcoming events only")
)

var leagueMap = map[string]string{
	mlb: bovada.MLBPath,
	nba: bovada.NBAPath,
	nhl: bovada.NHLPath,
	nfl: bovada.NFLPath,
	cfb: bovada.CFBPath,
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
		opts.TomorrowOnly()
	case week:
		opts.ThisWeek()
	default:
		// No-Op.
	}

	opts.UpcomingOnly(*upcoming)
	opts.GamesOnly()

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

	var ev = r.Events[:*num]
	sort.Slice(ev, func(i, j int) bool { return ev[i].StartTime < ev[j].StartTime })

	var w = tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	_, _ = fmt.Fprintln(w, strings.Join([]string{"Start Time", "Matchup", "Winner"}, "\t"))

	for _, event := range ev {
		rand.Seed(time.Now().UnixNano())
		_, _ = fmt.Fprintln(w, strings.Join([]string{
			time.Unix(event.StartTime/1000, 0).Format(time.RFC822),
			event.Description,
			event.Competitors[rand.Intn(2)].Name,
		}, "\t"))
	}

	_ = w.Flush()
}
