package profile

import (
	"context"

	"github.com/hakantaskin/onelol/app/usecases/dto"
	"github.com/hakantaskin/onelol/app/usecases/profile"
)

// ProfileCreateService ..
type ProfileCreateService interface {
	Create(ctx context.Context, discountForm profile.Form) (*dto.ProfileDTO, error)
}

// ProfileGetService ..
type ProfileGetService interface {
	Get(ctx context.Context, discountID int) (*dto.ProfileDTO, error)
}

// ProfileDeleteService ..
type ProfileDeleteService interface {
	Delete(ctx context.Context, discountID int, vendorIdentifier string) error
}

// ProfileUpdateService ..
type ProfileUpdateService interface {
	Update(ctx context.Context, discountID int, discountForm profile.UpdateForm) error
}

// ProfileListGetService ..
type ProfileListGetService interface {
	GetProfileList(ctx context.Context, filters profile.FiltersForm) (dto.PaginatedProfileList, error)
}
