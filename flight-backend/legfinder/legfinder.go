package legfinder

import (
	"strings"
	"time"
)

type LegFinder interface {
	Find(LegSpec) ([]Leg, error)
}

type Leg struct {
	Price    int // In pennies USD
	Segments []Segment
	//	fake     string
}

type Segment struct {
	Airlines      string
	FlightNumber  string
	ArrivalTime   time.Time
	DepartureTime time.Time
	Origin        string
	Destination   string
}

type LegSpec struct {
	Origin      string
	Destination string
	Dates       []string // "YYYY-MM-DD"
}

type TripSpec []LegSpec

func (l LegSpec) Hash() string {
	return strings.Join(append(l.Dates, l.Origin, l.Destination), "$")
}
