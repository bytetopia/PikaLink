import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  IconButton,
  Menu,
  MenuItem
} from '@mui/material';
import { 
  ArrowBack, 
  ExitToApp, 
  AccountCircle, 
  Settings,
  Assessment 
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
  const navigate = useNavigate();

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

  const handleImportExport = () => {
    handleMenuClose();
    navigate('/admin/import-export');
  };

  const handleAnalysis = () => {
    handleMenuClose();
    navigate('/admin/analysis');
  };

  const handleLogout = () => {
    handleMenuClose();
    onLogout();
  };

  const handleBackClick = () => {
    navigate(backDestination);
  };

  return (
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
          <MenuItem onClick={handleLogout}>
            <ExitToApp sx={{ mr: 1 }} />
            Logout
          </MenuItem>
        </Menu>
      </Toolbar>
    </AppBar>
  );
}

export default Header;