#!/usr/bin/env python3
import os
import subprocess
import sys
import datetime

# Iron Box Snapshot Tool
# Purpose: Create a safety restore point in .git_safety before a mutation.

SAFETY_DIR = ".git_safety"
WORK_TREE = "."

def run_git_safety(args):
    cmd = ["git", "--git-dir=" + SAFETY_DIR, "--work-tree=" + WORK_TREE] + args
    result = subprocess.run(cmd, capture_output=True, text=True)
    return result

def ensure_safety_repo():
    if not os.path.exists(SAFETY_DIR):
        print(f"🛡️ Initializing Safety Repo: {SAFETY_DIR}")
        os.makedirs(SAFETY_DIR, exist_ok=True)
        subprocess.run(["git", "init", "--bare", SAFETY_DIR], check=True)
        # Add to exclude if not already there
        exclude_path = ".git/info/exclude"
        if os.path.exists(".git"):
            os.makedirs(".git/info", exist_ok=True)
            with open(exclude_path, "a+") as f:
                f.seek(0)
                if SAFETY_DIR not in f.read():
                    f.write(f"\n{SAFETY_DIR}/\n")

def create_snapshot(label):
    ensure_safety_repo()
    
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    message = f"[{timestamp}] SNAPSHOT: {label}"
    
    print(f"📸 Creating Snapshot: {label}...")
    
    # Add all files to the safety index
    run_git_safety(["add", "."])
    
    # Commit
    res = run_git_safety(["commit", "--allow-empty", "-m", message])
    
    if res.returncode == 0:
        print(f"✅ Snapshot saved: {message}")
        return True
    else:
        print(f"❌ Snapshot failed: {res.stderr}")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 2:
        label = "Manual Snapshot"
    else:
        label = " ".join(sys.argv[1:])
    
    if create_snapshot(label):
        sys.exit(0)
    else:
        sys.exit(1)
