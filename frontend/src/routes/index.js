import { Routes, Route } from 'react-router-dom';
import AuthRedirect from './AuthRedirect';

import AdminLogin from '../pages/admin/Login';

const Home = () => <div>Home</div>;
const Dashboard = () => <div>Dashboard</div>;

function AppRoutes() {
  return (
    <Routes>
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
      <Route path="/admin/dashboard" element={<Dashboard />} />
    </Routes>
  );
}

export default AppRoutes;
