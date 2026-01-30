# VibeSync: Failure Modes & Recovery Protocols

This document defines the canonical taxonomy of failures within the VibeSync cluster and the required response for each.

---

## 🟥 1. Terminal Failures (Immediate PANIC)
*Definition: Any failure that compromises state integrity or trust boundaries. Non-recoverable without manual intervention.*

| Failure | Cause | Protocol | Recoverability |
| :--- | :--- | :--- | :--- |
| **Generation Drift** | Engine reload without handshake. | Instant **PANIC**; Lock hierarchies. | **Manual** |
| **Contract Violation** | Adapter sends malformed JSON/Illegal fields. | Instant **PANIC**; Revoke tokens. | **Manual** |
| **Security Breach** | Blocked pattern (e.g., `os.system`) detected. | Instant **PANIC**; Quarantine adapter. | **Manual** |
| **Deadlock** | Heartbeat timeout (>10s). | Trigger **Circuit Breaker**; Halt. | **Manual** |
| **Orchestrator Crash** | OS/Process failure. | Replay WAL on restart. | **Automatic** |

---

## 🟨 2. Recoverable Failures (Auto-Rollback)
*Definition: Operational errors that can be automatically reverted to a last known good state.*

| Failure | Cause | Protocol | Recoverability |
| :--- | :--- | :--- | :--- |
| **Hash Mismatch** | Post-import validation fails. | Execute `rollback` on target. | **Automatic** |
| **Resource Limit** | Vertex/Texture limit exceeded. | Block transfer; Return error. | **Automatic** |
| **Numerical Error** | NaN/Inf detected in payload. | Drop command; Re-fetch state. | **Automatic** |
| **Engine Busy** | Mutation during compilation/load. | Retry after 2s backoff. | **Automatic** |

---

## 🟦 3. Determinism Escalation
*Definition: Protocol for when the "Mathematical Truth" of the scene is in doubt.*

1.  **Halt**: If a hash mismatch occurs, the Orchestrator **MUST NOT** attempt to "fix" the data.
2.  **Quarantine**: The affected object is marked as `DIRTY` in the global ID map.
3.  **Mandatory Re-Snapshot**: The Orchestrator requests a fresh full-scene snapshot from the source before any further mutations are allowed on that branch.

---

## 👑 4. Authority & Override Policy
1. **Human > AI**: All AI-proposed mutations are held in a `PENDING` state until human validation or a "Trust-Mode" threshold is met.
2. **Emergency Override**: In a `PANIC` state, the Orchestrator provides an `emergency_unlock` command that forces engine hierarchies to open, acknowledging that this may result in permanent state loss.
3. **Brain > Limb**: If an engine (Limb) reports a state that contradicts the Orchestrator (Brain), the Orchestrator **MUST** force the engine to resync from the last known good Snapshot.

---

## 🛠️ Human Intervention Policy
Human intervention is mandatory **ONLY** when:
- An engine is in a `PANIC` state (Requires manual token reset).
- A `Terminal Failure` is logged in the WAL.
- The `DRIFT_TAXONOMY` indicates `Semantic Drift` that the AI cannot reconcile.

