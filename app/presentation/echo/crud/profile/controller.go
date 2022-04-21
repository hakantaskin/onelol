package profile

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

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

// NewController ..
func NewController(profileCreateService profile_crud.ProfileCreateService,

/*profileGetService profile_crud.ProfileGetService,
  profileDeleteService profile_crud.ProfileDeleteService,
  profileUpdateService profile_crud.ProfileUpdateService,
  profileListGetService profile_crud.ProfileListGetService,*/
) *Controller {
	return &Controller{
		profileCreateService: profileCreateService,
		/*profileGetService:     profileGetService,
		profileDeleteService:  profileDeleteService,
		profileUpdateService:  profileUpdateService,
		profileListGetService: profileListGetService,*/
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
		return cc.JSON(http.StatusBadRequest, err)
	}

	profileDTO, err := c.profileCreateService.Create(ctx, form)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, err)
	}

	return cc.JSON(http.StatusOK, profileDTO)
}

func (c *Controller) GetProfile(ctx echo.Context) error {
	return nil
}

func (c *Controller) GetProfiles(ctx echo.Context) error {
	return nil
}

func (c *Controller) DeleteProfileByID(ctx echo.Context) error {
	return nil
}

func (c *Controller) UpdateProfile(ctx echo.Context) error {
	return nil
}
