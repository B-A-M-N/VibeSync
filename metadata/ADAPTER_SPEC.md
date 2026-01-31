# VibeSync: Engine Adapter Specification (v0.4)

This document defines the technical contract any engine (Adapter) must satisfy to join the VibeSync cluster.

---

## 📡 Connectivity
- **Communication**: HTTP/1.1 JSON REST.
- **Port Mapping**:
  - `8085`: Unity
  - `22000`: Blender
  - `30000+`: Future Engines (Maya, Unreal, etc.)

---

## 🔒 Security Headers
Every request from the Orchestrator will include:
| Header | Description |
| :--- | :--- |
| `X-Vibe-Token` | The current session token (rotated during handshake). |
| `X-Vibe-Session` | The unique UUID for the current orchestrator session. |
| `X-Vibe-Generation` | Monotonic counter to detect engine reloads/drift. |
| `X-Vibe-Transaction` | The current atomic Transaction ID (`tid`). |

---

## 🛠️ Required Endpoints

### 1. **Lifecycle & Health**
- `GET /health`: Returns `{"status": "ok", "generation": X}`. Allowed without token.
- `POST /handshake`: Expects `{"version": "...", "new_token": "...", "challenge": "..."}`.
- `POST /panic`: Instant hierarchy lock. Rejects all subsequent mutations until restart.
- `GET /metrics`: Returns memory and engine load data.

### 2. **Atomic Sync**
- `POST /preflight/run`: (Source Only) Generates a scene/asset hash and checks resource limits.
- `POST /export`: (Source Only) Saves asset to a sandboxed `.vibesync/tmp` directory.
- `POST /import`: (Target Only) Loads asset into a temporary sandbox.
- `POST /validate`: (Target Only) Generates post-import hash for comparison.
- `POST /commit`: Finalizes mutation (moves asset to project folder).
- `POST /rollback`: Purges sandbox and reverts to last snapshot.

### 3. **Mutations**
- `POST /transform/set`: Sets position, rotation, and scale.
- `POST /material/update`: Updates shader properties.
- `POST /object/lock`: Enables/disables hierarchy-aware locking.
- `GET /state/get`: Returns a scene-wide or object-specific hash for verification.

---

## 🧵 Threading Model (Crucial)
Adapters operating in single-threaded engines (Unity, Blender) **MUST** use a **Split Architecture**:
1.  **Listener Thread**: Handles incoming HTTP requests immediately to avoid network timeouts.
2.  **Main Thread Queue**: Commands are enqueued and executed on the engine's main update loop to safely interact with the engine API.

## ⚖️ The Law of Reality
Every mutation MUST be followed by an independent `/state/get` call from the Orchestrator to verify the result. Adapters must ensure that `/state/get` returns the *actual* engine state, not the last commanded state.

---

## 📐 4. Coordinate & Unit Translation (The "Two Gods" Protocol)
To ensure seamless sync, adapters **MUST** perform the following basis conversions:

| Property | Blender Basis (Right-Hand, Z-Up) | Unity Basis (Left-Hand, Y-Up) | Conversion Formula |
| :--- | :--- | :--- | :--- |
| **Up Vector** | `[0, 0, 1]` | `[0, 1, 0]` | `Y_unity = Z_blender`, `Z_unity = Y_blender` |
| **Scale** | 1.0 (Metric) | 1.0 (Metric) | Ensure 1:1 unit matching. |
| **Rotation** | Euler / Quaternion | Quaternion | Invert `X` and `W` components for hand-flip. |

---

## 🔒 5. Privacy & "Human-Only" Marking
Adapters must support **AI-Invisibility** via the following mechanisms:

1.  **UUID Prefixing**: Any object with a UUID starting with `h-` (e.g., `h-1234...`) is **STRICTLY FORBIDDEN** from being sent to the AI. The adapter must redact these from telemetry.
2.  **Tagging**:
    - **Unity**: Objects tagged with `HumanOnly` are ignored by the bridge.
    - **Blender**: Objects in a Collection named `HUMAN_ONLY` are ignored by the bridge.

---
*Copyright (C) 2026 B-A-M-N*
