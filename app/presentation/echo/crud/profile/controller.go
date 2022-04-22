package profile

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/hakantaskin/onelol/app/domain/params"
	profile_crud "github.com/hakantaskin/onelol/app/domain/usecases/profile"
	"github.com/hakantaskin/onelol/app/usecases/profile"
)

// Controller ..
type Controller struct {
	profileCreateService  profile_crud.ProfileCreateService
	profileGetService     profile_crud.ProfileGetService
	profileDeleteService  profile_crud.ProfileDeleteService
	profileUpdateService  profile_crud.ProfileUpdateService
	profileListGetService profile_crud.ProfileListGetService
}

func NewController(profileCreateService profile_crud.ProfileCreateService,
	profileGetService profile_crud.ProfileGetService,
	profileListGetService profile_crud.ProfileListGetService,
	/*
	   profileDeleteService profile_crud.ProfileDeleteService,
	   profileUpdateService profile_crud.ProfileUpdateService,
	*/
) *Controller {
	return &Controller{
		profileCreateService:  profileCreateService,
		profileGetService:     profileGetService,
		profileListGetService: profileListGetService,
		/*profileDeleteService:  profileDeleteService,
		profileUpdateService:  profileUpdateService,
		*/
	}
}

// Init ..
func (c *Controller) Init(r *echo.Group) {
	r.POST("/profile", c.CreateProfile)
	r.GET("/profile/:id", c.GetProfile)
	r.GET("/profile", c.GetProfiles)
	r.DELETE("/profile/:id", c.DeleteProfileByID)
	r.PUT("/profile/:id", c.UpdateProfile)
}

func (c *Controller) CreateProfile(cc echo.Context) error {
	ctx := context.Background()
	var form profile.Form
	if err := cc.Bind(&form); err != nil {
		return cc.JSON(http.StatusBadRequest, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	profileDTO, err := c.profileCreateService.Create(ctx, form)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	return cc.JSON(http.StatusOK, params.JsonResponse{Data: profileDTO})
}

func (c *Controller) GetProfile(cc echo.Context) error {
	ctx := context.Background()
	profileParamID := cc.Param("id")
	profileID, err := strconv.Atoi(profileParamID)
	if err != nil {
		return cc.JSON(http.StatusBadRequest, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	profileDTO, err := c.profileGetService.Get(ctx, profileID)
	if err != nil {
		return cc.JSON(http.StatusNotFound, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	return cc.JSON(http.StatusOK, params.JsonResponse{Data: profileDTO})
}

func (c *Controller) GetProfiles(cc echo.Context) error {
	ctx := context.Background()
	form := cc.QueryParams()
	profileFilterForm, err := buildProfileListFilterFromRequestArgs(form)
	if err != nil {
		return cc.JSON(http.StatusBadRequest, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	result, err := c.profileListGetService.GetList(ctx, profileFilterForm)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, params.JsonResponse{Error: true, ErrorMsg: err.Error()})
	}

	return cc.JSON(http.StatusOK, params.JsonResponse{Data: result})
}

func (c *Controller) DeleteProfileByID(ctx echo.Context) error {
	return nil
}

func (c *Controller) UpdateProfile(ctx echo.Context) error {
	return nil
}

func buildProfileListFilterFromRequestArgs(form url.Values) (profile.FiltersForm, error) {
	var (
		err error

		limit   int
		offset  int
		active  *bool
		showAll *bool
	)
	if form.Get("limit") != "" {
		limit, err = strconv.Atoi(form.Get("limit"))
		if err != nil {
			return profile.FiltersForm{}, err
		}
	}

	if form.Get("offset") != "" {
		offset, err = strconv.Atoi(form.Get("offset"))
		if err != nil {
			return profile.FiltersForm{}, err
		}
	}

	if form.Get("active") != "" {
		isActive, err := strconv.ParseBool(form.Get("active"))
		if err != nil {
			return profile.FiltersForm{}, err
		}
		active = &isActive
	}

	if form.Get("show_all") != "" {
		isShowAll, err := strconv.ParseBool(form.Get("show_all"))
		if err != nil {
			return profile.FiltersForm{}, err
		}
		showAll = &isShowAll
	}

	return profile.FiltersForm{
		Limit:     limit,
		Offset:    offset,
		LatestID:  form.Get("latest_id"),
		Active:    active,
		StartDate: form.Get("start_date"),
		EndDate:   form.Get("end_date"),
		FullName:  form.Get("full_name"),
		ShowAll:   showAll,
	}, nil
}
