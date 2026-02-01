# 🛠️ VibeSync Blender Adapter: Technical Engineering & Installation Guide

This guide is intended for **Engineers, Technical Artists, and Power Users**. It provides a deep dive into the adapter architecture, security invariants, and advanced orchestration patterns for the VibeSync system.

---

## 🏗️ 1. System Architecture Overview

The VibeSync Blender Adapter operates as a **Split-Thread Execution Plane**:

1.  **The Control Plane (Go Orchestrator)**: The central kernel (`vibe-mcp-server`) that handles AST auditing, transaction ordering, and HMAC-signed protocol translation.
2.  **The Execution Plane (Python/Bpy)**: A local-loopback server (`bridge_server.py`) running within Blender that consumes commands via a thread-safe queue and `bpy.app.timers`.

### Architectural Invariants:
- **Zero-Trust Communication**: All requests MUST include `X-Vibe-Token`, `X-Vibe-Signature`, and `X-Vibe-Generation`.
- **Main-Thread Dispatch**: All `bpy` mutations are executed on the Blender main thread to prevent C-level memory corruption and deadlocks.
- **Strict Serializability**: Commands are processed in the total order established by the Go Orchestrator's monotonic counters.

---

## 🚀 2. Advanced Installation & Hardening

### Prerequisites
- **Blender**: 3.6 LTS or newer (Recommended).
- **Python**: 3.10+ (Internal to Blender, but external required for Orchestrator development).
- **Go**: 1.23+ (Required to run the `vibe-mcp-server` Orchestrator).

### Deployment Steps
1.  **Adapter Injection**:
    - Move or link the `blender-bridge/` directory into a location accessible by Blender.
    - Run `bridge_server.py` within Blender's Python environment (via text editor or as a background service).
2.  **Orchestrator Setup**:
    Initialize the Go kernel:
    ```bash
    cd mcp-server
    go run main.go
    ```
3.  **Security Handshake**:
    Establish the initial trust boundary:
    - The Orchestrator rotates the `VIBE_BLENDER_BOOTSTRAP_SECRET` into a unique session token upon the first `/handshake`.

---

## 🧠 3. Advanced Epistemic Governance

To prevent "AI Psychosis" (pattern-loop hallucination), the system implements **Referee-Model Verification**.

### A. Preflight Hashing
The adapter exposes `/preflight/run` to generate a binary state hash.
- **Verification**: The Go Orchestrator compares the pre-mutation hash with the post-mutation readback. If drift is detected, the system enters a `PANIC` state.

### B. Forensic Journaling
Every operation is journaled to the Write-Ahead Log (`wal.jsonl`) with its associated `Rationale`.
- If an AI command fails, the Orchestrator provides the raw Blender traceback to force the AI to reason through the failure mode.

---

## 🛠️ 4. Transactional Mutation Mastery

### Atomic Handshakes & State
The bridge uses generation counters (`X-Vibe-Generation`) to ensure the AI's "View of Reality" matches the engine's current state.
- If a user manually modifies the scene in Blender, the generation increments.
- The Orchestrator will reject any AI command based on an "Outdated" generation.

---

## 🛡️ 5. Security & AST Auditing

The `security_gate.py` in the root directory uses **Recursive AST Parsing** to block malicious payloads before they ever reach the Blender adapter.

### Whitelisted ISA Tools
The adapter only responds to a fixed set of whitelisted paths (see `_path_whitelist` in `bridge_server.py`).
- Any attempt to access a path not in the whitelist results in a `403 Forbidden` at the network layer.

---

## 📊 6. Performance & Liveness

### Resource Monitoring
The adapter monitors engine readiness via the `/health` endpoint.
- **Busy State**: If Blender is rendering or compiling shaders, the adapter reports `busy`, and the Orchestrator throttles outbound commands.

---

## 💀 7. Debugging the Bridge

| Symptom | Diagnosis | Solution |
| :--- | :--- | :--- |
| **401 Unauthorized** | Token Mismatch | Re-run `handshake_init` to rotate session keys. |
| **409 Conflict** | Generation Drift | AI is acting on stale scene data. Refresh telemetry. |
| **403 Invalid Signature** | Clock Skew or Replay | Ensure system clocks are synced within 5 seconds. |
| **Main Thread Hang** | Blocked Operator | Check for active Blender popups or heavy modal operators. |

## 🧠 8. VibeSync Glossary: Speaking the Language

| Term | Definition |
| :--- | :--- |
| **Vibe** | The collective state (transforms, materials, hashes) of a synchronized object. |
| **Orchestrator**| The central Go-based "Brain" that manages all engine communication. |
| **Adapter** | The "Dumb Limb" plugin (Unity/Blender) that executes raw commands. |
| **Monotonic ID** | A logical counter that ensures commands happen in the correct sequence. |
| **Epistemic Refusal**| When the AI stops an action because the outcome is "unknowable" or risky. |
| **Iron Box** | The security layer that wraps every mutation in a verified transaction. |

---

## 🆘 9. Troubleshooting & FAQ

**Q: Why is my engine state showing as "BUSY"?**
**A:** This happens if Unity is recompiling scripts or Blender is rendering. VibeSync throttles commands to prevent crashes. Wait a few seconds for the status to return to "READY."

**Q: Handshake failed with "AUTH_FAILED". What do I do?**
**A:** Ensure you haven't changed the `BOOTSTRAP_TOKEN` in the adapter source code. If you did, update the corresponding token in `mcp-server/main.go` and restart the Orchestrator.

**Q: My object moved in Blender but didn't update in Unity!**
**A:** Check the console for a "Hash Mismatch." This means VibeSync detected a potential data corruption and blocked the sync. Use the tool `reconcile_sync_state` to force a fresh match.

---

## 🛠️ 10. Manual Recovery: Fixing a Broken Scene

If VibeSync enters a `PANIC` or `DESYNC` state, follow these steps to restore your project:

1.  **Stop the Orchestrator**: Press `Ctrl+C` in the terminal.
2.  **Inspect the WAL**: Look at `.vibesync/wal.jsonl`. Find the last successful `COMMIT` entry.
3.  **Manual Rollback**:
    - In Blender: Go to the "VibeBridge" panel and click **"Force Global Re-Index."**
    - In Unity: Select **VibeSync > Emergency > Restore Last Snapshot.**
4.  **Clear Persistence**: If corruption persists, delete the `.vibesync/state.json` file and re-run `handshake_init`.

---
**Copyright (C) 2026 B-A-M-N**
