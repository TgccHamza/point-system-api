package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// WorkDayService defines the interface for workday-related operations.
type WorkDayService interface {
	CreateWorkDay(ctx context.Context, workday *models.WorkDay) error
	GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error)
	ListWorkDays(ctx context.Context) ([]*models.WorkDay, error)
	UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error
	DeleteWorkDay(ctx context.Context, id uint) error
}

// workDayService implements the WorkDayService interface.
type workDayService struct {
	workDayRepo repositories.WorkDayRepository
}

// NewWorkDayService creates a new instance of WorkDayService.
func NewWorkDayService(workDayRepo repositories.WorkDayRepository) WorkDayService {
	return &workDayService{
		workDayRepo: workDayRepo,
	}
}

// CreateWorkDay creates a new workday in the database.
func (s *workDayService) CreateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil {
		return errors.New("workday is nil")
	}

	if workday.Date.IsZero() {
		return errors.New("date is required")
	}

	if workday.DayType == "" {
		return errors.New("day type is required")
	}

	return s.workDayRepo.CreateWorkDay(ctx, workday)
}

// GetWorkDayByID retrieves a workday by its ID.
func (s *workDayService) GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error) {
	if id == 0 {
		return nil, errors.New("invalid workday ID")
	}

	return s.workDayRepo.GetWorkDayByID(ctx, id)
}

// ListWorkDays retrieves all workdays from the database.
func (s *workDayService) ListWorkDays(ctx context.Context) ([]*models.WorkDay, error) {
	return s.workDayRepo.ListWorkDays(ctx)
}

// UpdateWorkDay updates an existing workday in the database.
func (s *workDayService) UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil || workday.ID == 0 {
		return errors.New("invalid workday data")
	}

	return s.workDayRepo.UpdateWorkDay(ctx, workday)
}

// DeleteWorkDay deletes a workday by its ID.
func (s *workDayService) DeleteWorkDay(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid workday ID")
	}

	return s.workDayRepo.DeleteWorkDay(ctx, id)
}