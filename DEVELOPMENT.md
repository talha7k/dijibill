# Development Workflow

## Quick Start

### Option 1: Use the integrated development script (Recommended)
```bash
./dev.sh
```
This script will:
- Start the Wails development server
- Automatically monitor and fix TypeScript model generation issues
- Handle cleanup when you stop the server

### Option 2: Manual development
```bash
wails dev
```
If you encounter TypeScript model issues, run:
```bash
./fix-models.sh
```

## Build Process

When building for production:
```bash
wails build
```
The `fix-models.sh` script will automatically run after the frontend build completes (via npm postbuild hook).

## Scripts Explanation

### `dev.sh`
- Integrated development script that combines `wails dev` with automatic model fixing
- Monitors the generated models file and fixes issues automatically
- Provides better development experience with visual feedback

### `fix-models.sh`
- Fixes the "classs" typo in generated Wails TypeScript models
- Adds proper null checks for constructor calls
- Can be run manually or automatically via hooks
- Provides detailed feedback about what was fixed

## Integration Details

The fix-models script is integrated into the development workflow in multiple ways:

1. **Frontend postbuild hook**: Automatically runs after `npm run build` in the frontend
2. **Development script**: The `dev.sh` script monitors and fixes models automatically
3. **Manual execution**: Can be run manually anytime with `./fix-models.sh`

This ensures that TypeScript model issues are automatically resolved regardless of how you run the development server.