import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  CircularProgress,
  Alert
} from '@mui/material';
import { Info } from '@mui/icons-material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

function About({ user, onLogout }) {
  const [version, setVersion] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchVersion();
  }, []);

  const fetchVersion = async () => {
    try {
      const response = await axios.get(`${config.API_BASE_URL}/version`);
      setVersion(response.data.version);
    } catch (err) {
      setError('Failed to load version information');
      console.error('Error fetching version:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="About PikaLink"
        showBackButton={true}
      />
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Info sx={{ fontSize: 40, mr: 2, color: 'primary.main' }} />
            <Typography variant="h5" component="h1">
              About PikaLink
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box sx={{ mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Version Information
            </Typography>
            {loading ? (
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                <CircularProgress size={20} sx={{ mr: 2 }} />
                <Typography>Loading version...</Typography>
              </Box>
            ) : (
              <Typography variant="body1" sx={{ mt: 1 }}>
                <strong>Version:</strong> {version || 'Unknown'}
              </Typography>
            )}
          </Box>

          <Box sx={{ mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Open Source
            </Typography>
            <Typography variant="body1">
              PikaLink is open source and available on{' '}
              <a 
                href="https://github.com/bytetopia/PikaLink" 
                target="_blank" 
                rel="noopener noreferrer"
              >
                GitHub
              </a>.
            </Typography>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom>
              Contact
            </Typography>
            <Typography variant="body1" component="div">
              <ul>
                <li>Email: hi@idealland.app</li>
                <li>File an issue: <a href="https://github.com/bytetopia/PikaLink/issues" target="_blank" rel="noopener noreferrer">GitHub Issues</a></li>
              </ul>
            </Typography>
          </Box>
        </Paper>
      </Container>
    </div>
  );
}

export default About;
