import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
  TextField,
  Tab,
  Tabs,
  Alert,
  CircularProgress
} from '@mui/material';
import { Web as WebIcon, Preview as PreviewIcon, Save as SaveIcon } from '@mui/icons-material';
import Header from './Header';
import config from '../config';
import axios from 'axios';

const PageContentManager = ({ user, onLogout }) => {
  const [homeContent, setHomeContent] = useState('');
  const [notFoundContent, setNotFoundContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [activeTab, setActiveTab] = useState('home');
  const [message, setMessage] = useState('');

  const API_BASE_URL = config.API_BASE_URL;

  useEffect(() => {
    loadPageContent();
  }, []);  // eslint-disable-line react-hooks/exhaustive-deps

  const loadPageContent = async () => {
    try {
      setLoading(true);
      const [homeResponse, notFoundResponse] = await Promise.all([
        axios.get(`${API_BASE_URL}/config/home_page_html`),
        axios.get(`${API_BASE_URL}/config/404_page_html`)
      ]);
      
      setHomeContent(homeResponse.data.config_value || '');
      setNotFoundContent(notFoundResponse.data.config_value || '');
    } catch (error) {
      console.error('Error loading page content:', error);
      setMessage('Failed to load page content');
    } finally {
      setLoading(false);
    }
  };

  const saveContent = async (configKey, content) => {
    try {
      setSaving(true);
      await axios.put(`${API_BASE_URL}/config/${configKey}`, {
        config_value: content
      });
      setMessage('Content saved successfully!');
      setTimeout(() => setMessage(''), 3000);
    } catch (error) {
      console.error('Error saving content:', error);
      setMessage('Failed to save content');
      setTimeout(() => setMessage(''), 3000);
    } finally {
      setSaving(false);
    }
  };

  const handleSave = () => {
    if (activeTab === 'home') {
      saveContent('home_page_html', homeContent);
    } else {
      saveContent('404_page_html', notFoundContent);
    }
  };

  const previewContent = () => {
    const content = activeTab === 'home' ? homeContent : notFoundContent;
    const newWindow = window.open();
    newWindow.document.write(content);
    newWindow.document.close();
  };

  if (loading) {
    return (
      <>
        <Header 
          user={user} 
          onLogout={onLogout} 
          title="Customize Page Content"
          showBackButton={true}
          backButtonText="Back"
        />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Box display="flex" justifyContent="center" alignItems="center" height="200px">
            <CircularProgress />
          </Box>
        </Container>
      </>
    );
  }

  return (
    <>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="Customize Page Content"
        showBackButton={true}
        backButtonText="Back"
      />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Paper sx={{ p: 3 }}>
          <Box sx={{ mb: 3 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <WebIcon sx={{ mr: 1, fontSize: 28 }} />
              <Typography variant="h4" component="h1">
                Customize Page Content
              </Typography>
            </Box>
            <Typography variant="body1" color="text.secondary">
              Customize your home page and 404 page HTML content
            </Typography>
          </Box>

          {message && (
            <Alert 
              severity={message.includes('successfully') ? 'success' : 'error'}
              sx={{ mb: 3 }}
            >
              {message}
            </Alert>
          )}

          <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
            <Tabs value={activeTab} onChange={(e, newValue) => setActiveTab(newValue)}>
              <Tab label="Home Page" value="home" />
              <Tab label="404 Page" value="404" />
            </Tabs>
          </Box>

          <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">
              {activeTab === 'home' ? 'Home Page HTML' : '404 Page HTML'}
            </Typography>
            <Box sx={{ display: 'flex', gap: 2 }}>
              <Button
                variant="outlined"
                startIcon={<PreviewIcon />}
                onClick={previewContent}
              >
                Preview
              </Button>
              <Button
                variant="contained"
                startIcon={<SaveIcon />}
                onClick={handleSave}
                disabled={saving}
              >
                {saving ? 'Saving...' : 'Save Changes'}
              </Button>
            </Box>
          </Box>

          <TextField
            fullWidth
            multiline
            rows={20}
            value={activeTab === 'home' ? homeContent : notFoundContent}
            onChange={(e) => 
              activeTab === 'home' 
                ? setHomeContent(e.target.value)
                : setNotFoundContent(e.target.value)
            }
            placeholder="Enter your HTML content here..."
            sx={{
              '& .MuiInputBase-input': {
                fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
                fontSize: '14px'
              }
            }}
          />

          <Box sx={{ mt: 3, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
            <Typography variant="subtitle2" sx={{ mb: 1, fontWeight: 'bold' }}>
              Tips:
            </Typography>
            <Typography variant="body2" component="ul" sx={{ pl: 2, m: 0 }}>
              <li>Use complete HTML documents including &lt;!DOCTYPE html&gt;</li>
              <li>Include all necessary CSS and JavaScript inline or via CDN</li>
              <li>Test your changes using the Preview button before saving</li>
              <li>Changes will be applied immediately after saving</li>
            </Typography>
          </Box>
        </Paper>
      </Container>
    </>
  );
};

export default PageContentManager;