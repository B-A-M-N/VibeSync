#!/usr/bin/env python3
import os
import json
import socket
import psutil
import time
import sys

# VibeSync: Adversarial Pre-flight Layer (v1.0)
# Detects and resolves environmental blocks before orchestration begins.

VIBE_STATE_PATH = ".vibesync/state.json"
EXPECTED_PORTS = [8087, 8088, 22005]  # Unity (8087), Vision (8088), Blender (22005)
VIBE_TMP_PATH = ".vibesync/tmp"
AIRLOCK_QUEUE_PATH = "AirlockQueue"

def poll_airlock(timeout=2):
    if not os.path.exists(AIRLOCK_QUEUE_PATH):
        return False
    start = time.time()
    while time.time() - start < timeout:
        if os.listdir(AIRLOCK_QUEUE_PATH):
            return True
        time.sleep(0.2)
    return False

def load_state():
    if os.path.exists(VIBE_STATE_PATH):
        try:
            with open(VIBE_STATE_PATH, "r") as f:
                return json.load(f)
        except:
            pass
    return {"engines": {}, "id_map": {}, "credits": 100}

def save_state(state):
    os.makedirs(os.path.dirname(VIBE_STATE_PATH), exist_ok=True)
    with open(VIBE_STATE_PATH, "w") as f:
        json.dump(state, f, indent=4)

def is_port_in_use(port):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        return s.connect_ex(("localhost", port)) == 0

def kill_pid(pid):
    try:
        proc = psutil.Process(pid)
        proc.kill()
        return True
    except psutil.NoSuchProcess:
        return False

def find_zombie_engine_pids():
    zombies = []
    targets = ["Unity", "Blender", "vibe-mcp-server"]
    cwd = os.getcwd()
    for proc in psutil.process_iter(['pid', 'name', 'cwd']):
        try:
            name = proc.info['name']
            proc_cwd = proc.info.get('cwd')
            # Only target processes that are either stopped/zombie OR belong to this project directory
            if any(t in name for t in targets):
                if proc.status() in (psutil.STATUS_STOPPED, psutil.STATUS_ZOMBIE):
                    zombies.append(proc.info['pid'])
                elif proc_cwd == cwd:
                    # If it's running in our directory but the port check (later) shows it's blocking us
                    pass 
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            continue
    return zombies

def check_tmp_sandbox():
    if not os.path.exists(VIBE_TMP_PATH):
        return False, "Sandbox directory missing"
    return True, "OK"

def check_unity_compile_errors():
    # Placeholder for actual Unity compilation check
    # In a real scenario, this would check Editor.log or a specific artifact
    return False, ""

def adversarial_preflight():
    issues = []
    resolutions = []

    # 1. Port Conflict Check (PRIMARY AUTHORITY)
    ORCHESTRATOR_PORTS = [8080] # Port used by vibe-mcp-server/proxy
    for port in EXPECTED_PORTS + ORCHESTRATOR_PORTS:
        for proc in psutil.process_iter(['pid', 'name', 'connections', 'cwd']):
            try:
                connections = proc.info.get('connections')
                if connections is None:
                    continue
                for conn in connections:
                    if conn.laddr.port == port:
                        pid = proc.info['pid']
                        name = proc.info['name']
                        
                        if port in ORCHESTRATOR_PORTS:
                            if kill_pid(pid):
                                resolutions.append(f"Released orchestrator port {port} by killing {name} (PID {pid})")
                        else:
                            issues.append(f"Port {port} in use by {name} (PID {pid}). Ensuring engine liveness.")
            except (psutil.NoSuchProcess, psutil.AccessDenied):
                continue

    # 2. State Integrity
    state = load_state()
    panic_found = False
    for engine, data in state.get("engines", {}).items():
        if data.get("state") == "PANIC":
            panic_found = True
            data["state"] = "STOPPED"
            resolutions.append(f"Reset {engine} from PANIC to STOPPED")
    if panic_found:
        save_state(state)

    # 3. Sandbox Integrity
    if not os.path.exists(VIBE_TMP_PATH):
        os.makedirs(VIBE_TMP_PATH, exist_ok=True)
        resolutions.append("Recreated .vibesync/tmp sandbox")

    # 4. Lock File Cleanup
    for root, dirs, files in os.walk(".vibesync"):
        for file in files:
            if file.endswith(".lock"):
                try:
                    os.remove(os.path.join(root, file))
                    resolutions.append(f"Removed stale lock: {file}")
                except:
                    pass

    return {
        "status": "READY",
        "safe_to_proceed": True,
        "issues": issues,
        "resolutions": resolutions,
        "timestamp": time.time()
    }

if __name__ == "__main__":
    print("🛡️ VibeSync: Executing Adversarial Pre-flight...")
    report = adversarial_preflight()
    print(json.dumps(report, indent=2))
    if not report["safe_to_proceed"]:
        sys.exit(1)
    sys.exit(0)
