# VibeSync: Adapter Contract & Extensibility (v0.4)

Adapters are the "dumb limbs" of the VibeSync system. To be a legal participant in the cluster, an adapter must uphold these strict invariants.

---

## 🏗️ 1. The Dumb Limb Principle
Adapters **MUST NOT** possess autonomous logic, local state persistence (outside the session), or the authority to reject valid Orchestrator commands. Intelligence lives exclusively in the **Go Brain**.

## 🛡️ 2. Security Invariants
- **No Self-Authorization**: Adapters cannot generate their own tokens. They must receive them via a handshake.
- **Strict Whitelisting**: Adapters must only expose the minimum necessary API surface to the Orchestrator.
- **Pattern Rejection**: Adapters must self-audit for dangerous patterns (Reflection, Shell, External Network) as a second layer of defense.

## ⏱️ 3. State & Causality
- **Monotonic Reporting**: Every state report MUST include the last processed `MonotonicID`.
- **Atomic Execution**: Adapters must support "Preflight" checks where they report if a mutation *would* succeed without actually applying it.
- **Fail-Fast**: If an internal engine error occurs, the adapter must report it immediately and enter a `STOPPED` state until the Orchestrator intervenes.

## 🔌 4. Extensibility Requirements
Third-party adapters (e.g., Unreal, Maya) are supported if they implement the following:
1. **JSON-RPC Over Loopback**: The communication protocol is strictly JSON-RPC 2.0 over a localhost socket.
2. **Heartbeat Handler**: Response to a `health_check` ping within <100ms.
3. **Snapshot API**: Ability to generate a deterministic hierarchy and hash of the current scene.
4. **Sandboxed Import**: A mechanism to import assets into a temporary/hidden namespace for validation before committing to the main scene.

---
*VibeSync: Modular Reality.*
