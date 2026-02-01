# Iron Box Save-Game Protocol (Git Isolation)

This protocol defines how to use the `.git_safety` repository to create high-fidelity snapshots of project state without bloating the primary codebase repository.

---

## 🏛️ 1. The Separation of Concerns
- **Logic Kernel (`.git`)**: Stores source code, tool logic, and configuration. Pushes to remote. **Prohibited**: `.fbx`, `.png`, `.mat`, `.unity`, `.blend`, or any large binary asset.
- **Save-Game Cluster (`.git_safety`)**: Stores the current state of the 3D scene, materials, and baked assets. **Local-only**. **Purpose**: Immediate rollback if an engine crash or AI hallucination corrupts the scene.

---

## 🔄 2. The Snapshot Procedure
Before performing any **High-Risk Operation** (Baking, Exporting, Batch Renaming, Shader Application), the AI MUST:

1.  **Stage All Assets**:
    ```bash
    git --git-dir=.git_safety --work-tree=. add .
    ```
2.  **Commit Snapshot**:
    ```bash
    git --git-dir=.git_safety --work-tree=. commit -m "[PHASE] Pre-Operation Snapshot: [Description]"
    ```
3.  **Verify Integrity**: Ensure the snapshot is recorded before proceeding with the mutation.

---

## 🔁 3. The Rollback Procedure
If the engine crashes or the "Verify Engine State" tool fails:
1.  **Halt All Operations**.
2.  **Restore State**:
    ```bash
    git --git-dir=.git_safety --work-tree=. restore .
    ```
3.  **Audit Logs**: Check `metadata/LOG_TROUBLESHOOTING_MAPPING.md` to understand why the failure occurred.

---

## 🏗️ 4. Unity ↔ Blender Safe Bridge Workflow
- **Sandbox Import**: Import assets into a `Sandbox_` folder or scene first.
- **Material Verification**: Check naming parity (`Blender Material Name` == `Unity Material Slot`) before overwriting.
- **Poiyomi Check**: If using Poiyomi shaders, ensure the "Lock" state is accounted for during sync.

---
*Copyright (C) 2026 B-A-M-N*
