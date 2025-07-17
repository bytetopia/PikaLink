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
  - URL redirection

- **Frontend (React)**:
  - Material-UI components
  - User authentication
  - Link management (create, edit, delete, list)
  - Responsive design
  - Copy to clipboard functionality

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
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js 16 or later
- npm or yarn

### Backend Setup

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

### Frontend Setup

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

The frontend will start on `http://localhost:3000`

## Default Credentials

- **Username**: admin
- **Password**: admin123

## API Endpoints

- `POST /api/login` - User login
- `GET /api/links` - Get user's links
- `POST /api/links` - Create new link
- `PUT /api/links/:id` - Update link
- `DELETE /api/links/:id` - Delete link
- `GET /s/:code` - Redirect to original URL

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

## Technologies Used

### Backend
- Go
- Gin (HTTP framework)
- SQLite
- JWT for authentication
- bcrypt for password hashing

### Frontend
- React
- Material-UI
- React Router
- Axios for API calls

## License

This project is licensed under the MIT License.