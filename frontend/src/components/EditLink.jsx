import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Alert
} from '@mui/material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

const API_BASE_URL = config.API_BASE_URL;

function EditLink({ user, onLogout }) {
  const [originalUrl, setOriginalUrl] = useState('');
  const [title, setTitle] = useState('');
  const [shortCode, setShortCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(true);
  const navigate = useNavigate();
  const { id } = useParams();

  const fetchLink = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/links/${id}`);
      const link = response.data;
      
      setOriginalUrl(link.original_url);
      setTitle(link.title || '');
      setShortCode(link.short_code || '');
    } catch (err) {
      if (err.response?.status === 404) {
        setError('Link not found');
      } else {
        setError('Failed to fetch link details');
      }
    } finally {
      setFetchLoading(false);
    }
  };

  useEffect(() => {
    fetchLink();
  }, [id]);  // eslint-disable-line react-hooks/exhaustive-deps

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await axios.put(`${API_BASE_URL}/links/${id}`, {
        original_url: originalUrl,
        title: title,
        short_code: shortCode
      });

      navigate('/admin');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update link');
    } finally {
      setLoading(false);
    }
  };

  if (fetchLoading) return <div>Loading...</div>;

  return (
    <>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="Edit Link"
        showBackButton={true}
        backButtonText="Back"
        backDestination="/admin"
      />

      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Paper elevation={3} sx={{ padding: 4 }}>
          <Typography variant="h4" gutterBottom>
            Edit Short Link
          </Typography>

          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="originalUrl"
              label="Original URL"
              name="originalUrl"
              value={originalUrl}
              onChange={(e) => setOriginalUrl(e.target.value)}
              type="url"
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="shortCode"
              label="Short Code (Slug)"
              name="shortCode"
              value={shortCode}
              onChange={(e) => setShortCode(e.target.value)}
              helperText="1-50 characters, only letters, numbers, hyphens, and underscores allowed"
            />
            <TextField
              margin="normal"
              fullWidth
              id="title"
              label="Title (Optional)"
              name="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
              <Button
                type="submit"
                variant="contained"
                disabled={loading || !originalUrl || !shortCode}
              >
                {loading ? 'Updating...' : 'Update Link'}
              </Button>
              <Button
                variant="outlined"
                onClick={() => navigate('/admin')}
              >
                Cancel
              </Button>
            </Box>
          </Box>
        </Paper>
      </Container>
    </>
  );
}

export default EditLink;