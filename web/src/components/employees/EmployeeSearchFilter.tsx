import React, { useCallback } from 'react';
import { Search, X } from 'lucide-react';
import { Input } from '../ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select';

interface EmployeeSearchFilterProps {
  searchInput: string;
  pageSize: number;
  onSearchChange: (value: string) => void;
  onPageSizeChange: (value: number) => void;
}

const EmployeeSearchFilter: React.FC<EmployeeSearchFilterProps> = ({
  searchInput,
  pageSize,
  onSearchChange,
  onPageSizeChange,
}) => {
  const handleClearSearch = useCallback(() => {
    onSearchChange('');
  }, [onSearchChange]);

  const handlePageSizeChange = useCallback(
    (value: string) => {
      onPageSizeChange(Number(value));
    },
    [onPageSizeChange]
  );

  return (
    <div className="flex flex-col sm:flex-row gap-4">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search by name, email, job title, country, or department..."
          value={searchInput}
          onChange={(e) => {
            console.log('[SearchFilter] input changed:', e.target.value);
            onSearchChange(e.target.value);
          }}
          className="pl-10 pr-10"
        />
        {searchInput && (
          <button
            onClick={handleClearSearch}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
      <Select value={String(pageSize)} onValueChange={handlePageSizeChange}>
        <SelectTrigger className="w-[140px]">
          <SelectValue placeholder="Page size" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="5">5 per page</SelectItem>
          <SelectItem value="10">10 per page</SelectItem>
          <SelectItem value="25">25 per page</SelectItem>
          <SelectItem value="50">50 per page</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
};

export default EmployeeSearchFilter;
