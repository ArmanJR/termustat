import { Routes, Route } from 'react-router-dom';

import AdminLogin from '../pages/admin/Login';

const Home = () => <div>Home</div>;
const Dashboard = () => <div>Dashboard</div>;

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/admin/login" element={<AdminLogin />} />
      <Route path="/admin/dashboard" element={<Dashboard />} />
    </Routes>
  );
}

export default AppRoutes;
