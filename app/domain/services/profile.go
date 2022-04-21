package services

import (
	"context"

	"github.com/hakantaskin/onelol/app/domain/entities"
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

type ProfileService interface {
	ProfileCreator
	ProfileGetter
	ProfileUpdater
	ProfileDeleter
}

type ProfileValidatorFn func(ctx context.Context, profile entities.Profile) error
