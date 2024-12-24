package models

import (
	"time"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Webpage Model -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type Webpage struct {
	ID      int    `gorm:"primaryKey" bson:"_id"`
	URL     string `gorm:"not null" bson:"url"`
	Title   string `gorm:"not null" bson:"title"`
	Content string `bson:"content"`
	Version int    `bson:"version"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Deleted   bool      `bson:"deleted"`
}
