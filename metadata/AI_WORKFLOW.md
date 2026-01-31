# VibeSync AI Workflow Instructions (Pain-Point Focused)

**Goal:** Ensure the bridge never corrupts scene state, preserves identities, survives crashes, and always reconciles Blender ↔ Unity deterministically.

---

## PHASE 0 — BOOT & ATTACH

1. Confirm both editors (Unity + Blender) are running and idle.
2. Detect editor states:
   * Unity: Play Mode, domain reload, scene loaded
   * Blender: File loaded, undo depth, auto-save state
3. Load persistent registries:
   * UUID registry for all synced objects
   * Sidecar backup files
   * Logs from previous sessions
4. Validate registry schemas. Abort if corrupted.
5. If previous session crashed → enter **Safe Mode**, disable automation. (**Edge Case 17: Crash Recovery**)

---

## PHASE 1 — INDEXING & UUID ENFORCEMENT

6. Scan all assets and objects in both editors:
   * Unity: prefabs, materials, meshes, objects
   * Blender: objects, meshes, armatures, materials, collections
7. Build mapping:
   ```
   UUID → Editor Object Reference
   UUID → External Mapping
   ```
8. Detect missing UUIDs:
   * Generate UUID
   * Assign to object/datablock
   * Persist immediately
9. Detect duplicate UUIDs: (**Edge Case 1: Duplicate UUIDs**)
   * Freeze automation temporarily
   * Regenerate for newest instance
   * Persist immediately

---

## PHASE 2 — SNAPSHOT & LOG CONSULTATION

10. Take a full snapshot of scene state:
    * Object transforms
    * Scene hierarchy
    * Mesh/material references
11. Timestamp snapshot. Abort if incomplete.
12. Consult logs for prior failures, crash triggers, incomplete transactions. (**Edge Case 10 & 15**)
13. **No operation executes without log acknowledgment.**

---

## PHASE 3 — EVENT & RATE CONTROL

14. Queue all incoming changes from Unity or Blender.
15. Debounce and throttle: (**Edge Case 7 & 16**)
    * Max 1 operation per tick per object
    * Limit simultaneous mutations
16. Maintain event depth counter; abort if threshold exceeded.

---

## PHASE 4 — SYNC OPERATIONS (UUID-FIRST)

17. Determine required sync operations:
    * Object creation, deletion, modification
    * Scene hierarchy changes
    * Asset updates (mesh/material/animation)
18. Resolve all targets by UUID only. (**Edge Case 8: Renames/Moves**)
19. Wrap each operation in transaction: (**Edge Case 10**)
    ```
    BEGIN → mutate → validate → COMMIT / ROLLBACK
    ```
20. Respect editor restrictions:
    * Unity: snapshot editor state pre-Play Mode; restore post-Play. (**Edge Case 5: Play Mode**)
    * Blender: avoid persistent Python handles across undo. (**Edge Case 4: Undo/Redo**)

---

## PHASE 5 — DUPLICATION & DELETION HANDLING

21. On duplication: (**Edge Case 11**)
    * Preserve source UUID
    * Regenerate duplicate UUID immediately
22. On deletion: (**Edge Case 2: Dangling References**)
    * Tombstone UUID
    * Archive mapping
    * Notify other editor

---

## PHASE 6 — HOT RELOAD / UNDO / FILE RELOAD

23. Detect hot reloads or undo/redo. (**Edge Case 6**)
24. Drop stale handles.
25. Rescan all assets/datablocks.
26. Rebuild `UUID → Object` mapping.
27. Rebind external references by UUID.
28. Quarantine ambiguous objects; log for review.

---

## PHASE 7 — LOGGING & FAILURE ESCALATION

29. Every operation must log:
    ```
    timestamp, process_id, operation, uuid(s), phase, outcome, error_code
    ```
30. Consult logs before each operation: (**Edge Case 13: Failed Hotfix / AI Retry Loops**)
    * Repeated failures → degrade / pause / notify
    * Promote log entries into operational memory
31. Maintain failure counters per UUID/operation.
32. On threshold: pause automation, preserve state, enter safe mode.

---

## PHASE 8 — EXTERNAL SYNC VERIFICATION

33. Validate external payloads: (**Edge Case 9 & 18**)
    * Project ID, schema version, UUID namespace
    * Reject stale or mismatched data
34. Lock state during sync; unlock after validation
35. Ensure order-independent resolution; buffer events if necessary

---

## PHASE 9 — PERFORMANCE & WATCHDOG

36. Limit: (**Edge Case 16**)
    * Max scans per tick
    * Max fix/repair operations per tick
37. Yield execution to prevent editor freeze
38. Resume monitoring only during idle

---

## PHASE 10 — USER OVERRIDE DETECTION

39. Detect direct user edits (selection, typing, manual modifications). (**Edge Case 14**)
40. Pause automation immediately
41. Resume only after editor is idle

---

## PHASE 11 — SAFE IDLE / MONITORING

42. Enter watch mode for:
    * Asset/datablock changes
    * Scene/collection modifications
    * External sync requests
43. Throttle operations
44. Maintain ephemeral handle discipline — resolve all references fresh by UUID
45. Snapshot & rollback on demand

---

## PHASE 12 — ABSOLUTE INVARIANTS FOR VIBESYNC AI

* UUID authority overrides semantic names.
* Logs are inputs, not just outputs.
* Transactions are atomic; no partial commits.
* Editor integrity prioritized over automation.
* Cross-system version/schema checks are mandatory.
* Automation knows when to pause or stop.
* No stale object handles ever persist beyond safe boundaries.

---

## ❄️ FREEZE-PROOF DISCIPLINE (STRICT)

1. **Never block the main thread**: All heartbeat, network I/O, and file transfers MUST be backgrounded or async.
2. **Timeouts & Watchdogs**: Wrap every external call in `try/catch` with strict 1000ms timeouts.
3. **Queue-Based Mutation**: Background threads must NEVER modify Unity/Blender objects directly. Use `MainThreadDispatcher` (Unity) or `bpy.app.timers` (Blender) via thread-safe queues.
4. **File Integrity**: Use **Atomic Swap** (write temp -> verify -> rename) for asset transfers to avoid engine file locks.
5. **Log Throttling**: Max 1 log/sec for repetitive tasks (heartbeat) to prevent console freeze.
6. **Conflict Checking**: Only push changes if the source version is newer than the target.

---

## 🛑 COMPREHENSIVE EDGE CASE CHECKLIST

1. **Duplicate UUIDs**: Detect immediately, regenerate for newest, persist.
2. **Dangling References**: Tombstone UUIDs, notify other editor, archive mapping.
3. **TOCTOU (Time of Check → Time of Use)**: Re-resolve UUID immediately before mutation.
4. **Blender Undo/Redo**: Treat all handles as ephemeral; resolve fresh on every access.
5. **Unity Play Mode**: Snapshot pre-Play, restore post-Play, block persistence during Play.
6. **Hot/Script Reload**: Drop all cached handles and rebuild UUID mapping.
7. **Simultaneous Changes**: Debounce events and track depth to prevent infinite loops.
8. **Asset Renames/Moves**: Resolve by UUID only; ignore name/order changes for sync.
9. **External Tool Mismatch**: Validate Project ID, Schema, and Namespace.
10. **Partial Operations**: Wrap in transactions; consult logs for incomplete steps.
11. **Duplication Events**: Preserve source UUID, regenerate for duplicate, rebind.
12. **Multi-Scene Collisions**: Namespace UUIDs by project/scene; block ambiguous merges.
13. **AI Retry Loops**: Track failure counters; escalate/pause after threshold.
14. **State Drift**: Detect user edits and pause automation immediately.
15. **Registry Desync**: Verify against editor state before mutation; quarantine inconsistencies.
16. **Performance Drift**: Yield execution; set max scans/repairs per tick.
17. **Crash Recovery**: Safe mode on boot; preserve snapshots for recovery.
18. **Schema Mismatch**: Enforce version checks; block operations on mismatch.

---
*Copyright (C) 2026 B-A-M-N*
