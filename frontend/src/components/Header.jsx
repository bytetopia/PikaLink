import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  IconButton,
  Menu,
  MenuItem,
  Alert,
  Box,
  Collapse
} from '@mui/material';
import { 
  ArrowBack, 
  ExitToApp, 
  AccountCircle, 
  Settings,
  Assessment,
  Close,
  Web
} from '@mui/icons-material';

function Header({ 
  user, 
  onLogout, 
  title = "PikaLink",
  showBackButton = false,
  backButtonText = "Back",
  backDestination = "/admin"
}) {
  const [anchorEl, setAnchorEl] = useState(null);
  const [showPasswordAlert, setShowPasswordAlert] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // Check if user needs to change password
    const needsPasswordChange = localStorage.getItem('needs_password_change') === 'true';
    setShowPasswordAlert(needsPasswordChange);

    // Listen for password change events
    const handlePasswordChanged = (event) => {
      setShowPasswordAlert(event.detail.isDefaultPassword);
    };

    window.addEventListener('passwordChanged', handlePasswordChanged);

    return () => {
      window.removeEventListener('passwordChanged', handlePasswordChanged);
    };
  }, []);

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleChangePassword = () => {
    handleMenuClose();
    navigate('/admin/change-password');
  };

  const handleAnalysis = () => {
    handleMenuClose();
    navigate('/admin/analysis');
  };

  const handlePageContent = () => {
    handleMenuClose();
    navigate('/admin/page-content');
  };

  const handleLogout = () => {
    handleMenuClose();
    onLogout();
  };

  const handleBackClick = () => {
    navigate(backDestination);
  };

  const handleDismissPasswordAlert = () => {
    setShowPasswordAlert(false);
    // Don't remove from localStorage so it shows again on page refresh
    // Only remove when password is actually changed
  };

  const handleChangePasswordFromAlert = () => {
    setShowPasswordAlert(false);
    navigate('/admin/change-password');
  };

  return (
    <Box>
      <AppBar position="static">
        <Toolbar>
          {showBackButton && (
            <Button
              color="inherit"
              onClick={handleBackClick}
              startIcon={<ArrowBack />}
              sx={{ mr: 2 }}
            >
              {backButtonText}
            </Button>
          )}
          
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            {title}
          </Typography>
          
          <Typography variant="body2" sx={{ mr: 2 }}>
            Welcome, {user.username}
          </Typography>
          
          <IconButton
            color="inherit"
            onClick={handleMenuOpen}
            sx={{ mr: 1 }}
          >
            <AccountCircle />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            <MenuItem onClick={handleChangePassword}>
              <Settings sx={{ mr: 1 }} />
              Change Password
            </MenuItem>
            <MenuItem onClick={handleAnalysis}>
              <Assessment sx={{ mr: 1 }} />
              Access Analysis
            </MenuItem>
            <MenuItem onClick={handlePageContent}>
              <Web sx={{ mr: 1 }} />
              Customize Pages
            </MenuItem>
            <MenuItem onClick={handleLogout}>
              <ExitToApp sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      
      {/* Default Password Warning Banner */}
      <Collapse in={showPasswordAlert}>
        <Alert 
          severity="warning" 
          action={
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button 
                color="inherit" 
                size="small" 
                onClick={handleChangePasswordFromAlert}
              >
                Change Now
              </Button>
              <IconButton
                aria-label="close"
                color="inherit"
                size="small"
                onClick={handleDismissPasswordAlert}
              >
                <Close fontSize="inherit" />
              </IconButton>
            </Box>
          }
          sx={{ borderRadius: 0 }}
        >
          You are using the default password. Please change it for security reasons.
        </Alert>
      </Collapse>
    </Box>
  );
}

export default Header;