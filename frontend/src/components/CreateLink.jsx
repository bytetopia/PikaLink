import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Alert,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Collapse
} from '@mui/material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

const API_BASE_URL = config.API_BASE_URL;

function CreateLink({ user, onLogout }) {
  const [originalUrl, setOriginalUrl] = useState('');
  const [title, setTitle] = useState('');
  const [codeType, setCodeType] = useState('random');
  const [customCode, setCustomCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const getAuthHeaders = () => ({
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`
    }
  });

  const validateCustomCode = (code) => {
    if (!code || code.trim() === '') {
      return 'Custom code cannot be empty';
    }
    
    if (code.length < 3 || code.length > 50) {
      return 'Custom code must be 3-50 characters long';
    }
    
    if (!/^[a-zA-Z0-9_-]+$/.test(code)) {
      return 'Custom code can only contain letters, numbers, hyphens, and underscores';
    }
    
    return null;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    // Validate custom code if selected
    if (codeType === 'custom') {
      const validationError = validateCustomCode(customCode);
      if (validationError) {
        setError(validationError);
        setLoading(false);
        return;
      }
    }

    try {
      await axios.post(`${API_BASE_URL}/links`, {
        original_url: originalUrl,
        title: title,
        use_custom_code: codeType === 'custom',
        custom_code: codeType === 'custom' ? customCode.trim() : ''
      }, getAuthHeaders());

      navigate('/admin');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create link');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="Create New Link"
        showBackButton={true}
        backButtonText="Back to Links"
        backDestination="/admin"
        showSettings={false}
      />

      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Paper elevation={3} sx={{ padding: 4 }}>
          <Typography variant="h4" gutterBottom>
            Create New Short Link
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
              placeholder="https://example.com"
              value={originalUrl}
              onChange={(e) => setOriginalUrl(e.target.value)}
              type="url"
            />
            
            <TextField
              margin="normal"
              fullWidth
              id="title"
              label="Title (Optional)"
              name="title"
              placeholder="Enter a title for your link"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />

            <FormControl component="fieldset" sx={{ mt: 2, mb: 2 }}>
              <FormLabel component="legend">Short Code Type</FormLabel>
              <RadioGroup
                value={codeType}
                onChange={(e) => setCodeType(e.target.value)}
                row
              >
                <FormControlLabel
                  value="random"
                  control={<Radio />}
                  label="Random Generated"
                />
                <FormControlLabel
                  value="custom"
                  control={<Radio />}
                  label="Custom"
                />
              </RadioGroup>
            </FormControl>

            <Collapse in={codeType === 'custom'}>
              <TextField
                margin="normal"
                fullWidth
                id="customCode"
                label="Custom Short Code"
                name="customCode"
                placeholder="my-awesome-link"
                value={customCode}
                onChange={(e) => setCustomCode(e.target.value)}
                helperText="3-50 characters, letters, numbers, hyphens, and underscores only"
                required={codeType === 'custom'}
              />
            </Collapse>

            <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
              <Button
                type="submit"
                variant="contained"
                disabled={loading || !originalUrl || (codeType === 'custom' && !customCode.trim())}
              >
                {loading ? 'Creating...' : 'Create Link'}
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

export default CreateLink;