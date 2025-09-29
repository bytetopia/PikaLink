import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Alert,
  CircularProgress
} from '@mui/material';
import { Edit, Delete, Add, ContentCopy, ImportExport, Analytics } from '@mui/icons-material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

const API_BASE_URL = config.API_BASE_URL;

function LinkList({ user, onLogout }) {
  const [links, setLinks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  // Get the base URL for short redirects - now at root path
  const getRedirectBaseUrl = () => {
    return config.SHORT_LINK_BASE_URL;
  };

  const getAuthHeaders = () => ({
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`
    }
  });

  const fetchLinks = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/links`, getAuthHeaders());
      setLinks(response.data || []);
    } catch (err) {
      setError('Failed to fetch links');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLinks();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this link?')) {
      try {
        await axios.delete(`${API_BASE_URL}/links/${id}`, getAuthHeaders());
        setLinks(links.filter(link => link.id !== id));
      } catch (err) {
        setError('Failed to delete link');
      }
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="PikaLink Dashboard"
        showBackButton={false}
      />

      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h4">Your Links</Typography>
          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button
              variant="outlined"
              startIcon={<Analytics />}
              onClick={() => navigate('/admin/analysis')}
            >
              Analytics
            </Button>
            <Button
              variant="outlined"
              startIcon={<ImportExport />}
              onClick={() => navigate('/admin/import-export')}
            >
              Import / Export
            </Button>
            <Button
              variant="contained"
              startIcon={<Add />}
              onClick={() => navigate('/admin/create')}
            >
              Create New Link
            </Button>
          </Box>
        </Box>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Title</TableCell>
                <TableCell>Short Link</TableCell>
                <TableCell>Original URL</TableCell>
                <TableCell>Created</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center" sx={{ py: 8 }}>
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
                      <CircularProgress />
                      <Typography variant="body2" color="text.secondary">
                        Loading your links...
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : links.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    No links found. Create your first link!
                  </TableCell>
                </TableRow>
              ) : (
                links.map((link) => (
                  <TableRow key={link.id}>
                    <TableCell>{link.title || 'Untitled'}</TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Typography variant="body2" sx={{ mr: 1 }}>
                          {`${getRedirectBaseUrl()}/${link.short_code}`}
                        </Typography>
                        <IconButton
                          size="small"
                          onClick={() => copyToClipboard(`${getRedirectBaseUrl()}/${link.short_code}`)}
                        >
                          <ContentCopy fontSize="small" />
                        </IconButton>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" sx={{ 
                        maxWidth: 200, 
                        overflow: 'hidden', 
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap'
                      }}>
                        {link.original_url}
                      </Typography>
                    </TableCell>
                    <TableCell>{formatDate(link.created_at)}</TableCell>
                    <TableCell>
                      <IconButton
                        onClick={() => navigate(`/admin/edit/${link.id}`)}
                        color="primary"
                      >
                        <Edit />
                      </IconButton>
                      <IconButton
                        onClick={() => handleDelete(link.id)}
                        color="error"
                      >
                        <Delete />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Container>
    </>
  );
}

export default LinkList;