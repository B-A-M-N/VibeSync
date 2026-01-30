# VibeSync: Non-Goals & Doctrine

This document defines the intentional limitations of the VibeSync system. Silence on a feature does not imply intent to implement; if it contradicts these non-goals, it is prohibited.

---

## 🚫 1. No Creative Autonomy
VibeSync will **never** attempt to "improve," "clean up," or "optimize" a scene creatively. 
- It will not fix normals unless explicitly commanded.
- It will not "suggest" better material placements.
- It will not auto-generate geometry to "complete" a scene.

## 🚫 2. No Silent Self-Healing
If state reality diverges from mathematical intent, the system will **never** "average it out" or silently fix the drift.
- Failure is a first-class outcome.
- The system would rather halt and require human intervention than proceed with 99.9% accuracy.

## 🚫 3. No Guessing Artist Intent
The bridge will **never** infer intent from partial state. 
- If a command is ambiguous (e.g., "Delete active"), it is rejected.
- Context is always explicit (UUID-based).

## 🚫 4. No Implicit Scope Expansion
VibeSync will not automatically track or sync objects that have not been explicitly mapped to the `global_id_map`.
- We do not "mirror the world"; we "synchronize the intent."

## 🚫 5. No Trust in Editor UI
The system will **never** rely on the editor's UI state or "Success" popups as proof of work.
- Only telemetry read-back and binary hashes constitute proof.
- If the editor UI says "Done" but the hash hasn't updated, the operation is a failure.
