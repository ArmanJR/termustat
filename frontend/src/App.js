import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import AppRoutes from './routes/index.js';
import { ThemeProvider } from '@mui/material';
import muiTheme from './themes/muiTheme.js';

function App() {
  // get the current route
  const location = useLocation();
  // only show nav on HomePage
  const showNav = location.pathname === '/';

  return (
    <ThemeProvider theme={muiTheme}>
      {showNav && (
        <nav>
          <ul>
            <li><Link to="/">Home</Link></li>
            <li><Link to="/admin/login">Admin Login</Link></li>
          </ul>
        </nav>
      )}
      <AppRoutes />
    </ThemeProvider>
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
