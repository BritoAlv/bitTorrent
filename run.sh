#!/bin/bash
set -e  # Exit on any error

# Run build script
bash build.sh || { echo "Build failed! Exiting..."; exit 1; }

# Navigate to bin directory
cd "$(dirname "$0")/bin" || { echo "Failed to enter bin directory"; exit 1; }

# Start server and log output
./server > api_log.txt 2>&1 & 
echo "Server started with PID $!"

# Navigate to frontend
cd frontend || { echo "Failed to enter bin/frontend"; exit 1; }

# Start frontend host script in the background
python3 host.py & 
echo $! > frontend_pid.txt
echo "Frontend started with PID $(cat frontend_pid.txt)"
