# VibeSync: System Design & Architecture (v0.4)

VibeSync uses a central "Brain" (the Orchestrator) to coordinate two "Dumb Limbs" (the Adapters in Unity and Blender). This ensures that one tool doesn't accidentally break the other.

---

## 🏛️ Orchestrator vs. Adapter Architecture

The system is split into two distinct layers to ensure absolute pipeline safety:

### 1. The Orchestrator (Go)
The central authority and "Source of Reality." 
- **Language**: Go (Golang) for high-concurrency and runtime stability.
- **Responsibilities**: 
    - **IPC (Inter-Process Communication)**: Managing the high-speed data bus between engines.
    - **Strict Serializability**: Ensuring commands happen in a perfect, one-at-a-time order.
    - **Write-Ahead Log (WAL)**: An append-only journal of every mutation for forensic replay and crash recovery.
    - **Trust Management**: Rotating session tokens and auditing payloads.

### 2. The Adapters (C#/Python)
Isolated, untrusted endpoints that live within the creative engines.
- **Unity Adapter**: A C# kernel that marshals intents to the Unity Main Thread.
- **Blender Adapter**: A Python server using `bpy` to execute mutations.
- **Responsibilities**:
    - Executing raw mutations (transform sets, material updates).
    - Returning binary state hashes for verification.
    - Reporting engine readiness (Busy/Ready/Panic).

---

## 🔗 Internal Documentation
*   [**Ecosystem Summary**](../HUMAN_ONLY/ECOSYSTEM_SUMMARY.md): Professional overview of the VibeSync project family.
*   [**Architecture Blueprint**](ARCHITECTURE_BLUEPRINT.md): Visual sequence and state diagrams.
*   [**Adapter Contract**](ADAPTER_CONTRACT.md): The formal interface rules for engine adapters.

---
*Copyright (C) 2026 B-A-M-N*
