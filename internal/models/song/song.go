package song

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Song struct {
	ID        uuid.UUID      `gorm:"types:uuid;primaryKey"`
	Group     string         `gorm:"column:group;type:varchar"`
	Song      string         `gorm:"column:song;type:varchar"`
	Date      *time.Time     `gorm:"column:date;type:timestamp"`
	Lyrics    *string        `gorm:"column:lyrics;type:text"`
	Link      *string        `gorm:"column:link;type:varchar"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (song *Song) BeforeCreate(*gorm.DB) (err error) {
	song.ID = uuid.New()
	return nil
}

type Lyrics struct {
	Data  []string `json:"data"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int      `json:"total"`
}
