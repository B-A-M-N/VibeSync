# VibeSync Master Procedural Flow (v0.4)

**Goal:** Ensure deterministic, crash-proof, live syncing of Blender ↔ Unity scenes with full AI awareness.

---

## PHASE 0 — BRIDGE BOOT / ATTACH

1. Initialize bridge process.
2. Detect editor states:
   * Unity: Play Mode, domain reload, scene loaded
   * Blender: File loaded, undo depth, auto-save state
3. Load persistent state:
   * UUID registries from Unity & Blender
   * Sidecar files (backup)
   * Logs from prior sessions
4. Validate schemas. Abort if corrupted.
5. Enter **safe mode** if previous session crashed.

---

## PHASE 1 — INITIAL ASSET / SCENE INDEXING

6. Scan all assets and objects in both editors:
   * Unity: prefabs, materials, meshes, objects, collections
   * Blender: objects, meshes, armatures, materials, collections, actions
7. Build temporary in-memory map:
   ```
   UUID → Object Reference
   UUID → External Mapping
   ```
8. Validate UUID uniqueness; regenerate missing UUIDs. Persist immediately.

---

## PHASE 2 — SNAPSHOT & LOG CONSULTATION

9. Snapshot both editor states:
   * Object transforms, hierarchy, mesh data references
   * Scene collections
   * External mapping references
10. Timestamp snapshot. Abort if incomplete.
11. Consult logs for:
    * Prior failures
    * Crash triggers
    * Incomplete operations
12. No mutation occurs without log acknowledgment.

---

## PHASE 3 — EVENT & RATE CONTROL

13. Queue all incoming changes from either editor.
14. Debounce/throttle:
    * Max 1 mutation per tick per object
    * Rate-limit heavy operations
15. Monitor event depth to prevent infinite loops. Abort if exceeded.

---

## PHASE 4 — SYNC OPERATION EXECUTION

16. Determine required sync operations:
    * Object creation, deletion, modification
    * Scene hierarchy changes
    * Asset updates (mesh, material, animation)
17. Resolve all objects **by UUID only**.
18. Wrap each operation in a transaction:
    ```
    BEGIN → mutate → validate → COMMIT / ROLLBACK
    ```
19. Respect editor-specific restrictions:
    * Unity: snapshot pre-Play, restore post-Play
    * Blender: avoid Python handle persistence across undo

---

## PHASE 5 — DUPLICATION & DELETION HANDLING

20. Detect duplication events:
    * Preserve original UUID
    * Regenerate duplicate UUID immediately
21. Detect deletion events:
    * Tombstone UUID
    * Archive mapping
    * Notify external system

---

## PHASE 6 — HOT RELOAD / UNDO / FILE RELOAD

22. Detect editor reloads or undo/redo events.
23. Drop stale cached handles.
24. Re-scan all assets / datablocks.
25. Rebuild `UUID → Object` mapping.
26. Rebind external references (Unity ↔ Blender) by UUID.
27. Validate state; quarantine ambiguous objects.

---

## PHASE 7 — LOGGING & FAILURE ESCALATION

28. Every operation must log:
    ```
    timestamp, process_id, operation, uuid(s), phase, outcome, error_code
    ```
29. Consult logs pre-operation:
    * Repeated failures → degrade / pause / notify user
    * Promote log entries into operational memory
30. Maintain failure counters per UUID / operation.
31. On thresholds:
    * Pause automation
    * Enter safe mode
    * Preserve state

---

## PHASE 8 — EXTERNAL SYNC VERIFICATION

32. Validate incoming/outgoing payloads:
    * Project ID, schema version, UUID namespace
    * Reject stale or mismatched data
33. Lock state during sync, unlock after validation
34. Ensure operations are order-independent; buffer events as necessary

---

## PHASE 9 — PERFORMANCE & WATCHDOG

35. Set operation limits:
    * Max scans per tick
    * Max repair/fix operations per tick
36. Yield execution to prevent editor freeze
37. Resume monitoring only in idle states

---

## PHASE 10 — USER OVERRIDE & EDIT PROTECTION

38. Detect direct user edits (selection, typing, manual object changes)
39. Pause automation immediately
40. Resume only after idle

---

## PHASE 11 — SAFE IDLE / MONITORING

41. Enter watch mode for:
    * Asset/datablock changes
    * Scene/collection modifications
    * External sync requests
42. Throttle all operations
43. Maintain ephemeral handle discipline — always resolve by UUID
44. Ready to snapshot and rollback on demand

---

## PHASE 12 — ABSOLUTE INVARIANTS

* UUID authority is always the source of truth
* Logs are consulted **before every mutation**
* Transactions are atomic — no partial state commits
* Editor integrity is prioritized over AI/bridge automation
* Cross-system versioning/schema checks are mandatory
* Automation knows when to pause or stop
* No Python/Unity object references are persisted beyond safe boundaries

---
*Copyright (C) 2026 B-A-M-N*
