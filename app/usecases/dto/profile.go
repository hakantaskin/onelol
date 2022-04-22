package dto

import (
	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
)

type ProfileDTO struct {
	ID             int                       `json:"id"`
	FullName       string                    `json:"full_name"`
	Name           string                    `json:"name"`
	LastName       string                    `json:"last_name"`
	BirthDay       params.CustomDate         `json:"birth_day"`
	Information    string                    `json:"information"`
	SocialAccounts []ProfileSocialAccountDTO `json:"social_accounts"`
}

type ProfileSocialAccountDTO struct {
	Nickname string `json:"nickname"`
	Link     string `json:"link"`
	Social   string `json:"social"`
}

func BuildProfileDTO(
	profileEntity entities.Profile,
) *ProfileDTO {
	var socialAccounts []ProfileSocialAccountDTO
	for _, data := range profileEntity.SocialAccounts {
		socialAccounts = append(socialAccounts, ProfileSocialAccountDTO{
			Nickname: data.Nickname,
			Link:     data.Link,
			Social:   data.Social,
		})
	}
	return &ProfileDTO{
		ID:             profileEntity.ID.Int(),
		FullName:       profileEntity.BuildFullName(),
		Name:           profileEntity.Name,
		LastName:       profileEntity.LastName,
		BirthDay:       profileEntity.BirthDay,
		Information:    profileEntity.Information,
		SocialAccounts: socialAccounts,
	}
}
