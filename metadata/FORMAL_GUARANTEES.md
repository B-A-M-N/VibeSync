# VibeSync: Formal Guarantees & Non-Guarantees (v0.4)

This document defines the contractual behavior of the VibeSync Orchestrator. These are not aspirations; they are the foundational rules of the system.

---

## ⚖️ 1. Causality & Ordering Guarantees

- **Single-Writer Semantics**: The Orchestrator is the sole authority for state mutation. No engine-to-engine direct mutations are permitted.
- **Total Order of Intents**: All intents are strictly linearized via a global monotonic counter. Two intents NEVER interleave.
- **Causal Consistency**: Operation B is guaranteed to see the effects of Operation A if B was issued after A in the Orchestrator's timeline.
- **Conflict Resolution**: In the event of divergent state between Unity and Blender, the **Orchestrator's Source of Truth** (usually the last committed Snapshot) wins. Engines are treated as eventually-consistent caches that must be purged on conflict.

---

## 🏚️ 2. Failure Domains & Persistence

| Component Failure | Recoverability | Persistence Impact |
| :--- | :--- | :--- |
| **Adapter Crash** | **High** | Session revoked. State survives in the engine; must re-handshake. |
| **Engine Crash** | **Partial** | WAL replay can recover intent, but unsaved engine state is lost. |
| **Orchestrator Crash** | **High** | WAL (Write-Ahead Log) replay restores all committed transactions. |
| **Network Partition** | **Zero** | Triggers **Circuit Breaker**. Mutation is aborted/rolled back. |
| **Power Loss** | **Medium** | Committed WAL entries survive. Mid-commit transactions are rolled back. |

- **Persistence Guarantee**: Any mutation that has received a `commit_atomic_operation` confirmation is persistent in the Orchestrator's WAL.

---

## ⏱️ 3. The Time Model

- **Authoritative Clock**: The Orchestrator's internal monotonic counter is the **sole source of causality**.
- **Advisory Timestamps**: Wall-clock timestamps (UTC) are recorded for forensic logging but are **non-authoritative** for ordering or logic.
- **Clock Skew**: VibeSync is immune to system clock skew because it does not rely on wall-clock time for state synchronization.

---

## 👑 4. Authority Hierarchy

All operations follow a strict hierarchy of command:
1. **Emergency Override (Human)**: Can bypass any lock or state gate.
2. **Human Intent**: Explicit user commands always override AI-suggested mutations.
3. **Orchestrator (The Brain)**: The final arbiter of state between engines.
4. **AI Intent**: Proposed mutations that must pass Orchestrator validation.
5. **Engine (The Limb)**: Passive execution environment; has zero authority to reject valid Orchestrator intents.

---

## 🛡️ 5. Security & Trust Boundaries

- **Untrusted Adapters**: The Orchestrator assumes all adapters (Unity/Blender) are compromised or hostile. They are bound by a strict **[Adapter Contract](ADAPTER_CONTRACT.md)** but have zero authority.
- **Session Revocation**: Any deviation from the protocol (malformed JSON, unauthorized API calls) results in immediate session revocation.
- **Isolation**: A compromise of the Blender adapter does not grant authority over the Unity adapter; each is isolated by unique session tokens.
- **No Persistence in Adapters**: Adapters are "dumb limbs." They do not store state, logic, or secrets beyond the current session token.

---

## 🚫 6. Non-Guarantees (What We Refuse to Promise)

- **No Real-time Latency**: VibeSync prioritizes **Integrity over Speed**. We do not guarantee sub-millisecond sync.
- **No Engine EULA Enforcement**: Users are responsible for their own compliance with Unity/Blender licensing.
- **No Magic Recovery**: We cannot recover data that was never committed to the WAL or saved in the source engine.
- **No Non-Deterministic Safety**: If an engine's internal physics or solver is non-deterministic, VibeSync can only guarantee the *input* state was synced, not the *simulated result*.

---
*VibeSync: Engineering Reality.*
