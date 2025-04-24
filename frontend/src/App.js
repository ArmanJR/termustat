import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import AppRoutes from './routes/index.js';

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
            <li><Link to="/admin/login">Admin Login</Link></li>
          </ul>
        </nav>
      )}
      <AppRoutes />
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
