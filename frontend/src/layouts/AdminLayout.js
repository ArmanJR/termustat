import { useState } from "react";
import { Link as RouterLink, Outlet } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import styles from "./AdminLayout.module.css";
import logo_full from "../images/logo-full-white.png";
import logo from "../images/logo.png";

import {
  IconButton,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  useTheme,
  useMediaQuery,
} from "@mui/material";

import {
  Dashboard,
  School,
  Menu,
  ExitToApp,
} from "@mui/icons-material";

const menuItems = [
  { label: "صفحه اصلی", icon: <Dashboard />, route: "/admin/dashboard" },
  { label: "دانشگاه‌ها", icon: <School />, route: "/admin/universities" },
];

export default function AdminLayout() {
  const { logout } = useAuth();

  // State to toggle the sidebar on mobile devices
  const [mobileOpen, setMobileOpen] = useState(false);
  const handleDrawerToggle = () => setMobileOpen(!mobileOpen);

  // Use theme breakpoints for responsive behavior
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  // Close the drawer on mobile when a menu item is selected
  const handleMenuItemClick = () => {
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  // Sidebar drawer content
  const drawer = (
    <List>
      {menuItems.map(({ label, icon, route }) => (
        <ListItemButton
          key={label}
          component={RouterLink}
          to={route}
          onClick={handleMenuItemClick}
        >
          <ListItemIcon>{icon}</ListItemIcon>
          <ListItemText primary={label} />
        </ListItemButton>
      ))}
      <hr />
      <ListItemButton key={"خروج"} onClick={logout}>
        <ListItemIcon>
          <ExitToApp />
        </ListItemIcon>
        <ListItemText primary={"خروج"} />
      </ListItemButton>
    </List>
  );

  return (
    <>
      <div className={styles.pageWrapper}>
        
        {/* Topbar */}
        <div className={styles.topbar}>
          {isMobile && (
            <IconButton onClick={handleDrawerToggle}>
              <Menu />
            </IconButton>
          )}
          <h1>
            {!isMobile && (
              <span style={{ color: "#309a9a" }}>&#9699; &nbsp;</span>
            )}
            پنل مدیریت
          </h1>
          <img src={isMobile ? logo : logo_full} className={styles.logo} />
        </div>

        {/* Main layout area: sidebar + content */}
        <div className={styles.main}>

          {/* Sidebar */}
          <div className={styles.sidebar}>
            <Drawer
              variant={isMobile ? "temporary" : "permanent"}
              open={isMobile ? mobileOpen : true}
              onClose={handleDrawerToggle}
              ModalProps={{ keepMounted: true }}
            >
              {drawer}
            </Drawer>
          </div>

          {/* Content */}
          <div className={styles.content}>
            <Outlet />
          </div>
        </div>
      </div>
    </>
  );
}
