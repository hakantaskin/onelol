package transformer

import (
	"context"
	"time"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/usecases/profile"
)

type ProfileTransformer struct {
	now func() time.Time
}

func NewProfileTransformer() *ProfileTransformer {
	return &ProfileTransformer{
		now: time.Now,
	}
}

func (t *ProfileTransformer) TransformCreate(ctx context.Context, form profile.Form) (entities.Profile, error) {
	p, err := t.transform(ctx, form)
	if err != nil {
		return entities.Profile{}, err
	}

	p.CreatedAt = t.now()
	p.UpdatedAt = t.now()

	return p, nil
}

func (t *ProfileTransformer) TransformUpdate(ctx context.Context, entityProfileID entities.ProfileID, form profile.UpdateForm) (entities.Profile, error) {
	p, err := t.transform(ctx, form.Form)
	if err != nil {
		return entities.Profile{}, err
	}

	p.ID = entityProfileID
	p.UpdatedAt = t.now()

	return p, nil
}

func (t *ProfileTransformer) transform(ctx context.Context, form profile.Form) (entities.Profile, error) {
	entity := entities.Profile{
		FullName:       "",
		Name:           form.Name,
		LastName:       form.LastName,
		BirthDay:       form.BirthDay,
		Information:    form.Information,
		SocialAccounts: form.SocialAccounts,
	}

	return entity, nil
}
