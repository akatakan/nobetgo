import axios from 'axios';
import type {
    Employee, ShiftType, Schedule, ScheduleRequest, EmployeeFormData,
    ShiftTypeFormData, Department, DepartmentFormData, Title, TitleFormData,
    TimeEntry, TimeEntryRequest, ClockInRequest, ClockOutRequest,
    Leave, LeaveRequest, LeaveType, LeaveBalance,
    OvertimeRule, PublicHoliday, OvertimeSummary,
    PendingApprovals, AuditLog,
    WorkHoursReport, AbsenceReport, EmployeeSummaryReport, TrendData,
    CostAnalysisReport, Notification, AuthResponse, PaginationParams, PaginationResult
} from '../types';

const API_BASE_URL = 'http://localhost:8081/api/v1';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Auth interceptor
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

// Standard Error format handler
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export const authApi = {
    login: (data: any) => api.post<AuthResponse>('/auth/login', data),
};

export const departmentApi = {
    list: () => api.get<Department[]>('/departments'),
    getById: (id: number) => api.get<Department>(`/departments/${id}`),
    create: (data: DepartmentFormData) => api.post<Department>('/departments', data),
    update: (id: number, data: Partial<DepartmentFormData>) => api.put<Department>(`/departments/${id}`, data),
    delete: (id: number) => api.delete(`/departments/${id}`),
};

export const titleApi = {
    list: () => api.get<Title[]>('/titles'),
    getById: (id: number) => api.get<Title>(`/titles/${id}`),
    create: (data: TitleFormData) => api.post<Title>('/titles', data),
    update: (id: number, data: Partial<TitleFormData>) => api.put<Title>(`/titles/${id}`, data),
    delete: (id: number) => api.delete(`/titles/${id}`),
};

export const employeeApi = {
    list: (params?: PaginationParams) => api.get<PaginationResult<Employee>>('/employees', { params }),
    getById: (id: number) => api.get<Employee>(`/employees/${id}`),
    create: (data: EmployeeFormData) => api.post<Employee>('/employees', data),
    update: (id: number, data: Partial<EmployeeFormData>) => api.put<Employee>(`/employees/${id}`, data),
    delete: (id: number) => api.delete(`/employees/${id}`),
    import: (file: File) => {
        const formData = new FormData();
        formData.append('file', file);
        return api.post('/employees/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    }
};

export const shiftTypeApi = {
    list: () => api.get<ShiftType[]>('/shift-types'),
    getById: (id: number) => api.get<ShiftType>(`/shift-types/${id}`),
    create: (data: ShiftTypeFormData) => api.post<ShiftType>('/shift-types', data),
    update: (id: number, data: Partial<ShiftTypeFormData>) => api.put<ShiftType>(`/shift-types/${id}`, data),
    delete: (id: number) => api.delete(`/shift-types/${id}`),
};

export const scheduleApi = {
    generate: (data: ScheduleRequest) => api.post<Schedule[]>('/schedules/generate', data),
    create: (data: Partial<Schedule>) => api.post<Schedule>('/schedules', data),
    getMonthly: (month: number, year: number, departmentId?: number) =>
        api.get<Schedule[]>(`/schedules?month=${month}&year=${year}${departmentId ? `&department_id=${departmentId}` : ''}`),
    update: (id: number, data: Partial<Schedule>) => api.put<Schedule>(`/schedules/${id}`, data),
    delete: (id: number) => api.delete(`/schedules/${id}`),
    clear: (month: number, year: number) => api.delete(`/schedules/clear?month=${month}&year=${year}`),
};

// ===== Time Entry (Puantaj) =====
export const timeEntryApi = {
    clockIn: (data: ClockInRequest) => api.post<TimeEntry>('/time-entries/clock-in', data),
    clockOut: (data: ClockOutRequest) => api.post<TimeEntry>('/time-entries/clock-out', data),
    create: (data: TimeEntryRequest) => api.post<TimeEntry>('/time-entries', data),
    list: (params: PaginationParams & { employee_id?: number; department_id?: number; start?: string; end?: string }) =>
        api.get<PaginationResult<TimeEntry>>('/time-entries', { params }),
    getById: (id: number) => api.get<TimeEntry>(`/time-entries/${id}`),
    update: (id: number, data: TimeEntryRequest) => api.put<TimeEntry>(`/time-entries/${id}`, data),
    delete: (id: number) => api.delete(`/time-entries/${id}`),
};

// ===== Leave (İzin) =====
export const leaveApi = {
    request: (data: LeaveRequest) => api.post<Leave>('/leaves', data),
    list: (params: PaginationParams & { employee_id?: number; department_id?: number; start?: string; end?: string }) =>
        api.get<PaginationResult<Leave>>('/leaves', { params }),
    getById: (id: number) => api.get<Leave>(`/leaves/${id}`),
    approve: (id: number, approver_id: number) =>
        api.post<Leave>(`/leaves/${id}/approve`, { approver_id }),
    reject: (id: number, approver_id: number) =>
        api.post<Leave>(`/leaves/${id}/reject`, { approver_id }),
    getBalance: (employee_id: number, year: number) =>
        api.get<LeaveBalance[]>('/leaves/balance', { params: { employee_id, year } }),
};

export const leaveTypeApi = {
    list: () => api.get<LeaveType[]>('/leave-types'),
    create: (data: Partial<LeaveType>) => api.post<LeaveType>('/leave-types', data),
    update: (id: number, data: Partial<LeaveType>) => api.put<LeaveType>(`/leave-types/${id}`, data),
    delete: (id: number) => api.delete(`/leave-types/${id}`),
};

// ===== Overtime (Mesai) =====
export const overtimeApi = {
    calculate: (employee_id: number, month: number, year: number) =>
        api.get<OvertimeSummary>('/overtime/calculate', { params: { employee_id, month, year } }),
    departmentSummary: (department_id: number, month: number, year: number) =>
        api.get<OvertimeSummary[]>('/overtime/summary', { params: { department_id, month, year } }),
};

export const overtimeRuleApi = {
    list: () => api.get<OvertimeRule[]>('/overtime-rules'),
    create: (data: Partial<OvertimeRule>) => api.post<OvertimeRule>('/overtime-rules', data),
    update: (id: number, data: Partial<OvertimeRule>) => api.put<OvertimeRule>(`/overtime-rules/${id}`, data),
    delete: (id: number) => api.delete(`/overtime-rules/${id}`),
};

export const publicHolidayApi = {
    list: (year: number) => api.get<PublicHoliday[]>('/public-holidays', { params: { year } }),
    create: (data: Partial<PublicHoliday>) => api.post<PublicHoliday>('/public-holidays', data),
    update: (id: number, data: Partial<PublicHoliday>) => api.put<PublicHoliday>(`/public-holidays/${id}`, data),
    delete: (id: number) => api.delete(`/public-holidays/${id}`),
};

// ===== Approval (Onay) =====
export const approvalApi = {
    getPending: () => api.get<PendingApprovals>('/approvals/pending'),
    approveTimeEntry: (id: number, approver_id: number) =>
        api.post(`/approvals/time-entry/${id}/approve`, { approver_id }),
    rejectTimeEntry: (id: number, approver_id: number) =>
        api.post(`/approvals/time-entry/${id}/reject`, { approver_id }),
    getAuditLogs: (entity_type: string, entity_id?: number) =>
        api.get<AuditLog[]>('/audit-logs', { params: { entity_type, entity_id } }),
};

// ===== Reports =====
export const reportApi = {
    workHours: (month: number, year: number, department_id?: number) =>
        api.get<WorkHoursReport>('/reports/work-hours', { params: { month, year, department_id } }),
    absences: (month: number, year: number, department_id?: number) =>
        api.get<AbsenceReport>('/reports/absences', { params: { month, year, department_id } }),
    employeeSummary: (employee_id: number, month: number, year: number) =>
        api.get<EmployeeSummaryReport>('/reports/employee-summary', { params: { employee_id, month, year } }),
    trends: (start_month: number, end_month: number, year: number, department_id?: number) =>
        api.get<TrendData[]>('/reports/trends', { params: { start_month, end_month, year, department_id } }),
    costAnalysis: (month: number, year: number, department_id?: number) =>
        api.get<CostAnalysisReport>('/reports/cost-analysis', { params: { month, year, department_id } }),
};

// ===== Notifications =====
export const notificationApi = {
    getUnread: (employee_id: number) =>
        api.get<Notification[]>('/notifications/unread', { params: { employee_id } }),
    markAsRead: (id: number) =>
        api.post(`/notifications/${id}/read`),
    markAllAsRead: (employee_id: number) =>
        api.post('/notifications/read-all', { employee_id }),
};

export default api;
