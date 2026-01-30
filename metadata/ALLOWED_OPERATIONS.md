# VibeSync: Allowed Operations & Boundary Classifications

This document defines what the Orchestrator is permitted to do. If an operation's end-state cannot be independently verified, it is **Outlawed**.

---

## ✅ Class A: Pure State Targets (Permitted)
*Operations where the Orchestrator defines the result, not the method.*

| Operation | Verification Method |
| :--- | :--- |
| **Transform Set** | Telemetry read-back of world coordinates. |
| **Mesh Transfer** | Binary Hash comparison (Source vs Target). |
| **Material Sync** | Property value dump (e.g., Albedo == [1,0,0]). |
| **Hierarchy Sync** | UUID-based tree traversal and parent matching. |
| **Visibility/Lock** | State-bit verification in engine metadata. |

## ⚠️ Class B: Atomic Procedural (Conditional)
*Permitted only if followed by an immediate Class A re-verification.*

| Operation | Constraint |
| :--- | :--- |
| **Object Creation** | Must be followed by a UUID existence check. |
| **Modifier Application** | Must be followed by a full Mesh Hash re-generation. |
| **Animation Baking** | Must be verified via sample-point checksums. |

## 🚫 Class C: Prohibited Operations (Outlawed)
*Operations that are context-dependent, non-deterministic, or unverifiable.*

1.  **Live Mirroring**: No keystroke or mouse-movement replication.
2.  **Undo Sync**: The Orchestrator does not manage the editor's internal undo stack.
3.  **Side-Effect Logic**: No triggering of "Black Box" scripts that mutate untracked state.
4.  **Implicit Context**: No commands like "Delete Active Object" (Must be "Delete Object with UUID X").
5.  **Edit Mode Mirroring**: No sub-object (vertex/edge/face) live-syncing during active editing.

---

## ⚖️ The Litmus Test
Before implementing any new tool, the developer must answer:
> **"Can I verify the end state independently of how the editor claims it got there?"**
- If **NO**: The tool is rejected.
- If **YES**: Proceed to Class A or B implementation.

---
*Copyright (C) 2026 B-A-M-N*
