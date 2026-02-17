package services

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
)

type TimekeepingService struct {
	repo         repositories.AttendanceRepositoryInterface
	scheduleRepo ScheduleRepositoryInterface
}

func NewTimekeepingService(repo repositories.AttendanceRepositoryInterface, scheduleRepo ScheduleRepositoryInterface) *TimekeepingService {
	return &TimekeepingService{repo: repo, scheduleRepo: scheduleRepo}
}

// getShiftDurationHours calculates the planned duration from a ShiftType's StartTime/EndTime strings (HH:mm format)
func getShiftDurationHours(st core.ShiftType) float64 {
	start, errS := time.Parse("15:04", st.StartTime)
	end, errE := time.Parse("15:04", st.EndTime)
	if errS != nil || errE != nil {
		slog.Warn("Could not parse shift type times, defaulting to 8h", "start", st.StartTime, "end", st.EndTime)
		return 8.0
	}

	duration := end.Sub(start).Hours()
	if duration <= 0 {
		// Overnight shift (e.g. 16:00 -> 08:00)
		duration += 24.0
	}
	return duration
}

func (s *TimekeepingService) LogAttendance(req core.AttendanceRequest) (*core.Attendance, error) {
	// Get the schedule entry to find the shift type's planned hours
	plannedHours := 8.0 // fallback

	if s.scheduleRepo != nil {
		schedule, err := s.scheduleRepo.GetByID(req.ScheduleID)
		if err == nil && schedule != nil {
			plannedHours = getShiftDurationHours(schedule.ShiftType)
			slog.Info("Overtime calculation",
				"scheduleID", req.ScheduleID,
				"shiftType", schedule.ShiftType.Name,
				"plannedHours", plannedHours,
			)
		} else {
			slog.Warn("Could not fetch schedule for overtime calc, using 8h default", "scheduleID", req.ScheduleID, "error", err)
		}
	}

	// Calculate actual hours worked
	actualHours := req.ActualEndTime.Sub(req.ActualStartTime).Hours()
	if actualHours < 0 {
		return nil, fmt.Errorf("çıkış saati giriş saatinden önce olamaz")
	}

	overtime := 0.0
	isOvertime := false
	if actualHours > plannedHours {
		overtime = actualHours - plannedHours
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

func (s *TimekeepingService) UpdateAttendance(id uint, req core.AttendanceRequest) (*core.Attendance, error) {
	// Get existing attendance
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("kayıt bulunamadı: %w", err)
	}

	// Calculate planned hours from the schedule's shift type
	plannedHours := 8.0
	if s.scheduleRepo != nil {
		schedule, err := s.scheduleRepo.GetByID(existing.ScheduleID)
		if err == nil && schedule != nil {
			plannedHours = getShiftDurationHours(schedule.ShiftType)
		}
	}

	// Recalculate overtime
	actualHours := req.ActualEndTime.Sub(req.ActualStartTime).Hours()
	if actualHours < 0 {
		return nil, fmt.Errorf("çıkış saati giriş saatinden önce olamaz")
	}

	overtime := 0.0
	isOvertime := false
	if actualHours > plannedHours {
		overtime = actualHours - plannedHours
		isOvertime = true
	}

	existing.ActualStartTime = req.ActualStartTime
	existing.ActualEndTime = req.ActualEndTime
	existing.Notes = req.Notes
	existing.IsOvertime = isOvertime
	existing.OvertimeHours = overtime

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *TimekeepingService) GetPayrollReport(month int, year int) ([]core.Attendance, error) {
	return s.repo.GetCombinedReport(month, year)
}
