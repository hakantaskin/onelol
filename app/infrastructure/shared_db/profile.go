package shared_db

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/hakantaskin/onelol/app/domain/entities"
	"github.com/hakantaskin/onelol/app/domain/params"
	"github.com/hakantaskin/onelol/app/infrastructure/pctx"
)

type ProfileEntry struct {
	ID             int         `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	FullName       string      `json:"full_name" gorm:"type:varchar(150);index;"`
	Name           string      `json:"name" gorm:"type:varchar(75);"`
	LastName       string      `json:"last_name" gorm:"type:varchar(75);"`
	BirthDay       time.Time   `json:"birth_day"`
	Information    string      `json:"information" gorm:"type:text;"`
	SocialAccounts params.JSON `json:"social_accounts" gorm:"type:json;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ProfileSocialAccountEntry struct {
	Nickname string `json:"nickname"`
	Link     string `json:"link"`
	Social   string `json:"social"`
}

func (ProfileEntry) TableName() string {
	return "profiles"
}

type ProfileRepository struct {
	db      *gorm.DB
	now     func() time.Time // todo: add local timezone by geid
	timeout time.Duration
	doneFn  func()
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db:      db,
		now:     time.Now,
		timeout: time.Second * 2,
		doneFn:  func() {},
	}
}

func (s *ProfileRepository) Create(ctx context.Context, discount entities.Profile) (entities.ProfileID, error) {
	db := pctx.DBTransaction(ctx)
	if db == nil {
		var err error
		db = s.db
		if err != nil {
			return 0, err
		}
	}

	entry := ToProfileEntry(discount)

	if err := db.Create(&entry).Error; err != nil {
		return 0, err
	}

	return entities.ProfileID(entry.ID), nil
}

func ToProfileEntry(profile entities.Profile) ProfileEntry {
	jsonSocialAccounts, _ := json.Marshal(profile.SocialAccounts)

	return ProfileEntry{
		FullName:       profile.BuildFullName(),
		Name:           profile.Name,
		LastName:       profile.LastName,
		BirthDay:       profile.BirthDay.Time,
		Information:    profile.Information,
		SocialAccounts: params.JSON{RawMessage: jsonSocialAccounts},
		CreatedAt:      profile.CreatedAt,
		UpdatedAt:      profile.UpdatedAt,
	}
}
