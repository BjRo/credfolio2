#!/bin/bash
# Stop all dev server processes
#
# Usage: .claude/scripts/stop-dev.sh
#
# This script safely stops all development server processes:
# - Turbo orchestrator (turbo run dev)
# - Go backend (go run cmd/server)
# - Next.js frontend (next dev)
# - Any processes on ports 8080 and 3000
#
# The script is idempotent - safe to run even if no servers are running.

# Colors for output (matching reset-db.sh pattern)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Stopping dev servers...${NC}"
echo ""

# Step 1: Kill turbo orchestrator
echo -e "${GREEN}[1/6]${NC} Stopping Turbo orchestrator..."
if pkill -f "turbo run dev" 2>/dev/null; then
    echo "      Killed Turbo orchestrator"
else
    echo "      No Turbo orchestrator running"
fi

# Step 2: Kill Go backend
echo -e "${GREEN}[2/6]${NC} Stopping Go backend..."
if pkill -f "go run cmd/server" 2>/dev/null; then
    echo "      Killed Go backend"
else
    echo "      No Go backend running"
fi

# Step 3: Kill Next.js frontend
echo -e "${GREEN}[3/6]${NC} Stopping Next.js frontend..."
if pkill -f "next dev" 2>/dev/null; then
    echo "      Killed Next.js frontend"
else
    echo "      No Next.js frontend running"
fi

# Step 4: Kill anything on ports 8080 and 3000
echo -e "${GREEN}[4/6]${NC} Killing processes on ports 8080 and 3000..."
if fuser -k 8080/tcp 3000/tcp 2>/dev/null; then
    echo "      Killed processes on ports 8080 and 3000"
else
    echo "      No processes on ports 8080 and 3000"
fi

# Step 5: Wait for processes to terminate
echo -e "${GREEN}[5/6]${NC} Waiting for graceful shutdown..."
sleep 2
echo "      Done"

# Step 6: Verify ports are free
echo -e "${GREEN}[6/6]${NC} Verifying ports are free..."
if lsof -i :8080 -i :3000 2>/dev/null; then
    echo ""
    echo -e "${RED}ERROR: Processes still running on ports 8080 or 3000${NC}"
    echo -e "${YELLOW}Run 'lsof -i :8080 -i :3000' for details${NC}"
    exit 1
else
    echo "      Ports 8080 and 3000 are free"
fi

echo ""
echo -e "${GREEN}All dev servers stopped successfully!${NC}"
