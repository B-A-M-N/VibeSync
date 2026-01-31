# ⚖️ VibeSync: Absolute Invariance Contract (v0.4)

> **VibeSync is not a messenger. It is a transactional state machine with memory.**

This document defines the mechanical, contextual, and semantic locks that ensure "Mathematical Truth" is preserved across the Unity ↔ Blender bridge.

---

## 🏛️ Core Principle: Proof over Intent
VibeSync never forwards intent. It forwards **Proof**. If proof cannot be produced (via hashes or WAL entries), the bridge **refuses to act**.

---

## 🧩 1. The Three Orders of Invariance

| Order | Guards Against | Domain | Definition |
| :--- | :--- | :--- | :--- |
| **1st** | False State | Physics / Bytes | reality anchoring (hashes, WAL, handshakes) |
| **2nd** | False Causality | Control Systems | behavioral sanity (idempotency, budgets, monotonicity) |
| **3rd** | False Belief | Epistemology | belief integrity (provenance, confidence decay, amnesia) |

---

## 🔒 2. The Invariance Poll Set

The following tools are provided for the AI to verify the "Ground Truth" of the bridge.

### 📡 `/bridge/heartbeat` (Epoch Invariance)
Ensures the bridge is alive and all engines are aligned to the same session "Epoch."
- **Hard Gate**: Failure halts all outbound mutations.
- **Goal**: Prevents silent desyncs after an engine crash.

### 🤝 `/bridge/handshake_state` (Reality Alignment)
Verifies if Unity and Blender are describing the exact same asset state.
- **Rule**: `hash_match == false` → **BLOCK ALL ACTIONS**.

### 📜 `/bridge/wal_state` (The Truth Chain)
Exposes the authoritative history from the Write-Ahead Log.
- **Goal**: Eliminates "Tail-Chasing" by making time and causality explicit.

### ⚛️ `/bridge/transaction_state` (Atomicity Lock)
Ensures only one transaction is in-flight and that specific assets are locked.
- **Rule**: Timeout triggers automatic rollback.

### 📐 `/bridge/delta_state` (Mutation Proof)
Explains exactly what changed and provides a hash of the delta.
- **Goal**: Prevents "Ghost Changes" or implicit mutations.

---

## 🚀 3. Invariance Amplifiers (The 12 Mechanical Tricks)

1.  **Hash-Chaining Everything**: Every output is cryptographically chained to the previous.
2.  **Force-Fed Context**: Orchestrator injects verified hashes into every AI turn.
3.  **Proof-of-Work Gate**: Commit requires technical rationale matching verified hashes.
4.  **Dual-Witness Verification**: Facts must be independently observed by ≥2 sources.
5.  **Idempotency Keys**: Operations cannot be reapplied blindly; rejects same key with different hash.
6.  **Monotonic Time**: Actions blocked if tick ≤ last committed tick.
7.  **Entropy Budgets**: Mutation/retry limit per session; exhaustion = hard stop.
8.  **Cross-Layer Reconciliation**: Unity, Blender, and Bridge hashes must match at commit.
9.  **Belief Provenance & Expiry**: AI beliefs decay unless re-validated by fresh hashes.
10. **Preflight/Postflight Checks**: Mechanical verification before and after every operation.
11. **Silence as Error**: Missing heartbeats/signals trigger an immediate block.
12. **Observer Relativity**: Single-point observations are marked as `UNCONFIRMED`.

---

## 🛡️ 4. The Meta-Invariant

> **The AI is not allowed to “fix” invariance violations. Only machines may.**

If an invariant breaks, the AI may explain, summarize, or escalate—but it **cannot bypass**.

---

## 🛡️ 5. The Anti-Tail-Chasing Rule

> **If `wal_hash`, `delta_hash`, and `handshake_state` are unchanged, the AI is FORBIDDEN from retrying, re-diagnosing, or re-syncing.**

---

## 🔒 6. The Triple-Lock Summary

| Layer | Locked By | Purpose |
| :--- | :--- | :--- |
| **Mechanical** | Go Orchestrator + WAL | Prevents illegal state transitions. |
| **Contextual** | Force-fed hashes in Tool Outputs | Anchors the AI's mental model to reality. |
| **Semantic** | Commit Proof Gate | Forces the AI to reason through the evidence. |

---
*VibeSync: Engineering Reality.*
*Copyright (C) 2026 B-A-M-N*