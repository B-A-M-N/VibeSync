#!/usr/bin/env python3
import os
import sys

# Log Rotation & Context Truncation Tool
# Keeps logs manageable for AI context.

MAX_LOG_SIZE = 1024 * 1024 # 1MB

def rotate_log(file_path):
    if not os.path.exists(file_path):
        return
        
    size = os.path.getsize(file_path)
    if size > MAX_LOG_SIZE:
        print(f"🔄 Rotating Log: {file_path} ({size} bytes)")
        backup = file_path + ".old"
        if os.path.exists(backup):
            os.remove(backup)
        os.rename(file_path, backup)
        # Keep last 100 lines in new file
        with open(backup, "r") as f:
            lines = f.readlines()
        with open(file_path, "w") as f:
            f.writelines(lines[-100:])

def run_rotation():
    logs_to_rotate = [
        "mcp-server/vibe_server.log",
        "blender_bridge.log",
        ".vibesync/events.jsonl",
        ".vibesync/wal.jsonl"
    ]
    
    for log in logs_to_rotate:
        rotate_log(log)

if __name__ == "__main__":
    run_rotation()
    print("✅ Log rotation complete.")
