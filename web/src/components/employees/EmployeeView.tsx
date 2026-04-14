import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft,
  User,
  Mail,
  Briefcase,
  Globe,
  DollarSign,
  Building,
  Calendar,
  Clock,
  Edit,
  Trash2,
  Loader2,
  TrendingUp,
} from 'lucide-react';
import { employeeAPI, analyticsAPI } from '../../services/api';
import { Employee, CountrySalaryStats, JobTitleSalaryStats } from '../../types';
import { Button } from '../ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';

/**
 * EmployeeView Component
 *
 * Displays read-only details of a single employee including:
 * - Personal information (name, email)
 * - Job details (title, department, country)
 * - Salary information
 * - Employment dates
 * - Salary insights for their country and job title
 */
const EmployeeView: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [employee, setEmployee] = useState<Employee | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [countryStats, setCountryStats] = useState<CountrySalaryStats | null>(null);
  const [jobTitleStats, setJobTitleStats] = useState<JobTitleSalaryStats | null>(null);

  useEffect(() => {
    if (id) {
      fetchEmployeeData();
    }
  }, [id]);

  const fetchEmployeeData = async () => {
    try {
      setLoading(true);
      setError(null);

      const employeeData = await employeeAPI.getEmployee(Number(id));
      setEmployee(employeeData);

      // Fetch salary insights for the employee's country
      if (employeeData.country) {
        await fetchSalaryInsights(employeeData.country, employeeData.jobTitle);
      }
    } catch (err) {
      setError('Failed to fetch employee details');
      console.error('Error fetching employee:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchSalaryInsights = async (country: string, jobTitle: string) => {
    try {
      const [countryData, jobTitleData] = await Promise.all([
        analyticsAPI.getSalaryByCountry(),
        analyticsAPI.getSalaryByJobTitleInCountry(country),
      ]);

      // Find country stats
      const countryStat = countryData.find((stat) => stat.country === country);
      setCountryStats(countryStat || null);

      // Find job title stats
      const jobStat = jobTitleData.find((stat) => stat.jobTitle === jobTitle);
      setJobTitleStats(jobStat || null);
    } catch (err) {
      console.error('Error fetching salary insights:', err);
    }
  };

  const handleDelete = async () => {
    if (!employee) return;

    if (!window.confirm('Are you sure you want to delete this employee?')) {
      return;
    }

    try {
      await employeeAPI.deleteEmployee(employee.id);
      navigate('/');
    } catch (err) {
      setError('Failed to delete employee');
      console.error('Error deleting employee:', err);
    }
  };

  const formatSalary = (salary: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(salary);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="flex items-center gap-2">
          <Loader2 className="h-6 w-6 animate-spin" />
          <span>Loading employee details...</span>
        </div>
      </div>
    );
  }

  if (error || !employee) {
    return (
      <div className="min-h-screen bg-background">
        <div className="container mx-auto py-6">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-destructive mb-4">{error || 'Employee not found'}</p>
              <Button onClick={() => navigate('/')}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Dashboard
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto py-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
          <div className="flex items-center gap-4">
            <Link
              to="/"
              className="flex items-center text-sm text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="h-4 w-4 mr-1" />
              Back to Employees
            </Link>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => navigate(`/edit-employee/${employee.id}`)}
            >
              <Edit className="mr-2 h-4 w-4" />
              Edit
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </Button>
          </div>
        </div>

        {/* Employee Profile Card */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex items-start gap-4">
              <div className="h-20 w-20 rounded-full bg-primary/10 flex items-center justify-center">
                <User className="h-10 w-10 text-primary" />
              </div>
              <div className="flex-1">
                <CardTitle className="text-2xl">
                  {employee.firstName} {employee.lastName}
                </CardTitle>
                <CardDescription className="text-base mt-1">
                  {employee.jobTitle}
                </CardDescription>
                <div className="flex flex-wrap gap-2 mt-3">
                  <Badge variant="secondary">{employee.department}</Badge>
                  <Badge variant="outline">{employee.country}</Badge>
                </div>
              </div>
            </div>
          </CardHeader>
        </Card>

        <div className="grid gap-6 md:grid-cols-2">
          {/* Personal Information */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Personal Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Email</p>
                  <p className="font-medium">{employee.email}</p>
                </div>
              </div>
              <div className="border-t" />
              <div className="flex items-center gap-3">
                <Building className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Department</p>
                  <p className="font-medium">{employee.department}</p>
                </div>
              </div>
              <div className="border-t" />
              <div className="flex items-center gap-3">
                <Globe className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Country</p>
                  <p className="font-medium">{employee.country}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Employment Details */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Employment Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <Briefcase className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Job Title</p>
                  <p className="font-medium">{employee.jobTitle}</p>
                </div>
              </div>
              <div className="border-t" />
              <div className="flex items-center gap-3">
                <DollarSign className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Salary</p>
                  <p className="font-medium text-green-600">
                    {formatSalary(employee.salary)}
                  </p>
                </div>
              </div>
              <div className="border-t" />
              <div className="flex items-center gap-3">
                <Calendar className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Hire Date</p>
                  <p className="font-medium">{formatDate(employee.hireDate)}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Salary Insights */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Salary Insights
            </CardTitle>
            <CardDescription>
              Salary comparison for {employee.country} and {employee.jobTitle} role
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-6 md:grid-cols-2">
              {/* Country Insights */}
              {countryStats && (
                <div className="space-y-4">
                  <h4 className="font-medium text-muted-foreground">
                    {employee.country} Salary Statistics
                  </h4>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Minimum</p>
                      <p className="font-semibold">{formatSalary(countryStats.min)}</p>
                    </div>
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Maximum</p>
                      <p className="font-semibold">{formatSalary(countryStats.max)}</p>
                    </div>
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Average</p>
                      <p className="font-semibold">{formatSalary(countryStats.average)}</p>
                    </div>
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Employees</p>
                      <p className="font-semibold">{countryStats.count}</p>
                    </div>
                  </div>
                  <div className="p-3 bg-primary/5 rounded-lg">
                    <p className="text-sm text-muted-foreground">
                      {employee.firstName}&apos;s salary vs country average
                    </p>
                    <p className={`font-semibold ${
                      employee.salary >= countryStats.average ? 'text-green-600' : 'text-amber-600'
                    }`}>
                      {employee.salary >= countryStats.average ? '+' : ''}
                      {formatSalary(employee.salary - countryStats.average)}
                      {' '}
                      ({((employee.salary / countryStats.average - 1) * 100).toFixed(1)}%)
                    </p>
                  </div>
                </div>
              )}

              {/* Job Title Insights */}
              {jobTitleStats && (
                <div className="space-y-4">
                  <h4 className="font-medium text-muted-foreground">
                    {employee.jobTitle} in {employee.country}
                  </h4>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Average Salary</p>
                      <p className="font-semibold">{formatSalary(jobTitleStats.average)}</p>
                    </div>
                    <div className="p-3 bg-muted rounded-lg">
                      <p className="text-sm text-muted-foreground">Total Positions</p>
                      <p className="font-semibold">{jobTitleStats.count}</p>
                    </div>
                  </div>
                  <div className="p-3 bg-primary/5 rounded-lg">
                    <p className="text-sm text-muted-foreground">
                      {employee.firstName}&apos;s salary vs role average
                    </p>
                    <p className={`font-semibold ${
                      employee.salary >= jobTitleStats.average ? 'text-green-600' : 'text-amber-600'
                    }`}>
                      {employee.salary >= jobTitleStats.average ? '+' : ''}
                      {formatSalary(employee.salary - jobTitleStats.average)}
                      {' '}
                      ({((employee.salary / jobTitleStats.average - 1) * 100).toFixed(1)}%)
                    </p>
                  </div>
                </div>
              )}

              {!countryStats && !jobTitleStats && (
                <div className="col-span-2 text-center py-8 text-muted-foreground">
                  <Clock className="h-8 w-8 mx-auto mb-2" />
                  <p>No salary insights available for this employee&apos;s country or job title.</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* System Information */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle className="text-lg">System Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3 text-sm">
              <div>
                <p className="text-muted-foreground">Employee ID</p>
                <p className="font-mono">#{employee.id}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Created At</p>
                <p>{formatDate(employee.createdAt)}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Last Updated</p>
                <p>{formatDate(employee.updatedAt)}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default EmployeeView;
