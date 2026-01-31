# VibeSync: Zero-Trust Unity ↔ Blender Orchestrator
# Copyright (C) 2026 B-A-M-N
#
# This project is distributed under a DUAL-LICENSING MODEL:
# 1. Open-Source Path: GNU Affero General Public License v3
# 2. Commercial Path: "Work-or-Pay" Model
#
# See the LICENSE file in the project root for the full terms and conditions
# of both licensing paths.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

import bpy
import json
import http.server
import threading
from queue import Queue
import time
import hmac
import hashlib

# Configuration
HOST = "127.0.0.1"
PORT = 22000
BOOTSTRAP_TOKEN = "VIBE_BLENDER_BOOTSTRAP_SECRET" # Unique token for Blender

# State
_session_token = ""
_current_generation = 0
_command_queue = Queue()
_state_lock = threading.Lock()

_path_whitelist = {
    "/health", "/handshake", "/metrics", "/object/lock", "/panic", 
    "/preflight/run", "/export", "/camera/set", "/camera/get",
    "/selection/set", "/material/update", "/mesh/mutate", "/state/get",
    "/playback/control"
}

def compute_hmac(key, data):
    return hmac.new(key.encode(), data.encode(), hashlib.sha256).hexdigest()

class VibeRequestHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        # Path Whitelist Check
        if self.path not in _path_whitelist:
            self.send_response(403)
            self.end_headers()
            return

        if self.path == "/health":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            with _state_lock:
                gen = _current_generation
            self.wfile.write(json.dumps({"status": "ok", "generation": gen}).encode())
        elif self.path == "/camera/get":
            # Simple camera telemetry (mocked for this turn)
            response = {"status": "OK", "pos": [0, 5, -10], "rot": [0, 0, 0]}
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_error(404)

    def do_POST(self):
        # 1. Path Whitelist Check
        if self.path not in _path_whitelist:
            self.send_response(403)
            self.end_headers()
            return

        # Token & Security Header Validation
        received_token = self.headers.get("X-Vibe-Token")
        received_gen = self.headers.get("X-Vibe-Generation")
        received_sig = self.headers.get("X-Vibe-Signature")
        received_time = self.headers.get("X-Vibe-Timestamp")
        
        is_handshake = self.path == "/handshake"
        
        with _state_lock:
            curr_token = _session_token if _session_token != "" else BOOTSTRAP_TOKEN
            curr_gen = _current_generation

        if received_token != curr_token:
            self.send_response(401)
            self.end_headers()
            self.wfile.write(json.dumps({"error": "Unauthorized"}).encode())
            return

        # 2. Anti-Replay: Timestamp Check (5s window)
        try:
            if received_time:
                ts = int(received_time)
                now = int(time.time())
                if abs(now - ts) > 5:
                    self.send_response(403)
                    self.end_headers()
                    self.wfile.write(json.dumps({"error": "Request Expired"}).encode())
                    return
            elif not is_handshake:
                self.send_response(400)
                self.end_headers()
                return
        except ValueError:
            self.send_response(400)
            self.end_headers()
            return

        # 3. Generation Validation
        if received_gen and not is_handshake:
            try:
                if int(received_gen) != curr_gen:
                    self.send_response(409)
                    self.end_headers()
                    self.wfile.write(json.dumps({"error": "Generation Drift"}).encode())
                    return
            except ValueError:
                self.send_response(400)
                self.end_headers()
                return

        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length).decode() if content_length > 0 else ""

        # 4. HMAC Signature Verification
        if not is_handshake:
            sig_data = f"{received_time}|{self.command}|{self.path}|{body}"
            expected_sig = compute_hmac(curr_token, sig_data)
            if received_sig != expected_sig:
                self.send_response(403)
                self.end_headers()
                self.wfile.write(json.dumps({"error": "Invalid Signature"}).encode())
                return

        if is_handshake:
            with _state_lock:
                global _current_generation, _session_token
                _current_generation += 1
                try:
                    data = json.loads(body)
                    if "new_token" in data:
                        _session_token = data["new_token"]
                        print("🛡️ VibeSync: New Session Token Established")
                except:
                    data = {}

            response = {
                "status": "OK",
                "engine_version": bpy.app.version_string,
                "capabilities": ["mesh", "transform", "cycles", "eevee", "locking", "metrics"],
                "response": "VIBE_HASH_" + data.get("challenge", "UNKNOWN")
            }
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
            return

        if self.path == "/metrics":
            # Basic metrics for now
            response = {"status": "OK", "memory_usage": 0, "engine_busy": False}
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
            return

        if self.path == "/object/lock":
            _command_queue.put((self.path, body))
            self.send_response(200)
            self.end_headers()
            self.wfile.write(json.dumps({"status": "ok"}).encode())
            return

        if self.path == "/panic":
            # Logic to lock Blender (e.g. disable operators)
            print("🚨 VIBESYNC PANIC | Signal Received")
            self.send_response(200)
            self.end_headers()
            self.wfile.write(json.dumps({"status": "locked"}).encode())
            return

        if self.path == "/preflight/run":
            # Simple hash for now
            response = {"status": "OK", "hash": "BLENDER_HASH_" + str(time.time())}
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
            return

        if self.path == "/export":
            response = {"status": "OK", "meta": {"exporter": "VibeSync"}}
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
            return

        if self.path in ["/camera/set", "/selection/set", "/material/update"]:
            _command_queue.put((self.path, body))
            self.send_response(200)
            self.end_headers()
            self.wfile.write(json.dumps({"status": "ok"}).encode())
            return

        # Queue for other commands
        _command_queue.put((self.path, body))
        self.send_response(202)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps({"status": "queued"}).encode())

    def log_message(self, format, *args):
        # Suppress logging to avoid cluttering Blender console
        pass

def run_server():
    server = http.server.HTTPServer((HOST, PORT), VibeRequestHandler)
    print(f"🛡️ VibeSync Blender Bridge: Listening on port {PORT}")
    server.serve_forever()

# Timer for main thread execution
def process_queue():
    while not _command_queue.empty():
        path, body = _command_queue.get()
        print(f"VibeSync Command received on Blender Main Thread: {path}")
        # TODO: Implement command dispatching to modules
    return 0.1 # Run every 0.1 seconds

def register():
    # Start server in background thread
    thread = threading.Thread(target=run_server, daemon=True)
    thread.start()
    
    # Register timer for queue processing
    if not bpy.app.timers.is_registered(process_queue):
        bpy.app.timers.register(process_queue)

def unregister():
    if bpy.app.timers.is_registered(process_queue):
        bpy.app.timers.unregister(process_queue)

if __name__ == "__main__":
    register()
