package services

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

type TimekeepingService struct {
	repo repositories.AttendanceRepositoryInterface
}

func NewTimekeepingService(repo repositories.AttendanceRepositoryInterface) *TimekeepingService {
	return &TimekeepingService{repo: repo}
}

func (s *TimekeepingService) LogAttendance(req core.AttendanceRequest) (*core.Attendance, error) {
	// Calculate if overtime
	duration := req.ActualEndTime.Sub(req.ActualStartTime).Hours()

	// Simplified overtime check: if > 8 hours (or whatever shift type says, but we lack shift type context here easily without query)
	// For now, let's assume > 8h is overtime for simplicity of MVP or trust frontend to flag it?
	// Better: Calculate it.

	overtime := 0.0
	isOvertime := false
	if duration > 8.0 {
		overtime = duration - 8.0
		isOvertime = true
	}

	attendance := &core.Attendance{
		ScheduleID:      req.ScheduleID,
		ActualStartTime: req.ActualStartTime,
		ActualEndTime:   req.ActualEndTime,
		Notes:           req.Notes,
		IsOvertime:      isOvertime,
		OvertimeHours:   overtime,
	}

	if err := s.repo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *TimekeepingService) GetPayrollReport(month int, year int) ([]core.Attendance, error) {
	return s.repo.GetCombinedReport(month, year)
}
