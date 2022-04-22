package profile

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/services"
	"github.com/hakantaskin/onelol/app/infrastructure/pctx"
	"github.com/hakantaskin/onelol/app/usecases/dto"
)

type profileCreateTransformer interface {
	TransformCreate(ctx context.Context, form Form) (entities.Profile, error)
}

type CreateService struct {
	profileCreator     services.ProfileCreator
	profileValidatorFn services.ProfileValidatorFn
	profileTransformer profileCreateTransformer
}

func NewCreateService(profileCreator services.ProfileCreator,
	profileValidatorFn services.ProfileValidatorFn,
	profileTransformer profileCreateTransformer) *CreateService {
	return &CreateService{profileCreator: profileCreator,
		profileValidatorFn: profileValidatorFn,
		profileTransformer: profileTransformer}
}

func (s *CreateService) Create(ctx context.Context, form Form) (*dto.ProfileDTO, error) {
	logEntry := pctx.LoggerEntry(ctx)
	logEntry = logEntry.WithFields(logrus.Fields{
		"scope":  "profile",
		"form":   form,
		"method": "Create",
	})

	profileEntity, err := s.profileTransformer.TransformCreate(ctx, form)
	if err != nil {
		logEntry.Warn(err)
		return nil, err
	}

	if err := s.profileValidatorFn(ctx, profileEntity); err != nil {
		logEntry.Warn(err)
		return nil, err
	}

	profileID, err := s.createProfile(ctx, profileEntity)
	if err != nil {
		logEntry.Error(err)
		return nil, err
	}

	if !entities.ProfileID(form.ID).Empty() {
		profileID = entities.ProfileID(form.ID)
	}

	profileEntity.ID = profileID
	logEntry.WithFields(logrus.Fields{
		"id":    profileEntity.ID,
		"event": "profile_created",
	}).Info("profile created")

	return dto.BuildProfileDTO(profileEntity), nil
}

func (s *CreateService) createProfile(ctx context.Context, profileEntity entities.Profile) (entities.ProfileID, error) {
	var (
		profileID entities.ProfileID
		err       error
	)
	profileID, err = s.profileCreator.Create(ctx, profileEntity)
	if err != nil {
		return 0, err
	}

	return profileID, err
}
