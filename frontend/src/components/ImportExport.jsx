import { useState } from 'react';
import {
  Paper,
  Typography,
  Button,
  Box,
  Alert,
  CircularProgress,
  Divider,
  List,
  ListItem,
  ListItemText,
  Chip,
  Card,
  CardContent
} from '@mui/material';
import { CloudUpload, Download, FileUpload, GetApp } from '@mui/icons-material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

const API_BASE_URL = config.API_BASE_URL;

function ImportExport({ user, onLogout }) {
  const [importFile, setImportFile] = useState(null);
  const [importLoading, setImportLoading] = useState(false);
  const [importResult, setImportResult] = useState(null);
  const [exportLoading, setExportLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    if (file && file.type === 'text/csv') {
      setImportFile(file);
      setError('');
      setImportResult(null);
    } else {
      setError('Please select a valid CSV file');
      setImportFile(null);
    }
  };

  const handleImport = async () => {
    if (!importFile) {
      setError('Please select a file first');
      return;
    }

    setImportLoading(true);
    setError('');
    setSuccess('');
    setImportResult(null);

    try {
      const formData = new FormData();
      formData.append('file', importFile);

      const response = await axios.post(
        `${API_BASE_URL}/import`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }
      );

      setImportResult(response.data);
      if (response.data.success_count > 0) {
        setSuccess(`Successfully imported ${response.data.success_count} links!`);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to import links');
    } finally {
      setImportLoading(false);
    }
  };

  const handleExport = async () => {
    setExportLoading(true);
    setError('');
    setSuccess('');

    try {
      const response = await axios.get(
        `${API_BASE_URL}/export`,
        {
          responseType: 'blob'
        }
      );

      // Create download link
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'pikalink_export.csv');
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      setSuccess('Links exported successfully!');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to export links');
    } finally {
      setExportLoading(false);
    }
  };

  const resetImport = () => {
    setImportFile(null);
    setImportResult(null);
    setError('');
    setSuccess('');
    // Reset file input
    const fileInput = document.getElementById('csv-file-input');
    if (fileInput) {
      fileInput.value = '';
    }
  };

  return (
    <Box sx={{ width: '100vw', minHeight: '100vh', margin: 0, padding: 0 }}>
      <Header 
        user={user} 
        onLogout={onLogout} 
        title="Import / Export Links"
        showBackButton={true}
      />
      
      <Box sx={{ p: 3 }}>
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        {/* Import Section */}
        <Paper sx={{ p: 3, mb: 4 }}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <CloudUpload />
            Import Links
          </Typography>
          
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Upload a CSV file to import multiple links at once. The CSV should have columns: original_url, short_code, title
          </Typography>

          <Box sx={{ mb: 3 }}>
            <input
              id="csv-file-input"
              type="file"
              accept=".csv"
              onChange={handleFileSelect}
              style={{ display: 'none' }}
            />
            <label htmlFor="csv-file-input">
              <Button
                variant="outlined"
                component="span"
                startIcon={<FileUpload />}
                sx={{ mr: 2 }}
              >
                Choose CSV File
              </Button>
            </label>
            
            {importFile && (
              <Chip 
                label={importFile.name} 
                onDelete={resetImport}
                sx={{ ml: 1 }}
              />
            )}
          </Box>

          <Box sx={{ mb: 2 }}>
            <Button
              variant="contained"
              onClick={handleImport}
              disabled={!importFile || importLoading}
              startIcon={importLoading ? <CircularProgress size={20} /> : <CloudUpload />}
            >
              {importLoading ? 'Importing...' : 'Import Links'}
            </Button>
          </Box>

          {/* Import Results */}
          {importResult && (
            <Card sx={{ mt: 2 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Import Results
                </Typography>
                
                <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                  <Chip 
                    label={`Total: ${importResult.total_count}`} 
                    color="default" 
                  />
                  <Chip 
                    label={`Success: ${importResult.success_count}`} 
                    color="success" 
                  />
                  <Chip 
                    label={`Failed: ${importResult.failed_count}`} 
                    color={importResult.failed_count > 0 ? "error" : "default"}
                  />
                </Box>

                {importResult.failed_count > 0 && (
                  <Box>
                    <Typography variant="subtitle2" sx={{ mb: 1 }}>
                      Failed Items:
                    </Typography>
                    <List dense>
                      {importResult.failed_items.map((item, index) => (
                        <ListItem key={index} sx={{ py: 0.5 }}>
                          <ListItemText 
                            primary={item}
                            sx={{ 
                              '& .MuiListItemText-primary': { 
                                fontSize: '0.875rem',
                                color: 'error.main'
                              }
                            }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </Box>
                )}
              </CardContent>
            </Card>
          )}
        </Paper>

        <Divider sx={{ my: 4 }} />

        {/* Export Section */}
        <Paper sx={{ p: 3 }}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Download />
            Export Links
          </Typography>
          
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Download all your links as a CSV file. This file can be imported into another PikaLink instance.
          </Typography>

          <Button
            variant="contained"
            onClick={handleExport}
            disabled={exportLoading}
            startIcon={exportLoading ? <CircularProgress size={20} /> : <GetApp />}
          >
            {exportLoading ? 'Exporting...' : 'Export All Links'}
          </Button>
        </Paper>

        {/* CSV Format Info */}
        <Paper sx={{ p: 3, mt: 4, bgcolor: 'grey.50' }}>
          <Typography variant="h6" gutterBottom>
            CSV Format Information
          </Typography>
          
          <Typography variant="body2" sx={{ mb: 2 }}>
            The CSV format uses the following columns:
          </Typography>
          
          <Box component="pre" sx={{ 
            bgcolor: 'grey.100', 
            p: 2, 
            borderRadius: 1, 
            fontSize: '0.875rem',
            overflow: 'auto'
          }}>
{`original_url,short_code,title
https://example.com,example1,Example Website
https://google.com,search,"Google Search Engine"
https://site.com/path,custom-url,"Site with, comma in title"`}
          </Box>
          
          <Typography variant="body2" sx={{ mt: 2 }}>
            • Fields containing commas, quotes, or newlines will be automatically escaped<br/>
            • short_code must be unique and follow the same rules as manual creation<br/>
            • Existing short codes will be skipped during import
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
}

export default ImportExport;