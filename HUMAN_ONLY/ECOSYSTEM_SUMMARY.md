# 🌐 The VibeSync Ecosystem: Distributed Creative Orchestration

VibeSync is a high-integrity orchestration layer designed to synchronize state between disparate creative engines. It treats **Unity** and **Blender** not as standalone tools, but as "Dumb Limbs" in a larger, governed distributed system.

---

## 🏗️ 1. The Core Components

The ecosystem is comprised of three primary pillars:

### 🧠 [VibeSync (The Orchestrator)](https://github.com/B-A-M-N/VibeSync)
The central "Brain" of the system. Written in **Go**, it acts as the authoritative control plane. 
- **Role**: Handles IPC, enforces security constraints, manages the Write-Ahead Log (WAL), and ensures **Strict Serializability** (linear command execution).
- **Key Tech**: Go, MCP (Model Context Protocol), HMAC-SHA256, Docker.

### 🎮 [UnityVibeBridge (The Unity Adapter)](https://github.com/B-A-M-N/unityvibebridge)
A specialized C# kernel that lives inside the Unity Editor.
- **Role**: Translates Orchestrator intents into Unity-native API calls.
- **Key Tech**: C#, Unity Editor Scripting, Time-Budgeted Main-Thread Marshalling.

### 🧊 [BlenderVibeBridge (The Blender Adapter)](https://github.com/B-A-M-N/BlenderVibeBridge)
A Python-based server that runs within Blender's process.
- **Role**: Translates Orchestrator intents into Blender `bpy` commands.
- **Key Tech**: Python, `bpy`, HTTP/JSON, AST-based security auditing.

---

## 🛡️ 2. Technical Pillars

### ⚛️ Atomic Synchronization
Unlike standard exporters, VibeSync uses a **Snapshot → Preflight → Commit** pipeline. Mutations are transactional; if a sync fails in Unity, the source in Blender is automatically rolled back to ensure cluster consistency.

### 🚧 Zero-Trust Security
The system assumes both the AI and the engines are potentially compromised. Every command is audited via **Recursive AST Parsing** before execution, and all traffic is signed with session-rotated HMAC tokens.

### 🧠 Epistemic Governance
VibeSync solves the "AI Psychosis" problem (hallucinated success) by implementing **Truth Reconciliation Loops**. The Orchestrator acts as a "Referee," forcing the AI to verify its intent against actual engine telemetry read-backs.

---

## 🔄 3. How It Works (The Lifecycle)

1. **Intent**: A user (or AI) proposes a change (e.g., "Update Material X").
2. **Audit**: VibeSync's Go Orchestrator validates the intent for safety and resource limits.
3. **Dispatch**: The Orchestrator signs the request and sends it to the respective Adapters via local loopback.
4. **Execution**: The Adapters marshal the command to the Main Thread of Unity/Blender.
5. **Verify**: The Go Orchestrator reads the post-mutation hash from both engines.
6. **Journal**: The success or failure is recorded in the Forensic Write-Ahead Log.

---

## 🚀 4. Use Cases
- **Real-time Pipeline**: Model in Blender, see it live in Unity with verified physics and material accuracy.
- **AI-Driven Automation**: Use LLMs to rig, light, or arrange scenes safely within production editors.
- **Distributed Collaboration**: A unified control plane for multi-tool creative workflows.

---
**Created by the Vibe Bridge Team.**
**Copyright (C) 2026 B-A-M-N**
