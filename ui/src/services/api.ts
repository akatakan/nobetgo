import axios from 'axios';
import type { Employee, ShiftType, Schedule, Attendance, ScheduleRequest, EmployeeFormData, ShiftTypeFormData, Department, DepartmentFormData, Title, TitleFormData } from '../types';

const API_BASE_URL = 'http://localhost:8081/api/v1';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

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
    list: () => api.get<Employee[]>('/employees'),
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

export const attendanceApi = {
    log: (data: { schedule_id: number; actual_start_time: string; actual_end_time: string; notes?: string }) =>
        api.post<Attendance>('/attendance', data),
    update: (id: number, data: any) => api.put(`/attendance/${id}`, data),
    getReport: (month: number, year: number, departmentId?: number) =>
        api.get<Attendance[]>(`/attendance/reports?month=${month}&year=${year}${departmentId ? `&department_id=${departmentId}` : ''}`),
};

export default api;
