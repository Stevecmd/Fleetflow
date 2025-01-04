import React from 'react';
import { Outlet, Link } from 'react-router-dom';
import { useAuth } from '../features/auth/AuthProvider';
import Logo from '../assets/logo.svg';

/**
 * The main layout component for the app.
 *
 * This component renders the main layout structure for the app, including the
 * navigation bar and the main content area.
 *
 * The navigation bar contains a link to the login page if the user is not
 * authenticated, and a logout button if the user is authenticated.
 *
 * The main content area renders the child route component, which is passed in
 * as a prop from the parent route.
 */
const Layout: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0 flex items-center">
                <img 
                  src={Logo} 
                  alt="FleetFlow Logo" 
                  className="h-10 w-auto sm:h-12 md:h-14 lg:h-16 object-contain 
                    rounded-lg 
                    shadow-sm 
                    hover:shadow-md 
                    transition-all 
                    duration-300 
                    ease-in-out 
                    transform 
                    hover:scale-105 
                    bg-white 
                    p-1"
                />
              </div>
            </div>
            <div className="flex items-center">
              {isAuthenticated ? (
                <button
                  onClick={logout}
                  className="ml-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Logout
                </button>
              ) : (
                <Link
                  to="/login"
                  className="ml-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Login
                </Link>
              )}
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;
