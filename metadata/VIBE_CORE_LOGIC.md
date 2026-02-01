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
