#!/bin/bash

# Script to automatically fix the "classs" typo in generated Wails models
# Run this after wails dev or wails build to fix the TypeScript generation bug

MODELS_FILE="frontend/wailsjs/go/models.ts"

echo "🔧 Running fix-models.sh..."

if [ -f "$MODELS_FILE" ]; then
    echo "📁 Found models file: $MODELS_FILE"
    
    # Check if the file contains the typo before attempting to fix
    if grep -q "classs" "$MODELS_FILE"; then
        echo "🐛 Found 'classs' typo, fixing..."
        
        # Replace 'classs' with 'clazz'
        sed -i '' 's/classs/clazz/g' "$MODELS_FILE"
        
        # Add proper null checks for constructor calls
        sed -i '' 's/return new clazz(a);/return clazz \&\& typeof clazz === '\''function'\'' ? new clazz(a) : a;/g' "$MODELS_FILE"
        sed -i '' 's/a\[key\] = new clazz(a\[key\]);/a[key] = clazz \&\& typeof clazz === '\''function'\'' ? new clazz(a[key]) : a[key];/g' "$MODELS_FILE"
        sed -i '' 's/return new clazz(a);/return clazz \&\& typeof clazz === '\''function'\'' ? new clazz(a) : a;/g' "$MODELS_FILE"
        
        echo "✅ Fixed TypeScript models successfully!"
    else
        echo "✅ No 'classs' typo found, models are already correct!"
    fi
else
    echo "⚠️  Models file not found: $MODELS_FILE"
    echo "   This is normal if models haven't been generated yet."
fi

echo "🏁 fix-models.sh completed"