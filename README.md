# VibeSync: The Ultimate Zero-Trust Unity ↔ Blender Orchestrator

> [!CAUTION]
> **EXPERIMENTAL RESEARCH PREVIEW**: VibeSync is currently in active development (v0.4). This software is highly experimental and provided "as-is" for research and community evaluation. Users should expect frequent breaking changes and are advised to back up all project data before use.

### ⚖️ Project Status & Access
VibeSync is presently in a **Research & Development phase**. During this v0.x "Crowbar" cycle, all core orchestration features and advanced AI tools are provided at **zero cost** to the general public to encourage stress-testing and technical feedback.

As the system matures toward a v1.0 production release, the Author reserves the right to transition specific high-compute, cloud-integrated, or advanced automation features to a tiered licensing or credit-based access model. Use of the current preview constitutes acknowledgement of this roadmap.

---

VibeSync is a high-fidelity orchestration system...

---

## 🛡️ Hardcore Architecture (Zero-Trust)
*Version 0.3 "Crowbar" is now powered by Go for maximum performance and type safety.*

- **Establishing New Trust**: Every session begins with a bootstrap handshake that **rotates tokens**. The Orchestrator generates high-entropy session keys, ensuring the "Bootstrap Secret" is never used for actual engine mutations.
- **Atomic Sync Workflow**: Mutations follow a strict **Snapshot → Preflight → Export → Import → Validate → Commit/Rollback** pipeline. If a post-import hash mismatch occurs, the system automatically purges the sandbox and reverts the source engine.
- **Semantic Firewall**: The Go server performs AST-based auditing and pattern matching on all payloads to block dangerous operations (e.g., `os.system`, `Reflection`) before they reach the adapters.
- **Circuit Breaker**: A background heartbeat routine monitors engines every 5s. Any communication deadlock or hang triggers an immediate system-wide **PANIC**, locking engine hierarchies to prevent state corruption.
- **Numerical Safety**: Integrated NaN/Infinity sanitization and precision thresholding (>0.0001 delta) eliminate "exploding physics" and network noise.

---

## 🚀 Key Features

### **Atomic Verified Sync**
- **1M Vertex Gating**: Automated resource limits to prevent engine lockups during massive transfers.
- **Sandboxed Imports**: Unity assets are imported into a hidden `.vibesync/tmp` sandbox for validation before being committed to the project.
- **Hash Integrity**: Binary-level verification ensures the mesh/material data in Unity is an exact clone of the Blender source.

### **Unified Lifecycle**
- **Strict State Machine**: Engines transition through `STOPPED` -> `STARTING` -> `RUNNING` -> `PANIC`.
- **Domain Reload Resilience**: Unity readiness checks prevent mutations during script recompilation.
- **Forensic WAL**: Every operation is journaled to a Write-Ahead Log with unique Transaction IDs (`tid`).

---

## 🕹️ Tool Library (v0.3)

| Category | Tool | Description |
| :--- | :--- | :--- |
| **Logistics** | `initiate_handshake` | Establishes trust and rotates session tokens. |
| **Stability** | `trigger_panic` | Broadcasts emergency hierarchy lock to all engines. |
| **Atomic** | `sync_asset_atomic` | Full validated transfer with auto-rollback. |
| **Creative** | `sync_transform` | Safe, delta-based transform synchronization. |
| **Creative** | `sync_material` | Real-time object material property syncing. |
| **Ops** | `lock_object` | Hierarchy-aware locking to prevent concurrent edits. |
| **Diagnostics** | `get_diagnostics` | Real-time health, uptime, and WAL telemetry. |

---

## ⚖️ Dual-License & Credits (v1.2)
VibeSync is distributed under a **Dual-Licensing Model**.

1. **Open-Source Path (AGPLv3)**: Free for non-commercial use.
   - **SOURCE ACCESS**: In accordance with Section 13 of the AGPLv3, if you run a modified version of this software over a network, you must provide a way for users to access your source code. You can find our official source at: `[Your GitHub Repository Link]`.
2. **Commercial Path ("Work-or-Pay")**: For revenue-generating entities. Requires either maintenance contributions or a license fee.
3. **Operation Credits**: Basic sync is free. Complex AI-driven operations (Optimization, Vision Audits) are gated via the **Go Hub Credit System**.

---

## 📦 Installation & Setup

### **Prerequisites**
- **Go 1.23+** (Orchestrator)
- **Python 3.10+** (Dev Watcher)
- **Unity 2022.3+** & **Blender 3.6+**

### **Execution**
```bash
# 1. Enter the server directory
cd mcp-server

# 2. Launch the Orchestrator (Automatic rebuilds on change)
go run watcher.go
```

---
**VibeSync** | Hardcore Distributed Systems for Real-time Content Creation.
