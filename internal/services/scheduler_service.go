package services

import (
	"log/slog"
	"math"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/internal/services/scheduler"
)

type ScheduleRepositoryInterface interface {
	Create(schedule *core.Schedule) error
	Update(schedule *core.Schedule) error
	GetCombinedSchedule(month int, year int) ([]core.Schedule, error)
	DeleteByMonthYear(month int, year int) error
	Delete(id uint) error
	GetByID(id uint) (*core.Schedule, error)
	ListPaginated(params core.PaginationParams, month, year int) ([]core.Schedule, int64, error)
}

type SchedulerService struct {
	repo         ScheduleRepositoryInterface
	employeeRepo EmployeeRepositoryInterface
	shiftRepo    ShiftTypeRepositoryInterface
	leaveRepo    repositories.LeaveRepositoryInterface
}

func NewSchedulerService(repo ScheduleRepositoryInterface, empRepo EmployeeRepositoryInterface, shiftRepo ShiftTypeRepositoryInterface, leaveRepo repositories.LeaveRepositoryInterface) *SchedulerService {
	return &SchedulerService{
		repo:         repo,
		employeeRepo: empRepo,
		shiftRepo:    shiftRepo,
		leaveRepo:    leaveRepo,
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

	// Fetch approved leaves for the given range to support LeaveOverlapConstraint
	start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	var leaves []core.Leave
	if req.DepartmentID > 0 {
		leaves, _ = s.leaveRepo.ListByDepartment(req.DepartmentID, start, end)
	} else {
		// If no department, we'd need to fetch all or per employee, but ListByDepartment handles it if dept is specific
		// For now, let's just use what we have or add a ListByStatus to repository if needed.
		leaves, _ = s.leaveRepo.ListByStatus("approved")
	}

	constraints := []scheduler.Constraint{
		&scheduler.NoConsecutiveShifts{},
		&scheduler.WeeklyHourLimit{LimitHours: threshold},
		&scheduler.AnnualLeaveOverlapConstraint{ApprovedLeaves: leaves},
		&scheduler.MinimumRestConstraint{MinRestHours: 24},
	}

	optimizer := scheduler.NewOptimizer(constraints)

	// Fetch department to support bed capacity mode
	var department *core.Department
	if req.DepartmentID > 0 {
		var d core.Department
		if err := s.employeeRepo.GetDB().Model(&core.Department{}).Where("id = ?", req.DepartmentID).First(&d).Error; err == nil {
			department = &d
		}
	}

	// 5. Generate optimized schedule
	bestSchedule := optimizer.OptimizeSchedule(req, department, employees, shiftTypes)

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

func (s *SchedulerService) GetPaginatedSchedules(params core.PaginationParams, month, year int) (*core.PaginationResult, error) {
	data, total, err := s.repo.ListPaginated(params, month, year)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}

	return &core.PaginationResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
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
