import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { 
  BuildingOfficeIcon, 
  ChartBarIcon, 
  UserGroupIcon,
  CurrencyDollarIcon 
} from '@heroicons/react/24/outline';
import { healthAPI } from './services/api';
import { HealthResponse } from './types';
import EmployeeList from './components/EmployeeList';
import EmployeeForm from './components/EmployeeForm';
import Analytics from './components/Analytics';
import './App.css';

function App() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const healthData = await healthAPI.getHealth();
        setHealth(healthData);
      } catch (error) {
        console.error('Health check failed:', error);
        setHealth({ status: 'error' });
      } finally {
        setLoading(false);
      }
    };

    checkHealth();
  }, []);

  const Navigation = () => (
    <nav className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              <BuildingOfficeIcon className="h-8 w-8 text-blue-600" />
              <span className="ml-2 text-xl font-bold text-gray-900">Salary Manager</span>
            </div>
            <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
              <Link
                to="/"
                className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                <UserGroupIcon className="h-5 w-5 mr-1" />
                Employees
              </Link>
              <Link
                to="/analytics"
                className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                <ChartBarIcon className="h-5 w-5 mr-1" />
                Analytics
              </Link>
              <Link
                to="/add-employee"
                className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                <CurrencyDollarIcon className="h-5 w-5 mr-1" />
                Add Employee
              </Link>
            </div>
          </div>
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                health?.status === 'ok' 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              }`}>
                {loading ? 'Checking...' : health?.status === 'ok' ? 'API Connected' : 'API Error'}
              </span>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Navigation />
        <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <Routes>
            <Route path="/" element={<EmployeeList />} />
            <Route path="/analytics" element={<Analytics />} />
            <Route path="/add-employee" element={<EmployeeForm />} />
            <Route path="/edit-employee/:id" element={<EmployeeForm />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
