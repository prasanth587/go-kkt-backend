#!/usr/bin/env bash

# Add Go to PATH
export PATH=$PATH:/usr/local/go/bin

export APP_NAME="kk-transport"
export ENV="dev"
#export DB_CONNECTION_CONN="root:3nRrF3rn5f@tcp(3.88.53.87:3306)/transport_hub"
export DB_CONNECTION_CONN="root:YourNewSecurePassword123!@tcp(localhost:3306)/transport_hub"
# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export EMPLOYEE_IMAGE_DIRECTORY="/t_hub_document/employee"
export BASE_DIRECTORY="${SCRIPT_DIR}/../uploads"
export IMAGE_DIRECTORY="/t_hub_document"

# Create upload directories if they don't exist
mkdir -p "${BASE_DIRECTORY}/t_hub_document/employee"
export LOG_LEVEL="Info"


#export PORT=9010
  
# Sensitive environment variables are added in dev.env
# source dev.env

#swag init
# Swag init done

go install
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
 

echo "Doing some cleaning ..."
go clean
echo "Done."

# Check if required commands are available
command -v go >/dev/null 2>&1 || { echo >&2 "Go is required but it's not installed. Aborting."; exit 1; }
# command -v godep >/dev/null 2>&1 || { echo >&2 "Godep is required but it's not installed. Aborting."; exit 1; }
# command -v goimports >/dev/null 2>&1 || { echo >&2 "Goimports is required but it's not installed. Aborting."; exit 1; }
command -v gofmt >/dev/null 2>&1 || { echo >&2 "Gofmt is required but it's not installed. Aborting."; exit 1; }


# echo "Running go vet ..."
# go vet ./internal/...
# if [ $? != 0 ]; then
#   exit
# fi
# echo "Done."

#echo "Running go generate ..."
#go generate ./internal/...
#echo "Done."

echo "Running go format ..."
gofmt -w .
echo "Done."

echo "Running go build ..."
go build -race
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
echo "Done."

#echo "Running unit test ..."
#go test -parallel 1 ./internal/...
#export PG_DATABASE_URL='ER9OaxWlwxYday7oZ7-Wecnq9HNvAHy4h8BW-0uShA3NcMajFtBMahDeO-XI_y92eaDpjH-bt9nItiLfIDsfFyzgLjxZfGn2qbA3WBd2PcztpGMtdCf6QNWbFp-glIY8f-tMVLGP-Gpl4LIue_pH-nh5QO-69eKmp2ORbB4OY_9VZ8tiZAexTGd3'

if [ $? == 0 ]; then
	echo "## Starting service ##"
    ./go-transport-hub
fi