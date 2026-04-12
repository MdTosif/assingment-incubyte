export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  jobTitle: string;
  country: string;
  salary: number;
  department: string;
  hireDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateEmployeeRequest {
  firstName: string;
  lastName: string;
  email: string;
  jobTitle: string;
  country: string;
  salary: number;
  department: string;
}

export interface UpdateEmployeeRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  jobTitle?: string;
  country?: string;
  salary?: number;
  department?: string;
}

export interface EmployeesResponse {
  employees: Employee[];
  total: number;
  page: number;
  limit: number;
  pages: number;
}

export interface CountrySalaryStats {
  country: string;
  min: number;
  max: number;
  average: number;
  count: number;
}

export interface JobTitleSalaryStats {
  jobTitle: string;
  average: number;
  count: number;
}

export interface DepartmentSalaryStats {
  department: string;
  min: number;
  max: number;
  average: number;
  count: number;
}

export interface HealthResponse {
  status: string;
}

// Authentication types
export interface User {
  id: number;
  email: string;
  role: string;
  firstName: string;
  lastName: string;
  isActive: boolean;
  lastLogin?: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expiresAt: string;
  user: User;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}
