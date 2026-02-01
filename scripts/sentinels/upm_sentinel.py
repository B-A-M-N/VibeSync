#!/usr/bin/env python3
import json
import os
import sys

# UPM (Unity Package Manager) Sentinel
# Audits manifest.json for required modules.

REQUIRED_PACKAGES = [
    "com.cysharp.memorypack",
    "com.cysharp.unitask"
]

def audit_upm():
    manifest_path = "unity-bridge/manifest.json" # Relative to workspace root for sync
    if not os.path.exists(manifest_path):
        # Check standard location
        manifest_path = "Packages/manifest.json"
        
    if not os.path.exists(manifest_path):
        print("⚠️ VibeSync: manifest.json not found. Skipping UPM audit.")
        return True

    print(f"🔍 Auditing UPM Manifest: {manifest_path}")
    
    with open(manifest_path, "r") as f:
        manifest = json.load(f)
        
    dependencies = manifest.get("dependencies", {})
    missing = []
    
    for pkg in REQUIRED_PACKAGES:
        if pkg not in dependencies:
            missing.append(pkg)
            
    if missing:
        print(f"🚨 VibeSync UPM Violation: Missing required packages: {', '.join(missing)}")
        return False
        
    print("✅ UPM Manifest Audit Passed.")
    return True

if __name__ == "__main__":
    if audit_upm():
        sys.exit(0)
    else:
        sys.exit(1)
