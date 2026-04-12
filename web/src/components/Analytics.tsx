import React, { useState, useEffect } from 'react';
import { 
  DollarSign, 
  Globe, 
  Briefcase,
  Building2,
  BarChart3,
  X
} from 'lucide-react';
import { analyticsAPI } from '../services/api';
import { CountrySalaryStats, JobTitleSalaryStats, DepartmentSalaryStats } from '../types';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

const Analytics: React.FC = () => {
  const [countryStats, setCountryStats] = useState<CountrySalaryStats[]>([]);
  const [jobTitleStats, setJobTitleStats] = useState<JobTitleSalaryStats[]>([]);
  const [departmentStats, setDepartmentStats] = useState<DepartmentSalaryStats[]>([]);
  const [selectedCountry, setSelectedCountry] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAnalyticsData();
  }, []);

  const fetchAnalyticsData = async () => {
    try {
      setLoading(true);
      const [countryData, departmentData] = await Promise.all([
        analyticsAPI.getSalaryByCountry(),
        analyticsAPI.getDepartmentInsights()
      ]);
      
      setCountryStats(countryData);
      setDepartmentStats(departmentData);
      setError(null);
    } catch (err) {
      setError('Failed to fetch analytics data');
      console.error('Error fetching analytics:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchJobTitleStats = async (country: string) => {
    try {
      const data = await analyticsAPI.getSalaryByJobTitleInCountry(country);
      setJobTitleStats(data);
    } catch (err) {
      console.error('Error fetching job title stats:', err);
    }
  };

  const handleCountryClick = (country: string) => {
    setSelectedCountry(country);
    fetchJobTitleStats(country);
  };

  const formatSalary = (salary: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(salary);
  };

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat('en-US').format(num);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div>
        <CardHeader className="px-0">
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-6 w-6" />
            Salary Analytics
          </CardTitle>
          <CardDescription>
            Comprehensive salary insights across countries, job titles, and departments.
          </CardDescription>
        </CardHeader>
      </div>

      {error && (
        <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Salary by Country */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Globe className="h-5 w-5 text-blue-600" />
              Salary by Country
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {countryStats.slice(0, 10).map((stat) => (
                <div key={stat.country} className="flex items-center space-x-4 p-3 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer" onClick={() => handleCountryClick(stat.country)}>
                  <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                    <DollarSign className="h-4 w-4 text-blue-600" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium hover:text-primary truncate">
                      {stat.country}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {formatNumber(stat.count)} employees
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-medium">
                      {formatSalary(stat.average)}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {formatSalary(stat.min)} - {formatSalary(stat.max)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Department Insights */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Building2 className="h-5 w-5 text-green-600" />
              Department Salary Insights
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {departmentStats.slice(0, 10).map((stat) => (
                <div key={stat.department} className="flex items-center space-x-4 p-3 rounded-lg">
                  <div className="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center">
                    <Briefcase className="h-4 w-4 text-green-600" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium truncate">
                      {stat.department}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {formatNumber(stat.count)} employees
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-medium">
                      {formatSalary(stat.average)}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {formatSalary(stat.min)} - {formatSalary(stat.max)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Job Title Stats by Country */}
      {selectedCountry && jobTitleStats.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Briefcase className="h-5 w-5 text-purple-600" />
                Job Title Salaries in {selectedCountry}
              </CardTitle>
              <Button variant="ghost" size="sm" onClick={() => setSelectedCountry('')}>
                <X className="h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {jobTitleStats.map((stat) => (
                <div key={stat.jobTitle} className="flex items-center space-x-4 p-3 rounded-lg">
                  <div className="h-8 w-8 rounded-full bg-purple-100 flex items-center justify-center">
                    <DollarSign className="h-4 w-4 text-purple-600" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium truncate">
                      {stat.jobTitle}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {formatNumber(stat.count)} employees
                    </div>
                  </div>
                  <div className="text-sm font-medium">
                    {formatSalary(stat.average)}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Summary Statistics */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <Globe className="h-6 w-6 text-blue-600" />
              <div className="ml-4">
                <p className="text-sm font-medium text-muted-foreground">
                  Countries
                </p>
                <p className="text-2xl font-bold">
                  {countryStats.length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <Building2 className="h-6 w-6 text-green-600" />
              <div className="ml-4">
                <p className="text-sm font-medium text-muted-foreground">
                  Departments
                </p>
                <p className="text-2xl font-bold">
                  {departmentStats.length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <DollarSign className="h-6 w-6 text-yellow-600" />
              <div className="ml-4">
                <p className="text-sm font-medium text-muted-foreground">
                  Highest Avg Salary
                </p>
                <p className="text-2xl font-bold">
                  {countryStats.length > 0 ? formatSalary(countryStats[0].average) : 'N/A'}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <Briefcase className="h-6 w-6 text-purple-600" />
              <div className="ml-4">
                <p className="text-sm font-medium text-muted-foreground">
                  Total Employees
                </p>
                <p className="text-2xl font-bold">
                  {formatNumber(countryStats.reduce((sum, stat) => sum + stat.count, 0))}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Analytics;
