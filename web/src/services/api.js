import axios from 'axios';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Employee API calls
export const employeeAPI = {
  // Get all employees with pagination
  getEmployees: async (page = 1, limit = 50) => {
    const response = await api.get('/employees', {
      params: { page, limit }
    });
    return response.data;
  },

  // Get single employee by ID
  getEmployee: async (id) => {
    const response = await api.get(`/employees/${id}`);
    return response.data;
  },

  // Create new employee
  createEmployee: async (employeeData) => {
    const response = await api.post('/employees', employeeData);
    return response.data;
  },

  // Update employee
  updateEmployee: async (id, employeeData) => {
    const response = await api.put(`/employees/${id}`, employeeData);
    return response.data;
  },

  // Delete employee
  deleteEmployee: async (id) => {
    await api.delete(`/employees/${id}`);
  },
};

// Analytics API calls
export const analyticsAPI = {
  // Get salary statistics by country
  getSalaryByCountry: async () => {
    const response = await api.get('/analytics/salary/by-country');
    return response.data;
  },

  // Get salary statistics by job title in a country
  getSalaryByJobTitleInCountry: async (country) => {
    const response = await api.get(`/analytics/salary/by-job-title/${encodeURIComponent(country)}`);
    return response.data;
  },

  // Get department salary insights
  getDepartmentInsights: async () => {
    const response = await api.get('/analytics/salary/department-insights');
    return response.data;
  },
};

// Health check
export const healthAPI = {
  getHealth: async () => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api;
