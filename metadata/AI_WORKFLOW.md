# VibeSync AI Workflow Instructions (Pain-Point Focused)

**Goal:** Ensure the bridge never corrupts scene state, preserves identities, survives crashes, and always reconciles Blender ↔ Unity deterministically.

---

## PHASE -3 — STATIC PRE-FLIGHT (NEW CODE ONLY)

If the AI has generated new code or scripts for an engine:

1.  **Blender**: Run `pyright` against the script using `fake-bpy-module` stubs. Fix all errors.
2.  **Unity**: Audit syntax via Roslyn-style checks and verify against `metadata/unity_api_map.json`.
3.  **Dual-Validator**: Execute `python3 scripts/validators/dual_validator.py` to confirm environmental stability.

---

## PHASE -2 — FORENSIC TRIGGER ANALYSIS (TRIGGERED)

If a tool output or engine response contains an error pattern (e.g., `NullReferenceException`, `ECONNREFUSED`):

1.  **Consult Mapping**: Resolve the trigger to its required log files using `metadata/LOG_TROUBLESHOOTING_MAPPING.md`.
2.  **Verify Hash**: Check if the target log has changed since the last read. If unchanged, skip to step 4.
3.  **Targeted Log Read**: Read the relevant sections of the forensic logs.
4.  **Inject Forensic Report**: Create a structured summary of the error context and add it to the active context.
5.  **Refusal/Recovery**: If the log indicates a terminal state (e.g., Kernel Deadlock), use `epistemic_refusal` or `trigger_panic`. DO NOT guess a fix.

---

## PHASE -1 — ADVERSARIAL PRE-FLIGHT (MANDATORY)



Before any turn involving connection or mutation, the AI **MUST** verify the environment is clean and reconcile local context.



1.  **Execute Pre-flight**: Run `python3 scripts/preflight.py`.

2.  **Path Discovery Gate**: Before executing any mutation in a subdirectory, you MUST first read the `.gemini` or `README.md` file within that specific directory to reconcile local invariants. 

3.  **Git Isolation Check**: Verify that `.git_safety` is initialized for project snapshots. Logic/Code changes must be targeted to the primary `.git`.

4.  **Performance Mode (High-Frequency)**: For high-frequency data (transforms, playback), use the "Fast Pipe" tools if available. These bypass heavy HMAC/Audit checks in exchange for session-level trust.

3.  **Verify Report**:

 
    - If `safe_to_proceed` is `false`, stop and report the specific issues (e.g., Unity compile errors) to the human.
    - If zombie processes were killed or ports released, acknowledge the cleanup in the rationale.
3.  **Clear Persistence**: If `get_diagnostics` continues to report `DESYNC` after cleanup, ask for human permission to purge `.vibesync/state.json`.

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

10. **Iron Box Snapshot**: Before any destructive mutation, take a local snapshot using the `.git_safety` repo (see `metadata/IRON_BOX_SAVE_GAME.md`).
    ```bash
    git --git-dir=.git_safety --work-tree=. add . && git --git-dir=.git_safety --work-tree=. commit -m "[AI_SYNC] Pre-mutation snapshot"
    ```
11. **Scene State Snapshot**: Take a full snapshot of scene state:
    * Object transforms
    * Scene hierarchy
    * Mesh/material references
12. Timestamp snapshot. Abort if incomplete.
13. Consult logs for prior failures, crash triggers, incomplete transactions. (**Edge Case 10 & 15**)
14. **No operation executes without log acknowledgment and safety snapshot.**

---

## PHASE 3 — EVENT & RATE CONTROL

14. Queue all incoming changes from Unity or Blender.
15. **Conflict Detection**: Check for simultaneous overlapping intents. Apply `metadata/CONFLICT_RESOLUTION_POLICY.md`.
    - **Cosmetic Override**: Promote intent with higher Monotonic ID.
    - **Merge**: Apply non-overlapping cosmetic changes.
    - **Quarantine**: If Structural/Destructive overlap, pause and emit `CONFLICT_EVENT`.
16. Debounce and throttle: (**Edge Case 7 & 16**)
    * Max 1 operation per tick per object
    * Limit simultaneous mutations
17. Maintain event depth counter; abort if threshold exceeded.

---

## PHASE 4 — SYNC OPERATIONS (UUID-FIRST)

18. **Real-Time Sentinel Check**: Verify engines are not busy (Compiling, Updating, Depsgraph recalculating).
19. **Path Selection**: Categorize mutation as **Fast Path** (Cosmetic/Transform) or **Slow Path** (Structural).
20. **Conflict Validation**: Ensure the intent does not violate existing provisional locks or the Conflict Resolution Policy.
21. **Graph Validation**: For parenting changes, compute **Ancestor Closure**. Ensure the new parent is not a descendant of the target.
22. **Closure Computation**: For destructive operations (Deletes), compute the **Delete Closure** (children, refs, constraints). Intent MUST target the full closure.
23. Resolve all objects **by UUID and Prefab Depth**.
24. Wrap each operation in a transaction:
    ```
    BEGIN → mutate → PROVISIONAL COMMIT (Fast Path) / BLOCK (Slow Path) → validate → COMMIT / ROLLBACK
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

39. **Human Active Lock**: Detect direct user edits (selection, typing, manual modifications). (**Edge Case 14**)
40. **Veto AI Actions**: Immediately transition all overlapping AI intents to `WAIT_HUMAN_LOCK`.
41. **Provisional Rollback**: If the human edit overlaps with a provisional AI commit, issue an authoritative `ROLLBACK` for the AI intent.
42. **Resume Policy**: AI may only resume sync after the `HUMAN_ACTIVE` lock is released and state hash is re-verified.

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
