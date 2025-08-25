package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID                    int       `gorm:"primaryKey;autoIncrement;type:integer" json:"id"`
	Name                  string    `json:"name"`
	Location              string    `json:"location"`
	TotalTicketsPurchased int64     `json:"total_tickets_purchased" gorm:"-"`
	TotalTicketsEntered   int64     `json:"total_tickets_entered" gorm:"-"`
	Date                  time.Time `json:"date"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type EventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventID uint) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, eventID uint, updateData map[string]interface{}) (*Event, error)
	DeleteOne(ctx context.Context, eventID uint) error
}

func (e *Event) AfterFind(db *gorm.DB) error {
	baseQuery := db.Model(&Ticket{}).Where("event_id = ?", e.ID)
	if err := baseQuery.Count(&e.TotalTicketsPurchased).Error; err != nil {
		return err
	}

	if err := baseQuery.Where("entered = ?", true).Count(&e.TotalTicketsEntered).Error; err != nil {
		return err
	}

	return nil
}
