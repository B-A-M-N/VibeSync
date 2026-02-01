# VibeSync: Formal Guarantees & Non-Guarantees (v0.4)

This document defines the contractual behavior of the VibeSync Orchestrator. These are not aspirations; they are the foundational rules of the system.

---

## ⚖️ 1. Causality & Ordering Guarantees

- **Single-Writer Semantics**: The Orchestrator is the sole authority for state mutation. No engine-to-engine direct mutations are permitted.
- **Total Order of Intents**: All intents are strictly linearized via a global monotonic counter. Two intents NEVER interleave.
- **Causal Hash-Chaining**: Every entry in the Write-Ahead Log (WAL) contains a cryptographic hash of the previous entry. The history of "Reality" is tamper-evident and immutable.
- **Intent Budgeting**: Every session is assigned a temporal budget. Exceeding the "Mutation-Per-Minute" (MPM) threshold triggers mandatory trust degradation.
- **Conflict Resolution**: VibeSync guarantees deterministic resolution of Cosmetic conflicts using Monotonic Intent ID tie-breaking. Structural and Destructive conflicts are guaranteed to be detected and **QUARANTINED** to prevent silent semantic loss.
- **Human Supremacy**: VibeSync guarantees that human intents always trump AI intents. Active human manipulation creates a `HUMAN_ACTIVE` lock that veteos all overlapping AI mutations.
- **The Golden Rule**: If resolving a conflict requires guessing user intent, the system must stop and escalate to the human arbiter.
- **Speculative Finality**: The system guarantees that "Fast Path" mutations (Cosmetic/Transforms) are visible immediately while being verified asynchronously. Finality is deferred until background hash verification is complete.
- **Atomic Rollback Guarantee**: The system guarantees that any speculative state which fails background verification will be automatically and deterministically rolled back to the last known authoritative state.
- **Deterministic Convergence**: The system guarantees that both engines will eventually reach the same hash state once the intent queue is empty.

---

## ⏱️ 3. The Time Model & Behavioral Throttling

- **Authoritative Clock**: The Orchestrator's internal monotonic counter is the **sole source of causality**.
- **Adaptive Throttling**: The system monitors the "entropy" of commands. Rapidly varying or structurally deep mutations trigger a "Cooldown" window where the engine is locked to read-only state.
- **Networking & Latency**: 
  - VibeSync is optimized for **Local Loopback (<1ms)**. 
  - **Jitter Buffer**: The Orchestrator maintains a 50ms jitter buffer for out-of-order monotonic ID arrival. 
  - **Congestion Policy**: If the engine `BUSY` signal persists >5s, the pending mutation queue is purged to prevent memory exhaustion.

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
