# VibeSync: Adapter Contract & Extensibility (v0.4)

Adapters are the "dumb limbs" of the VibeSync system. To be a legal participant in the cluster, an adapter must uphold these strict invariants.

---

## 🏗️ 1. The Dumb Limb Principle
Adapters **MUST NOT** possess autonomous logic, local state persistence (outside the session), or the authority to reject valid Orchestrator commands. Intelligence lives exclusively in the **Go Brain**.

## 🛡️ 2. Security Invariants
- **No Self-Authorization**: Adapters cannot generate their own tokens. They must receive them via a handshake.
- **Strict Whitelisting**: Adapters must only expose the minimum necessary API surface to the Orchestrator.
- **Pattern Rejection**: Adapters must self-audit for dangerous patterns (Reflection, Shell, External Network) as a second layer of defense.

## ⏱️ 3. State & Causality
- **Monotonic Reporting**: Every state report MUST include the last processed `MonotonicID`.
- **Atomic Execution**: Adapters must support "Preflight" checks where they report if a mutation *would* succeed without actually applying it.
- **Fail-Fast**: If an internal engine error occurs, the adapter must report it immediately and enter a `STOPPED` state until the Orchestrator intervenes.

## 🔌 4. Extensibility & Schema

VibeSync supports custom extensions via the **Payload Injection** system.



- **Custom Asset Types**: Users can register new UUID prefixes for proprietary data (e.g., custom physics volumes). These must implement a `to_intermediate_json()` method.

- **Hook Registry**: Adapters can expose `pre_sync` and `post_sync` hooks for local engine side-effects (e.g., updating a local UI element). These hooks are **non-blocking** and cannot reject the mutation.



---



## 🏗️ 5. Versioning & Evolution

- **Contract Pinning**: Every adapter must report its **Protocol Version** (e.g., `v0.4.2`). 

- **Breaking Changes**: Any change to the JSON schema results in a major version bump. Orchestrator will refuse connection to an adapter more than one minor version behind.

- **Schema Migration**: VibeSync does not support live schema migration. All engines must restart to adopt a new protocol version.



---

## 🛠️ 6. API Schema Specification

All adapters must implement these JSON endpoints via HTTP POST/GET.

### A. `/handshake` (POST)
**Request:**
```json
{ "version": "v0.4.0", "new_token": "uuid-string", "challenge": "nonce-string" }
```
**Response:**
```json
{ 
  "status": "OK", 
  "engine_version": "2022.3.x", 
  "capabilities": ["mesh", "transform"], 
  "unit_settings": { "system": "Metric", "scale_length": 1.0 },
  "response": "VIBE_HASH_[challenge]" 
}
```

### B. `/health` (GET)
**Response:**
```json
{ "status": "ok|busy", "generation": 12 }
```

### C. `/transform/set` (POST)
**Request:**
```json
{ "id": "uuid", "transform": { "pos": [x,y,z], "rot": [x,y,z,w], "sca": [x,y,z] }, "generation": 12 }
```

---

*VibeSync: Modular Reality.*
