import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Pencil,
  Trash2,
  Eye,
  MoreHorizontal,
} from 'lucide-react';
import { Button } from '../ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';

interface EmployeeActionsDropdownProps {
  employeeId: number;
  onDelete: (id: number) => void;
}

/**
 * EmployeeActionsDropdown Component
 *
 * Dropdown menu providing actions for an employee record:
 * - View/Edit: Navigate to employee edit page
 * - Delete: Trigger delete confirmation
 */
const EmployeeActionsDropdown: React.FC<EmployeeActionsDropdownProps> = ({
  employeeId,
  onDelete,
}) => {
  // ==================== Hooks ====================
  const navigate = useNavigate();

  // ==================== Handlers ====================

  /** Navigate to employee detail/edit page */
  const handleView = () => {
    navigate(`/edit-employee/${employeeId}`);
  };

  /** Navigate to employee edit page */
  const handleEdit = () => {
    navigate(`/edit-employee/${employeeId}`);
  };

  /** Trigger delete callback for this employee */
  const handleDelete = () => {
    onDelete(employeeId);
  };

  // ==================== Render ====================

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="bg-slate-800 border-slate-500 text-slate-300 hover:bg-slate-700 hover:text-white"
        >
          <MoreHorizontal className="h-4 w-4 mr-1" />
          Actions
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40 bg-slate-900 border-slate-700">
        <DropdownMenuItem
          onClick={handleView}
          className="cursor-pointer text-blue-400 focus:text-blue-400 focus:bg-blue-950"
        >
          <Eye className="h-4 w-4 mr-2" />
          View Details
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={handleEdit}
          className="cursor-pointer text-slate-300 focus:text-slate-300 focus:bg-slate-800"
        >
          <Pencil className="h-4 w-4 mr-2" />
          Edit
        </DropdownMenuItem>
        <DropdownMenuSeparator className="bg-slate-700" />
        <DropdownMenuItem
          onClick={handleDelete}
          className="cursor-pointer text-red-400 focus:text-red-400 focus:bg-red-950"
        >
          <Trash2 className="h-4 w-4 mr-2" />
          Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default EmployeeActionsDropdown;
