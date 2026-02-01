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
        
        # Load System Instructions (Local Gem)
        self.system_instruction = self.load_gem_prompt()
        
        # Directory paths...

    def load_gem_prompt(self):
        prompt_path = "metadata/GEMS_SYSTEM_PROMPTS.md"
        if not os.path.exists(prompt_path):
            return f"You are the {self.engine} {self.role}."
            
        with open(prompt_path, "r") as f:
            content = f.read()
            
        # Extract the specific role's block from the MD file
        # This is a simple parser looking for the Role header
        role_marker = f"Agent {'Beta' if self.engine == 'blender' else 'Gamma'}-{'1' if self.role == 'Foreman' else '2'}"
        try:
            role_block = content.split(role_marker)[1].split("```text")[1].split("```")[0].strip()
            return role_block
        except:
            return f"You are the {self.engine} {self.role}."

    def think(self, work_order):
        # Foreman: Convert High-level intent to strictly mapped Opcode
        prompt = {
            "contents": [{"parts": [{"text": json.dumps(work_order)}]}],
            "system_instruction": {"parts": [{"text": self.system_instruction}]}
        }
        res = requests.post(f"{API_URL}?key={self.api_key}", json=prompt)
        # ... logic to return opcode-mapped order ...
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
