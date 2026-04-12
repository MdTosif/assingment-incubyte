import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { 
  Building2, 
  BarChart3, 
  Users, 
  UserPlus,
  Menu,
  X,
  LogOut,
  Settings
} from 'lucide-react';
import { Button } from './ui/button';
import { cn } from '@/lib/utils';
import { healthAPI } from '../services/api';
import { HealthResponse } from '../types';

interface NavigationProps {
  user?: {
    firstName: string;
    lastName: string;
    email: string;
    role: string;
  };
  onLogout?: () => void;
}

const Navigation: React.FC<NavigationProps> = ({ user, onLogout }) => {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const location = useLocation();

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

  const navigation = [
    { name: 'Employees', href: '/', icon: Users },
    { name: 'Analytics', href: '/analytics', icon: BarChart3 },
    { name: 'Add Employee', href: '/add-employee', icon: UserPlus },
  ];

  const isActive = (href: string) => {
    if (href === '/') {
      return location.pathname === '/' || location.pathname.startsWith('/edit-employee');
    }
    return location.pathname === href;
  };

  return (
    <nav className="bg-background border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              <Link to="/" className="flex items-center space-x-3">
                <Building2 className="h-8 w-8 text-primary" />
                <span className="text-xl font-bold text-foreground">Salary Manager</span>
              </Link>
            </div>
            
            {/* Desktop Navigation */}
            <div className="hidden sm:ml-6 sm:flex sm:space-x-1">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <Link
                    key={item.name}
                    to={item.href}
                    className={cn(
                      "inline-flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors",
                      isActive(item.href)
                        ? "bg-primary text-primary-foreground"
                        : "text-muted-foreground hover:text-foreground hover:bg-accent"
                    )}
                  >
                    <Icon className="h-4 w-4 mr-2" />
                    {item.name}
                  </Link>
                );
              })}
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* Health Status */}
            <div className="hidden sm:flex items-center">
              <div className={cn(
                "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium",
                health?.status === 'ok' 
                  ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200" 
                  : "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
              )}>
                <div className={cn(
                  "w-2 h-2 rounded-full mr-2",
                  health?.status === 'ok' ? "bg-green-500" : "bg-red-500"
                )} />
                {loading ? 'Checking...' : health?.status === 'ok' ? 'API Connected' : 'API Error'}
              </div>
            </div>

            {/* User Menu */}
            {user && (
              <div className="hidden sm:flex items-center space-x-3">
                <div className="text-sm text-muted-foreground">
                  <span className="font-medium text-foreground">
                    {user.firstName} {user.lastName}
                  </span>
                  <span className="ml-2 text-xs bg-secondary px-2 py-1 rounded">
                    {user.role}
                  </span>
                </div>
                <Button variant="ghost" size="sm" onClick={onLogout}>
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            )}

            {/* Mobile menu button */}
            <div className="sm:hidden">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              >
                {mobileMenuOpen ? (
                  <X className="h-6 w-6" />
                ) : (
                  <Menu className="h-6 w-6" />
                )}
              </Button>
            </div>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div className="sm:hidden border-t">
            <div className="px-2 pt-2 pb-3 space-y-1">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <Link
                    key={item.name}
                    to={item.href}
                    className={cn(
                      "flex items-center px-3 py-2 rounded-md text-base font-medium transition-colors",
                      isActive(item.href)
                        ? "bg-primary text-primary-foreground"
                        : "text-muted-foreground hover:text-foreground hover:bg-accent"
                    )}
                    onClick={() => setMobileMenuOpen(false)}
                  >
                    <Icon className="h-5 w-5 mr-3" />
                    {item.name}
                  </Link>
                );
              })}
            </div>
            
            {user && (
              <div className="pt-4 pb-3 border-t">
                <div className="px-2 space-y-1">
                  <div className="flex items-center px-3 py-2">
                    <div className="text-sm text-muted-foreground">
                      <div className="font-medium text-foreground">
                        {user.firstName} {user.lastName}
                      </div>
                      <div className="text-xs">{user.email}</div>
                      <span className="inline-block mt-1 text-xs bg-secondary px-2 py-1 rounded">
                        {user.role}
                      </span>
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="w-full justify-start"
                    onClick={() => {
                      onLogout?.();
                      setMobileMenuOpen(false);
                    }}
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Logout
                  </Button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navigation;
