import axios from 'axios';
import type { Employee, ShiftType, Schedule, Attendance, ScheduleRequest, EmployeeFormData, ShiftTypeFormData } from '../types';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const employeeApi = {
    list: () => api.get<Employee[]>('/employees'),
    getById: (id: number) => api.get<Employee>(`/employees/${id}`),
    create: (data: EmployeeFormData) => api.post<Employee>('/employees', data),
    update: (id: number, data: Partial<EmployeeFormData>) => api.put<Employee>(`/employees/${id}`, data),
    delete: (id: number) => api.delete(`/employees/${id}`),
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
    getMonthly: (month: number, year: number) => api.get<Schedule[]>(`/schedules?month=${month}&year=${year}`),
    update: (id: number, data: Partial<Schedule>) => api.put<Schedule>(`/schedules/${id}`, data),
};

export const attendanceApi = {
    log: (data: { schedule_id: number; actual_start_time: string; actual_end_time: string; notes?: string }) =>
        api.post<Attendance>('/attendance', data),
    getReport: (month: number, year: number) =>
        api.get<Attendance[]>(`/attendance/reports?month=${month}&year=${year}`),
};

export default api;
