import axios, { AxiosResponse, AxiosInstance } from 'axios';
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

// Get API base URL from env or fallback
export const getApiBaseUrl = (): string => {
  return import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
};

// Create axios instance with auth headers
export const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: getApiBaseUrl(),
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Add request interceptor for JWT token
  client.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  // Add response interceptor for auth errors
  client.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );

  return client;
};

// Default api instance
const api = createApiClient();


// Employee API calls
export const employeeAPI = {
  // Get all employees with pagination and optional search
  getEmployees: async (page = 1, limit = 10, search?: string): Promise<EmployeesResponse> => {
    const params: Record<string, unknown> = { page, limit };
    if (search && search.trim()) {
      params.search = search.trim();
    }
    const response: AxiosResponse<EmployeesResponse> = await api.get('/employees', {
      params
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
    const response: AxiosResponse<HealthResponse> = await api.get('/api/health');
    return response.data;
  },
};

export default api;
