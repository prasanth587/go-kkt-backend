# Running the Backend Locally

This guide will help you run the Go backend locally on your machine.

## Prerequisites

1. **Go** (Version 1.23.4 or compatible)
   - Check: `go version`
   - Install: https://golang.org/dl/

2. **MySQL Database**
   - Make sure MySQL is installed and running
   - You need a database called `transport_hub` (or create it)
   - Default connection: `localhost:3306`

## Setup Steps

### 1. Update Database Connection

Edit the `build.sh` file and update the database connection string:

```bash
export DB_CONNECTION_CONN="root:YOUR_PASSWORD@tcp(localhost:3306)/transport_hub"
```

Replace:
- `YOUR_PASSWORD` with your MySQL root password
- `transport_hub` with your database name (if different)
- `localhost:3306` with your MySQL host/port (if different)

### 2. Make the Script Executable

```bash
chmod +x build.sh
```

### 3. Create Required Directories

The script will automatically create the upload directories, but you can verify:

```bash
mkdir -p ../uploads/t_hub_document/employee
```

### 4. Run the Application

**Option A: Using the build script (Recommended)**
```bash
cd go-kkt-backend
./build.sh
```

**Option B: Run directly with environment variables**
```bash
cd go-kkt-backend

# Set environment variables
export DB_CONNECTION_CONN="root:YOUR_PASSWORD@tcp(localhost:3306)/transport_hub"
export BASE_DIRECTORY="$(pwd)/../uploads"
export IMAGE_DIRECTORY="/t_hub_document"
export EMPLOYEE_IMAGE_DIRECTORY="/t_hub_document/employee"

# Build
go build -o go-transport-hub

# Run
./go-transport-hub
```

### 5. Verify it's Running

The server should start on port **9005**. You should see:
```
Listening on: 9005
Connected to database successfully!
```

Test the API:
```bash
curl http://localhost:9005/v1/ping
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_CONNECTION_CONN` | MySQL connection string | `root:password@tcp(localhost:3306)/transport_hub` |
| `BASE_DIRECTORY` | Base directory for file uploads | `/path/to/uploads` |
| `IMAGE_DIRECTORY` | Image storage directory | `/t_hub_document` |
| `EMPLOYEE_IMAGE_DIRECTORY` | Employee images directory | `/t_hub_document/employee` |

## Troubleshooting

### Database Connection Error
- Make sure MySQL is running: `sudo systemctl status mysql` (Linux) or `brew services list` (Mac)
- Verify credentials in the connection string
- Check if the database exists: `mysql -u root -p -e "SHOW DATABASES;"`
- Create database if needed: `mysql -u root -p -e "CREATE DATABASE transport_hub;"`

### Port Already in Use
If port 9005 is already in use:
```bash
# Find what's using the port
sudo lsof -i :9005

# Or change the port in main.go (line 51)
Addr: fmt.Sprintf(":%d", YOUR_PORT),
```

### Go Not Found
Make sure Go is in your PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
# Add this to your ~/.bashrc or ~/.zshrc for persistence
```

### Permission Denied
Make the build script executable:
```bash
chmod +x build.sh
```

## Development vs Production

**Current Setup (Local Development):**
- Database: `localhost:3306`
- Base Directory: `../uploads` (relative to project)
- Port: `9005`

**Production Setup:**
- Database: Remote server
- Base Directory: `/home/ubuntu/t_hub_document`
- Port: `9005` (usually behind a reverse proxy)

## Frontend Connection

If you're also running the frontend locally, make sure it points to:
```
http://localhost:9005
```

The frontend config is in: `hub/frontend/src/api/apiUrls.js`

