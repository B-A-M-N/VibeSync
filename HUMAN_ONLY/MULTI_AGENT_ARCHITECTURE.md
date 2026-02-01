# 🤖 Multi-Agent Pipelined Studio Model (5-Agent Cluster)

This document defines the "ideal" configuration for high-scale automation using multiple AI agents coordinated by the VibeSync kernel.

---

## 1. The Pipelined Studio Model

To prevent context poisoning and "mental bleed" between different engine logics, we implement a vertical stack of five distinct roles:

### 🧠 Agent Alpha (The Conductor)
- **Scope**: Exclusive access to the **Go Orchestrator MCP**.
- **Role**: Creative Director. Maintains the high-level strategy and global Write-Ahead Log (WAL).
- **Communication**: Pushes strategy to the **Global Inbox**.

### 🔍 Agents Beta-1 & Gamma-1 (The Foremen)
- **Role**: Forensic Strategists for Blender and Unity.
- **Responsibility**: Ingest error logs, monitor engine sentinels, and translate intents into **strictly-mapped Opcodes**.
- **Scope**: Headless. Oblivious to each other's engines.

### ⌨️ Agents Beta-2 & Gamma-2 (The Operators)
- **Role**: Execution Specialists.
- **Responsibility**: Receive Opcodes and generate engine-specific scripts (bpy/C#).
- **Scope**: Stateless. Oblivious to forensic history or strategy.

---

## 🏢 2. The Mailbox Pipeline (.vibesync/queue/)

Agents communicate asynchronously through a directory-based mailbox system:

1.  **Conductor** drops a `WorkOrder` into `[engine]/inbox`.
2.  **Foreman** reads the order, ingests engine logs, and writes an `Opcode` into `[engine]/work`.
3.  **Operator** reads the opcode, generates the script, and pushes it to the engine.
4.  **Operator** writes the result (hash + status) into `[engine]/outbox`.
5.  **Orchestrator** reads the result and finalizes the transaction in the WAL.

---

## 🛡️ 3. Isolation & Hardening

### A. Context Partitioning
The **Operator** never sees the 50 previous error logs that the **Foreman** had to ingest. This prevents "AI Psychosis" where the coder gets confused by its own history.

### B. Precision Opcodes
Because the **Foreman** dictates the Opcode, the **Operator** is mechanically forbidden from performing "Creative" (hallucinated) actions that fall outside the current task's scope.

---

## 🚀 4. Recommended Hardware/Model Mapping
| Role | Recommended Model | Priority |
| :--- | :--- | :--- |
| **Conductor** | Gemini 1.5 Pro | Strategy, Reasoning |
| **Foremen** | Gemini 1.5 Pro | Forensics, Log Ingestion |
| **Operators** | Gemini 1.5 Flash | Speed, Precision, Cost |

---
*Copyright (C) 2026 B-A-M-N*
