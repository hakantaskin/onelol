package profile

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
	"github.com/hakantaskin/onelol/app/domain/services"
	"github.com/hakantaskin/onelol/app/domain/valueobjects"
	"github.com/hakantaskin/onelol/app/infrastructure/pctx"
	"github.com/hakantaskin/onelol/app/usecases/dto"
)

const (
	defaultListSize = 20
)

type ListService struct {
	profileListGetter services.ProfileListGetter
}

func NewListService(
	discountListGetter services.ProfileListGetter) *ListService {
	return &ListService{
		profileListGetter: discountListGetter,
	}
}

func (s *ListService) GetList(ctx context.Context, filters FiltersForm) (dto.PaginatedProfileList, error) {
	l := pctx.LoggerEntry(ctx)
	l = l.WithFields(
		logrus.Fields{
			"scope":   "profile",
			"method":  "GetList",
			"filters": filters,
		},
	)

	profileListParams, err := toProfileListParams(filters)
	if err != nil {
		l.Error(err)
		return dto.PaginatedProfileList{}, err
	}

	profiles := make([]entities.Profile, 0)
	profiles, err = s.profileListGetter.GetList(ctx, profileListParams)
	if err != nil {
		l.Error(err)
		return dto.PaginatedProfileList{}, err
	}

	latestID := ""
	if len(profiles) > 0 {
		latestID = profiles[len(profiles)-1].ID.String()
	}

	more := true
	if len(profiles) < profileListParams.Limit {
		more = false
	}

	pagination := valueobjects.Pagination{
		Limit:    profileListParams.Limit,
		Offset:   profileListParams.Offset,
		LatestID: latestID,
		More:     more,
	}

	return dto.BuildProfileListPaginated(profiles, pagination), nil
}

func toProfileListParams(filters FiltersForm) (params.GetProfileListParams, error) {
	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)

	if filters.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", filters.StartDate)
		if err != nil {
			return params.GetProfileListParams{}, err
		}
	}

	if filters.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", filters.EndDate)
		if err != nil {
			return params.GetProfileListParams{}, err
		}
	}

	if filters.Limit <= 0 {
		filters.Limit = defaultListSize
	}

	return params.GetProfileListParams{
		Limit:     filters.Limit,
		Offset:    filters.Offset,
		LatestID:  filters.LatestID,
		StartDate: startDate,
		EndDate:   endDate,
		FullName:  filters.FullName,
	}, nil
}
