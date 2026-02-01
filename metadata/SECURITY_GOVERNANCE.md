# VibeSync: Security Governance & "The Iron Box"

We treat your creative project like a high-security vault. VibeSync assumes that both the engines and the AI operators are potentially hostile or unreliable.

---

## 🛡️ Security Capabilities

### 1. Iron Handshake (Zero-Trust)
Communication is secured via **Token Rotation** (keys change every session) and **HMAC-SHA256 Request Signing**. This provides cryptographic proof that commands were issued by the authoritative Orchestrator and were not tampered with in transit.

### 2. Atomic Sync (Transactional Pipeline)
A formal **Snapshot → Preflight → Commit** pipeline ensures that state changes are all-or-nothing. If a sync fails in Unity, the source in Blender is automatically rolled back to prevent "split-brain" divergence.

### 3. Semantic Firewall (AST Auditing)
**AST-based auditing** (scanning the abstract syntax tree of code) blocks dangerous payloads like `os.system`, `Reflection`, or unauthorized network calls before they ever reach the execution engine.

### 4. Deadman Switch (Circuit Breaker)
A **5000ms Heartbeat monitor** tracks engine liveness. If any engine freezes, deadlocks, or becomes unresponsive, the Orchestrator triggers an immediate **Global PANIC** lock, freezing all hierarchies to prevent data corruption.

### 5. Adversarial Pre-flight Layer
A specialized pre-boot auditor (`scripts/preflight.py`) detects and resolves environmental blocks before orchestration begins. It automatically kills zombie engine processes, releases port conflicts, and validates Unity compilation status.

### 6. OS & Process Isolation
- **Docker Isolation**: Minimal Alpine-based containerization for the Go Orchestrator ensures environment isolation.
- **OS Hardening**: Standard security tools like `ufw` (firewall) and `AppArmor` are used to restrict the bridge's system access.

---

## 🔒 7. Triple-Lock Invariance System (Absolute Governance)

To achieve absolute invariance and prevent "Optimistic Bypass," VibeSync implements three nested safety locks as defined in the **[Absolute Invariance Contract](INVARIANCE_CONTRACT.md)**:

### 1. Mechanical Invariance (The Ground Truth Lock)
The `execute_governed_mutation` tool doesn't just send a command; it automatically waits for the engine response and performs an independent state read-back before returning success to the AI.

### 2. Contextual Invariance (The Forensic Feed)
Every tool response is "Force-Fed" with a **Forensic Report**, including the last 3 lines of the WAL, engine health flags, and generation counters. This ensures errors are always in the AI's immediate context.

### 3. Semantic Invariance (The Proof of Work)
The `commit_atomic_operation` tool is "Hard Gated." It will mechanically refuse to execute unless the AI provides a valid `ProofOfWork` string summarizing the evidence found in the Forensic Report.

---

### 🛡️ 8. Enterprise Stability & Hallucination Prevention
To achieve commercial-grade reliability, VibeSync employs:
- **Static Pre-Flight**: Every AI-generated script is audited via **Pyright** (Blender) and **Roslyn** (Unity) before execution.
- **API Intelligence**: `fake-bpy-module` and `ApiMetadataMapper` provide exact tool signatures to the AI, eliminating "hallucinated" attributes.
- **Automated Integrity Testing**: Headless tests run via `pytest-blender` and Unity UTF verify mutations inside the actual application processes.

### 🛡️ 9. Real-Time State Sentinels
- **Unity Compilation Sentinel**: Mutate-before-verify loops are blocked if `EditorApplication.isCompiling` is true.
- **Blender Depsgraph Sentinel**: Mutations wait for `evaluated_depsgraph_get()` to stabilize.
- **Global Traffic Controller**: The Go Orchestrator locks the "Creation Perimeter" on both apps during sync to prevent race conditions.

### 🛡️ 10. Infrastructure Hardening
- **Auto-Snapshot System**: `snap_commit.py` creates a safety restore point in `.git_safety` before every high-risk mutation.
- **Log Rotation**: Automated truncation of forensic logs keeps AI context within stable token limits.
- **License Compliance**: Continuous auditing ensures all dependencies (MIT/Apache 2.0) are ready for commercial redistribution.

---

## 🔗 Security Resources
- [**Security Gate Script**](../security_gate.py): The pre-execution auditor.
- [**Iron Box Constraints**](../AI_ENGINEERING_CONSTRAINTS.md): The formal rules governing all codebase changes.
- [**Threat Model**](THREAT_MODEL.md): Detailed analysis of potential attack vectors.

---
*Copyright (C) 2026 B-A-M-N*