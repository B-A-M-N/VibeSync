# VibeSync: Pathological Edge Case Handling

This document defines the mechanical defenses against three classes of pathological system failures: Nested Prefab Drift, Cyclic Parenting, and Partial Deletes.

---

## 🏗️ 1. Nested Prefabs / Multi-Layer Authority
**Problem**: Divergence between *Definition* vs *Instance* and *Local Override* vs *Source Asset* at different nesting depths.

**Defense**: **Identity-Depth Metadata**
- **Kind Enforcement**: All identity reports MUST distinguish between `PREFAB_DEF`, `PREFAB_INSTANCE`, and `OBJECT`.
- **Depth Mask**: Verification hashes must include the `prefab_depth`. A mutation at depth 2 cannot be legally verified against state at depth 0.
- **Rule**: If `prefab_depth` mismatch is detected during verification, the transaction enters `QUARANTINE`.

---

## 🌳 2. Cyclic Parenting / Graph Invalidity
**Problem**: Creating a loop (A → B → C → A) during distributed mutation where engines may resolve the conflict inconsistently.

**Defense**: **Engine-Independent Central DAG Validation**
- **Central Authority**: The Orchestrator MUST compute the **Ancestor Closure** before issuing any parenting command.
- **Invariant**: `∀ intent(parent=A, child=B): A ∉ ancestors(B)`.
- **Outcome**: Cycles are detected at the **Validate** phase; the mutation is never issued to the engines, and the intent is `ROLLED_BACK` instantly.

---

## 💀 3. Partial Deletes / Identity Holes
**Problem**: Deleting a parent or reference without cleaning up children or dependencies, leading to dangling UUIDs and "ghost" objects.

**Defense**: **Closure-Aware Destructive Intents**
- **Closure Computation**: Before issuing a `DELETE`, the agent MUST compute the **Delete Closure** (Children + Constraints + Modifiers).
- **Atomic Deletion**: The `DELETE` command MUST target the entire closure list.
- **Safety Net**: All Destructive intents require a `.git_safety` snapshot. If any part of the closure fails to delete, the entire graph is restored.

---

## ⚖️ Summary of Stress Test Results

| Pathology | Survives? | Primary Enforcement |
| :--- | :--- | :--- |
| **Nested Prefabs** | ✅ | Depth-aware Identity Metadata |
| **Cyclic Parenting** | ✅ | Central DAG Validation (Pre-mutation) |
| **Partial Deletes** | ✅ | Closure-aware Snapshots |
| **Human-AI Racing** | ✅ | `HUMAN_ACTIVE` Veto Lock |

---
*VibeSync: Engineering Reality.*
*Copyright (C) 2026 B-A-M-N*
