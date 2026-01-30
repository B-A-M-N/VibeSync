# VibeSync: AI Security Threat Acceptance (v1.1)

This document formalizes the accepted risks and non-negotiable defensive mandates for the VibeSync cluster. Use of the system constitutes acceptance of these terms.

---

## 🛑 1. The Zero-Trust Mandate
- **Untrusted AI**: All AI-generated code, intents, and payloads are considered potentially malicious.
- **Untrusted Adapters**: Adapters are "dumb limbs." They are assumed to be compromised or lying.
- **Orchestrator Authority**: The Go Orchestrator is the **sole source of truth** and authority.

## 🛡️ 2. Core Defensive Posture

### A. Strict Isolation
- **Network Isolation**: Orchestrator and Adapters communicate strictly over `localhost`. No external network calls are permitted.
- **Process Isolation**: Every engine (Unity/Blender) runs in its own memory space. Compromise of one must not lead to the compromise of the other.

### B. Payload Auditing
- All incoming mutations MUST pass the **Semantic Firewall** (AST/Pattern-based auditing) to block:
  - Reflection and dynamic code execution.
  - Shell/Process spawning (`os.system`, `exec`).
  - Unauthorized file system access outside `.vibesync/`.

### C. State Verification (Ref-Mental Model)
- **Read-Before-Write**: No mutation is permitted without a pre-check of the target state.
- **Post-Verification**: Every mutation MUST be followed by an independent state read-back.
- **Refusal on Drift**: Any detected drift or hash mismatch results in an immediate **PANIC** and rollback.

## ⚖️ 3. Accepted Risks
- **Non-Deterministic Physics**: Minor floating-point variations in engine-specific solvers are accepted, provided the *input* state remains deterministic.
- **Human-in-the-Loop Delay**: We accept operational latency in favor of human verification for low-confidence mutations.
- **Project Corruption (Legacy)**: Use of experimental bridges (v0.x) carries an inherent risk of data loss. **MANDATORY BACKUPS** are the user's responsibility.

## 🛠️ 4. Security Enforcement (The Gauntlet)
- **Token Rotation**: Every session uses rotated ephemeral tokens.
- **Monotonic Causality**: All operations are serialized and stamped with a monotonic ID. Out-of-order execution is a terminal failure.
- **Circuit Breaker**: Heartbeat failures (>5s) trigger a global lock.

---
*VibeSync: Engineering Trust in an Untrusted World.*
