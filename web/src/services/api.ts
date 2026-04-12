import axios, { AxiosResponse } from 'axios';
import {
  Employee,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  EmployeesResponse,
  CountrySalaryStats,
  JobTitleSalaryStats,
  DepartmentSalaryStats,
  HealthResponse,
  LoginRequest,
  LoginResponse,
  ChangePasswordRequest,
  User
} from '../types';

const API_BASE_URL = '/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, clear local storage and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Employee API calls
export const employeeAPI = {
  // Get all employees with pagination
  getEmployees: async (page = 1, limit = 50): Promise<EmployeesResponse> => {
    const response: AxiosResponse<EmployeesResponse> = await api.get('/employees', {
      params: { page, limit }
    });
    return response.data;
  },

  // Get single employee by ID
  getEmployee: async (id: number): Promise<Employee> => {
    const response: AxiosResponse<Employee> = await api.get(`/employees/${id}`);
    return response.data;
  },

  // Create new employee
  createEmployee: async (employeeData: CreateEmployeeRequest): Promise<Employee> => {
    const response: AxiosResponse<Employee> = await api.post('/employees', employeeData);
    return response.data;
  },

  // Update employee
  updateEmployee: async (id: number, employeeData: UpdateEmployeeRequest): Promise<Employee> => {
    const response: AxiosResponse<Employee> = await api.put(`/employees/${id}`, employeeData);
    return response.data;
  },

  // Delete employee
  deleteEmployee: async (id: number): Promise<void> => {
    await api.delete(`/employees/${id}`);
  },
};

// Analytics API calls
export const analyticsAPI = {
  // Get salary statistics by country
  getSalaryByCountry: async (): Promise<CountrySalaryStats[]> => {
    const response: AxiosResponse<CountrySalaryStats[]> = await api.get('/analytics/salary/by-country');
    return response.data;
  },

  // Get salary statistics by job title in a country
  getSalaryByJobTitleInCountry: async (country: string): Promise<JobTitleSalaryStats[]> => {
    const response: AxiosResponse<JobTitleSalaryStats[]> = await api.get(
      `/analytics/salary/by-job-title/${encodeURIComponent(country)}`
    );
    return response.data;
  },

  // Get department salary insights
  getDepartmentInsights: async (): Promise<DepartmentSalaryStats[]> => {
    const response: AxiosResponse<DepartmentSalaryStats[]> = await api.get('/analytics/salary/department-insights');
    return response.data;
  },
};

// Authentication API calls
export const authAPI = {
  // Login
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response: AxiosResponse<LoginResponse> = await api.post('/auth/login', credentials);
    return response.data;
  },

  // Get current user
  getMe: async (): Promise<User> => {
    const response: AxiosResponse<User> = await api.get('/auth/me');
    return response.data;
  },

  // Logout
  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },

  // Change password
  changePassword: async (passwordData: ChangePasswordRequest): Promise<void> => {
    await api.post('/auth/change-password', passwordData);
  },
};

// Health check
export const healthAPI = {
  getHealth: async (): Promise<HealthResponse> => {
    const response: AxiosResponse<HealthResponse> = await api.get('/health');
    return response.data;
  },
};

export default api;
