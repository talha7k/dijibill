#!/bin/bash

# Development script that runs wails dev and automatically fixes models
# Usage: ./dev.sh

echo "🚀 Starting Dijibill development server..."
echo "📝 This will automatically fix TypeScript models when they're generated"
echo ""

# Make sure fix-models.sh is executable
chmod +x fix-models.sh

# Function to run fix-models in the background
fix_models_loop() {
    while true; do
        sleep 5  # Check every 5 seconds
        if [ -f "frontend/wailsjs/go/models.ts" ]; then
            ./fix-models.sh > /dev/null 2>&1
        fi
    done
}

# Start the fix-models loop in the background
fix_models_loop &
FIX_PID=$!

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping development server..."
    kill $FIX_PID 2>/dev/null
    exit 0
}

# Set up trap to cleanup on exit
trap cleanup SIGINT SIGTERM

# Run wails dev
echo "🔧 Running: wails dev"
wails dev

# Cleanup when wails dev exits
cleanup