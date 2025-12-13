#!/bin/bash

# Development script to run the backend with Railway database connection and log to file

# Load environment variables from .env.local if it exists
if [ -f .env.local ]; then
    source .env.local
fi

# Check if DB_CONNECTION_CONN is set
if [ -z "$DB_CONNECTION_CONN" ]; then
    echo "Error: DB_CONNECTION_CONN environment variable is not set!"
    echo "Please set it or create a .env.local file with:"
    echo 'export DB_CONNECTION_CONN="your_connection_string"'
    exit 1
fi

# Set BASE_DIRECTORY and IMAGE_DIRECTORY for local development
# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set environment variables if not already set
if [ -z "$BASE_DIRECTORY" ]; then
    export BASE_DIRECTORY="${SCRIPT_DIR}/../uploads"
    echo "Setting BASE_DIRECTORY to: $BASE_DIRECTORY"
fi

if [ -z "$IMAGE_DIRECTORY" ]; then
    export IMAGE_DIRECTORY="/t_hub_document"
    echo "Setting IMAGE_DIRECTORY to: $IMAGE_DIRECTORY"
fi

if [ -z "$EMPLOYEE_IMAGE_DIRECTORY" ]; then
    export EMPLOYEE_IMAGE_DIRECTORY="/t_hub_document/employee"
    echo "Setting EMPLOYEE_IMAGE_DIRECTORY to: $EMPLOYEE_IMAGE_DIRECTORY"
fi

# Twilio SMS Configuration (optional - set these if you want SMS functionality)
# Get these from your Twilio Console: https://console.twilio.com/
# Set these as environment variables or update with your actual credentials
export TWILIO_ACCOUNT_SID="${TWILIO_ACCOUNT_SID:-your_account_sid_here}"
export TWILIO_AUTH_TOKEN="${TWILIO_AUTH_TOKEN:-your_auth_token_here}"
export TWILIO_PHONE_NUMBER="${TWILIO_PHONE_NUMBER:-your_phone_number_here}"  # Your Twilio phone number (with country code)

# Create upload directories if they don't exist
mkdir -p "${BASE_DIRECTORY}/t_hub_document/employee"
echo "Created upload directories at: ${BASE_DIRECTORY}/t_hub_document"

# Create logs directory if it doesn't exist
LOG_DIR="${SCRIPT_DIR}/logs"
mkdir -p "$LOG_DIR"

# Log file with timestamp
LOG_FILE="${LOG_DIR}/backend-$(date +%Y%m%d-%H%M%S).log"
echo "Logs will be written to: $LOG_FILE"
echo "You can also view logs in real-time with: tail -f $LOG_FILE"

echo ""
echo "Starting backend server..."
echo "Database: Railway (maglev.proxy.rlwy.net:26072)"
echo "Server will run on: http://localhost:9005"
echo "BASE_DIRECTORY: $BASE_DIRECTORY"
echo "IMAGE_DIRECTORY: $IMAGE_DIRECTORY"
echo "Log file: $LOG_FILE"
echo ""

# Run the Go application and redirect both stdout and stderr to log file
# Also display on terminal using tee
go run main.go 2>&1 | tee "$LOG_FILE"

