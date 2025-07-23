// Configuration for different environments
const config = {
  // Default to current window location in production, localhost in development
  API_BASE_URL: process.env.REACT_APP_API_URL || 
    (process.env.NODE_ENV === 'production' 
      ? `${window.location.protocol}//${window.location.host}/api`
      : 'http://localhost:8080/api'),
  
  // Base URL for short links
  SHORT_LINK_BASE_URL: process.env.REACT_APP_SHORT_LINK_URL ||
    (process.env.NODE_ENV === 'production'
      ? `${window.location.protocol}//${window.location.host}`
      : 'http://localhost:8080')
};

export default config;
