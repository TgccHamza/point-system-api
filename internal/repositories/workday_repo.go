package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// WorkDayRepository defines the interface for workday-related database operations.
type WorkDayRepository interface {
	CreateWorkDay(ctx context.Context, workday *models.WorkDay) error
	GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error)
	ListWorkDays(ctx context.Context) ([]*models.WorkDay, error)
	UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error
	DeleteWorkDay(ctx context.Context, id uint) error
}

// workDayRepository implements the WorkDayRepository interface.
type workDayRepository struct {
	db *gorm.DB
}

// NewWorkDayRepository creates a new instance of WorkDayRepository.
func NewWorkDayRepository(db *gorm.DB) WorkDayRepository {
	return &workDayRepository{
		db: db,
	}
}

// CreateWorkDay inserts a new workday into the database.
func (r *workDayRepository) CreateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil {
		return errors.New("workday is nil")
	}

	if err := r.db.WithContext(ctx).Create(workday).Error; err != nil {
		return fmt.Errorf("failed to create workday: %w", err)
	}

	return nil
}

// GetWorkDayByID retrieves a workday by its ID.
func (r *workDayRepository) GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error) {
	var workday models.WorkDay
	if err := r.db.WithContext(ctx).First(&workday, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No workday found
		}
		return nil, fmt.Errorf("failed to retrieve workday by ID: %w", err)
	}
	return &workday, nil
}

// ListWorkDays retrieves all workdays from the database.
func (r *workDayRepository) ListWorkDays(ctx context.Context) ([]*models.WorkDay, error) {
	var workdays []*models.WorkDay
	if err := r.db.WithContext(ctx).Find(&workdays).Error; err != nil {
		return nil, fmt.Errorf("failed to list workdays: %w", err)
	}
	return workdays, nil
}

// UpdateWorkDay updates an existing workday in the database.
func (r *workDayRepository) UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil || workday.ID == 0 {
		return errors.New("invalid workday data")
	}

	if err := r.db.WithContext(ctx).Save(workday).Error; err != nil {
		return fmt.Errorf("failed to update workday: %w", err)
	}

	return nil
}

// DeleteWorkDay deletes a workday by its ID.
func (r *workDayRepository) DeleteWorkDay(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid workday ID")
	}

	if err := r.db.WithContext(ctx).Delete(&models.WorkDay{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete workday: %w", err)
	}

	return nil
}