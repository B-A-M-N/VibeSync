# VibeSync: Multiplexer Specification (The Kernel-Driver Model)

This document defines how sub-MCPs (Sensors) interact with the Go Orchestrator to access the engine HTTP sockets.

---

## 🏗️ 1. The Kernel-Driver Model
In the 4-MCP stack, the **VibeSync Go Orchestrator** is the **Kernel**. It owns the HTTP sockets, security tokens, and state telemetry. Other MCPs (Vision, Selection) are **Drivers**.

### **The Flow**:
1.  **Sensor** generates an intent (e.g., "Selection changed").
2.  **Sensor** calls `vibe_multiplex` on the Kernel.
3.  **Kernel** audits the command against the **Driver Registry**.
4.  **Kernel** executes the HTTP call, attaches security headers, and assigns a `MonotonicID`.
5.  **Kernel** returns the verified telemetry to the **Sensor**.

---

## 📡 2. Driver Registration
Before a Sensor can use the multiplexer, it must declare its "Capability Set" to the Kernel.

| Driver ID | Allowed Endpoints | Knowledge Domain |
| :--- | :--- | :--- |
| **vision_mcp** | `render/capture`, `material/get`, `light/get` | Visual Data & Rendering |
| **selection_mcp** | `selection/set`, `hierarchy/get`, `camera/frame` | Hierarchy & Focus |

---

## 🔒 3. Multiplex Protocol
Every `vibe_multiplex` call must include:
- `sensor_id`: The registered ID of the calling MCP.
- `endpoint`: The target engine-adapter path.
- `payload`: The JSON data for the command.

**Kernel Knowledge Requirement**: The Kernel does not own the Sensor's logic, but it **MUST** maintain a local schema for every "Allowed Endpoint" to perform the **Law of Independent Verification**.

---
*Copyright (C) 2026 B-A-M-N*
