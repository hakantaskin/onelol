package dto

type PaginationDTO struct {
	// limit
	// example: 20
	Limit int `json:"limit"`

	// offset
	// example: 0
	Offset int `json:"offset"`

	// size of collection
	// example: 1203
	Size int `json:"size"`

	// latest id of a batch
	// example: 123
	LatestID string `json:"latest_id"`

	// show if there is more data
	// example: true
	More bool `json:"more"`
}
