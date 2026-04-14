import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  Building2,
  Users,
  TrendingUp,
  DollarSign,
  LogOut,
  Menu,
  Plus,
  Loader2,
  BarChart3
} from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Avatar, AvatarFallback } from './ui/avatar';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { Sheet, SheetContent, SheetTrigger } from './ui/sheet';
import { employeeAPI, analyticsAPI, authAPI } from '../services/api';
import { Employee, DepartmentSalaryStats, CountrySalaryStats, User, EmployeesResponse } from '../types';
import { useDebounce } from '../hooks/useDebounce';
import { EmployeeSearchFilter, EmployeeTable, EmployeePagination } from './employees';

interface Analytics {
  totalEmployees: number;
  avgSalary: number;
  totalSalaryExpense: number;
  departments: DepartmentSalaryStats[];
  countryStats: CountrySalaryStats[];
}

const Dashboard: React.FC = () => {
  // ==================== State ====================
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [analytics, setAnalytics] = useState<Analytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  // Employee list state with server-side pagination
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [totalItems, setTotalItems] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [searchInput, setSearchInput] = useState('');
  const [tableLoading, setTableLoading] = useState(false);
  const debouncedSearch = useDebounce(searchInput, 500);

  // ==================== Data Fetching ====================

  /** Fetch employees with server-side pagination and optional search */
  const fetchEmployees = useCallback(async (page: number, limit: number, search?: string) => {
    console.log('[Dashboard] fetchEmployees called:', { page, limit, search });
    try {
      setTableLoading(true);
      const response: EmployeesResponse = await employeeAPI.getEmployees(page, limit, search);
      console.log('[Dashboard] API response:', { total: response.total, pages: response.pages, count: response.employees.length });
      setEmployees(response.employees);
      setTotalPages(response.pages);
      setTotalItems(response.total);
    } catch (error) {
      console.error('Failed to fetch employees:', error);
    } finally {
      setTableLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    const loadData = async () => {
      try {
        const userData = localStorage.getItem('user');
        if (!userData) {
          navigate('/login');
          return;
        }
        setUser(JSON.parse(userData));

        const [employeesData, salaryData, departmentData] = await Promise.all([
          employeeAPI.getEmployees(1, 10),
          analyticsAPI.getSalaryByCountry(),
          analyticsAPI.getDepartmentInsights()
        ]);

        setEmployees(employeesData.employees || []);
        setTotalPages(employeesData.pages || 0);
        setTotalItems(employeesData.total || 0);

        // Calculate analytics from API responses
        const totalEmployees = employeesData.total || 0;
        const avgSalary = salaryData.length > 0 ?
          salaryData.reduce((sum: number, country: CountrySalaryStats) => sum + country.average, 0) / salaryData.length : 0;
        const totalSalaryExpense = employeesData.employees?.reduce((sum: number, emp: Employee) => sum + emp.salary, 0) || 0;

        setAnalytics({
          totalEmployees,
          avgSalary,
          totalSalaryExpense,
          departments: departmentData || [],
          countryStats: salaryData || []
        });
      } catch (error) {
        console.error('Failed to load data:', error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [navigate]);

  // Fetch when pagination/search changes
  useEffect(() => {
    console.log('[Dashboard] useEffect triggered:', { loading, currentPage, pageSize, debouncedSearch });
    if (!loading) {
      fetchEmployees(currentPage, pageSize, debouncedSearch);
    }
  }, [currentPage, pageSize, debouncedSearch, fetchEmployees, loading]);

  // Reset to first page when search changes
  useEffect(() => {
    setCurrentPage(1);
  }, [debouncedSearch]);

  // ==================== Handlers ====================

  /** Handle user logout - calls API and clears local storage */
  const handleLogout = async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      navigate('/login');
    }
  };

  /** Handle employee deletion with confirmation dialog */
  const handleDelete = useCallback(async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this employee?')) {
      return;
    }
    try {
      await employeeAPI.deleteEmployee(id);
      fetchEmployees(currentPage, pageSize, debouncedSearch);
    } catch (error) {
      console.error('Failed to delete employee:', error);
    }
  }, [currentPage, pageSize, debouncedSearch, fetchEmployees]);

  /** Handle pagination page changes */
  const handlePageChange = useCallback((page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  }, [totalPages]);

  /** Handle page size changes - resets to page 1 */
  const handlePageSizeChange = useCallback((value: number) => {
    setPageSize(value);
    setCurrentPage(1);
  }, []);

  /** Handle search input changes */
  const handleSearchChange = useCallback((value: string) => {
    setSearchInput(value);
  }, []);

  // ==================== Render ====================

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Building2 className="h-8 w-8 text-primary" />
            <h1 className="text-xl font-semibold">Salary Manager</h1>
          </div>
          
          {/* Mobile Menu */}
          <Sheet open={mobileMenuOpen} onOpenChange={setMobileMenuOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden">
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent side="right">
              <div className="flex flex-col space-y-4 mt-8">
                <div className="flex items-center space-x-3">
                  <Avatar>
                    <AvatarFallback>
                      {user?.firstName?.[0]}{user?.lastName?.[0]}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium">{user?.firstName} {user?.lastName}</p>
                    <p className="text-sm text-muted-foreground">{user?.email}</p>
                  </div>
                </div>
                <Button variant="outline" className="w-full" asChild>
                  <Link to="/analytics">
                    <BarChart3 className="mr-2 h-4 w-4" />
                    Analytics
                  </Link>
                </Button>
                <Button onClick={handleLogout} variant="outline" className="w-full">
                  <LogOut className="mr-2 h-4 w-4" />
                  Logout
                </Button>
              </div>
            </SheetContent>
          </Sheet>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center space-x-4">
            <Button variant="ghost" asChild>
              <Link to="/analytics">
                <BarChart3 className="mr-2 h-4 w-4" />
                Analytics
              </Link>
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback>
                      {user?.firstName?.[0]}{user?.lastName?.[0]}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <div className="flex items-center justify-start gap-2 p-2">
                  <div className="flex flex-col space-y-1 leading-none">
                    <p className="font-medium">{user?.firstName} {user?.lastName}</p>
                    <p className="w-[200px] truncate text-sm text-muted-foreground">
                      {user?.email}
                    </p>
                  </div>
                </div>
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        {/* Stats Cards */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Employees</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{analytics?.totalEmployees || 0}</div>
              <p className="text-xs text-muted-foreground">Active employees</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Average Salary</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">${analytics?.avgSalary?.toFixed(0) || 0}</div>
              <p className="text-xs text-muted-foreground">Per employee</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Expense</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">${analytics?.totalSalaryExpense?.toFixed(0) || 0}</div>
              <p className="text-xs text-muted-foreground">Monthly payroll</p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Departments</CardTitle>
              <Building2 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{analytics?.departments?.length || 0}</div>
              <p className="text-xs text-muted-foreground">Active departments</p>
            </CardContent>
          </Card>
        </div>

        {/* Quick Salary Insights Preview */}
        <Card className="mb-8">
          <CardHeader>
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div className="flex items-center gap-2">
                <BarChart3 className="h-5 w-5 text-primary" />
                <CardTitle>Quick Insights</CardTitle>
              </div>
              <Button variant="outline" size="sm" asChild>
                <Link to="/analytics">
                  View Full Analytics
                  <TrendingUp className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
            <CardDescription>
              Top countries by average salary
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-muted">
                  <tr>
                    <th className="px-4 py-3 text-left font-medium">Country</th>
                    <th className="px-4 py-3 text-right font-medium">Employees</th>
                    <th className="px-4 py-3 text-right font-medium">Min</th>
                    <th className="px-4 py-3 text-right font-medium">Max</th>
                    <th className="px-4 py-3 text-right font-medium">Average</th>
                  </tr>
                </thead>
                <tbody>
                  {analytics?.countryStats?.slice(0, 5).map((stat) => (
                    <tr key={stat.country} className="border-t">
                      <td className="px-4 py-3 font-medium">{stat.country}</td>
                      <td className="px-4 py-3 text-right">{stat.count}</td>
                      <td className="px-4 py-3 text-right">${stat.min.toLocaleString()}</td>
                      <td className="px-4 py-3 text-right">${stat.max.toLocaleString()}</td>
                      <td className="px-4 py-3 text-right font-semibold text-primary">
                        ${stat.average.toFixed(0)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        {/* Employee Table */}
        <Card>
          <CardHeader className="pb-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div className="space-y-1">
                <CardTitle className="text-2xl">Employees</CardTitle>
                <CardDescription className="text-base">
                  Manage your team members
                </CardDescription>
              </div>
              <div className="flex items-center gap-2">
                <Button onClick={() => navigate('/add-employee')}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Employee
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="pt-0 space-y-4">
            {/* Search and Filters */}
            <EmployeeSearchFilter
              searchInput={searchInput}
              pageSize={pageSize}
              onSearchChange={handleSearchChange}
              onPageSizeChange={handlePageSizeChange}
            />

            {/* Loading overlay */}
            {tableLoading && (
              <div className="flex justify-center items-center py-4">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
              </div>
            )}

            {/* Employee Table */}
            <EmployeeTable
              employees={employees}
              isSearching={!!debouncedSearch}
              onDelete={handleDelete}
            />

            {/* Pagination */}
            <EmployeePagination
              currentPage={currentPage}
              totalPages={totalPages}
              totalItems={totalItems}
              pageSize={pageSize}
              loading={tableLoading}
              onPageChange={handlePageChange}
            />
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default Dashboard;
