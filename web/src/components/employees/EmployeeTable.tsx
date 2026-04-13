import React, { useCallback } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../ui/table';
import { Employee } from '../../types';
import EmployeeActionsDropdown from './EmployeeActionsDropdown';

interface EmployeeTableProps {
  employees: Employee[];
  isSearching: boolean;
  onDelete: (id: number) => void;
}

const EmployeeTable: React.FC<EmployeeTableProps> = ({
  employees,
  isSearching,
  onDelete,
}) => {
  const formatSalary = useCallback((salary: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(salary);
  }, []);

  if (employees.length === 0) {
    return (
      <div className="rounded-md border">
        <Table>
          <TableBody>
            <TableRow>
              <TableCell colSpan={7} className="text-center py-12">
                <div className="text-muted-foreground">
                  {isSearching
                    ? 'No employees found matching your search.'
                    : 'No employees found.'}
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    );
  }

  return (
    <div className="rounded-md border overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="min-w-[150px]">Name</TableHead>
            <TableHead className="min-w-[200px]">Email</TableHead>
            <TableHead className="min-w-[180px]">Job Title</TableHead>
            <TableHead className="min-w-[120px]">Country</TableHead>
            <TableHead className="min-w-[120px]">Salary</TableHead>
            <TableHead className="min-w-[120px]">Department</TableHead>
            <TableHead className="text-right w-[120px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {employees.map((employee) => (
            <TableRow key={employee.id}>
              <TableCell className="font-medium">
                {employee.firstName} {employee.lastName}
              </TableCell>
              <TableCell>{employee.email}</TableCell>
              <TableCell>{employee.jobTitle}</TableCell>
              <TableCell>{employee.country}</TableCell>
              <TableCell>{formatSalary(employee.salary)}</TableCell>
              <TableCell>{employee.department}</TableCell>
              <TableCell className="text-right whitespace-nowrap">
                <EmployeeActionsDropdown
                  employeeId={employee.id}
                  onDelete={onDelete}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default EmployeeTable;
