/**
 * App Component
 * 
 * Root application component that handles:
 * - Authentication state management
 * - Route protection based on auth status
 * - Token validation on app load
 */

import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { User } from './types';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import EmployeeForm from './components/EmployeeForm';
import Analytics from './components/Analytics';
import { EmployeeView } from './components/employees';
import './App.css';

function App() {
  // ==================== State ====================
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // ==================== Authentication ====================
  const checkAuth = () => {
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (token && userData) {
      try {
        const parsedUser = JSON.parse(userData);
        setUser(parsedUser);
        return parsedUser;
      } catch (error) {
        console.error('Failed to parse user data:', error);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setUser(null);
        return null;
      }
    } else {
      setUser(null);
      return null;
    }
  };

  useEffect(() => {
    checkAuth();
    setLoading(false);
  }, []);

  // ==================== Effects ====================

  // Listen for storage changes (login from other tabs)
  useEffect(() => {
    const handleStorageChange = () => {
      checkAuth();
    };
    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  /** Logout handler - clears local storage and resets user state */
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  };

  /** Login handler - stores token/user and updates state */
  const handleLogin = (token: string, userData: User) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
  };

  // ==================== Render ====================

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <Router>
      <div className="min-h-screen bg-background">
        {user ? (
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/analytics" element={<Analytics />} />
            <Route path="/add-employee" element={<EmployeeForm />} />
            <Route path="/edit-employee/:id" element={<EmployeeForm />} />
            <Route path="/employee/:id" element={<EmployeeView />} />
            <Route path="/login" element={<Navigate to="/" replace />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        ) : (
          <Routes>
            <Route path="/login" element={<Login onLogin={handleLogin} />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        )}
      </div>
    </Router>
  );
}

export default App;
