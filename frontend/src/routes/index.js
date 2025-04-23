import { Routes, Route } from 'react-router-dom';

import AdminLogin from '../pages/admin/Login';

const Home = () => <div>Home</div>;

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/admin/login" element={<AdminLogin />} />
    </Routes>
  );
}

export default AppRoutes;
