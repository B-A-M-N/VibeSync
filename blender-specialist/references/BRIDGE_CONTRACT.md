# 📜 VibeSync: The Bridge Contract (v0.4)

This document defines the absolute authority boundaries and operational guarantees for the Unity ↔ Blender synchronization cluster.

---

## 🏛️ 1. Authority Boundaries
- **Orchestrator (Go)**: The sole source of truth and ordering. If an engine contradicts the Orchestrator, the engine is wrong.
- **Human**: Absolute override authority.
- **AI**: An operator bound by the Iron Box; has zero authority to bypass protocol checks.
- **Engines (Unity/Blender)**: Passive execution environments.

## ⚛️ 2. Sync Guarantees
- **Strict Serializability**: Commands are processed in a total order. No interleaving.
- **Eventual Consistency**: All engines will reach the same state once the intent queue is drained.
- **Atomicity**: Mutations are all-or-nothing (Snapshot -> Commit -> Rollback).

## 🛡️ 3. State Ownership
- **Single-Writer**: Only one entity may own a specific UUID at a time during a transaction.
- **No Implicit Merges**: Desync is handled via Rollback, never via "best-guess" merging.

## 🚫 4. Explicit Non-Goals
- **No Real-time Speed over Integrity**: We would rather be slow and correct than fast and corrupt.
- **No Autonomous Cleanup**: We do not "fix" scenes unless explicitly commanded.

---
*VibeSync: Engineering Reality.*
*Copyright (C) 2026 B-A-M-N*
