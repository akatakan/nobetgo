package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeNotificationRepo struct {
	notifications map[uint]*core.Notification
	nextID        uint
}

func (r *fakeNotificationRepo) Create(notification *core.Notification) error {
	r.nextID++
	notification.ID = r.nextID
	r.notifications[notification.ID] = notification
	return nil
}

func (r *fakeNotificationRepo) GetUnreadByEmployee(employeeID uint) ([]core.Notification, error) {
	var result []core.Notification
	for _, notification := range r.notifications {
		if notification.EmployeeID == employeeID && !notification.IsRead {
			result = append(result, *notification)
		}
	}
	return result, nil
}

func (r *fakeNotificationRepo) MarkAsReadForEmployee(id uint, employeeID uint) (bool, error) {
	notification, ok := r.notifications[id]
	if !ok || notification.EmployeeID != employeeID {
		return false, nil
	}
	notification.IsRead = true
	return true, nil
}

func (r *fakeNotificationRepo) MarkAllAsRead(employeeID uint) error {
	for _, notification := range r.notifications {
		if notification.EmployeeID == employeeID {
			notification.IsRead = true
		}
	}
	return nil
}

type fakeTimeEntryRepo struct {
	entries map[uint]*core.TimeEntry
	nextID  uint
}

func (r *fakeTimeEntryRepo) Create(entry *core.TimeEntry) error {
	r.nextID++
	entry.ID = r.nextID
	r.entries[entry.ID] = entry
	return nil
}

func (r *fakeTimeEntryRepo) Update(entry *core.TimeEntry) error {
	r.entries[entry.ID] = entry
	return nil
}

func (r *fakeTimeEntryRepo) GetByID(id uint) (*core.TimeEntry, error) {
	entry, ok := r.entries[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return entry, nil
}

func (r *fakeTimeEntryRepo) Delete(id uint) error {
	delete(r.entries, id)
	return nil
}

func (r *fakeTimeEntryRepo) GetOpenEntry(employeeID uint) (*core.TimeEntry, error) {
	for _, entry := range r.entries {
		if entry.EmployeeID == employeeID && entry.ClockOut == nil {
			return entry, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeTimeEntryRepo) ListByEmployee(employeeID uint, start, end time.Time) ([]core.TimeEntry, error) {
	var result []core.TimeEntry
	for _, entry := range r.entries {
		if entry.EmployeeID == employeeID {
			result = append(result, *entry)
		}
	}
	return result, nil
}

func (r *fakeTimeEntryRepo) ListByDepartment(departmentID uint, start, end time.Time) ([]core.TimeEntry, error) {
	var result []core.TimeEntry
	for _, entry := range r.entries {
		if entry.Employee.DepartmentID == departmentID {
			result = append(result, *entry)
		}
	}
	return result, nil
}

func (r *fakeTimeEntryRepo) ListByDateRange(start, end time.Time) ([]core.TimeEntry, error) {
	var result []core.TimeEntry
	for _, entry := range r.entries {
		result = append(result, *entry)
	}
	return result, nil
}

func (r *fakeTimeEntryRepo) ListByStatus(status string, start, end time.Time) ([]core.TimeEntry, error) {
	var result []core.TimeEntry
	for _, entry := range r.entries {
		if entry.Status == status {
			result = append(result, *entry)
		}
	}
	return result, nil
}

func (r *fakeTimeEntryRepo) ListPaginated(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) ([]core.TimeEntry, int64, error) {
	var result []core.TimeEntry
	for _, entry := range r.entries {
		if employeeID != 0 && entry.EmployeeID != employeeID {
			continue
		}
		if departmentID != 0 && entry.Employee.DepartmentID != departmentID {
			continue
		}
		result = append(result, *entry)
	}
	return result, int64(len(result)), nil
}

type fakeLeaveRepo struct {
	leaves    map[uint]*core.Leave
	balances  map[uint][]core.LeaveBalance
	nextID    uint
	leaveType core.LeaveType
}

func (r *fakeLeaveRepo) Create(leave *core.Leave) error {
	r.nextID++
	leave.ID = r.nextID
	r.leaves[leave.ID] = leave
	return nil
}

func (r *fakeLeaveRepo) Update(leave *core.Leave) error {
	r.leaves[leave.ID] = leave
	return nil
}

func (r *fakeLeaveRepo) GetByID(id uint) (*core.Leave, error) {
	leave, ok := r.leaves[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return leave, nil
}

func (r *fakeLeaveRepo) Delete(id uint) error {
	delete(r.leaves, id)
	return nil
}

func (r *fakeLeaveRepo) ListByEmployee(employeeID uint, start, end time.Time) ([]core.Leave, error) {
	var result []core.Leave
	for _, leave := range r.leaves {
		if employeeID != 0 && leave.EmployeeID != employeeID {
			continue
		}
		result = append(result, *leave)
	}
	return result, nil
}

func (r *fakeLeaveRepo) ListByDepartment(departmentID uint, start, end time.Time) ([]core.Leave, error) {
	var result []core.Leave
	for _, leave := range r.leaves {
		if leave.Employee.DepartmentID == departmentID {
			result = append(result, *leave)
		}
	}
	return result, nil
}

func (r *fakeLeaveRepo) ListByStatus(status string) ([]core.Leave, error) {
	var result []core.Leave
	for _, leave := range r.leaves {
		if leave.Status == status {
			result = append(result, *leave)
		}
	}
	return result, nil
}

func (r *fakeLeaveRepo) HasOverlap(employeeID uint, start, end time.Time, excludeID uint) (bool, error) {
	for _, leave := range r.leaves {
		if leave.EmployeeID != employeeID || leave.ID == excludeID || leave.Status == "rejected" {
			continue
		}
		if leave.StartDate.Before(end) && leave.EndDate.After(start) {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeLeaveRepo) CreateLeaveType(lt *core.LeaveType) error { return nil }
func (r *fakeLeaveRepo) UpdateLeaveType(lt *core.LeaveType) error { return nil }
func (r *fakeLeaveRepo) DeleteLeaveType(id uint) error            { return nil }
func (r *fakeLeaveRepo) ListLeaveTypes() ([]core.LeaveType, error) {
	return []core.LeaveType{r.leaveType}, nil
}
func (r *fakeLeaveRepo) GetLeaveTypeByID(id uint) (*core.LeaveType, error) { return &r.leaveType, nil }

func (r *fakeLeaveRepo) GetBalance(employeeID uint, leaveTypeID uint, year int) (*core.LeaveBalance, error) {
	for _, balance := range r.balances[employeeID] {
		if balance.LeaveTypeID == leaveTypeID && balance.Year == year {
			b := balance
			return &b, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeLeaveRepo) GetAllBalances(employeeID uint, year int) ([]core.LeaveBalance, error) {
	var result []core.LeaveBalance
	for _, balance := range r.balances[employeeID] {
		if balance.Year == year {
			result = append(result, balance)
		}
	}
	return result, nil
}

func (r *fakeLeaveRepo) UpsertBalance(balance *core.LeaveBalance) error {
	r.balances[balance.EmployeeID] = append(r.balances[balance.EmployeeID], *balance)
	return nil
}

func (r *fakeLeaveRepo) ListPaginated(params core.PaginationParams, employeeID, departmentID uint, start, end time.Time) ([]core.Leave, int64, error) {
	var result []core.Leave
	for _, leave := range r.leaves {
		if employeeID != 0 && leave.EmployeeID != employeeID {
			continue
		}
		if departmentID != 0 && leave.Employee.DepartmentID != departmentID {
			continue
		}
		result = append(result, *leave)
	}
	return result, int64(len(result)), nil
}

type fakeOvertimeRuleRepo struct {
	active *core.OvertimeRule
}

func (r *fakeOvertimeRuleRepo) Create(rule *core.OvertimeRule) error            { return nil }
func (r *fakeOvertimeRuleRepo) Update(rule *core.OvertimeRule) error            { return nil }
func (r *fakeOvertimeRuleRepo) GetByID(id uint) (*core.OvertimeRule, error)     { return r.active, nil }
func (r *fakeOvertimeRuleRepo) GetActive() (*core.OvertimeRule, error)          { return r.active, nil }
func (r *fakeOvertimeRuleRepo) List() ([]core.OvertimeRule, error)              { return nil, nil }
func (r *fakeOvertimeRuleRepo) Delete(id uint) error                            { return nil }
func (r *fakeOvertimeRuleRepo) CreateHoliday(holiday *core.PublicHoliday) error { return nil }
func (r *fakeOvertimeRuleRepo) UpdateHoliday(holiday *core.PublicHoliday) error { return nil }
func (r *fakeOvertimeRuleRepo) GetHolidayByID(id uint) (*core.PublicHoliday, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeOvertimeRuleRepo) ListHolidays(year int) ([]core.PublicHoliday, error) {
	return nil, nil
}
func (r *fakeOvertimeRuleRepo) DeleteHoliday(id uint) error            { return nil }
func (r *fakeOvertimeRuleRepo) IsHoliday(date time.Time) (bool, error) { return false, nil }

type fakeAuditLogRepo struct{}

func (r *fakeAuditLogRepo) Create(log *core.AuditLog) error { return nil }
func (r *fakeAuditLogRepo) ListByEntity(entityType string, entityID uint) ([]core.AuditLog, error) {
	return nil, nil
}
func (r *fakeAuditLogRepo) ListByDateRange(start, end time.Time) ([]core.AuditLog, error) {
	return nil, nil
}
func (r *fakeAuditLogRepo) ListByPerformer(performerID uint, start, end time.Time) ([]core.AuditLog, error) {
	return nil, nil
}

func newAuthedRouter(userID uint, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	})
	return router
}

func TestNotificationGetUnreadScopesToActingUser(t *testing.T) {
	repo := &fakeNotificationRepo{
		notifications: map[uint]*core.Notification{
			1: {Model: gorm.Model{ID: 1}, EmployeeID: 1, Title: "mine"},
			2: {Model: gorm.Model{ID: 2}, EmployeeID: 2, Title: "other"},
		},
		nextID: 2,
	}
	handler := NewNotificationHandler(services.NewNotificationService(repo))
	router := newAuthedRouter(1, "user")
	router.GET("/notifications/unread", handler.GetUnread)

	req := httptest.NewRequest(http.MethodGet, "/notifications/unread?employee_id=1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var notifications []core.Notification
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &notifications))
	require.Len(t, notifications, 1)
	require.Equal(t, uint(1), notifications[0].EmployeeID)
}

func TestNotificationMarkAsReadRejectsForeignNotification(t *testing.T) {
	repo := &fakeNotificationRepo{
		notifications: map[uint]*core.Notification{
			2: {Model: gorm.Model{ID: 2}, EmployeeID: 2, Title: "other"},
		},
		nextID: 2,
	}
	handler := NewNotificationHandler(services.NewNotificationService(repo))
	router := newAuthedRouter(1, "user")
	router.POST("/notifications/:id/read", handler.MarkAsRead)

	req := httptest.NewRequest(http.MethodPost, "/notifications/2/read", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
	require.False(t, repo.notifications[2].IsRead)
}

func TestClockInUsesActingUser(t *testing.T) {
	repo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	handler := NewTimeEntryHandler(services.NewTimekeepingService(repo, nil))
	router := newAuthedRouter(7, "user")
	router.POST("/time-entries/clock-in", handler.ClockIn)

	req := httptest.NewRequest(http.MethodPost, "/time-entries/clock-in", strings.NewReader(`{"employee_id":99,"notes":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)
	require.Len(t, repo.entries, 1)
	for _, entry := range repo.entries {
		require.Equal(t, uint(7), entry.EmployeeID)
	}
}

func TestGetTimeEntryRejectsForeignOwner(t *testing.T) {
	repo := &fakeTimeEntryRepo{
		entries: map[uint]*core.TimeEntry{
			1: {Model: gorm.Model{ID: 1}, EmployeeID: 2},
		},
		nextID: 1,
	}
	handler := NewTimeEntryHandler(services.NewTimekeepingService(repo, nil))
	router := newAuthedRouter(1, "user")
	router.GET("/time-entries/:id", handler.GetTimeEntry)

	req := httptest.NewRequest(http.MethodGet, "/time-entries/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestListTimeEntriesRejectsForeignEmployeeFilter(t *testing.T) {
	repo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	handler := NewTimeEntryHandler(services.NewTimekeepingService(repo, nil))
	router := newAuthedRouter(1, "user")
	router.GET("/time-entries", handler.ListTimeEntries)

	req := httptest.NewRequest(http.MethodGet, "/time-entries?employee_id=2", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestRequestLeaveUsesActingUser(t *testing.T) {
	repo := &fakeLeaveRepo{
		leaves:    map[uint]*core.Leave{},
		balances:  map[uint][]core.LeaveBalance{},
		leaveType: core.LeaveType{Model: gorm.Model{ID: 1}},
	}
	handler := NewLeaveHandler(services.NewLeaveService(repo))
	router := newAuthedRouter(5, "user")
	router.POST("/leaves", handler.RequestLeave)

	body := `{"employee_id":99,"leave_type_id":1,"start_date":"2026-02-01T00:00:00Z","end_date":"2026-02-02T00:00:00Z","reason":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/leaves", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)
	require.Len(t, repo.leaves, 1)
	for _, leave := range repo.leaves {
		require.Equal(t, uint(5), leave.EmployeeID)
	}
}

func TestGetLeaveRejectsForeignOwner(t *testing.T) {
	repo := &fakeLeaveRepo{
		leaves: map[uint]*core.Leave{
			1: {Model: gorm.Model{ID: 1}, EmployeeID: 2},
		},
		balances:  map[uint][]core.LeaveBalance{},
		leaveType: core.LeaveType{},
	}
	handler := NewLeaveHandler(services.NewLeaveService(repo))
	router := newAuthedRouter(1, "user")
	router.GET("/leaves/:id", handler.GetLeave)

	req := httptest.NewRequest(http.MethodGet, "/leaves/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestGetLeaveBalanceScopesToActingUser(t *testing.T) {
	repo := &fakeLeaveRepo{
		leaves: map[uint]*core.Leave{},
		balances: map[uint][]core.LeaveBalance{
			1: {{EmployeeID: 1, LeaveTypeID: 1, Year: 2026}},
			2: {{EmployeeID: 2, LeaveTypeID: 1, Year: 2026}},
		},
		leaveType: core.LeaveType{},
	}
	handler := NewLeaveHandler(services.NewLeaveService(repo))
	router := newAuthedRouter(1, "user")
	router.GET("/leaves/balance", handler.GetLeaveBalance)

	req := httptest.NewRequest(http.MethodGet, "/leaves/balance?employee_id=1&year=2026", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var balances []core.LeaveBalance
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &balances))
	require.Len(t, balances, 1)
	require.Equal(t, uint(1), balances[0].EmployeeID)
}

func TestApprovalUsesActingApprover(t *testing.T) {
	repo := &fakeTimeEntryRepo{
		entries: map[uint]*core.TimeEntry{
			1: {Model: gorm.Model{ID: 1}, EmployeeID: 2, Status: "pending"},
		},
		nextID: 1,
	}
	service := services.NewApprovalService(&fakeAuditLogRepo{}, repo, &fakeLeaveRepo{leaves: map[uint]*core.Leave{}, balances: map[uint][]core.LeaveBalance{}})
	handler := NewApprovalHandler(service)
	router := newAuthedRouter(42, "admin")
	router.POST("/approvals/time-entry/:id/approve", handler.ApproveTimeEntry)

	req := httptest.NewRequest(http.MethodPost, "/approvals/time-entry/1/approve", strings.NewReader(`{"approver_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.NotNil(t, repo.entries[1].ApprovedBy)
	require.Equal(t, uint(42), *repo.entries[1].ApprovedBy)
}

func TestCalculateOvertimeRejectsForeignEmployee(t *testing.T) {
	timeRepo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	ruleRepo := &fakeOvertimeRuleRepo{active: &core.OvertimeRule{WeeklyHourLimit: 45}}
	handler := NewOvertimeHandler(services.NewOvertimeService(ruleRepo, timeRepo))
	router := newAuthedRouter(1, "user")
	router.GET("/overtime/calculate", handler.CalculateOvertime)

	req := httptest.NewRequest(http.MethodGet, "/overtime/calculate?employee_id=2&month=2&year=2026", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestDepartmentOvertimeSummaryRequiresAdmin(t *testing.T) {
	timeRepo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	ruleRepo := &fakeOvertimeRuleRepo{active: &core.OvertimeRule{WeeklyHourLimit: 45}}
	handler := NewOvertimeHandler(services.NewOvertimeService(ruleRepo, timeRepo))
	router := newAuthedRouter(1, "user")
	router.GET("/overtime/summary", handler.GetDepartmentSummary)

	req := httptest.NewRequest(http.MethodGet, "/overtime/summary?department_id=1&month=2&year=2026", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestEmployeeSummaryRejectsForeignEmployee(t *testing.T) {
	timeRepo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	leaveRepo := &fakeLeaveRepo{leaves: map[uint]*core.Leave{}, balances: map[uint][]core.LeaveBalance{}}
	ruleRepo := &fakeOvertimeRuleRepo{active: &core.OvertimeRule{WeeklyHourLimit: 45}}
	handler := NewReportingHandler(services.NewReportingService(timeRepo, leaveRepo, ruleRepo))
	router := newAuthedRouter(1, "user")
	router.GET("/reports/employee-summary", handler.GetEmployeeSummary)

	req := httptest.NewRequest(http.MethodGet, "/reports/employee-summary?employee_id=2&month=2&year=2026", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestWorkHoursReportRequiresAdmin(t *testing.T) {
	timeRepo := &fakeTimeEntryRepo{entries: map[uint]*core.TimeEntry{}}
	leaveRepo := &fakeLeaveRepo{leaves: map[uint]*core.Leave{}, balances: map[uint][]core.LeaveBalance{}}
	ruleRepo := &fakeOvertimeRuleRepo{active: &core.OvertimeRule{WeeklyHourLimit: 45}}
	handler := NewReportingHandler(services.NewReportingService(timeRepo, leaveRepo, ruleRepo))
	router := newAuthedRouter(1, "user")
	router.GET("/reports/work-hours", handler.GetWorkHoursReport)

	req := httptest.NewRequest(http.MethodGet, "/reports/work-hours?month=2&year=2026", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusForbidden, resp.Code)
}

var _ repositories.NotificationRepositoryInterface = (*fakeNotificationRepo)(nil)
var _ repositories.TimeEntryRepositoryInterface = (*fakeTimeEntryRepo)(nil)
var _ repositories.LeaveRepositoryInterface = (*fakeLeaveRepo)(nil)
var _ repositories.OvertimeRuleRepositoryInterface = (*fakeOvertimeRuleRepo)(nil)
var _ repositories.AuditLogRepositoryInterface = (*fakeAuditLogRepo)(nil)
