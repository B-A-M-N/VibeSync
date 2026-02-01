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

## 🏎️ 3. Fast Path vs. Slow Path

| Category | Fast Path (Speculative) | Slow Path (Blocking) |
| :--- | :--- | :--- |
| **Criteria** | Cosmetic changes, Transforms, Material parameters, Scalar values. | Mesh topology, Parenting changes, Deletions, Prefab instancing. |
| **Behavior** | Instant Provisional Commit + Background Verify. | Sequential validation; Verification MUST pass before UI unlock. |
| **Verification** | Asynchronous (Deferred). | Synchronous (Immediate). |

---

## 📦 4. Intent Batching (Coalescing)
To reduce verification churn, small intents are coalesced:
- **Time-based**: Group changes within a 250ms window.
- **Semantic**: Group changes affecting the same set of UUIDs.
- **Atomic Batch**: A single `VERIFY` call covers the entire batch.

---

## 🚨 5. Conflict & Panic Handling
- **Speculation Halt**: If a **Panic Lock** is triggered (heartbeat failure, critical desync), all speculation stops immediately.
- **No Persistence**: Provisional state is NEVER saved to disk or persistent storage until `FINALIZED`.
- **Conflict Resolution**: If a user manually edits a provisionally-held object, the Orchestrator immediately aborts the speculation and snapshots the user's edit as the new source of truth.

---
*VibeSync: Speed without Compromise.*
