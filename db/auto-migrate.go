package db

import (
	"github.com/Zhandos28/ticket-booking/model"
	"gorm.io/gorm"
)

func Migrator(db *gorm.DB) error {
	return db.AutoMigrate(&model.Event{}, &model.Ticket{}, &model.User{})
}
