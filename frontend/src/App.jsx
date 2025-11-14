import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Login from './components/Login';
import LinkList from './components/LinkList';
import CreateLink from './components/CreateLink';
import EditLink from './components/EditLink';
import ChangePassword from './components/ChangePassword';
import ImportExport from './components/ImportExport';
import Analysis from './components/Analysis';
import PageContentManager from './components/PageContentManager';
import About from './components/About';
import { setupAxiosInterceptors } from './utils/axiosConfig';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Handle automatic logout when 401 is received
  const handleUnauthorized = () => {
    setUser(null);
    // The localStorage cleanup is already handled in the interceptor
  };

  useEffect(() => {
    // Setup axios interceptors for automatic 401 handling
    setupAxiosInterceptors(handleUnauthorized);

    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (token && userData) {
      setUser(JSON.parse(userData));
    }
    setLoading(false);
  }, []);

  const login = (token, userData) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
    
    // Store default password status
    if (userData.is_default_password) {
      localStorage.setItem('needs_password_change', 'true');
    } else {
      localStorage.removeItem('needs_password_change');
    }
    
    setUser(userData);
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route 
            path="/admin/login" 
            element={
              user ? <Navigate to="/admin" /> : <Login onLogin={login} />
            } 
          />
          <Route 
            path="/admin" 
            element={
              user ? <LinkList user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/create" 
            element={
              user ? <CreateLink user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/edit/:id" 
            element={
              user ? <EditLink user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/change-password" 
            element={
              user ? <ChangePassword user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/import-export" 
            element={
              user ? <ImportExport user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/analysis" 
            element={
              user ? <Analysis user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/page-content" 
            element={
              user ? <PageContentManager user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          <Route 
            path="/admin/about" 
            element={
              user ? <About user={user} onLogout={logout} /> : <Navigate to="/admin/login" />
            } 
          />
          {/* Redirect root to admin */}
          <Route path="/" element={<Navigate to="/admin" />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;