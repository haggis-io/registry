package model

import (
	"github.com/jinzhu/gorm"
)

type Query func(*gorm.DB) *gorm.DB
