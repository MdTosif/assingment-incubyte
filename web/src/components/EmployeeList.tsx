import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Link } from 'react-router-dom';
import { Users, Loader2 } from 'lucide-react';
import { employeeAPI } from '../services/api';
import { Employee, EmployeesResponse } from '../types';
import { useDebounce } from '../hooks/useDebounce';
import { Button } from './ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import {
  EmployeeSearchFilter,
  EmployeeTable,
  EmployeePagination,
} from './employees';

const EmployeeList: React.FC = () => {
  // State management
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [totalItems, setTotalItems] = useState(0);
  const [pageSize, setPageSize] = useState(10);

  // Search state
  const [searchInput, setSearchInput] = useState('');
  const debouncedSearch = useDebounce(searchInput, 500);

  // Abort controller ref for cancelling pending requests
  const abortControllerRef = useRef<AbortController | null>(null);

  // Fetch employees with server-side pagination and search
  const fetchEmployees = useCallback(async (page: number, limit: number, search?: string) => {
    // Cancel any pending request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    abortControllerRef.current = new AbortController();

    try {
      setLoading(true);
      setError(null);

      const response: EmployeesResponse = await employeeAPI.getEmployees(page, limit, search);

      setEmployees(response.employees);
      setTotalPages(response.pages);
      setTotalItems(response.total);
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      setError('Failed to fetch employees. Please try again.');
      console.error('Error fetching employees:', err);
    } finally {
      setLoading(false);
      setInitialLoading(false);
    }
  }, []);

  // Effect to fetch data when page, pageSize, or debouncedSearch changes
  useEffect(() => {
    fetchEmployees(currentPage, pageSize, debouncedSearch);

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [currentPage, pageSize, debouncedSearch, fetchEmployees]);

  // Reset to first page when search changes
  useEffect(() => {
    setCurrentPage(1);
  }, [debouncedSearch]);

  // Handle delete
  const handleDelete = useCallback(async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this employee?')) {
      return;
    }

    try {
      await employeeAPI.deleteEmployee(id);
      fetchEmployees(currentPage, pageSize, debouncedSearch);
    } catch (err) {
      setError('Failed to delete employee');
      console.error('Error deleting employee:', err);
    }
  }, [currentPage, pageSize, debouncedSearch, fetchEmployees]);

  // Handle page change
  const handlePageChange = useCallback((page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  }, [totalPages]);

  // Handle page size change
  const handlePageSizeChange = useCallback((value: number) => {
    setPageSize(value);
    setCurrentPage(1);
  }, []);

  // Handle search change
  const handleSearchChange = useCallback((value: string) => {
    setSearchInput(value);
  }, []);

  if (initialLoading) {
    return (
      <div className="container mx-auto py-6">
        <Card>
          <CardContent className="flex justify-center items-center h-64">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6">
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-6 w-6" />
                Employees
              </CardTitle>
              <CardDescription>
                Manage your employees with server-side pagination and search.
              </CardDescription>
            </div>
            <Link to="/add-employee">
              <Button>Add Employee</Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Search and Filters */}
            <EmployeeSearchFilter
              searchInput={searchInput}
              pageSize={pageSize}
              onSearchChange={handleSearchChange}
              onPageSizeChange={handlePageSizeChange}
            />

            {/* Error Message */}
            {error && (
              <div className="flex items-center justify-between text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                <span>{error}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => fetchEmployees(currentPage, pageSize, debouncedSearch)}
                >
                  Retry
                </Button>
              </div>
            )}

            {/* Loading overlay */}
            {loading && (
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
              loading={loading}
              onPageChange={handlePageChange}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default EmployeeList;
