package services

import (
	"log/slog"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services/scheduler"
)

type ScheduleRepositoryInterface interface {
	Create(schedule *core.Schedule) error
	Update(schedule *core.Schedule) error
	GetCombinedSchedule(month int, year int) ([]core.Schedule, error)
	DeleteByMonthYear(month int, year int) error
	Delete(id uint) error
	GetByID(id uint) (*core.Schedule, error)
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
	// 1. Clear existing schedule for this month
	if err := s.repo.DeleteByMonthYear(req.Month, req.Year); err != nil {
		return nil, err
	}

	// 2. Fetch employees — filter by department if specified
	var employees []core.Employee
	var err error
	if req.DepartmentID > 0 {
		employees, err = s.employeeRepo.ListByDepartment(req.DepartmentID)
	} else {
		employees, err = s.employeeRepo.List()
	}
	if err != nil {
		return nil, err
	}

	// Filter by specific employee IDs if provided
	if len(req.EmployeeIDs) > 0 {
		selectedSet := make(map[uint]bool)
		for _, id := range req.EmployeeIDs {
			selectedSet[id] = true
		}
		var filtered []core.Employee
		for _, e := range employees {
			if selectedSet[e.ID] {
				filtered = append(filtered, e)
			}
		}
		employees = filtered
	}

	// 3. Fetch shift types — filter by selected IDs if specified
	allShiftTypes, err := s.shiftRepo.List()
	if err != nil {
		return nil, err
	}

	var shiftTypes []core.ShiftType
	if len(req.ShiftTypeIDs) > 0 {
		// Build lookup set
		selected := make(map[uint]bool)
		for _, id := range req.ShiftTypeIDs {
			selected[id] = true
		}
		for _, st := range allShiftTypes {
			if selected[st.ID] {
				shiftTypes = append(shiftTypes, st)
			}
		}
	} else {
		shiftTypes = allShiftTypes
	}

	if len(employees) == 0 || len(shiftTypes) == 0 {
		slog.Warn("Schedule generation skipped", "employees", len(employees), "shiftTypes", len(shiftTypes))
		return []core.Schedule{}, nil
	}

	slog.Info("Generating schedule",
		"month", req.Month,
		"year", req.Year,
		"employees", len(employees),
		"shiftTypes", len(shiftTypes),
		"departmentID", req.DepartmentID,
	)

	// 4. Initialize optimizer with constraints
	threshold := req.OvertimeThreshold
	if threshold == 0 {
		threshold = 45.0
	}

	constraints := []scheduler.Constraint{
		&scheduler.NoConsecutiveShifts{},
		&scheduler.WeeklyHourLimit{LimitHours: threshold},
	}

	optimizer := scheduler.NewOptimizer(constraints)

	// 5. Generate optimized schedule
	bestSchedule := optimizer.OptimizeSchedule(employees, shiftTypes, req.Month, req.Year)

	// 6. Save to database
	for _, sched := range bestSchedule {
		toSave := sched
		if err := s.repo.Create(&toSave); err != nil {
			return nil, err
		}
	}

	slog.Info("Schedule generated", "assignments", len(bestSchedule))
	return bestSchedule, nil
}

func (s *SchedulerService) GetMonthlySchedule(month, year int) ([]core.Schedule, error) {
	return s.repo.GetCombinedSchedule(month, year)
}

func (s *SchedulerService) UpdateSchedule(id uint, req core.Schedule) (*core.Schedule, error) {
	req.ID = id
	if err := s.repo.Update(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *SchedulerService) ClearMonthlySchedule(month, year int) error {
	slog.Info("Clearing schedule", "month", month, "year", year)
	return s.repo.DeleteByMonthYear(month, year)
}

func (s *SchedulerService) CreateSingleSchedule(schedule *core.Schedule) error {
	slog.Info("Creating single schedule", "date", schedule.Date, "employeeID", schedule.EmployeeID, "shiftTypeID", schedule.ShiftTypeID)
	return s.repo.Create(schedule)
}

func (s *SchedulerService) DeleteSchedule(id uint) error {
	slog.Info("Deleting single schedule", "id", id)
	return s.repo.Delete(id)
}
