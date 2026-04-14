import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import {
  Building2,
  Users,
  TrendingUp,
  DollarSign,
  LogOut,
  Menu,
  MapPin,
  Briefcase,
  BarChart3,
  ArrowLeft,
  Loader2,
  Globe,
  PieChart,
  ChevronDown,
} from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Avatar, AvatarFallback } from './ui/avatar';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { Sheet, SheetContent, SheetTrigger } from './ui/sheet';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { analyticsAPI, authAPI } from '../services/api';
import { DepartmentSalaryStats, CountrySalaryStats, JobTitleSalaryStats, User } from '../types';

interface GlobalAnalytics {
  totalEmployees: number;
  avgSalary: number;
  totalSalaryExpense: number;
  departments: DepartmentSalaryStats[];
  countryStats: CountrySalaryStats[];
}

const Analytics: React.FC = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [analytics, setAnalytics] = useState<GlobalAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  // Salary insights state
  const [selectedCountry, setSelectedCountry] = useState<string>('');
  const [jobTitleStats, setJobTitleStats] = useState<JobTitleSalaryStats[]>([]);
  const [insightsLoading, setInsightsLoading] = useState(false);

  // Department by country state
  const [deptCountry, setDeptCountry] = useState<string>('');
  const [deptStatsByCountry, setDeptStatsByCountry] = useState<DepartmentSalaryStats[]>([]);
  const [deptLoading, setDeptLoading] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      try {
        const userData = localStorage.getItem('user');
        if (!userData) {
          navigate('/login');
          return;
        }
        setUser(JSON.parse(userData));

        const [salaryData, departmentData] = await Promise.all([
          analyticsAPI.getSalaryByCountry(),
          analyticsAPI.getDepartmentInsights()
        ]);

        // Calculate global stats
        const totalEmployees = salaryData.reduce((sum: number, c: CountrySalaryStats) => sum + c.count, 0);
        const avgSalary = salaryData.length > 0 ?
          salaryData.reduce((sum: number, c: CountrySalaryStats) => sum + c.average, 0) / salaryData.length : 0;
        const totalSalaryExpense = salaryData.reduce((sum: number, c: CountrySalaryStats) => {
          return sum + (c.average * c.count);
        }, 0);

        setAnalytics({
          totalEmployees,
          avgSalary,
          totalSalaryExpense,
          departments: departmentData || [],
          countryStats: salaryData || []
        });

        // Set first country as selected if available
        if (salaryData && salaryData.length > 0) {
          setSelectedCountry(salaryData[0].country);
          setDeptCountry(salaryData[0].country);
        }

      } catch (error) {
        console.error('Failed to load analytics:', error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [navigate]);

  // Fetch job title stats when country changes
  useEffect(() => {
    const fetchJobTitleStats = async () => {
      if (!selectedCountry) return;
      try {
        setInsightsLoading(true);
        const stats = await analyticsAPI.getSalaryByJobTitleInCountry(selectedCountry);
        setJobTitleStats(stats);
      } catch (error) {
        console.error('Failed to fetch job title stats:', error);
      } finally {
        setInsightsLoading(false);
      }
    };

    fetchJobTitleStats();
  }, [selectedCountry]);

  // Fetch department stats when country changes
  useEffect(() => {
    const fetchDeptStatsByCountry = async () => {
      if (!deptCountry) return;
      try {
        setDeptLoading(true);
        const stats = await analyticsAPI.getDepartmentInsightsByCountry(deptCountry);
        setDeptStatsByCountry(stats);
      } catch (error) {
        console.error('Failed to fetch department stats by country:', error);
      } finally {
        setDeptLoading(false);
      }
    };

    fetchDeptStatsByCountry();
  }, [deptCountry]);

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

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 0,
    }).format(value);
  };

  const formatNumber = (value: number) => {
    return new Intl.NumberFormat('en-US').format(value);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="flex items-center gap-3">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <span className="text-lg">Loading analytics...</span>
        </div>
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
            <span className="text-muted-foreground">|</span>
            <span className="text-muted-foreground flex items-center gap-2">
              <BarChart3 className="h-4 w-4" />
              Analytics
            </span>
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
                  <Link to="/">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Dashboard
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
              <Link to="/">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Dashboard
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
        {/* Global Stats Cards */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Employees</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(analytics?.totalEmployees || 0)}</div>
              <p className="text-xs text-muted-foreground">Across all countries</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Global Average Salary</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatCurrency(analytics?.avgSalary || 0)}</div>
              <p className="text-xs text-muted-foreground">Per employee</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Payroll</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatCurrency(analytics?.totalSalaryExpense || 0)}</div>
              <p className="text-xs text-muted-foreground">Monthly expense</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Countries</CardTitle>
              <Globe className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{analytics?.countryStats?.length || 0}</div>
              <p className="text-xs text-muted-foreground">Active locations</p>
            </CardContent>
          </Card>
        </div>

        {/* Analytics Tabs */}
        <Tabs defaultValue="countries" className="space-y-6">
          <TabsList className="grid w-full md:w-auto grid-cols-3 md:inline-flex">
            <TabsTrigger value="countries" className="flex items-center gap-2">
              <Globe className="h-4 w-4" />
              <span className="hidden sm:inline">By Country</span>
            </TabsTrigger>
            <TabsTrigger value="jobtitles" className="flex items-center gap-2">
              <Briefcase className="h-4 w-4" />
              <span className="hidden sm:inline">By Job Title</span>
            </TabsTrigger>
            <TabsTrigger value="departments" className="flex items-center gap-2">
              <PieChart className="h-4 w-4" />
              <span className="hidden sm:inline">By Department</span>
            </TabsTrigger>
          </TabsList>

          {/* Countries Tab */}
          <TabsContent value="countries" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Globe className="h-5 w-5" />
                  Salary by Country
                </CardTitle>
                <CardDescription>
                  Minimum, maximum, and average salaries across all countries
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border overflow-hidden">
                  <table className="w-full text-sm">
                    <thead className="bg-muted">
                      <tr>
                        <th className="px-4 py-3 text-left font-medium">Country</th>
                        <th className="px-4 py-3 text-right font-medium">Employees</th>
                        <th className="px-4 py-3 text-right font-medium">Min Salary</th>
                        <th className="px-4 py-3 text-right font-medium">Max Salary</th>
                        <th className="px-4 py-3 text-right font-medium">Average</th>
                        <th className="px-4 py-3 text-right font-medium">Total Payroll</th>
                      </tr>
                    </thead>
                    <tbody>
                      {analytics?.countryStats?.map((stat) => (
                        <tr key={stat.country} className="border-t">
                          <td className="px-4 py-3 font-medium">{stat.country}</td>
                          <td className="px-4 py-3 text-right">{stat.count}</td>
                          <td className="px-4 py-3 text-right">{formatCurrency(stat.min)}</td>
                          <td className="px-4 py-3 text-right">{formatCurrency(stat.max)}</td>
                          <td className="px-4 py-3 text-right font-semibold text-primary">
                            {formatCurrency(stat.average)}
                          </td>
                          <td className="px-4 py-3 text-right text-muted-foreground">
                            {formatCurrency(stat.average * stat.count)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Job Titles Tab */}
          <TabsContent value="jobtitles" className="space-y-6">
            <Card>
              <CardHeader>
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      <Briefcase className="h-5 w-5" />
                      Salary by Job Title
                    </CardTitle>
                    <CardDescription>
                      Average salaries for each job title within a selected country
                    </CardDescription>
                  </div>
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" className="w-[200px] justify-between bg-slate-900 border-slate-700">
                          {selectedCountry || 'Select country'}
                          <ChevronDown className="h-4 w-4 ml-2 opacity-50" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent className="w-[200px] bg-slate-900 border-slate-700">
                        {analytics?.countryStats?.map((stat) => (
                          <DropdownMenuItem
                            key={stat.country}
                            onClick={() => setSelectedCountry(stat.country)}
                            className={`cursor-pointer ${selectedCountry === stat.country ? 'bg-slate-800' : ''}`}
                          >
                            {stat.country}
                          </DropdownMenuItem>
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {insightsLoading ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                  </div>
                ) : jobTitleStats.length > 0 ? (
                  <div className="rounded-md border overflow-hidden">
                    <table className="w-full text-sm">
                      <thead className="bg-muted">
                        <tr>
                          <th className="px-4 py-3 text-left font-medium">Job Title</th>
                          <th className="px-4 py-3 text-right font-medium">Positions</th>
                          <th className="px-4 py-3 text-right font-medium">Average Salary</th>
                          <th className="px-4 py-3 text-right font-medium">Total Cost</th>
                        </tr>
                      </thead>
                      <tbody>
                        {jobTitleStats.map((stat) => (
                          <tr key={stat.jobTitle} className="border-t">
                            <td className="px-4 py-3 font-medium">{stat.jobTitle}</td>
                            <td className="px-4 py-3 text-right">{stat.count}</td>
                            <td className="px-4 py-3 text-right font-semibold text-primary">
                              {formatCurrency(stat.average)}
                            </td>
                            <td className="px-4 py-3 text-right text-muted-foreground">
                              {formatCurrency(stat.average * stat.count)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : selectedCountry ? (
                  <div className="text-center py-12 text-muted-foreground">
                    <Briefcase className="h-12 w-12 mx-auto mb-3" />
                    <p className="text-lg">No job title data available for {selectedCountry}</p>
                  </div>
                ) : (
                  <div className="text-center py-12 text-muted-foreground">
                    <MapPin className="h-12 w-12 mx-auto mb-3" />
                    <p className="text-lg">Select a country to view job title insights</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Departments Tab */}
          <TabsContent value="departments" className="space-y-6">
            <Card>
              <CardHeader>
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      <PieChart className="h-5 w-5" />
                      Salary by Department
                    </CardTitle>
                    <CardDescription>
                      Salary statistics filtered by country
                    </CardDescription>
                  </div>
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" className="w-[160px] justify-between bg-slate-900 border-slate-700">
                          {deptCountry || 'Select country'}
                          <ChevronDown className="h-4 w-4 ml-2 opacity-50" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent className="w-[160px] bg-slate-900 border-slate-700">
                        {analytics?.countryStats?.map((stat) => (
                          <DropdownMenuItem
                            key={stat.country}
                            onClick={() => setDeptCountry(stat.country)}
                            className={`cursor-pointer ${deptCountry === stat.country ? 'bg-slate-800' : ''}`}
                          >
                            {stat.country}
                          </DropdownMenuItem>
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {deptLoading ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                  </div>
                ) : deptStatsByCountry.length > 0 ? (
                  <>
                    <div className="rounded-md border overflow-hidden">
                      <table className="w-full text-sm">
                        <thead className="bg-muted">
                          <tr>
                            <th className="px-4 py-3 text-left font-medium">Department</th>
                            <th className="px-4 py-3 text-right font-medium">Employees</th>
                            <th className="px-4 py-3 text-right font-medium">Min Salary</th>
                            <th className="px-4 py-3 text-right font-medium">Max Salary</th>
                            <th className="px-4 py-3 text-right font-medium">Average</th>
                            <th className="px-4 py-3 text-right font-medium">Total Payroll</th>
                          </tr>
                        </thead>
                        <tbody>
                          {deptStatsByCountry.map((dept) => (
                            <tr key={dept.department} className="border-t">
                              <td className="px-4 py-3 font-medium">{dept.department}</td>
                              <td className="px-4 py-3 text-right">{dept.count}</td>
                              <td className="px-4 py-3 text-right">{formatCurrency(dept.min)}</td>
                              <td className="px-4 py-3 text-right">{formatCurrency(dept.max)}</td>
                              <td className="px-4 py-3 text-right font-semibold text-primary">
                                {formatCurrency(dept.average)}
                              </td>
                              <td className="px-4 py-3 text-right text-muted-foreground">
                                {formatCurrency(dept.average * dept.count)}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </>
                ) : deptCountry ? (
                  <div className="text-center py-12 text-muted-foreground">
                    <PieChart className="h-12 w-12 mx-auto mb-3" />
                    <p className="text-lg">No department data available for {deptCountry}</p>
                  </div>
                ) : (
                  <div className="text-center py-12 text-muted-foreground">
                    <MapPin className="h-12 w-12 mx-auto mb-3" />
                    <p className="text-lg">Select a country to view department insights</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
};

export default Analytics;
