package services

import (
	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services/scheduler"
)

type ScheduleRepositoryInterface interface {
	Create(schedule *core.Schedule) error
	Update(schedule *core.Schedule) error
	GetCombinedSchedule(month int, year int) ([]core.Schedule, error)
	DeleteByMonthYear(month int, year int) error
}

type SchedulerService struct {
	repo         ScheduleRepositoryInterface
	employeeRepo EmployeeRepositoryInterface
	shiftRepo    ShiftTypeRepositoryInterface
}

func NewSchedulerService(repo ScheduleRepositoryInterface, empRepo EmployeeRepositoryInterface, shiftRepo ShiftTypeRepositoryInterface) *SchedulerService {
	return &SchedulerService{
		repo:         repo,
		employeeRepo: empRepo,
		shiftRepo:    shiftRepo,
	}
}

func (s *SchedulerService) GenerateSchedule(req core.ScheduleRequest) ([]core.Schedule, error) {
	// 1. Clear existing schedule
	if err := s.repo.DeleteByMonthYear(req.Month, req.Year); err != nil {
		return nil, err
	}

	// 2. Fetch resources
	employees, err := s.employeeRepo.List()
	if err != nil {
		return nil, err
	}

	shiftTypes, err := s.shiftRepo.List()
	if err != nil {
		return nil, err
	}

	if len(employees) == 0 || len(shiftTypes) == 0 {
		return []core.Schedule{}, nil
	}

	// 3. Output Optimization
	// Set defaults if not provided
	threshold := req.OvertimeThreshold
	if threshold == 0 {
		threshold = 45.0
	}
	multiplier := req.OvertimeMultiplier
	if multiplier == 0 {
		multiplier = 1.5
	}

	// Initialize constraints
	constraints := []scheduler.Constraint{
		&scheduler.NoConsecutiveShifts{},
		&scheduler.WeeklyHourLimit{LimitHours: threshold},
	}

	optimizer := scheduler.NewOptimizer(constraints)

	// Run 100 iterations
	bestSchedule := optimizer.OptimizeSchedule(employees, shiftTypes, req.Month, req.Year, 100, threshold, multiplier)

	// Save best schedule
	for _, sched := range bestSchedule {
		// sched is a value copy, we need to pass pointer to Create
		toSave := sched
		if err := s.repo.Create(&toSave); err != nil {
			return nil, err
		}
	}

	return bestSchedule, nil
}

func (s *SchedulerService) UpdateSchedule(id uint, req core.Schedule) (*core.Schedule, error) {
	req.ID = id
	if err := s.repo.Update(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
