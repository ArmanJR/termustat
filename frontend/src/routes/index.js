import { Routes, Route } from 'react-router-dom';
import AuthListener from '../services/AuthListener';
import AuthRedirect from './AuthRedirect';
import AdminRoute from './AdminRoute';

import AdminLogin from '../pages/admin/Login';
import AdminLayout from '../layouts/AdminLayout';
import Dashboard from '../pages/admin/Dashboard';
import Universities from '../pages/admin/Universities';

const Home = () => <div>Home</div>;

function AppRoutes() {
  return (
    <>
      <AuthListener />
      <Routes>
        {/* Public Routes */}
        <Route path="/" element={<Home />} />
        <Route
          path="/admin/login"
          element={
            /* Redirect to /admin/dashboard if logged in */
            <AuthRedirect redirectTo="/admin/dashboard">
              <AdminLogin />
            </AuthRedirect>
          }
        />

        {/* Admin Protected Routes */}
        <Route element={<AdminRoute />}>
          <Route path="/admin" element={<AdminLayout />}>
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="universities" element={<Universities />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default AppRoutes;
