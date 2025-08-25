package model

import (
	"context"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:integer;autoIncrement"`
	EventID   uint      `json:"eventID" gorm:"type:integer;"`
	UserID    uint      `json:"userID" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Event     *Event    `json:"event" gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketRepository interface {
	GetMany(ctx context.Context, userID uint) ([]*Ticket, error)
	GetOne(ctx context.Context, userID, ticketID uint) (*Ticket, error)
	CreateOne(ctx context.Context, userID uint, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, userID, ticketID uint, updateData map[string]interface{}) (*Ticket, error)
}

type ValidateTicket struct {
	TicketID uint `json:"ticketID"`
	OwnerID  uint `json:"ownerID"`
}
