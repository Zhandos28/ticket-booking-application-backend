package repository

import (
	"context"
	"time"

	"github.com/Zhandos28/ticket-booking/model"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func (r *EventRepository) GetMany(ctx context.Context) ([]*model.Event, error) {
	var events []*model.Event

	if err := r.db.WithContext(ctx).Model(model.Event{}).Order("updated_at desc").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) GetOne(ctx context.Context, eventID uint) (*model.Event, error) {
	var event model.Event

	if err := r.db.WithContext(ctx).Model(model.Event{}).Where("id = ?", eventID).First(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) CreateOne(ctx context.Context, event *model.Event) (*model.Event, error) {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) UpdateOne(ctx context.Context, eventID uint, updateData map[string]interface{}) (*model.Event, error) {
	updateData["updated_at"] = time.Now()
	if err := r.db.WithContext(ctx).Model(model.Event{}).Where("id = ?", eventID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	event, err := r.GetOne(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) DeleteOne(ctx context.Context, eventID uint) error {
	if err := r.db.WithContext(ctx).Delete(model.Event{}, "id = ?", eventID).Error; err != nil {
		return err
	}

	return nil
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}
