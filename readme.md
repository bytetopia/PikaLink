# PikaLink - URL Shortener

> [!NOTE]
> 🤖 **Created by GitHub Copilot** ✨  
> This awesome project was generated and designed by your friendly AI programming assistant! 🚀

A full-stack URL shortener application built with Go backend and React frontend.

## Features

- **Backend (Go)**:
  - RESTful API with Gin framework
  - SQLite database for data storage
  - JWT authentication
  - CRUD operations for links
  - Click tracking
  - URL redirection at root path
  - Admin interface serving

- **Frontend (React)**:
  - Material-UI components
  - User authentication
  - Link management (create, edit, delete, list)
  - Responsive design
  - Copy to clipboard functionality
  - Admin interface accessible at `/admin`

## Project Structure

```
pikalink/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── models/
│   ├── handlers/
│   ├── database/
│   └── middleware/
├── frontend/
│   ├── package.json
│   ├── src/
│   │   ├── components/
│   │   ├── App.jsx
│   │   └── index.js
│   └── public/
├── build.sh          # Build script for Linux/macOS
├── build.ps1         # Build script for Windows
└── README.md
```

## URL Structure

The application uses a clean URL structure:

- **Short URL redirects**: `https://yourdomain.com/abc` → redirects to the target URL
- **Admin interface**: `https://yourdomain.com/admin` → management portal
- **API endpoints**: `https://yourdomain.com/api/*` → backend API

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js 16 or later
- npm or yarn

### Development Setup

#### Backend Setup

1. Navigate to the backend directory:
   ```powershell
   cd backend
   ```

2. Install dependencies:
   ```powershell
   go mod tidy
   ```

3. Run the server:
   ```powershell
   go run main.go
   ```

The backend will start on `http://localhost:8080`

#### Frontend Setup

1. Navigate to the frontend directory:
   ```powershell
   cd frontend
   ```

2. Install dependencies:
   ```powershell
   npm install
   ```

3. Start the development server:
   ```powershell
   npm start
   ```

The frontend development server will start on `http://localhost:3000/admin`

## Production Build & Deployment

### Quick Build (Automated)

Use the provided build scripts for easy production builds:

**Windows (PowerShell):**
```powershell
.\build.ps1
```

**Linux/macOS (Bash):**
```bash
chmod +x build.sh
./build.sh
```

This will create a `dist` directory with all production files ready for deployment.

### Manual Build Process

If you prefer to build manually or need custom configurations:

#### Step 1: Build Frontend
```powershell
cd frontend
npm install
npm run build
```

#### Step 2: Prepare Backend with Frontend Assets
```powershell
cd ../backend
# Remove any existing frontend folder
rm -rf frontend  # or Remove-Item frontend -Recurse -Force on Windows
mkdir frontend
# Copy built frontend to backend
cp -r ../frontend/build/* frontend/  # or Copy-Item on Windows
```

#### Step 3: Build Go Binary
```powershell
go mod tidy
go build -o pikalink main.go  # On Windows: go build -o pikalink.exe main.go
```

#### Step 4: Prepare Distribution
```powershell
# Create distribution directory
mkdir ../dist
# Copy binary
cp pikalink ../dist/  # or Copy-Item on Windows
# Copy frontend assets
cp -r frontend ../dist/  # or Copy-Item on Windows
# Copy database (if exists)
cp pikalink.db ../dist/  # Optional, will be created on first run
```

### Cross-Platform Builds

To build for different platforms:

```powershell
# Build for Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o pikalink-linux main.go

# Build for macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o pikalink-macos main.go

# Build for Windows (from Linux/macOS)
GOOS=windows GOARCH=amd64 go build -o pikalink.exe main.go
```

### Deployment Structure

After building, your `dist` directory should contain:

```
dist/
├── pikalink          # or pikalink.exe on Windows
├── frontend/          # Built React app
│   ├── index.html
│   ├── static/
│   └── ...
└── pikalink.db       # SQLite database (created on first run)
```

### Running in Production

1. Copy the `dist` directory to your server
2. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x pikalink
   ```
3. Run the application:
   ```bash
   ./pikalink          # Linux/macOS
   pikalink.exe        # Windows
   ```

The application will start on port 8080 by default.

### Production Environment Variables

You can customize the production deployment with environment variables:

```bash
export PORT=8080                    # Server port (default: 8080)
export GIN_MODE=release            # Set Gin to release mode
export DATA_PATH=./data              # Data directory path (for database and logs)
```

### Systemd Service (Linux)

Create a systemd service for automatic startup:

```ini
# /etc/systemd/system/pikalink.service
[Unit]
Description=PikaLink URL Shortener
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/pikalink
ExecStart=/opt/pikalink/pikalink
Restart=always
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl enable pikalink
sudo systemctl start pikalink
```

### Reverse Proxy Setup (Nginx)

Example Nginx configuration:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Default Credentials

- **Username**: admin
- **Password**: admin123

## API Endpoints

### Authentication
- `POST /api/login` - User login

### Link Management (Protected)
- `GET /api/links` - Get user's links
- `POST /api/links` - Create new link
- `PUT /api/links/:id` - Update link
- `DELETE /api/links/:id` - Delete link
- `POST /api/change-password` - Change user password

### Public Endpoints
- `GET /:code` - Redirect to original URL (short link)
- `GET /admin/*` - Admin interface routes

## Database Schema

### Users Table
- `id` - Primary key
- `username` - Unique username
- `password` - Hashed password

### Links Table
- `id` - Primary key
- `original_url` - Original URL
- `short_code` - Generated short code
- `title` - Optional title
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `user_id` - Foreign key to users
- `click_count` - Number of clicks

## Usage Examples

### Creating Short Links
1. Access the admin interface at `http://yourdomain.com/admin`
2. Log in with your credentials
3. Click "Create New Link"
4. Enter the original URL and optional title
5. Choose between auto-generated or custom short code
6. The short link will be available at `http://yourdomain.com/{shortcode}`

### Custom Short Codes
- Must be 3-50 characters long
- Only alphanumeric characters, hyphens, and underscores allowed
- Must be unique

## Technologies Used

### Backend
- Go
- Gin (HTTP framework)
- SQLite
- JWT for authentication
- bcrypt for password hashing
- CORS middleware

### Frontend
- React
- Material-UI
- React Router
- Axios for API calls

## Development Notes

- The backend serves the built React app in production
- CORS is configured for `http://localhost:3000` in development
- Frontend homepage is set to `/admin` for correct asset loading
- Short URL redirects have priority over other routes
- SQLite database is created automatically on first run
- All frontend assets are embedded in the Go binary deployment

## Troubleshooting

### Common Issues

1. **Frontend not loading**: Ensure the `frontend` directory exists in the same location as the binary
2. **Database errors**: Check write permissions for the database file
3. **Port already in use**: Change the port using environment variables
4. **CORS issues in development**: Ensure the frontend dev server is running on port 3000

### Logs

The application logs to stdout. In production, redirect to a file:
```bash
./pikalink > pikalink.log 2>&1
```

## License

This project is licensed under the MIT License.