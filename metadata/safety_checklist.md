# VibeSync MCP Hardcore Safety Checklist

## 1. Engine Lifecycle & State Persistence
- [ ] **Domain Reload Recovery**: Implement `OnEnable` state restoration to survive Unity C# recompiles.
- [ ] **Undo Stack Integration**: Ensure server-side mutations are registered in Unity `Undo.RecordObject` and Blender `bpy.ops.ed.undo_push`.
- [ ] **Stale Reference Cleanup**: Auto-invalidate object pointers after engine-side scene loads.

## 2. Asset & Data Integrity
- [ ] **Nested Prefab Resolution**: Track parent-child GUIDs to prevent breaking Unity Prefab Variants.
- [ ] **Constraint Translation**: Warning system for non-mappable Blender Constraints (e.g., Limit Distance).
- [ ] **Import Setting Enforcement**: Auto-configure FBX `Mesh Compression` and `Read/Write Enabled` on push.

## 3. Concurrency & Performance
- [ ] **Hierarchy-Aware Locking**: Lock parent nodes when children are undergoing structural changes.
- [ ] **Delta Transform Sync**: Only broadcast changes exceeding a precision threshold (e.g., 0.0001).
- [ ] **Operation Journaling**: Maintain a write-ahead log (WAL) to replay missed commands after a crash.

## 4. Technical Validation
- [ ] **Temporal Alignment**: Use master timecodes (SMPTE) to prevent floating-point animation drift.
- [ ] **Version Gatekeeping**: Hard-fail handshake if engine minor versions are incompatible.
- [ ] **Metadata Sanity Check**: Periodic verification of cross-engine UUID maps to detect duplicates.

## 5. Error & Recovery
- [ ] **Error Propagation**: Broadcast "Engine B Failed" alerts back to "Engine A" UI.
- [ ] **Automatic Re-Reconciliation**: Trigger a full-state check after a heartbeat timeout recovery.

---
*Copyright (C) 2026 B-A-M-N*
