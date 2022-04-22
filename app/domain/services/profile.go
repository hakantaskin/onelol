package services

import (
	"context"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
)

type ProfileCreator interface {
	Create(ctx context.Context, profile entities.Profile) (entities.ProfileID, error)
}

type ProfileGetter interface {
	Get(ctx context.Context, profileID entities.ProfileID) (entities.Profile, error)
}

type ProfileDeleter interface {
	Delete(ctx context.Context, profile entities.Profile) error
}

type ProfileUpdater interface {
	Update(ctx context.Context, profile entities.Profile) error
}

type ProfileListGetter interface {
	GetList(ctx context.Context, filters params.GetProfileListParams) (
		profiles []entities.Profile, err error,
	)
}

type ProfileService interface {
	ProfileCreator
	ProfileGetter
	ProfileUpdater
	ProfileDeleter
	ProfileListGetter
}

type ProfileValidatorFn func(ctx context.Context, profile entities.Profile) error
