import React, { useState, useEffect } from 'react';
import { Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import {
    Container,
    Paper,
    Typography,
    Button,
    Box,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Grid,
    Card,
    CardContent,
    Alert,
    CircularProgress,
    Divider
} from '@mui/material';
import { Analytics as AnalyticsIcon, TrendingUp } from '@mui/icons-material';
import Header from './Header';
import axios from 'axios';
import config from '../config';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const Analysis = ({ user, onLogout }) => {
    const [months, setMonths] = useState([]);
    const [selectedMonth, setSelectedMonth] = useState('');
    const [analysisData, setAnalysisData] = useState(null);
    const [selectedShortURL, setSelectedShortURL] = useState('');
    const [activeShortURL, setActiveShortURL] = useState(''); // The URL that was used in the last analysis
    const [selectedStatusCode, setSelectedStatusCode] = useState('');
    const [activeStatusCode, setActiveStatusCode] = useState(''); // The status code that was used in the last analysis
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);
    const [analyzing, setAnalyzing] = useState(false);

    useEffect(() => {
        const fetchMonths = async () => {
            try {
                setLoading(true);
                const response = await axios.get(`${config.API_BASE_URL}/analysis/months`);
                const data = response.data;
                console.log('Months data received:', data); // Debug log
                const monthsArray = Array.isArray(data) ? data.filter(month => month != null).map(month => String(month)) : [];
                setMonths(monthsArray);
                if (monthsArray.length > 0) {
                    setSelectedMonth(String(monthsArray[0]));
                }
                setError('');
            } catch (error) {
                setError(error.response?.data?.error || error.message || 'Failed to fetch months');
            } finally {
                setLoading(false);
            }
        };
        fetchMonths();
    }, []);

    const handleAnalyze = async () => {
        if (!selectedMonth) {
            setError('Please select a month.');
            return;
        }
        try {
            setAnalyzing(true);
            setError('');
            let url = `${config.API_BASE_URL}/analyze?month=${selectedMonth}`;
            if (selectedShortURL) {
                url += `&short_url=${selectedShortURL}`;
            }
            if (selectedStatusCode) {
                url += `&status_code=${selectedStatusCode}`;
            }
            const response = await axios.get(url);
            const data = response.data;
            console.log('Raw analysis data received:', JSON.stringify(data, null, 2)); // Enhanced debug log
            
            // Validate and sanitize the data more thoroughly
            const sanitizedData = {
                total_clicks: Number(data?.total_clicks) || 0,
                status_distribution: data?.status_distribution && typeof data.status_distribution === 'object' && !Array.isArray(data.status_distribution) ? data.status_distribution : {},
                ua_distribution: data?.ua_distribution && typeof data.ua_distribution === 'object' && !Array.isArray(data.ua_distribution) ? data.ua_distribution : {},
                referer_distribution: data?.referer_distribution && typeof data.referer_distribution === 'object' && !Array.isArray(data.referer_distribution) ? data.referer_distribution : {},
                browser_distribution: data?.browser_distribution && typeof data.browser_distribution === 'object' && !Array.isArray(data.browser_distribution) ? data.browser_distribution : {},
                os_distribution: data?.os_distribution && typeof data.os_distribution === 'object' && !Array.isArray(data.os_distribution) ? data.os_distribution : {},
                device_distribution: data?.device_distribution && typeof data.device_distribution === 'object' && !Array.isArray(data.device_distribution) ? data.device_distribution : {},
                bot_distribution: data?.bot_distribution && typeof data.bot_distribution === 'object' && !Array.isArray(data.bot_distribution) ? data.bot_distribution : {},
                short_urls: Array.isArray(data?.short_urls) ? data.short_urls.filter(url => url != null).map(url => String(url)) : [],
                status_codes: Array.isArray(data?.status_codes) ? data.status_codes.filter(code => code != null).map(code => String(code)) : []
            };
            
            console.log('Sanitized analysis data:', JSON.stringify(sanitizedData, null, 2)); // Debug sanitized data
            setAnalysisData(sanitizedData);
            setActiveShortURL(selectedShortURL); // Update the active URL after successful analysis
            setActiveStatusCode(selectedStatusCode); // Update the active status code after successful analysis
        } catch (error) {
            console.error('Analysis error:', error); // Debug log
            setError(error.response?.data?.error || error.message || 'Failed to fetch analysis data');
            setAnalysisData(null);
        } finally {
            setAnalyzing(false);
        }
    };

    const getChartData = (distribution, color = 'rgba(25, 118, 210, 0.6)') => {
        if (!distribution || typeof distribution !== 'object' || Array.isArray(distribution)) {
            return { labels: [], datasets: [] };
        }
        
        try {
            // Process entries to handle special cases
            const entries = Object.entries(distribution).filter(([key, value]) => key != null && value != null);
            
            // Create safe labels and data
            const labels = entries.map(([key]) => {
                let label = String(key).trim();
                // Handle empty or dash-only labels
                if (!label || label === '-' || label === 'null' || label === 'undefined') {
                    label = 'Direct/Unknown';
                }
                // Truncate very long labels (like user agents) for better chart display
                if (label.length > 30) {
                    label = label.substring(0, 27) + '...';
                }
                // Remove any potential problematic characters that might cause React issues
                return label.replace(/[\x00-\x1F\x7F-\x9F]/g, '');
            });
            
            const data = entries.map(([, value]) => {
                const num = Number(value);
                return isNaN(num) || num < 0 ? 0 : Math.floor(num);
            });
            
            // Ensure we have valid data
            if (labels.length === 0 || data.length === 0) {
                return { labels: [], datasets: [] };
            }
            
            return {
                labels: labels,
                datasets: [
                    {
                        label: 'Count',
                        data: data,
                        backgroundColor: color,
                        borderColor: typeof color === 'string' ? color.replace('0.6', '1') : color,
                        borderWidth: 1,
                    },
                ],
            };
        } catch (error) {
            console.error('Error processing chart data:', error, distribution);
            return { labels: [], datasets: [] };
        }
    };

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
            intersect: false,
            mode: 'index'
        },
        plugins: {
            legend: {
                position: 'top',
                labels: {
                    generateLabels: function(chart) {
                        const original = ChartJS.overrides.bar.plugins.legend.labels.generateLabels;
                        const labels = original.call(this, chart);
                        return labels.map(label => ({
                            ...label,
                            text: String(label.text || 'Count')
                        }));
                    }
                }
            },
            title: {
                display: false,
            },
            tooltip: {
                callbacks: {
                    title: function(context) {
                        return String(context[0].label || '');
                    },
                    label: function(context) {
                        return `Count: ${String(context.parsed.y || 0)}`;
                    }
                }
            }
        },
        scales: {
            x: {
                ticks: {
                    maxRotation: 45,
                    minRotation: 0,
                    callback: function(value, index) {
                        const label = this.getLabelForValue(value);
                        return String(label || '').length > 15 ? String(label || '').substring(0, 12) + '...' : String(label || '');
                    }
                }
            },
            y: {
                beginAtZero: true,
                ticks: {
                    precision: 0,
                    callback: function(value) {
                        return String(Math.floor(Number(value) || 0));
                    }
                }
            }
        }
    };

    if (loading) {
        return (
            <>
                <Header
                    user={user} 
                    onLogout={onLogout} 
                    title="Analysis"
                    showBackButton={true}
                    backButtonText="Back"
                    backDestination="/admin"
                />
                <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
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
                title="Analysis"
                showBackButton={true}
                backButtonText="Back"
                backDestination="/admin"
            />
            <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                <Paper elevation={3} sx={{ p: 3 }}>
                    <Box display="flex" alignItems="center" mb={3}>
                        <AnalyticsIcon sx={{ mr: 2, fontSize: 32, color: 'primary.main' }} />
                        <Typography variant="h4" component="h1" fontWeight="bold">
                            Log Analysis
                        </Typography>
                    </Box>

                    {error && (
                        <Alert severity="error" sx={{ mb: 3 }}>
                            {String(error)}
                        </Alert>
                    )}

                    {/* Controls Section */}
                    <Paper elevation={1} sx={{ p: 3, mb: 3, bgcolor: 'grey.50' }}>
                        <Typography variant="h6" gutterBottom>
                            Analysis Controls
                        </Typography>
                        <Grid container spacing={3} alignItems="center">
                            <Grid item xs={12} sm={3}>
                                <FormControl fullWidth>
                                    <InputLabel>Select Month</InputLabel>
                                    <Select
                                        value={selectedMonth || ''}
                                        label="Select Month"
                                        onChange={(e) => setSelectedMonth(String(e.target.value || ''))}
                                        disabled={analyzing}
                                    >
                                        {months.map((month, index) => (
                                            <MenuItem key={month || index} value={month}>
                                                {String(month || '')}
                                            </MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                            </Grid>
                            
                            {analysisData && analysisData.short_urls && (
                                <Grid item xs={12} sm={3}>
                                    <FormControl fullWidth>
                                        <InputLabel>Short URL (Optional)</InputLabel>
                                        <Select
                                            value={selectedShortURL || ''}
                                            label="Short URL (Optional)"
                                            onChange={(e) => setSelectedShortURL(String(e.target.value || ''))}
                                            disabled={analyzing}
                                        >
                                            <MenuItem value="">All URLs</MenuItem>
                                            {analysisData.short_urls.map((url, index) => (
                                                <MenuItem key={url || index} value={url}>
                                                    {String(url || '')}
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                </Grid>
                            )}
                            
                            {analysisData && analysisData.status_codes && (
                                <Grid item xs={12} sm={3}>
                                    <FormControl fullWidth>
                                        <InputLabel>Status Code (Optional)</InputLabel>
                                        <Select
                                            value={selectedStatusCode || ''}
                                            label="Status Code (Optional)"
                                            onChange={(e) => setSelectedStatusCode(String(e.target.value || ''))}
                                            disabled={analyzing}
                                        >
                                            <MenuItem value="">All Status Codes</MenuItem>
                                            {analysisData.status_codes.map((code, index) => (
                                                <MenuItem key={code || index} value={code}>
                                                    {String(code || '')}
                                                </MenuItem>
                                            ))}
                                        </Select>
                                    </FormControl>
                                </Grid>
                            )}
                            
                            <Grid item xs={12} sm={3}>
                                <Button
                                    variant="contained"
                                    size="large"
                                    onClick={handleAnalyze}
                                    disabled={!selectedMonth || analyzing}
                                    startIcon={analyzing ? <CircularProgress size={20} color="inherit" /> : <TrendingUp />}
                                    fullWidth
                                >
                                    {analyzing ? 'Analyzing...' : 'Analyze'}
                                </Button>
                            </Grid>
                        </Grid>
                    </Paper>

                    {/* Results Section */}
                    {analysisData && typeof analysisData === 'object' && !Array.isArray(analysisData) && (
                        <>
                            <Divider sx={{ my: 3 }} />
                            
                            {/* Summary Card */}
                            <Card elevation={2} sx={{ mb: 3, bgcolor: 'primary.main', color: 'white' }}>
                                <CardContent>
                                    <Box display="flex" alignItems="center" justifyContent="center">
                                        <TrendingUp sx={{ mr: 2, fontSize: 32 }} />
                                        <Typography variant="h4" component="div" fontWeight="bold">
                                            {String(Number(analysisData?.total_clicks || 0).toLocaleString())}
                                        </Typography>
                                        <Typography variant="h6" sx={{ ml: 1 }}>
                                            Total Clicks
                                        </Typography>
                                    </Box>
                                    {(activeShortURL || activeStatusCode) && (
                                        <Typography variant="body2" align="center" sx={{ mt: 1, opacity: 0.9 }}>
                                            {activeShortURL && `for ${String(activeShortURL)}`}
                                            {activeShortURL && activeStatusCode && ' • '}
                                            {activeStatusCode && `status ${String(activeStatusCode)}`}
                                        </Typography>
                                    )}
                                </CardContent>
                            </Card>

                            {/* Charts Grid */}
                            <Grid container spacing={3}>
                                {analysisData.status_distribution && Object.keys(analysisData.status_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    HTTP Status Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    {/* Temporarily replaced with table - Chart code commented out
                                                    <Bar 
                                                        data={getChartData(analysisData.status_distribution, 'rgba(76, 175, 80, 0.6)')} 
                                                        options={chartOptions}
                                                    />
                                                    */}
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>Status Code</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.status_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px' }}>{String(key)}</td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#4caf50' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.ua_distribution && Object.keys(analysisData.ua_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    User-Agent Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    {/* Temporarily replaced with table - Chart code commented out
                                                    <Bar 
                                                        data={getChartData(analysisData.ua_distribution, 'rgba(255, 152, 0, 0.6)')} 
                                                        options={{
                                                            ...chartOptions,
                                                            scales: {
                                                                ...chartOptions.scales,
                                                                x: {
                                                                    ticks: {
                                                                        maxRotation: 45,
                                                                        minRotation: 45,
                                                                        callback: function(value, index) {
                                                                            const label = this.getLabelForValue(value);
                                                                            return label.length > 20 ? label.substring(0, 20) + '...' : label;
                                                                        }
                                                                    }
                                                                }
                                                            }
                                                        }}
                                                    />
                                                    */}
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>User Agent</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.ua_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px', wordBreak: 'break-all', fontSize: '0.85em' }}>
                                                                        {String(key).length > 80 ? String(key).substring(0, 80) + '...' : String(key)}
                                                                    </td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#ff9800' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.referer_distribution && Object.keys(analysisData.referer_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    Referer Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    {/* Temporarily replaced with table - Chart code commented out
                                                    <Bar 
                                                        data={getChartData(analysisData.referer_distribution, 'rgba(156, 39, 176, 0.6)')} 
                                                        options={chartOptions}
                                                    />
                                                    */}
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>Referer</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.referer_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px', wordBreak: 'break-all' }}>
                                                                        {key === '-' ? 'Direct/Unknown' : String(key)}
                                                                    </td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#9c27b0' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.browser_distribution && Object.keys(analysisData.browser_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    Browser Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>Browser</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.browser_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px' }}>{String(key)}</td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#2196f3' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.os_distribution && Object.keys(analysisData.os_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    Operating System Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>OS</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.os_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px' }}>{String(key)}</td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#009688' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.device_distribution && Object.keys(analysisData.device_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    Device Type Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>Device Type</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.device_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px' }}>{String(key)}</td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: '#795548' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                                
                                {analysisData.bot_distribution && Object.keys(analysisData.bot_distribution).length > 0 && (
                                    <Grid item xs={12} lg={4}>
                                        <Card elevation={2} sx={{ height: 400 }}>
                                            <CardContent sx={{ height: '100%' }}>
                                                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">
                                                    Bot vs Human Distribution
                                                </Typography>
                                                <Box sx={{ height: 300, overflow: 'auto' }}>
                                                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                                                        <thead>
                                                            <tr style={{ borderBottom: '2px solid #e0e0e0' }}>
                                                                <th style={{ padding: '8px', textAlign: 'left' }}>Traffic Type</th>
                                                                <th style={{ padding: '8px', textAlign: 'right' }}>Count</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {Object.entries(analysisData.bot_distribution)
                                                                .sort(([, a], [, b]) => Number(b) - Number(a))
                                                                .map(([key, value]) => (
                                                                <tr key={key} style={{ borderBottom: '1px solid #f0f0f0' }}>
                                                                    <td style={{ padding: '8px' }}>{String(key)}</td>
                                                                    <td style={{ padding: '8px', textAlign: 'right', fontWeight: 'bold', color: key === 'Bot' ? '#f44336' : '#4caf50' }}>{String(value)}</td>
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                </Box>
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                )}
                            </Grid>

                            {/* No data message */}
                            {(!analysisData.status_distribution || Object.keys(analysisData.status_distribution).length === 0) &&
                             (!analysisData.ua_distribution || Object.keys(analysisData.ua_distribution).length === 0) &&
                             (!analysisData.referer_distribution || Object.keys(analysisData.referer_distribution).length === 0) &&
                             (!analysisData.browser_distribution || Object.keys(analysisData.browser_distribution).length === 0) &&
                             (!analysisData.os_distribution || Object.keys(analysisData.os_distribution).length === 0) &&
                             (!analysisData.device_distribution || Object.keys(analysisData.device_distribution).length === 0) &&
                             (!analysisData.bot_distribution || Object.keys(analysisData.bot_distribution).length === 0) && (
                                <Alert severity="info" sx={{ mt: 3 }}>
                                    No distribution data available for the selected period.
                                </Alert>
                            )}
                        </>
                    )}

                    {/* Initial state message */}
                    {!analysisData && !error && (
                        <Alert severity="info" sx={{ mt: 3 }}>
                            Select a month and click "Analyze" to view click statistics and distribution charts.
                        </Alert>
                    )}
                </Paper>
            </Container>
        </>
    );
};

export default Analysis;
