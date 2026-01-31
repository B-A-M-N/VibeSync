# ⚖️ VibeSync: Absolute Invariance Contract (v0.4)

> **VibeSync is not a messenger. It is a transactional state machine with memory.**

This document defines the mechanical, contextual, and semantic locks that ensure "Mathematical Truth" is preserved across the Unity ↔ Blender bridge.

---

## 🏛️ Core Principle: Proof over Intent
VibeSync never forwards intent. It forwards **Proof**. If proof cannot be produced (via hashes or WAL entries), the bridge **refuses to act**.

---

## 🔒 1. The Invariance Poll Set

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

## 🛡️ 2. The Anti-Tail-Chasing Rule

> **If `wal_hash`, `delta_hash`, and `handshake_state` are unchanged, the AI is FORBIDDEN from retrying, re-diagnosing, or re-syncing.**

No state movement → No action allowed.

---

## 🔒 3. The Triple-Lock Summary

| Layer | Locked By | Purpose |
| :--- | :--- | :--- |
| **Mechanical** | Go Orchestrator + WAL | Prevents illegal state transitions. |
| **Contextual** | Force-fed hashes in Tool Outputs | Anchors the AI's mental model to reality. |
| **Semantic** | Commit Proof Gate | Forces the AI to reason through the evidence. |

---
*VibeSync: Engineering Reality.*
*Copyright (C) 2026 B-A-M-N*
