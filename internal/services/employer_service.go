package services

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type EmployeeRepositoryInterface interface {
	Create(employee *core.Employee) error
	GetByID(id uint) (*core.Employee, error)
	List() ([]core.Employee, error)
	ListByDepartment(departmentID uint) ([]core.Employee, error)
	ListPaginated(ctx context.Context, params core.PaginationParams) ([]core.Employee, int64, error)
	Update(employee *core.Employee) error
	Delete(id uint) error
	GetDB() *gorm.DB
}

type EmployeeService struct {
	repo      EmployeeRepositoryInterface
	deptRepo  DepartmentRepositoryInterface
	titleRepo TitleRepositoryInterface
}

func NewEmployeeService(repo EmployeeRepositoryInterface, deptRepo DepartmentRepositoryInterface, titleRepo TitleRepositoryInterface) *EmployeeService {
	return &EmployeeService{
		repo:      repo,
		deptRepo:  deptRepo,
		titleRepo: titleRepo,
	}
}

func (s *EmployeeService) CreateEmployee(employee *core.Employee) error {
	return s.repo.Create(employee)
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*core.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) GetAllEmployees() ([]core.Employee, error) {
	return s.repo.List()
}

func (s *EmployeeService) GetPaginatedEmployees(ctx context.Context, params core.PaginationParams) (*core.PaginationResult, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	employees, total, err := s.repo.ListPaginated(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return &core.PaginationResult{
		Data:       employees,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *EmployeeService) UpdateEmployee(employee *core.Employee) error {
	return s.repo.Update(employee)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	return s.repo.Delete(id)
}

func (s *EmployeeService) ImportEmployees(reader io.Reader) error {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return err
	}
	defer f.Close()

	// Assume first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return err
	}

	// Skip header (row 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 7 {
			continue // Skip invalid rows
		}

		// Expected Columns: FirstName, LastName, Email, Phone, TitleName, DeptName, HourlyRate
		firstName := row[0]
		lastName := row[1]
		email := row[2]
		phone := row[3]
		titleName := row[4]
		deptName := row[5]
		hourlyRateStr := row[6]

		hourlyRate, _ := strconv.ParseFloat(hourlyRateStr, 64)

		// Lookup Title
		var titleID uint
		if title, err := s.titleRepo.GetByName(titleName); err == nil {
			titleID = title.ID
		} else {
			// Optional: Create title if not exists? For now, skip or set 0
			// fmt.Printf("Title %s not found\n", titleName)
		}

		// Lookup Department
		var deptID uint
		if dept, err := s.deptRepo.GetByName(deptName); err == nil {
			deptID = dept.ID
		} else {
			// fmt.Printf("Dept %s not found\n", deptName)
		}

		emp := &core.Employee{
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			Phone:        phone,
			TitleID:      titleID,
			DepartmentID: deptID,
			HourlyRate:   hourlyRate,
		}

		if err := s.repo.Create(emp); err != nil {
			// Continue or error? Just log error and continue for bulk import usually
			// or return partial error. Let's return error on first failure for now or valid imports.
			// Ideally we should collect errors.
			fmt.Printf("Failed to import %s %s: %v\n", firstName, lastName, err)
		}
	}

	return nil
}
