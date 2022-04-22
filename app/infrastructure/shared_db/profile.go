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

type ProfileEntryWithCreatedAndUpdatedBy struct {
	ProfileEntry
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
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

func (s *ProfileRepository) Get(ctx context.Context, profileID entities.ProfileID) (entities.Profile, error) {
	var err error
	db := pctx.DBTransaction(ctx)
	if db == nil {
		db = s.db
		if err != nil {
			return entities.Profile{}, err
		}
	}

	var entry ProfileEntry

	db = db.Select("id, name, last_name, full_name, birth_day, information, social_accounts, created_at, updated_at")

	if err = db.Take(&entry, profileID.Int()).Error; err != nil {
		return entities.Profile{}, err
	}

	return ToProfileEntity(entry), nil
}

func (s *ProfileRepository) GetList(ctx context.Context, filters params.GetProfileListParams) ([]entities.Profile, error) {
	var err error
	db := pctx.DBTransaction(ctx)
	if db == nil {
		db = s.db
		if err != nil {
			return nil, err
		}
	}

	filters.Limit += 1 // we get one more, to know if we have more data
	if filters.LatestID != "" {
		filters.Offset = 0
	}

	scopes, err := profileGetListCommonScopes(filters)
	if err != nil {
		return nil, err
	}

	scopes = append(scopes,
		discountLimitScope(filters.Limit),
		discountOffsetScope(filters.Offset),
	)

	var entries []ProfileEntry
	if err := db.Table(ProfileEntry{}.TableName() + " p").
		Select("id, name, last_name, full_name, birth_day, information, social_accounts, created_at, updated_at").
		Scopes(scopes...).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return []entities.Profile{}, nil
	}

	if len(entries) == filters.Limit {
		// we also want to remove the latest record from result
		entries = entries[:len(entries)-1]
	}

	return toProfileEntries(entries), nil
}

func toProfileEntries(entries []ProfileEntry) []entities.Profile {
	profiles := make([]entities.Profile, len(entries))
	for idx, entry := range entries {
		profileEntity := ToProfileEntity(entry)
		profiles[idx] = profileEntity
	}
	return profiles
}

func profileGetListCommonScopes(filters params.GetProfileListParams) ([]func(*gorm.DB) *gorm.DB, error) {
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if filters.LatestID != "" {
		scopes = append(
			scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where("p.id < ?", filters.LatestID)
			},
		)
	}

	scopes = append(scopes,
		discountStartDateScope(filters.StartDate),
		discountEndDateScope(filters.EndDate),
		discountFullNameScope(filters.FullName),
		func(db *gorm.DB) *gorm.DB {
			return db.Order("p.id DESC").
				Group("p.id")
		},
	)

	return scopes, nil
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

func ToProfileEntity(profile ProfileEntry) entities.Profile {
	var socialAccounts []entities.ProfileSocialAccount
	_ = json.Unmarshal(profile.SocialAccounts.RawMessage, &socialAccounts)

	return entities.Profile{
		ID:             entities.ProfileID(profile.ID),
		FullName:       profile.FullName,
		Name:           profile.Name,
		LastName:       profile.LastName,
		BirthDay:       params.CustomDate{Time: profile.BirthDay},
		Information:    profile.Information,
		SocialAccounts: socialAccounts,
		CreatedAt:      profile.CreatedAt,
		UpdatedAt:      profile.UpdatedAt,
	}
}

func discountFullNameScope(title string) scopeFn {
	return func(db *gorm.DB) *gorm.DB {
		if title != "" {
			return db.Where("p.full_name LIKE ?", "%"+title+"%")
		}

		return db
	}
}

func discountEndDateScope(endDate time.Time) scopeFn {
	return func(db *gorm.DB) *gorm.DB {
		if endDate.IsZero() {
			return db
		}

		return db.Where("p.created_at <= ?", endDate.Format("2006-01-02"))
	}
}

func discountStartDateScope(startDate time.Time) scopeFn {
	return func(db *gorm.DB) *gorm.DB {
		if startDate.IsZero() {
			return db
		}

		return db.Where("p.created_at >= ?", startDate.Format("2006-01-02"))
	}
}

func discountLimitScope(limit int) scopeFn {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}

		return db
	}
}

func discountOffsetScope(offset int) scopeFn {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Offset(offset)
		}

		return db
	}
}
