#!/usr/bin/env python3
import subprocess
import sys
import os

# Pip-Blender Orchestration Script
# Purpose: Manage isolated Python dependencies for the Blender Bridge.

def run_in_blender(script_content):
    # This assumes blender is in the path. In a real environment, 
    # we would use the discovered blender path.
    temp_script = ".vibesync/tmp/pip_orchestrator.py"
    os.makedirs(".vibesync/tmp", exist_ok=True)
    with open(temp_script, "w") as f:
        f.write(script_content)
    
    cmd = ["blender", "--background", "--python", temp_script]
    print(f"Running blender pip orchestrator...")
    subprocess.check_call(cmd)

def get_orchestration_script(packages):
    pkgs_str = ", ".join([f"'{p}'" for p in packages])
    return f"""
import subprocess
import sys
import ensurepip

def install_packages():
    print("📦 VibeSync: Initializing Blender Pip...")
    ensurepip.bootstrap()
    
    packages = [{pkgs_str}]
    print(f"📦 VibeSync: Installing/Updating packages: {{packages}}")
    
    subprocess.check_call([sys.executable, "-m", "pip", "install", "--upgrade"] + packages)
    print("✅ VibeSync: Blender dependencies synchronized.")

if __name__ == "__main__":
    install_packages()
"""

def sync_blender_dependencies():
    packages = ["numpy", "memorypack", "fake-bpy-module-3.6"]
    script = get_orchestration_script(packages)
    run_in_blender(script)

if __name__ == "__main__":
    try:
        sync_blender_dependencies()
    except Exception as e:
        print(f"❌ Dependency synchronization failed: {e}")
        sys.exit(1)
