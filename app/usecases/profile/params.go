package profile

import (
	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
)

// Form .
type Form struct {
	ID             int                             `json:"id"`
	Name           string                          `json:"name"`
	LastName       string                          `json:"last_name"`
	BirthDay       params.CustomDate               `json:"birth_day"`
	Information    string                          `json:"information"`
	IsActive       bool                            `json:"is_active"`
	SocialAccounts []entities.ProfileSocialAccount `json:"social_accounts"`
}

// UpdateForm .
type UpdateForm struct {
	Form Form `json:"form"`
}

// FiltersForm ..
type FiltersForm struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	LatestID  string `json:"latest_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	FullName  string `json:"full_name"`
	Active    *bool  `json:"active"`
	ShowAll   *bool  `json:"show_all"`
}
