# PikaLink ⚡️ - Lightweight URL Shortener

A full-stack URL shortener application built with Go backend and React frontend.

> [!NOTE]
> This project is actively iterating. Current release is ready for use with basic functionalities, but the API / DB schema might have siginificant change in the near future. It's not recommended to use it in production business yet.


## 📚 Features

- Lightweight app
   - Built with Go for great performance and easy deployment
   - SQLite database

- Admin interface with Material-UI
   - CRUD for links
   - click tracking


## 🔍 Demo site

// Demo site is WIP.

> Admin URL: `YOUR_DOMAIN/admin/`
>
> Default credential: user `admin`, pwd `admin123`


## 🤖 Selfhost 

We provide prebuilt docker image, available at [bytetopia/pikalink](https://hub.docker.com/r/bytetopia/pikalink)

```bash
# Using Docker Hub
docker run -d \
  --name pikalink \
  -p 8080:8080 \
  -v pikalink_data:/app/data \
  --restart unless-stopped \
  bytetopia/pikalink:latest
```

You can also use docker-compose for quick deployment, please refer to [docker-compose.yml](./blob/main/docker-compose.yml) for details.

```bash
docker-compose up -d
```


## 🔧 Build

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

```
dist/
├── pikalink          # or pikalink.exe on Windows
├── frontend/          # Built React app
│   ├── index.html
│   ├── static/
│   └── ...
└── pikalink.db       # SQLite database (created on first run)
```

Copy the `dist` directory to your server and run the executable, application will start on port 8080 by default.

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

The frontend development server will start on `http://localhost:3000/admin/`


### Logs

The application logs to stdout. In production, redirect to a file:
```bash
./pikalink > pikalink.log 2>&1
```

## License

This project is licensed under the MIT License.