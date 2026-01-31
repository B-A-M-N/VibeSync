# VibeSync: Non-Goals & Doctrine

This document defines the intentional limitations of the VibeSync system. Silence on a feature does not imply intent to implement; if it contradicts these non-goals, it is prohibited.

---

## 🚫 1. No Creative Autonomy
VibeSync will **never** attempt to "improve," "clean up," or "optimize" a scene creatively. 
- It will not fix normals unless explicitly commanded.
- It will not "suggest" better material placements.
- It will not auto-generate geometry to "complete" a scene.
- It will not "hallucinate" details to fill gaps in the `global_id_map`.

## 🚫 2. No Silent Self-Healing
If state reality diverges from mathematical intent, the system will **never** "average it out" or silently fix the drift.
- Failure is a first-class outcome.
- The system would rather halt and require human intervention than proceed with 99.9% accuracy.
- **No automatic retry** without a journaled entry in the Write-Ahead Log (WAL).

## 🚫 3. No Guessing Artist Intent
The bridge will **never** infer intent from partial state. 
- If a command is ambiguous (e.g., "Delete active"), it is rejected.
- Context is always explicit (UUID-based).
- **Epistemic Refusal**: If the AI cannot prove an intent is safe, it must refuse to execute it.

## 🚫 4. No Implicit Scope Expansion
VibeSync will not automatically track or sync objects that have not been explicitly mapped to the `global_id_map`.
- We do not "mirror the world"; we "synchronize the intent."
- **No "Viral" Syncing**: Syncing a parent does not automatically sync children unless explicitly requested.

## 🚫 5. No Trust in Editor UI
The system will **never** rely on the editor's UI state or "Success" popups as proof of work.
- Only telemetry read-back and binary hashes constitute proof.
- If the editor UI says "Done" but the hash hasn't updated, the operation is a failure.
- **No Screen Scraping**: We do not read pixels or UI text elements.

## 🚫 6. No Latency Hiding
VibeSync will **never** predictively render or hide network/processing lag.
- We show the raw state of the system.
- If the bridge is slow, the user sees the slowness.
- **No "Rubber-Banding"**: We do not interpolate between conflicting states to smooth out jitter.

## 🚫 7. No Hidden Persistence
VibeSync will **never** store state outside the `.vibesync` directory.
- **No Registry Keys**: We do not touch the OS registry.
- **No Hidden Dotfiles**: We do not create files in user home directories (other than the designated workspace).
- **No Analytics**: We do not track user behavior for product improvement.

---
*Copyright (C) 2026 B-A-M-N*