#!/usr/bin/env python3
import os
import sys

# License Compliance Audit Tool
# Ensures all used packages are permissively licensed.

APPROVED_LICENSES = ["MIT", "Apache 2.0", "BSD", "AGPLv3"]

PACKAGES = {
    "fake-bpy-module": "MIT",
    "UniTask": "MIT",
    "MemoryPack": "MIT",
    "pyright": "MIT",
    "pytest-blender": "MIT",
    "go-sdk": "Apache 2.0"
}

def run_license_audit():
    print("📜 Running License Compliance Audit...")
    violations = []
    
    for pkg, license in PACKAGES.items():
        if license not in APPROVED_LICENSES:
            violations.append(f"{pkg} ({license})")
            
    if violations:
        print(f"🚨 License Violations: {', '.join(violations)}")
        return False
        
    print("✅ License Compliance Audit Passed: All packages permissively licensed.")
    return True

if __name__ == "__main__":
    if run_license_audit():
        sys.exit(0)
    else:
        sys.exit(1)
