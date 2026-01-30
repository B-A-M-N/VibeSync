# VIBE MCP Pain Points Checklist

This checklist is designed to prevent common AI-driven failures in Blender and Unity integrations. Use this during development, PR reviews, and debugging to ensure the bridge remains stable and predictable.

## 🛡️ 1. Safety & Stability
- [x] **No Ghost Mutations:** Every operation is gated by `submit_intent` and `Rationale`.
- [ ] **Crash-Proofing:** Are long-running or heavy operations gated by performance checks?
- [ ] **Undo Reliability:** Can every operation be reverted with a single `Ctrl+Z`?
- [x] **Atomic Recovery:** Implemented `sync_asset_atomic` with auto-rollback on hash mismatch.

## 🚫 2. Overreach & Assumptions
- [x] **Hands-Off Cleanup:** Forbidden by `AI_ENGINEERING_CONSTRAINTS`.
- [ ] **Asset Protection:** Are renames/deletions blocked unless primary goal?
- [x] **Zero-Guess Generation:** Forbidden by `NON_GOALS.md` and `AI_ENGINEERING_CONSTRAINTS`.

## 🧠 3. Context Awareness
- [x] **Inspect-First Logic:** Enforced "Read-Before-Write" lifecycle in orchestrator.
- [x] **Post-Change Verification:** Implemented `verifyEngineState` loop (Referee Mental Model).
- [ ] **Translation Accuracy:** Are engine-specific concepts mapped correctly?

## ⚡ 4. Workflow & Performance
- [x] **Sandbox Purity:** All temporary files are stored in `.vibesync/tmp`.
- [x] **Hot-Reload Safety:** Generation tracking in `handshake_init` and heartbeats detect reloads.
- [x] **Resource Throttling:** Implemented `isRateLimited` in Go Orchestrator.

## 🏗️ 5. Project Scaling
- [ ] **Big Scene Resilience:** Does the bridge remain deterministic for 1000+ objects?
- [x] **Reference Integrity:** Using GUIDs/UUIDs in the `global_id_map`.
- [x] **Version Gatekeeping:** `handshake_init` and heartbeat verify version and generation.

## 🔒 6. Security & Verification
- [x] **Zero Trust Input:** Assets are imported into sandboxes and verified by hash.
- [x] **Execution Lock:** `auditPayload` blocks `exec()`, `Reflection`, and `os.system`.
- [x] **Pipeline Integrity:** Orchestrator acts as the central validator for all engine mutations.

## 🔍 7. Debugging & Diagnostics
- [x] **Reasoned Failures:** `dispatchVibeEvent` and WAL provide detailed failure reasons.
- [x] **Audit Trail:** Every operation is logged to the forensic WAL with `mid` and `tid`.
- [x] **Non-Silent Errors:** Orchestrator propagates explicit engine errors to the AI/Human.

---
*Copyright (C) 2026 B-A-M-N*
