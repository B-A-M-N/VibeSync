# VibeSync: Speculative Commit & Deferred Finality Protocol

To eliminate user-visible latency while maintaining 100% state integrity, VibeSync employs a **Speculative Commit** model.

---

## ⚡ 1. The Core Principle
**"Never remove verification. Remove blocking."**
The system acts as if a mutation succeeded (Provisional State) while verifying the result in the background (Deferred Finality).

---

## 🔄 2. The Speculative Lifecycle

### Phase A: Provisional Commit (Instant)
1. **Mutation**: The Orchestrator issues a command to the engine.
2. **Overlay**: The engine applies the change to a **Provisional Overlay** (e.g., Unity Undo stack or a temporary state buffer).
3. **Response**: The engine returns a `PROVISIONAL_OK` immediately after the command is queued/applied to the overlay.
4. **UI**: The change is visible in the editor without a blocking wait.

### Phase B: Deferred Verification (Asynchronous)
1. **Background Hash**: The engine calculates the new state hash without blocking the main thread.
2. **Reporting**: The engine sends the hash back to the Orchestrator with the associated `Monotonic Intent ID`.
3. **Matching**: The Orchestrator compares the reported hash against the `Expected Hash`.

### Phase C: Finalization or Rollback
- **Finalize**: If hashes match, the Orchestrator marks the WAL entry as `FINALIZED`. The engine clears the provisional tag.
- **Rollback**: If hashes mismatch or a timeout occurs, the Orchestrator issues a `ROLLBACK`. The engine reverts the overlay to the last authoritative state.

---

## 📜 3. WAL Metadata Schema (Production-Grade)

All entries in the Write-Ahead Log (WAL) MUST adhere to this structure to ensure formal auditability of speculative states.

```json
{
  "intent_id": "monotonic:uint64",
  "parent_hash": "sha256",
  "entry_hash": "sha256",
  "timestamp": "orchestrator_time_ns",
  "engine": "unity|blender",
  "actor": "ai|human|system",
  "scope": {
    "uuids": ["uuid-v4", "..."],
    "intent_class": "cosmetic|structural|destructive"
  },
  "phase": "PROVISIONAL|FINAL|ROLLED_BACK|QUARANTINED",
  "verification": {
    "expected_hash": "sha256",
    "observed_hash": "sha256|null",
    "epsilon": 1e-5,
    "verified_at": "timestamp|null"
  },
  "rollback": {
    "undo_token": "engine_specific",
    "snapshot_ref": "git_safety_ref|null"
  }
}
```

### State Machine Invariants
- **PROVISIONAL** entries may not mutate the `parent_hash` of the authoritative chain.
- **Only FINAL** entries advance "Reality" in the global state.
- **ROLLED_BACK** entries remain in the log as immutable evidence of failure.
- **QUARANTINED** entries halt all causality in the affected engine until manual intervention.

---

## 🔁 4. Rollback Protocol

Rollback is issued **exclusively by the Orchestrator**. Engines and AI agents are prohibited from self-initiating rollback.

### Trigger Conditions
- Hash mismatch beyond epsilon threshold.
- Missing UUID during identity parity check.
- Exception mapped to `metadata/LOG_TROUBLESHOOTING_MAPPING.md`.
- Heartbeat drift during the provisional window.

### Message Format
```json
{
  "type": "ROLLBACK",
  "intent_id": "uint64",
  "reason": "HASH_MISMATCH|IDENTITY_BREAK|ENGINE_FAULT",
  "severity": "SOFT|HARD"
}
```

---

## 🏎️ 5. Hard Classification Rules (Operation Semantics)

AI agents NEVER decide the classification. It is derived mechanically from the operation.

| Category | Fast Path (Cosmetic) | Slow Path (Structural) | Guarded Path (Destructive) |
| :--- | :--- | :--- | :--- |
| **Criteria** | No object graph change. No UUID creation/removal. No reference alteration. | Alters hierarchy, references, topology, or instancing. | Deletes UUIDs. Overwrites assets. Irreversible without snapshot. |
| **Validation** | Epsilon Check (`1e-5`). | **DAG Check**: No cycles. **Depth Check**: Prefab nesting parity. | **Closure Check**: All deps included. Snapshot required. |
| **Examples** | Transforms, Material params, Shader edits, Visibility toggles. | Parenting, Prefab instantiation, Modifier stack changes. | Object deletion, Mesh overwrite, Asset replacement. |
| **Commit Mode** | AUTO | AUTO | AUTO + SNAPSHOT |
| **Batch Window** | 250–500ms | 0–100ms | DISABLED |
| **Verification** | Asynchronous (Deferred) | Synchronous (Immediate) | Snapshot-Gated |

---

## 📦 6. Intent Batching & Conflict Detection
To reduce verification churn and ensure semantic integrity, micro-intents are coalesced before finalization.
- **Time-based Coalescing**: Group changes within a 250ms window.
- **Semantic Coalescing**: Group changes affecting the same set of UUIDs.
- **Conflict Enforcement**: During batching, the Orchestrator applies the `metadata/CONFLICT_RESOLUTION_POLICY.md`.
- **Atomic Batch Verification**: A single `VERIFY` call covers the entire batch. If a Structural/Destructive conflict is detected, the affected intents transition to `QUARANTINED`.

---

## 🚨 7. Conflict & Panic Handling
- **Speculation Halt**: If a **Panic Lock** is triggered (heartbeat failure, critical desync), all speculation stops immediately.
- **Human-Active Lock**: If a human actively manipulates an object (lock_type: `HUMAN_ACTIVE`), all speculative AI intents for that UUID enter a `WAIT_HUMAN_LOCK` state.
- **No Persistence**: Provisional state is NEVER saved to disk or persistent storage until `FINALIZED`.
- **Conflict Resolution**: If a user manually edits a provisionally-held object, the Orchestrator immediately aborts the speculation, rolls back the AI intent, and snapshots the user's edit as the new source of truth.

---
*VibeSync: Speed without Compromise.*