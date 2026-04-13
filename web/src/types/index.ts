// ==================== Entity Types ====================

/** Employee entity representing an employee in the system */
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

// ==================== Request Types ====================

/** Request body for creating a new employee */
export interface CreateEmployeeRequest {
  firstName: string;
  lastName: string;
  email: string;
  jobTitle: string;
  country: string;
  salary: number;
  department: string;
}

/** Request body for updating an existing employee (all fields optional) */
export interface UpdateEmployeeRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  jobTitle?: string;
  country?: string;
  salary?: number;
  department?: string;
}

// ==================== Response Types ====================

/** Paginated response for employee list endpoints */
export interface EmployeesResponse {
  employees: Employee[];
  total: number;
  page: number;
  limit: number;
  pages: number;
}

// ==================== Statistics Types ====================

/** Salary statistics aggregated by country */
export interface CountrySalaryStats {
  country: string;
  min: number;
  max: number;
  average: number;
  count: number;
}

/** Salary statistics aggregated by job title within a country */
export interface JobTitleSalaryStats {
  jobTitle: string;
  average: number;
  count: number;
}

/** Salary statistics aggregated by department */
export interface DepartmentSalaryStats {
  department: string;
  min: number;
  max: number;
  average: number;
  count: number;
}

/** Health check response from the backend */
export interface HealthResponse {
  status: string;
}

// ==================== User Types ====================

/** Authenticated user (HR or Admin) */
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

/** Login request with email and password */
export interface LoginRequest {
  email: string;
  password: string;
}

/** Login response with JWT token and user data */
export interface LoginResponse {
  token: string;
  expiresAt: string;
  user: User;
}

/** Request body for changing user password */
export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}
