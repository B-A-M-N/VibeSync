# VibeSync v0.4 Roadmap: "The Cathedral"

This document defines the essential additions required to move VibeSync from a hardened spine (v0.3) to a production-grade studio orchestrator.

---

## 🛰️ 1. Observability & Telemetry
- [x] **Metrics Engine**: Track memory and engine load per adapter via `/metrics`.
- [x] **Event Levels**: Implemented structured logging (INFO/DEBUG/WARN/ERROR) in the event bus (`events.jsonl`).
- [ ] **Real-time Dashboard**: A CLI or Web UI to visualize system health and transaction queues.

## ⚔️ 2. Conflict & Concurrency
- [x] **Locking Strategy**: Implemented per-object and per-collection hierarchy locking via `lock_object`.
- [ ] **Collision Policy**: Define "Last Write Wins" vs. "Human Manual Resolve" for AI/Human collisions.

## 🧠 3. AI Governance & Traceability
- [x] **Rationale Journaling**: AI must provide a "Reasoning String" for every intent via `submit_intent`.
- [x] **Confidence Scoring**: Implemented a `0.0 - 1.0` confidence threshold for automated mutations via `validate_intent`.
- [x] **Manual Override**: A dedicated "Human-in-the-Loop" gate for low-confidence AI actions via `human_approve_intent`.

## 📡 4. Networking & Distribution
- [x] **Resilience**: Implemented exponential backoff for failed engine requests in `sendToEngine`.
- [ ] **WebSocket Sync**: High-frequency telemetry channel for real-time transform feedback.
- [ ] **Encryption**: Mutual TLS or encrypted payloads for session data.

## 🧪 5. Testing & Verification
- [ ] **Integration Test Suite**: Automated "Headless" sync tests for Unity and Blender.
- [x] **Validation Hooks**: AST-based auditing (`auditPayload`) for security and numerical stability.

---
*Copyright (C) 2026 B-A-M-N*
