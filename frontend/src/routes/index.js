import { Routes, Route } from 'react-router-dom';
import AuthRedirect from './AuthRedirect';
import AdminRoute from './AdminRoute';

import AdminLogin from '../pages/admin/Login';

const Home = () => <div>Home</div>;
const Dashboard = () => <div>Dashboard</div>;

function AppRoutes() {
  return (
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
        <Route path="/admin/dashboard" element={<Dashboard />}/>
      </Route>
    </Routes>
  );
}

export default AppRoutes;
