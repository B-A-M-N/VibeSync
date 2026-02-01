# 👔 For Hiring Managers: Engineering Audit & Systems Design

If you are evaluating this repository for a technical role, this document provides a high-level audit of the architectural decisions, security invariants, and systems engineering principles implemented in **VibeSync**.

---

## 🏛️ 1. Orchestrator vs. Adapter Separation (Distributed Systems)
The core innovation of VibeSync is the move from a simple script-bridge to a **Governed Distributed System**.
*   **The Problem**: Creative engines (Unity/Blender) are non-deterministic, single-threaded environments that crash or hang easily. Relying on an AI to manage state across them directly is a recipe for corruption and "split-brain" scenarios.
*   **The Solution**: We implemented a **Split-Architecture Control Plane**.
    *   **The Brain (Go)**: A high-concurrency Go Orchestrator acts as the central source of truth, enforcing **Strict Serializability** and monotonic causality.
    *   **The Limbs (C#/Python)**: Untrusted, "dumb" adapters that execute raw mutations within the engines.
*   **Engineering Impact**: This ensures that even if an engine crashes or a network packet is lost, the **Write-Ahead Log (WAL)** in the Orchestrator can reconstruct the state and guarantee cluster consistency.

## 🛡️ 2. Zero-Trust Security & "Iron Box" Isolation
VibeSync treats both the AI and the engines as **hostile operators**.
*   **Semantic Firewall**: Every payload is audited via **Recursive AST Parsing** before execution. We block high-risk capabilities (Reflection, Shell execution, non-localhost network access) at the gateway.
*   **Iron Handshake**: Communication is secured via **HMAC-SHA256 request signing** and **Session Token Rotation**. Every transaction is verified against a 5-second anti-replay window.
*   **Circuit Breaker (Deadman Switch)**: A 5000ms heartbeat monitor tracks engine liveness. Any hang or deadlock triggers an immediate **Global PANIC** lock, freezing all hierarchies to prevent data loss.

## ⚛️ 3. Atomic Sync & Transactional Integrity
Synchronizing state between two different engines (with different coordinate systems and memory models) requires absolute atomicity.
*   **Mechanism**: Every sync operation follows a formal lifecycle: `submit_intent` -> `validate_intent` -> `begin_atomic_operation` -> [Mutation] -> **Provisional Commit (Fast Path)** -> `verify_engine_state` (Async) -> `commit_atomic_operation` (Finalize). This separates user-visible immediacy from ontological finality.
*   **Safety**: If background verification fails, the Orchestrator issues an authoritative **Rollback**, restoring the object graph to the last known-good state via engine undo tokens or snapshots.
*   **Outcome**: One AI Intent = One Clean, Verified Step across the entire cluster.

## 🏃 4. Performance Engineering & Main-Thread Safety
Managing the "Two Gods" problem (Blender and Unity both wanting main-thread control) required careful thread marshalling.
*   **Split-Thread Listener**: Both adapters use a background listener thread to handle immediate health checks and buffering, while marshalling mutations to the engine's main thread (via `bpy.app.timers` or `EditorApplication.update`) to prevent C-level race conditions.
*   **Numerical Stability**: All transform and material deltas pass an `auditPayload` check for NaN/Inf, preventing "Floating-Point Poisoning" of the engine's scene graph.

## 🧠 5. Epistemic Governance (Anti-Psychosis Protocol)
To solve the "Blind AI" problem and prevent "AI Psychosis" (hallucinated success):
*   **Forensic Write-Ahead Log (WAL)**: Every intent, rationale, and engine response is journaled with a unique `tid` (Transaction ID).
*   **Referee Mental Model**: The Orchestrator acts as an independent "Referee," forcing the AI to verify its own work through mandatory telemetry read-backs (`read_engine_state`).
*   **Rationale Journaling**: The AI is mandated to provide a `Rationale` for every mutation, ensuring that every state change is traceable and technically justified.

---

## 🛠️ Tech Stack Summary
- **Orchestrator**: Go (High-concurrency, hardened IPC, WAL persistence).
- **Adapters**: C# (Unity Editor SDK), Python (Blender `bpy` SDK).
- **Security**: HMAC-SHA256, AST Auditing, Seccomp profiles, AppArmor.
- **Architecture**: Distributed State Machine, Monotonic ID ordering, Atomic Transactions.

**Conclusion**: VibeSync demonstrates more than just "engine integration"—it demonstrates **Distributed Systems Governance**. It provides a hardened blueprint for how to synchronize complex, non-deterministic creative environments safely at scale.
