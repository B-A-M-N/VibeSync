# VibeSync: Conflict Resolution Policy (v0.4)

This document defines the formal policy for resolving simultaneous conflicting intents across the VibeSync cluster.

---

## ⚖️ 1. The Golden Rule
> **"If resolving a conflict requires guessing user intent, the system must stop."**

VibeSync prioritizes **Integrity over Automation**. We never "guess" which user was right; we provide the mechanical floor to prevent state corruption and escalate to the human when semantics are ambiguous.

---

## 🔍 2. Conflict Classification Matrix

Conflicts are detected during the **Phase 3 (Batching)** and **Phase 4 (Execution)** windows by the Go Orchestrator.

| Conflict Type | Auto-resolve | Strategy | Action |
| :--- | :--- | :--- | :--- |
| **Cosmetic vs Cosmetic** (Same Property) | ✅ | Deterministic Override | Last-Causal-Writer-Wins (Monotonic ID tie-break). |
| **Cosmetic vs Cosmetic** (Disjoint) | ✅ | Commutable Merge | Reorder and apply both. Single hash verification. |
| **Structural vs Cosmetic** | ⚠️ | Structural Precedence | Structural wins; Cosmetic is `ROLLED_BACK`. |
| **Structural vs Structural** | ❌ | Mandatory Quarantine | Both intents `QUARANTINED`. Manual resolution required. |
| **Destructive vs Any** | ❌ | Mandatory Quarantine | All overlapping intents `QUARANTINED`. Snapshot retained. |
| **Human vs AI** (Any Overlap) | 👑 | Human Supremacy | Human intent immediately trumps AI. AI intent transitions to `WAIT`. |
| **Identity Mismatch** | ❌ | Panic / Quarantine | Immediate engine hierarchy lock. |

---

## 🛠️ 3. Resolution Strategies

### A. Deterministic Override (Cosmetic Only)
For non-structural changes (Color, Roughness, Transforms), the Orchestrator uses the **Monotonic Intent ID** as the authoritative tie-breaker.
- The intent with the **higher ID** is promoted to `FINAL`.
- The intent with the **lower ID** is `ROLLED_BACK`.

### B. Commutable Merge
If two cosmetic intents target the same UUID but different property sets (e.g., Position vs Material Color), they are treated as non-conflicting.
- Both are applied in a single atomic batch.
- Verified via a single post-mutation hash.

### C. Mandatory Quarantine (Structural/Destructive)
If two structural changes (e.g., competing parentage) or a destructive change (Delete) overlap in the same batch:
1. **Freeze**: Both intents are marked as `QUARANTINED`.
2. **Revert**: Engines are ordered to revert to the last authoritative `FINAL` hash.
3. **Escalate**: A `CONFLICT_EVENT` is broadcast to the human with a forensic dump of both intents.

### D. Human Supremacy (The "Anti-Fighting" Lock)
If a human is actively manipulating an object (lock_type: `HUMAN_ACTIVE`):
1. **Veto**: Any AI intent touching the locked UUID is automatically deferred.
2. **Wait State**: AI intents transition to `WAIT_HUMAN_LOCK`.
3. **Rollback**: If an AI intent was provisionally applied during the lock window, it is instantly `ROLLED_BACK`.
4. **Adaptive Deferral**: The AI receives a `WAIT_HUMAN_LOCK` status and must exponentially back-off until the lock is released.

---

## 📜 4. Conflict Metadata (WAL)

Every resolved or quarantined conflict MUST include this metadata in its WAL entry:

```json
"conflict": {
  "type": "COSMETIC|STRUCTURAL|DESTRUCTIVE",
  "resolution": "DETERMINISTIC_OVERRIDE|MERGE|QUARANTINE",
  "winner_intent_id": "uint64|null",
  "reason": "PROPERTY_OVERLAP|GRAPH_INVALIDITY|DESTRUCTIVE_DOMINANCE"
}
```

---
*VibeSync: Engineering Truth.*
*Copyright (C) 2026 B-A-M-N*
