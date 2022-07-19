package models

import (
	"time"

	"gorm.io/gorm"
)

type Wishlists struct {
	ID                    uint           `gorm:"primary key;autoIncrement" json:"id"`
	USER_ID               *uint          `gorm:"index" json:"user_id"`
	CATALOG_ID            *uint          `gorm:"index" json:"catalog_id"`
	CATALOG_NAME          *string        `json:"catalog_name"`
	CATALOG_TYPE          *string        `json:"catalog_type"`
	CATALOG_IMAGE_URL     *string        `json:"catalog_image_url"`
	CATALOG_CONDITION     *string        `json:"catalog_condition"`
	CATALOG_STRIKE_PRICE  *float64       `json:"catalog_strike_price"`
	CATALOG_SELLING_PRICE *float64       `json:"catalog_selling_price"`
	CREATED_AT            time.Time      `json:"created_at"`
	DELETED_AT            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func MigrateWishlists(db *gorm.DB) error {
	err := db.AutoMigrate(&Wishlists{})
	return err
}
