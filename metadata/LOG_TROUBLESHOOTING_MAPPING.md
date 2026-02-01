# VibeSync: Forensic Log Mapping & Troubleshooting Triggers

This document defines the mandatory mapping between error triggers and their associated forensic logs.

---

## 🔍 Forensic Triggers (The "Must-Read" List)

| Trigger Pattern | Target Log(s) | Forensic Rationale |
| :--- | :--- | :--- |
| `NullReferenceException` | `unity-bridge/VibeBridge_Debug.log`, `mcp-server/vibe_server.log` | Identity drift or invalid object handle. |
| `Kernel deadlock` | `mcp-server/vibe_server.log`, `.vibesync/events.jsonl` | Lock contention or blocked main thread. |
| `Assembly trap` | `unity-bridge/UnityAssembly.log` | Compilation error blocking domain reload. |
| `Authentication Bypass` | `mcp-server/vibe_server.log`, `.vibesync/wal.jsonl` | Token rotation failure or HMAC mismatch. |
| `Drift Detected` | `.vibesync/wal.jsonl`, `metadata/bridge_activity.txt` | Desync between source and target hashes. |
| `ECONNREFUSED` / `Timeout` | `scripts/preflight.py` output, `mcp-server/vibe_server.log` | Port conflict or zombie engine process. |
| `NaN / Inf Detected` | `blender_bridge.log`, `mcp-server/vibe_server.log` | Numerical instability in transform/material data. |

---

## ⚖️ Forensic Enforcement Rules

1. **Trigger Recognition**: If any tool output or console message matches a **Trigger Pattern**, the agent MUST immediately halt and transition to **Forensic Phase**.
2. **Hash-Verified Reading**: Before reading a log, the agent MUST check the file's hash (or last modified time). If unchanged since the last read, use the existing context. If changed, the log MUST be read and structured.
3. **Structured Context Injection**: Raw log data must be summarized into a `Forensic Report` before being added to the active working context.
4. **Conditional Continuation**: No further mutations are permitted until the `Forensic Report` explains the root cause of the trigger.

---
*Copyright (C) 2026 B-A-M-N*
