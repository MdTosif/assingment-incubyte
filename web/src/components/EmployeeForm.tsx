import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, UserPlus } from 'lucide-react';
import { employeeAPI } from '../services/api';
import { Employee, CreateEmployeeRequest, UpdateEmployeeRequest } from '../types';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';

/**
 * EmployeeForm Component
 * 
 * Provides a form for creating new employees or editing existing ones.
 * Uses controlled inputs with validation and supports all employee fields.
 */
const EmployeeForm: React.FC = () => {
  // ==================== State & Refs ====================
  const { id } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const isEditing = !!id;

  const [employee, setEmployee] = useState<Employee | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState<CreateEmployeeRequest>({
    firstName: '',
    lastName: '',
    email: '',
    jobTitle: '',
    country: '',
    salary: 0,
    department: '',
  });

  // ==================== Data Constants ====================

  // Predefined list of job titles for the dropdown
  const jobTitles = [
    'Software Engineer',
    'Senior Software Engineer',
    'Lead Software Engineer',
    'Principal Software Engineer',
    'Software Architect',
    'DevOps Engineer',
    'Senior DevOps Engineer',
    'Site Reliability Engineer',
    'Product Manager',
    'Senior Product Manager',
    'UX Designer',
    'Senior UX Designer',
    'UI Designer',
    'Frontend Developer',
    'Senior Frontend Developer',
    'Backend Developer',
    'Senior Backend Developer',
    'Full Stack Developer',
    'Senior Full Stack Developer',
    'Data Scientist',
    'Senior Data Scientist',
    'Machine Learning Engineer',
    'QA Engineer',
    'Senior QA Engineer',
    'Technical Writer',
    'Project Manager',
    'Scrum Master',
    'Business Analyst',
    'Systems Administrator',
    'Network Engineer',
    'Security Engineer',
    'Cloud Engineer',
    'Database Administrator',
    'Mobile Developer',
    'Senior Mobile Developer',
    'Engineering Manager',
    'Director of Engineering',
    'CTO',
    'VP of Engineering',
    'HR Manager',
    'Senior HR Manager',
    'Recruiter',
    'Senior Recruiter',
    'Marketing Manager',
    'Senior Marketing Manager',
    'Sales Manager',
    'Senior Sales Manager',
    'Account Manager',
    'Senior Account Manager',
    'Financial Analyst',
    'Senior Financial Analyst',
    'Operations Manager',
    'Senior Operations Manager',
    'Customer Success Manager',
    'Technical Support Engineer',
    'Senior Technical Support Engineer',
  ];

  const countries = [
    'United States',
    'United Kingdom',
    'Canada',
    'Germany',
    'France',
    'India',
    'Japan',
    'China',
    'Australia',
    'Netherlands',
    'Sweden',
    'Norway',
    'Denmark',
    'Finland',
    'Switzerland',
    'Austria',
    'Belgium',
    'Ireland',
    'Spain',
    'Italy',
    'Poland',
    'Czech Republic',
    'Hungary',
    'Romania',
    'Bulgaria',
    'Greece',
    'Portugal',
    'Turkey',
    'Israel',
    'UAE',
    'Saudi Arabia',
    'South Africa',
    'Brazil',
    'Argentina',
    'Mexico',
    'Chile',
    'Colombia',
    'Peru',
    'South Korea',
    'Singapore',
    'Malaysia',
    'Thailand',
    'Indonesia',
    'Philippines',
    'Vietnam',
    'New Zealand',
    'Russia',
    'Ukraine',
    'Egypt',
    'Nigeria',
    'Kenya',
    'Ghana',
    'Morocco',
  ];

  // Predefined list of departments for the dropdown
  const departments = [
    'Engineering',
    'Product',
    'Design',
    'Marketing',
    'Sales',
    'HR',
    'Finance',
    'Operations',
    'Customer Success',
    'Legal',
    'IT',
    'Data Science',
    'Security',
    'Infrastructure',
    'Mobile',
    'Web',
    'Backend',
    'Frontend',
    'DevOps',
    'QA',
  ];

  // ==================== Effects ====================

  /** Fetch employee data when editing existing employee */
  useEffect(() => {
    if (isEditing) {
      fetchEmployee();
    }
  }, [id, isEditing]);

  // ==================== Data Fetching ====================

  /** Load employee data from API for editing */
  const fetchEmployee = async () => {
    try {
      setLoading(true);
      const employeeData = await employeeAPI.getEmployee(Number(id));
      setEmployee(employeeData);
      setFormData({
        firstName: employeeData.firstName,
        lastName: employeeData.lastName,
        email: employeeData.email,
        jobTitle: employeeData.jobTitle,
        country: employeeData.country,
        salary: employeeData.salary,
        department: employeeData.department,
      });
    } catch (err) {
      setError('Failed to fetch employee');
      console.error('Error fetching employee:', err);
    } finally {
      setLoading(false);
    }
  };

  // ==================== Form Handlers ====================

  /** Handle input field changes and update form state */
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'salary' ? parseFloat(value) || 0 : value
    }));
  };

  /** Handle form submission - creates or updates employee */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      setLoading(true);

      if (isEditing) {
        const updateData: UpdateEmployeeRequest = {
          firstName: formData.firstName,
          lastName: formData.lastName,
          email: formData.email,
          jobTitle: formData.jobTitle,
          country: formData.country,
          salary: formData.salary,
          department: formData.department,
        };
        await employeeAPI.updateEmployee(Number(id), updateData);
      } else {
        await employeeAPI.createEmployee(formData);
      }

      navigate('/');
    } catch (err) {
      setError(isEditing ? 'Failed to update employee' : 'Failed to create employee');
      console.error('Error saving employee:', err);
    } finally {
      setLoading(false);
    }
  };

  // ==================== Render ====================

  if (loading && isEditing) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6">
      <Card>
        <CardHeader>
          <div className="flex items-center gap-4 mb-4">
            <Link to="/" className="flex items-center text-sm text-muted-foreground hover:text-foreground">
              <ArrowLeft className="h-4 w-4 mr-1" />
              Back to Employees
            </Link>
          </div>
          <CardTitle className="flex items-center gap-2">
            <UserPlus className="h-6 w-6" />
            {isEditing ? 'Edit Employee' : 'Add New Employee'}
          </CardTitle>
          <CardDescription>
            {isEditing 
              ? 'Update the employee information below.'
              : 'Fill in the information below to add a new employee to the system.'
            }
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                {error}
              </div>
            )}

            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name *</Label>
                <Input
                  id="firstName"
                  name="firstName"
                  required
                  value={formData.firstName}
                  onChange={handleInputChange}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name *</Label>
                <Input
                  id="lastName"
                  name="lastName"
                  required
                  value={formData.lastName}
                  onChange={handleInputChange}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email *</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleInputChange}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="jobTitle">Job Title *</Label>
                <Select
                  value={formData.jobTitle}
                  onValueChange={(value) => handleInputChange({ target: { name: 'jobTitle', value } } as any)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a job title" />
                  </SelectTrigger>
                  <SelectContent>
                    {jobTitles.map(title => (
                      <SelectItem key={title} value={title}>{title}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="country">Country *</Label>
                <Select
                  value={formData.country}
                  onValueChange={(value) => handleInputChange({ target: { name: 'country', value } } as any)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a country" />
                  </SelectTrigger>
                  <SelectContent>
                    {countries.map(country => (
                      <SelectItem key={country} value={country}>{country}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="salary">Salary (USD) *</Label>
                <Input
                  id="salary"
                  name="salary"
                  type="number"
                  required
                  min="0"
                  step="0.01"
                  value={formData.salary}
                  onChange={handleInputChange}
                />
              </div>

              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="department">Department *</Label>
                <Select
                  value={formData.department}
                  onValueChange={(value) => handleInputChange({ target: { name: 'department', value } } as any)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a department" />
                  </SelectTrigger>
                  <SelectContent>
                    {departments.map(dept => (
                      <SelectItem key={dept} value={dept}>{dept}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="flex justify-end gap-3">
              <Link to="/">
                <Button variant="outline">
                  Cancel
                </Button>
              </Link>
              <Button type="submit" disabled={loading}>
                {loading ? 'Saving...' : (isEditing ? 'Update Employee' : 'Create Employee')}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default EmployeeForm;
