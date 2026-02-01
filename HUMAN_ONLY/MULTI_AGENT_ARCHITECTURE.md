# 🤖 Multi-Agent Isolation Architecture (Tri-Silo Model)

This document defines the "ideal" configuration for high-scale automation using multiple AI agents coordinated by the VibeSync kernel.

---

## 1. The Tri-Silo Model

To prevent context poisoning and "mental bleed" between different engine coordinate systems and logics, we implement three distinct sandboxes:

### 🧠 Agent Alpha (The Kernel Coordinator)
- **Scope**: Exclusive access to the **Go Orchestrator MCP**.
- **Rules**: Follows `BRIDGE_CONTRACT.md`. Deals only in **UUIDs** and **Intents**.
- **Role**: Maintains the Write-Ahead Log (WAL) and determines high-level scene state goals. It never sees raw C# or Python.

### 🧊 Agent Beta (The Blender Specialist)
- **Scope**: Exclusive access to the **BlenderVibeBridge MCP**.
- **Rules**: Follows Blender `FREEZE_PROOF_GUIDE.md`. 
- **Role**: Translates intents into precise `bpy` operations. It is oblivious to Unity's existence.

### 🎮 Agent Gamma (The Unity Specialist)
- **Scope**: Exclusive access to the **UnityVibeBridge MCP**.
- **Rules**: Follows Unity `AI_ENGINEERING_CONSTRAINTS.md`.
- **Role**: Ensures Unity state matches target hashes. Operates in Y-up space, oblivious to Blender's Z-up world.

---

## 🛡️ 2. Isolation Mechanisms (The Firewall)

### A. Protocol-Only IPC
Messages passed from the Coordinator to specialists are filtered through the **VibeSync Sanitizer**. This script strips engine-specific jargon (e.g., removing "GameObject" from Blender-bound messages) to prevent cross-contamination of the agent's mental model.

### B. Stateless Specialist Prompts
Specialists are treated as "ephemeral." Every invocation includes the **Current State Hash** and the **Target Intent**, ensuring they don't rely on long-term chat history that could contain stale or poisoned data.

### C. The UUID Firewall
All technical engine-internal handles (pointers, InstanceIDs) are stripped at the Orchestrator level. All agents communicate using the **Global ID Map** (UUIDs), ensuring they can reference the same entity without sharing technical "scaffolding" that leads to confusion.

---

## 🚀 3. Ideal Hardware/Model Mapping
| Role | Recommended Model | Priority |
| :--- | :--- | :--- |
| **Coordinator** | Gemini 1.5 Pro / Claude 3.5 Sonnet | Deep Reasoning, State Integrity |
| **Blender Specialist** | Gemini 1.5 Flash / Claude 3 Haiku | Speed, Literal Protocol Adherence |
| **Unity Specialist** | Gemini 1.5 Flash / Claude 3 Haiku | Speed, Literal Protocol Adherence |

---
*Copyright (C) 2026 B-A-M-N*
