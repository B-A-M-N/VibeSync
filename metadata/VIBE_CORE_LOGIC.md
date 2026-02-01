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
