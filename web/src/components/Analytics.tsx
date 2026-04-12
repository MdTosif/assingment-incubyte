import React, { useState, useEffect } from 'react';
import { 
  CurrencyDollarIcon, 
  GlobeAltIcon, 
  BriefcaseIcon,
  BuildingOfficeIcon
} from '@heroicons/react/24/outline';
import { analyticsAPI } from '../services/api';
import { CountrySalaryStats, JobTitleSalaryStats, DepartmentSalaryStats } from '../types';

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
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-semibold text-gray-900">Salary Analytics</h1>
          <p className="mt-2 text-sm text-gray-700">
            Comprehensive salary insights across countries, job titles, and departments.
          </p>
        </div>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4 mb-6">
          <div className="text-sm text-red-700">{error}</div>
        </div>
      )}

      <div className="mt-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Salary by Country */}
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <GlobeAltIcon className="h-8 w-8 text-blue-600" />
              <h3 className="ml-3 text-lg leading-6 font-medium text-gray-900">
                Salary by Country
              </h3>
            </div>
            <div className="mt-4">
              <div className="flow-root">
                <ul className="-my-5 divide-y divide-gray-200">
                  {countryStats.slice(0, 10).map((stat) => (
                    <li key={stat.country} className="py-4">
                      <div className="flex items-center space-x-4">
                        <div className="flex-shrink-0">
                          <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                            <CurrencyDollarIcon className="h-4 w-4 text-blue-600" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <button
                            onClick={() => handleCountryClick(stat.country)}
                            className="text-sm font-medium text-gray-900 hover:text-blue-600 truncate"
                          >
                            {stat.country}
                          </button>
                          <div className="text-sm text-gray-500">
                            {formatNumber(stat.count)} employees
                          </div>
                        </div>
                        <div className="flex-shrink-0 text-right">
                          <div className="text-sm font-medium text-gray-900">
                            {formatSalary(stat.average)}
                          </div>
                          <div className="text-xs text-gray-500">
                            {formatSalary(stat.min)} - {formatSalary(stat.max)}
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>

        {/* Department Insights */}
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <BuildingOfficeIcon className="h-8 w-8 text-green-600" />
              <h3 className="ml-3 text-lg leading-6 font-medium text-gray-900">
                Department Salary Insights
              </h3>
            </div>
            <div className="mt-4">
              <div className="flow-root">
                <ul className="-my-5 divide-y divide-gray-200">
                  {departmentStats.slice(0, 10).map((stat) => (
                    <li key={stat.department} className="py-4">
                      <div className="flex items-center space-x-4">
                        <div className="flex-shrink-0">
                          <div className="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center">
                            <BriefcaseIcon className="h-4 w-4 text-green-600" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-gray-900 truncate">
                            {stat.department}
                          </p>
                          <div className="text-sm text-gray-500">
                            {formatNumber(stat.count)} employees
                          </div>
                        </div>
                        <div className="flex-shrink-0 text-right">
                          <div className="text-sm font-medium text-gray-900">
                            {formatSalary(stat.average)}
                          </div>
                          <div className="text-xs text-gray-500">
                            {formatSalary(stat.min)} - {formatSalary(stat.max)}
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Job Title Stats by Country */}
      {selectedCountry && jobTitleStats.length > 0 && (
        <div className="mt-8 bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center">
                <BriefcaseIcon className="h-8 w-8 text-purple-600" />
                <h3 className="ml-3 text-lg leading-6 font-medium text-gray-900">
                  Job Title Salaries in {selectedCountry}
                </h3>
              </div>
              <button
                onClick={() => setSelectedCountry('')}
                className="text-sm text-gray-500 hover:text-gray-700"
              >
                Clear
              </button>
            </div>
            <div className="mt-4">
              <div className="flow-root">
                <ul className="-my-5 divide-y divide-gray-200">
                  {jobTitleStats.map((stat) => (
                    <li key={stat.jobTitle} className="py-4">
                      <div className="flex items-center space-x-4">
                        <div className="flex-shrink-0">
                          <div className="h-8 w-8 rounded-full bg-purple-100 flex items-center justify-center">
                            <CurrencyDollarIcon className="h-4 w-4 text-purple-600" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-gray-900 truncate">
                            {stat.jobTitle}
                          </p>
                          <div className="text-sm text-gray-500">
                            {formatNumber(stat.count)} employees
                          </div>
                        </div>
                        <div className="flex-shrink-0">
                          <div className="text-sm font-medium text-gray-900">
                            {formatSalary(stat.average)}
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Summary Statistics */}
      <div className="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <GlobeAltIcon className="h-6 w-6 text-blue-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Countries
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {countryStats.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <BuildingOfficeIcon className="h-6 w-6 text-green-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Departments
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {departmentStats.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <CurrencyDollarIcon className="h-6 w-6 text-yellow-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Highest Avg Salary
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {countryStats.length > 0 ? formatSalary(countryStats[0].average) : 'N/A'}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <BriefcaseIcon className="h-6 w-6 text-purple-600" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Employees
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {formatNumber(countryStats.reduce((sum, stat) => sum + stat.count, 0))}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Analytics;
