# VibeSync: Failure Modes & Recovery Protocols

This document defines the canonical taxonomy of failures within the VibeSync cluster and the required response for each.

---

## 🟥 1. Terminal Failures (Immediate PANIC)
*Definition: Any failure that compromises state integrity or trust boundaries.*

| Failure | Cause | Protocol |
| :--- | :--- | :--- |
| **Generation Drift** | Engine reload without handshake. | Instant **PANIC**; Lock all engine hierarchies. |
| **Contract Violation** | Adapter sends malformed JSON or illegal fields. | Instant **PANIC**; Invalidate session tokens. |
| **Security Breach** | Blocked pattern (e.g., `os.system`) detected. | Instant **PANIC**; Permanent quarantine of adapter. |
| **Deadlock** | Heartbeat timeout (>10s). | Trigger **Circuit Breaker**; Halt all mutations. |

---

## 🟨 2. Recoverable Failures (Auto-Rollback)
*Definition: Operational errors that can be reverted to a last known good state.*

| Failure | Cause | Protocol |
| :--- | :--- | :--- |
| **Hash Mismatch** | Post-import validation fails. | Execute `rollback` on target; Notify user. |
| **Resource Exhaustion** | Vertex/Texture limit exceeded. | Block transfer; Return error to orchestrator. |
| **Numerical Instability** | NaN/Inf detected in payload. | Drop command; Log to WAL; Re-fetch current state. |
| **Engine Busy** | Mutation attempted during compilation. | Queue operation or retry after 2s backoff. |

---

## 🟦 3. Determinism Escalation
*Definition: Protocol for when the "Mathematical Truth" of the scene is in doubt.*

1.  **Halt**: If a hash mismatch occurs, the Orchestrator **MUST NOT** attempt to "fix" the data.
2.  **Quarantine**: The affected object is marked as `DIRTY` in the global ID map.
3.  **Mandatory Re-Snapshot**: The Orchestrator requests a fresh full-scene snapshot from the source before any further mutations are allowed on that branch.

---

## 🛠️ Human Intervention Policy
Human intervention is mandatory **ONLY** when:
- An engine is in a `PANIC` state (Requires manual token reset).
- A `Terminal Failure` is logged in the WAL.
- The `DRIFT_TAXONOMY` indicates `Semantic Drift` that the AI cannot reconcile.
