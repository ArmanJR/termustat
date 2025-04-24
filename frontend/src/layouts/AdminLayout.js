import { useState } from "react";
import { Link as RouterLink, Outlet } from "react-router-dom";

import { adminLogout } from "../api/adminAuth.js";

import styles from "./AdminLayout.module.css";
import logo_full from "../images/logo-full-white.png";
import logo from "../images/logo.png";

// Material UI components
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

// Material UI icons
import {
  Dashboard,
  Menu,
  ExitToApp,
} from "@mui/icons-material";

// Sidebar menu items
const menuItems = [
  { label: "صفحه اصلی", icon: <Dashboard />, route: "/admin/dashboard" }
];

export default function AdminLayout() {

  // State to toggle the sidebar on mobile devices
  const [mobileOpen, setMobileOpen] = useState(false);
  const handleDrawerToggle = () => setMobileOpen(!mobileOpen);

  // Use theme breakpoints for responsive behavior
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  // Sidebar drawer content (menu and logout)
  const drawer = (
    <List>
      {menuItems.map(({ label, icon, route }) => (
        <ListItemButton key={label} component={RouterLink} to={route}>
          <ListItemIcon>{icon}</ListItemIcon>
          <ListItemText
            primary={label}
            primaryTypographyProps={{ fontFamily: "Vazirmatn, sans-serif" }}
          />
        </ListItemButton>
      ))}
      <hr />
      <ListItemButton key={"خروج"} onClick={adminLogout}>
        <ListItemIcon>
          <ExitToApp />
        </ListItemIcon>
        <ListItemText
          primary={"خروج"}
          primaryTypographyProps={{ fontFamily: "Vazirmatn, sans-serif" }}
        />
      </ListItemButton>
    </List>
  );

  return (
    <div className={styles.pageWrapper}>

      {/* Topbar: logo and title */}
      <div className={styles.topbar}>

        {/* Mobile hamburger menu button */}
        {isMobile && (
          <IconButton
            color="inherit"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2 }}
          >
            <Menu />
          </IconButton>
        )}
        {/* Panel title */}
        <h1>
          {!isMobile && (
            <span style={{ color: "#309a9a" }}>&#9699; &nbsp;</span>
          )}
          پنل مدیریت
        </h1>
        {/* Logo based on screen size */}
        <img src={isMobile ? logo : logo_full} className={styles.logo} />
      </div>

      {/* Main layout area: sidebar + content */}
      <div className={styles.main}>

        {/* Sidebar */}
        <div className={styles.sidebar}>
          <Drawer
            anchor={"right"}
            variant={isMobile ? "temporary" : "permanent"}
            open={isMobile ? mobileOpen : true}
            onClose={handleDrawerToggle}
            ModalProps={{ keepMounted: true }}
            sx={{
              "& .MuiDrawer-paper": {
                position: "relative",
                width: 240,
                border: "none",
                boxSizing: "border-box",
                direction: "rtl",
                ...(isMobile && {
                  right: 0,
                  left: "auto",
                  position: "absolute",
                }),
              },
              "& .MuiListItemText-root": {
                textAlign: "right",
              },
            }}
          >
            {drawer}
          </Drawer>
        </div>

        {/* Content */}
        <div className={styles.content}>
          <Outlet /> {/* Routed page content gets injected here */}
        </div>

      </div>
    </div>
  );
}
