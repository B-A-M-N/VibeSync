# VibeSync: Canonical Engine State Machine

All adapters must adhere to this formal state machine. The Orchestrator is the final authority on an engine's current state.

---

## 🛰️ State Definitions

| State | Description | Allowed Operations |
| :--- | :--- | :--- |
| **INIT** | Adapter starting; port open; handshake pending. | `handshake` |
| **READY** | Trust established; tokens rotated; healthy. | ALL (Sync, Mutate, Audit) |
| **BUSY** | Engine compiling or AssetDatabase updating. | `health`, `read` (NO MUTATIONS) |
| **DEGRADED** | Minor desync or resource warning. | `read`, `rollback`, `snapshot` |
| **DESYNC** | Hash mismatch or generation drift. | `panic`, `rollback`, `handshake` |
| **HALTED** | Panic triggered; hierarchy locked. | `NONE` (Requires manual reset) |

---

## 🔄 State Transitions

1.  **INIT → READY**: Triggered by successful `initiate_handshake`.
2.  **READY → BUSY**: Reported by engine `/health` during compilation.
3.  **ANY → DESYNC**: Triggered by hash mismatch or heartbeat failure.
4.  **ANY → HALTED**: Triggered by `trigger_panic` or `Terminal Failure`.
5.  **HALTED → INIT**: Mandatory path for recovery; requires manual user restart of the adapter.

---

## 🧵 State Enforcement
- The Orchestrator will **reject** any mutation command sent to an engine in `BUSY`, `DEGRADED`, or `HALTED` states.
- If an engine enters `DESYNC`, all active transactions are immediately aborted.

---
*Copyright (C) 2026 B-A-M-N*
