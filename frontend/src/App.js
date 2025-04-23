import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';

// Page components
import AdminLogin from './pages/admin/Login.js';
const Home = () => <div>Home</div>;

function App() {
  // get the current route
  const location = useLocation();
  // only show nav on HomePage
  const showNav = location.pathname === '/';

  return (
    <div>
      {showNav && (
        <nav>
          <ul>
            <li><Link to="/">Home</Link></li>
            <li><Link to="/admin-login">Admin Login</Link></li>
          </ul>
        </nav>
      )}

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/admin-login" element={<AdminLogin />} />
      </Routes>
    </div>
  );
}

function AppWithRouter() {
  return (
    <Router>
      <App />
    </Router>
  );
}

export default AppWithRouter;
