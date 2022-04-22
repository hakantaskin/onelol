package shared_db

import "gorm.io/gorm"

type scopeFn func(*gorm.DB) *gorm.DB
