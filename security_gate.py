#!/usr/bin/env python3
import os
import sys
import re

# VibeSync: Security Gate Auditor (v1.0)
# This script scans the codebase for prohibited patterns and security violations.

PROHIBITED_PATTERNS = [
    (r"os\.system", "DANGEROUS: Shell execution detected."),
    (r"eval\(", "DANGEROUS: Dynamic execution detected."),
    (r"exec\(", "DANGEROUS: Dynamic execution detected."),
    (r"subprocess\.", "DANGEROUS: Process spawning detected."),
    (r"Reflection", "DANGEROUS: C# Reflection detected."),
    (r"UnityEditorInternal", "DANGEROUS: Accessing internal Unity APIs."),
    (r"System\.IO\.File\.(?!ReadAllText|WriteAllText)", "DANGEROUS: Raw file system access outside sanctioned APIs."),
    (r"http://(?!localhost|127\.0\.0\.1)", "DANGEROUS: Non-localhost network access."),
    (r"https://(?!localhost|127\.0\.0\.1)", "DANGEROUS: Non-localhost network access."),
    (r"api_key\s*=\s*['\"][a-zA-Z0-9]{20,}['\"]", "DANGEROUS: Potential hardcoded API key."),
    (r"token\s*=\s*['\"][a-zA-Z0-9]{20,}['\"]", "DANGEROUS: Potential hardcoded secret token."),
]

def audit_file(file_path):
    violations = []
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            for i, line in enumerate(f, 1):
                if "skip-security-gate" in line:
                    continue
                for pattern, message in PROHIBITED_PATTERNS:
                    if re.search(pattern, line):
                        violations.append(f"Line {i}: {message}\n  -> {line.strip()}")
    except Exception as e:
        return [f"ERROR: Could not read {file_path}: {e}"]
    return violations

def main():
    print("🛡️ VibeSync Security Gate: Initiating Audit...")
    total_violations = 0
    
    # Directories to scan
    scan_dirs = ["mcp-server", "unity-bridge", "blender-bridge"]
    
    for scan_dir in scan_dirs:
        if not os.path.exists(scan_dir):
            continue
            
        print(f"Scanning {scan_dir}...")
        for root, _, files in os.walk(scan_dir):
            for file in files:
                if file.endswith((".go", ".cs", ".py")):
                    file_path = os.path.join(root, file)
                    violations = audit_file(file_path)
                    if violations:
                        print(f"\n🚨 VIOLATIONS FOUND in {file_path}:")
                        for v in violations:
                            print(v)
                        total_violations += len(violations)

    if total_violations > 0:
        print(f"\n❌ AUDIT FAILED: {total_violations} security violations detected.")
        sys.exit(1)
    else:
        print("\n✅ AUDIT PASSED: No security violations detected.")
        sys.exit(0)

if __name__ == "__main__":
    main()
