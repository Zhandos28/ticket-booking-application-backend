package repository

import (
	"context"
	"time"

	"github.com/Zhandos28/ticket-booking/model"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func (r *TicketRepository) GetMany(ctx context.Context, userID uint) ([]*model.Ticket, error) {
	tickets := make([]*model.Ticket, 0)

	if err := r.db.WithContext(ctx).Model(&model.Ticket{}).Preload("Event").Order("updated_at desc").
		Find(&tickets, "user_id", userID).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, userID, ticketID uint) (*model.Ticket, error) {
	ticket := model.Ticket{}

	if err := r.db.WithContext(ctx).Model(&ticket).Preload("Event").First(&ticket, "id = ? AND user_id = ?", ticketID, userID).
		Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, userID uint, ticket *model.Ticket) (*model.Ticket, error) {
	ticket.UserID = userID
	if err := r.db.WithContext(ctx).Create(&ticket).Error; err != nil {
		return nil, err
	}

	data, err := r.GetOne(ctx, userID, ticket.ID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TicketRepository) UpdateOne(ctx context.Context, userID, ticketID uint, updateData map[string]interface{}) (*model.Ticket, error) {
	ticket := model.Ticket{}
	updateData["updated_at"] = time.Now()

	if err := r.db.WithContext(ctx).Model(&ticket).Where("id = ? AND user_id = ?", ticketID, userID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	data, err := r.GetOne(ctx, userID, ticketID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}
