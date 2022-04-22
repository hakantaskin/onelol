package profile

import (
	"context"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/services"
	"github.com/hakantaskin/onelol/app/infrastructure/pctx"
	"github.com/hakantaskin/onelol/app/usecases/dto"
)

type GetService struct {
	profileGetter services.ProfileGetter
}

func NewGetService(profileGetter services.ProfileGetter) *GetService {
	return &GetService{profileGetter: profileGetter}
}

func (s *GetService) Get(ctx context.Context, profileID int) (*dto.ProfileDTO, error) {
	logEntry := pctx.LoggerEntry(ctx)
	l := logEntry.WithFields(map[string]interface{}{
		"scope":       "profile",
		"discount_id": profileID,
		"method":      "Get",
	})

	entityProfileID := entities.ProfileID(profileID)
	var profile entities.Profile
	var err error

	profile, err = s.profileGetter.Get(ctx, entityProfileID)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return dto.BuildProfileDTO(profile), nil
}
