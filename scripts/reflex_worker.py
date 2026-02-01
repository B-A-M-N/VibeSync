#!/usr/bin/env python3
import os
import time
import json
import requests
import sys

# VibeSync Reflex Worker (Template for B1, B2, G1, G2)
# Handles automated processing of Work Orders via the Mailbox system.

class ReflexWorker:
    def __init__(self, role, engine, queue_dir):
        self.role = role # Foreman | Operator
        self.engine = engine # blender | unity
        self.queue_dir = queue_dir
        self.api_key = os.getenv("VIBE_API_KEY")
        self.api_url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
        
        # Directory paths
        if self.role == "Foreman":
            self.in_dir = os.path.join(queue_dir, engine, "inbox")
            self.out_dir = os.path.join(queue_dir, engine, "work")
        else: # Operator
            self.in_dir = os.path.join(queue_dir, engine, "work")
            self.out_dir = os.path.join(queue_dir, engine, "outbox")

    def run(self):
        print(f"🤖 VibeSync {self.engine.capitalize()} {self.role} Online.")
        print(f"Watching: {self.in_dir}")
        
        while True:
            files = sorted([f for f in os.listdir(self.in_dir) if f.endswith(".json")])
            if files:
                self.process_file(os.path.join(self.in_dir, files[0]))
            time.sleep(0.5)

    def process_file(self, path):
        print(f"📦 Processing: {os.path.basename(path)}")
        with open(path, "r") as f:
            data = json.load(f)
            
        if self.role == "Foreman":
            result = self.think(data)
        else:
            result = self.execute(data)
            
        out_path = os.path.join(self.out_dir, os.path.basename(path))
        with open(out_path, "w") as f:
            json.dump(result, f, indent=2)
            
        os.remove(path)
        print(f"✅ Finished: {os.path.basename(path)}")

    def think(self, work_order):
        # Foreman: Convert High-level intent to strictly mapped Opcode
        # Mocking API call for now
        work_order["opcode"] = 0x03 # Transform (Default)
        return work_order

    def execute(self, work_order):
        # Operator: Generate code and push to bridge
        # Mocking implementation
        return {
            "work_order_id": work_order.get("id", "unk"),
            "status": "SUCCESS",
            "hash": "verified_sha256_mock"
        }

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: reflex_worker.py [Foreman|Operator] [blender|unity]")
        sys.exit(1)
        
    role = sys.argv[1]
    engine = sys.argv[2]
    queue_root = ".vibesync/queue"
    
    worker = ReflexWorker(role, engine, queue_root)
    worker.run()
