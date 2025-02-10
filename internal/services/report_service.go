package services

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ReportService interface {
	GenerateReport(ctx context.Context, companyID uint, startDate, endDate time.Time) ([]ReportResult, error)
}

type reportService struct {
	db *gorm.DB
}

type ReportResult struct {
	UserID       uint
	EmployeeName string
	WorkDays     float64
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportService{db: db}
}

func (s *reportService) GenerateReport(ctx context.Context, companyID uint, startDate, endDate time.Time) ([]ReportResult, error) {
	var results []ReportResult

	query := `
    WITH uniqueWorkDay AS (
        SELECT DISTINCT * 
        FROM work_days 
        WHERE work_days.date BETWEEN ? AND ?
        ORDER BY work_days.id DESC
    )
    SELECT 
        raw_attendances.user_id, 
        MIN(raw_attendances.employee_name) AS employee_name, 
        ROUND(
            SUM(
                IF(
                    ((raw_attendances.total_hours - raw_attendances.total_hour_out - 
                      IF(raw_attendances.calculate_lunch_hour, 1, 0))) > 9 
                    AND raw_attendances.calculate_over_time, 
                    (raw_attendances.total_hours - raw_attendances.total_hour_out - 
                     IF(raw_attendances.calculate_lunch_hour, 1, 0)), 
                    IF(
                        ((raw_attendances.total_hours - raw_attendances.total_hour_out - 
                          IF(raw_attendances.calculate_lunch_hour, 1, 0))) > 9, 
                        9, 
                        (raw_attendances.total_hours - raw_attendances.total_hour_out - 
                         IF(raw_attendances.calculate_lunch_hour, 1, 0))
                    )
                ) / 9
            ) * 2 
        ) / 2 AS work_days
    FROM uniqueWorkDay
    INNER JOIN raw_attendances 
        ON raw_attendances.work_day_id = uniqueWorkDay.id
    WHERE raw_attendances.company_id = ?
    GROUP BY raw_attendances.user_id;
    `

	err := s.db.WithContext(ctx).Raw(query, startDate, endDate, companyID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
