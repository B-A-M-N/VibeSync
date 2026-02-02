# VibeSync: Core Operational Logic

This document defines the essential invariants of the VibeSync cluster. These rules are enforced by the Go Orchestrator and are non-negotiable.

---

## ⚖️ 1. The Law of Reality (Determinism)
- **Hash Supremacy**: No mutation is valid until the target engine returns a hash matching the source.
- **The Silence Law**: Success messages without measurable telemetry deltas are treated as failures.
- **Independent Verification**: Every mutation is followed by a mandatory state read-back to prove intent matches reality.

## ⏱️ 2. The Law of Causality (Temporal Control)
- **Monotonicity**: Every event carries a monotonic counter. Out-of-order events are dropped.
- **Atomic Windows**: During a transaction, the target scope is locked. Overlapping mutations trigger a rollback.

## 🚨 3. The Law of Stability (Circuit Breaker)
- **Heartbeat Requirement**: Engines must respond to a health check every 5 seconds.
- **Auto-Panic**: Any heartbeat failure or unhandled exception triggers an immediate hierarchy lock in both engines.

## 🔍 4. The Law of Forensic Necessity
- **Mandatory Consultation**: If an error trigger (defined in `metadata/LOG_TROUBLESHOOTING_MAPPING.md`) occurs, the bridge MUST consult the associated logs before attempting recovery.
- **Hash-Gated Analysis**: Logs are only re-processed if their state hash has changed, ensuring deterministic troubleshooting context.
- **Context Primacy**: Forensic evidence from logs overrides "Success" status messages from adapters.

## 🏎️ 5. The Law of Speculative Finality
- **Non-Blocking Verification**: To eliminate UI latency, the Orchestrator may grant **Provisional Commit** status to "Fast Path" operations (Transforms, Materials) while verification happens asynchronously.
- **State Transitions**: Mutations MUST follow the monotonic path: `PROVISIONAL` -> `FINAL` (Verified) OR `ROLLED_BACK` (Mismatch).
- **Deterministic Rollback**: Any provisional state that fails deferred verification MUST be rolled back automatically using engine-level undo tokens.
- **Atomic Batching**: Micro-intents are coalesced into semantic batches to reduce verification overhead.

## 📜 6. The Law of Log-Driven Governance
- **Ingestion Precondition**: No intent may be submitted until the actor has ingested the recent forensic history.
- **History Anchoring**: Intents must be grounded in the `log_hash` provided by the Orchestrator.

## 🔢 7. The Law of Opcode Integrity
- **Strict Mapping**: All engine commands MUST be issued via strictly mapped Opcodes (0x01-0x11).
- **Intent Binding**: Commands that deviate from the declared intent (e.g., `RIG` performing a `LIGHT` operation) are mechanically rejected.

## 🛑 8. The Law of Failing Closed (Anti-Thrashing)
- **Terminal States**: If a specific failure signature (Hash of Engine+Opcode+Error) appears twice, the intent enters a **TERMINAL** state.
- **Authority Lock**: Terminal states cannot be bypassed by any AI model. They require an **External Reset** from a human operator.
- **Failure over Cleverness**: The system must fail closed and stop, rather than attempting "clever" but unverified workarounds.

---

## 🛰️ Engine State Machine
| State | Allowed Ops |
| :--- | :--- |
| **INIT** | handshake |
| **READY** | ALL |
| **BUSY** | health, read |
| **PANIC** | NONE (Hierarchy Locked) |

---
*Copyright (C) 2026 B-A-M-N*
