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

# Configuration
HOST = "localhost"
PORT = 22000
BOOTSTRAP_TOKEN = "VIBE_BLENDER_BOOTSTRAP_SECRET" # Unique token for Blender

# State
_session_token = ""
_current_generation = 0
_command_queue = Queue()
_state_lock = threading.Lock()

class VibeRequestHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        # ... (health and camera/get remain the same)
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
        # Token Validation
        received_token = self.headers.get("X-Vibe-Token")
        received_gen = self.headers.get("X-Vibe-Generation")
        
        with _state_lock:
            is_authenticated = (received_token == _session_token and _session_token != "") or \
                              (received_token == BOOTSTRAP_TOKEN)
            curr_gen = _current_generation

        if not is_authenticated:
            self.send_response(401)
            self.end_headers()
            self.wfile.write(json.dumps({"error": "Unauthorized"}).encode())
            return

        # Generation Validation
        if received_gen and self.path != "/handshake":
            try:
                if int(received_gen) != curr_gen:
                    self.send_response(409)
                    self.end_headers()
                    self.wfile.write(json.dumps({"error": "Generation Drift", "engine": curr_gen, "received": received_gen}).encode())
                    return
            except ValueError:
                self.send_response(400)
                self.end_headers()
                return

        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length).decode() if content_length > 0 else ""

        # Validate JSON if body exists
        if body:
            try:
                json.loads(body)
            except json.JSONDecodeError:
                self.send_response(400)
                self.end_headers()
                self.wfile.write(json.dumps({"error": "Invalid JSON"}).encode())
                return

        if self.path == "/handshake":
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
