#!/usr/bin/env python3
import subprocess
import sys
import os

# Blender API Intelligence Installer
# Installs fake-bpy-module for static analysis and pyright.

def run_cmd(cmd):
    print(f"Executing: {' '.join(cmd)}")
    subprocess.check_call(cmd)

def setup_blender_api_intel():
    print("🧊 Setting up Blender API Intelligence...")
    
    # Install fake-bpy-module for the current blender version (default 3.6 for v0.4)
    run_cmd([sys.executable, "-m", "pip", "install", "fake-bpy-module-3.6", "pyright", "pytest", "pytest-blender"])
    
    print("✅ Blender API Intelligence setup complete.")

if __name__ == "__main__":
    try:
        setup_blender_api_intel()
    except Exception as e:
        print(f"❌ Setup failed: {e}")
        sys.exit(1)
