package entities

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hakantaskin/onelol/app/domain/params"
)

type Profile struct {
	ID             ProfileID              `json:"id"`
	FullName       string                 `json:"full_name"`
	Name           string                 `json:"name"`
	LastName       string                 `json:"last_name"`
	BirthDay       params.CustomDate      `json:"birth_day"`
	Information    string                 `json:"information"`
	SocialAccounts []ProfileSocialAccount `json:"social_accounts"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *Profile) BuildFullName() string {
	s.FullName = fmt.Sprintf("%s %s", s.Name, s.LastName)
	return s.FullName
}

type ProfileSocialAccount struct {
	Nickname string `json:"nickname"`
	Link     string `json:"link"`
	Social   string `json:"social"`
}

type ProfileID int

// Int ..
func (pid ProfileID) Int() int {
	return int(pid)
}

// String ..
func (pid ProfileID) String() string {
	return strconv.Itoa(pid.Int())
}

// Empty returns true if it's empty or not valid
func (pid ProfileID) Empty() bool {
	return pid <= 0
}
