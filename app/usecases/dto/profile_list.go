package dto

import (
	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
	"github.com/hakantaskin/onelol/app/domain/valueobjects"
)

type PaginatedProfileList struct {
	Profiles   []ProfileListDTO `json:"profiles"`
	Pagination PaginationDTO    `json:"pagination"`
}

type ProfileListDTO struct {
	ID             int                       `json:"id"`
	FullName       string                    `json:"full_name"`
	Name           string                    `json:"name"`
	LastName       string                    `json:"last_name"`
	BirthDay       params.CustomDate         `json:"birth_day"`
	Information    string                    `json:"information"`
	SocialAccounts []ProfileSocialAccountDTO `json:"social_accounts"`
}

func BuildProfileListPaginated(profiles []entities.Profile, pagination valueobjects.Pagination) PaginatedProfileList {
	result := make([]ProfileListDTO, len(profiles))
	var socialAccounts []ProfileSocialAccountDTO
	for i, profileEntity := range profiles {
		for _, sa := range profileEntity.SocialAccounts {
			socialAccounts = append(socialAccounts, ProfileSocialAccountDTO{
				Nickname: sa.Nickname,
				Link:     sa.Link,
				Social:   sa.Social,
			})
		}

		profileDTO := ProfileListDTO{
			ID:             profileEntity.ID.Int(),
			FullName:       profileEntity.FullName,
			Name:           profileEntity.Name,
			LastName:       profileEntity.LastName,
			BirthDay:       profileEntity.BirthDay,
			Information:    profileEntity.Information,
			SocialAccounts: socialAccounts,
		}
		result[i] = profileDTO
	}

	return PaginatedProfileList{
		Profiles:   result,
		Pagination: buildPaginationDTO(pagination),
	}
}

func buildPaginationDTO(p valueobjects.Pagination) PaginationDTO {
	return PaginationDTO{
		Limit:    p.Limit,
		Offset:   p.Offset,
		Size:     p.Size,
		LatestID: p.LatestID,
		More:     p.More,
	}
}
