package valueobjects

type Pagination struct {
	Limit    int
	Offset   int
	Size     int
	LatestID string
	More     bool
}
