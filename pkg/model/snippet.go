package model

import "github.com/jinzhu/gorm"

type Snippet struct {
	gorm.Model
	Text            string
	TestCase        string
	DocumentName    string
	DocumentVersion Version
}
