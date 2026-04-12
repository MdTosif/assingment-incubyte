import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { employeeAPI } from '../services/api';
import { Employee, CreateEmployeeRequest, UpdateEmployeeRequest } from '../types';

const EmployeeForm: React.FC = () => {
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

  useEffect(() => {
    if (isEditing) {
      fetchEmployee();
    }
  }, [id, isEditing]);

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

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'salary' ? parseFloat(value) || 0 : value
    }));
  };

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

  if (loading && isEditing) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="sm:flex sm:items-center mb-6">
        <div className="sm:flex-auto">
          <Link
            to="/"
            className="inline-flex items-center text-sm text-gray-500 hover:text-gray-700 mb-4"
          >
            <ArrowLeftIcon className="h-4 w-4 mr-1" />
            Back to Employees
          </Link>
          <h1 className="text-2xl font-semibold text-gray-900">
            {isEditing ? 'Edit Employee' : 'Add New Employee'}
          </h1>
          <p className="mt-2 text-sm text-gray-700">
            {isEditing 
              ? 'Update the employee information below.'
              : 'Fill in the information below to add a new employee to the system.'
            }
          </p>
        </div>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4 mb-6">
          <div className="text-sm text-red-700">{error}</div>
        </div>
      )}

      <div className="bg-white shadow sm:rounded-lg">
        <form onSubmit={handleSubmit} className="px-4 py-5 sm:p-6">
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div>
              <label htmlFor="firstName" className="block text-sm font-medium text-gray-700">
                First Name *
              </label>
              <input
                type="text"
                id="firstName"
                name="firstName"
                required
                value={formData.firstName}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <div>
              <label htmlFor="lastName" className="block text-sm font-medium text-gray-700">
                Last Name *
              </label>
              <input
                type="text"
                id="lastName"
                name="lastName"
                required
                value={formData.lastName}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email *
              </label>
              <input
                type="email"
                id="email"
                name="email"
                required
                value={formData.email}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <div>
              <label htmlFor="jobTitle" className="block text-sm font-medium text-gray-700">
                Job Title *
              </label>
              <select
                id="jobTitle"
                name="jobTitle"
                required
                value={formData.jobTitle}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                <option value="">Select a job title</option>
                {jobTitles.map(title => (
                  <option key={title} value={title}>{title}</option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="country" className="block text-sm font-medium text-gray-700">
                Country *
              </label>
              <select
                id="country"
                name="country"
                required
                value={formData.country}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                <option value="">Select a country</option>
                {countries.map(country => (
                  <option key={country} value={country}>{country}</option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="salary" className="block text-sm font-medium text-gray-700">
                Salary (USD) *
              </label>
              <input
                type="number"
                id="salary"
                name="salary"
                required
                min="0"
                step="0.01"
                value={formData.salary}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <div className="sm:col-span-2">
              <label htmlFor="department" className="block text-sm font-medium text-gray-700">
                Department *
              </label>
              <select
                id="department"
                name="department"
                required
                value={formData.department}
                onChange={handleInputChange}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                <option value="">Select a department</option>
                {departments.map(dept => (
                  <option key={dept} value={dept}>{dept}</option>
                ))}
              </select>
            </div>
          </div>

          <div className="mt-6 flex justify-end">
            <Link
              to="/"
              className="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 mr-3"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {loading ? 'Saving...' : (isEditing ? 'Update Employee' : 'Create Employee')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EmployeeForm;
