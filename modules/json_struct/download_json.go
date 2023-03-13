package json_struct

import (
	"time"
)

type Cruises struct {
	CruiseName        string
	CuriseDescription string
	StartDate         time.Time
	EndDate           time.Time
	StartPort         string
	EndPort           string
	Sailor            []Sailor
}

type Sailor struct {
	Position   string
	First_name string
	Last_name  string
}
