import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const employeeApi = {
    list: () => api.get('/employees'),
    create: (data) => api.post('/employees', data),
};

export const shiftTypeApi = {
    list: () => api.get('/shift-types'),
    create: (data) => api.post('/shift-types', data),
};

export const scheduleApi = {
    generate: (data) => api.post('/schedules/generate', data),
    get: (month, year) => api.get(`/attendance/reports?month=${month}&year=${year}`), // Reusing attendance report for view
};

export default api;
