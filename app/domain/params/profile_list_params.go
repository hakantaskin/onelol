package params

import (
	"time"
)

type GetProfileListParams struct {
	Limit     int
	Offset    int
	LatestID  string
	StartDate time.Time
	EndDate   time.Time
	Active    *bool
	ShowAll   *bool
}
