import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';

// Page components
import AdminLoginPage from './pages/AdminLoginPage';
const HomePage = () => <div>Home Page</div>;

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
            <li><Link to="/admin-login">Admin Login Page</Link></li>
          </ul>
        </nav>
      )}

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/admin-login" element={<AdminLoginPage />} />
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
