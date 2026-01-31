# VibeSync: Formal Threat Model

This document defines the security surface of the VibeSync cluster. Security is enforced via **Distrust-by-Default**.

---

## 🏗️ I. Identity & Authentication
- **Threat**: Spoofed adapters, session hijacking, or request replay.
- **Mitigation**: 
  - **Version Pinning**: Refusal on minor build drift.
  - **Trust Rotation**: Bootstrap secrets are exchanged for ephemeral session tokens during handshake.
  - **HMAC-SHA256 Signing**: Every request is signed with the session token; tampering or spoofing results in immediate rejection.
  - **Anti-Replay Timestamps**: Strict 5s TTL on all requests to prevent interception and replay.
  - **Mutual Auth**: Orchestrator challenges the Engine; Engine must return a hashed nonce.

## ⚖️ II. Authorization & Authority
- **Threat**: Privilege creep, unauthorized "while-I'm-here" mutations, or endpoint injection.
- **Mitigation**:
  - **Path Whitelisting**: Adapters reject any HTTP path not explicitly defined in the protocol.
  - **Intent Scoping**: Permissions are limited to the UUIDs declared in the `SubmitIntent` envelope.
  - **Signed Overrides**: Human "sudo" access requires a cryptographic key and a reason string.

## 📦 III. Input & Payload Safety
- **Threat**: Malicious FBX/Assets or prompt smuggling.
- **Mitigation**:
  - **Zero-Trust Ingestion**: All assets arrive in a hidden `.vibesync/tmp` sandbox for validation.
  - **Semantic Firewall**: AST-based auditing of all incoming payloads.

## ⏱️ IV. Execution & Isolation
- **Threat**: TOCTOU (Time-of-Check/Time-of-Use) attacks or process memory corruption.
- **Mitigation**:
  - **Read-Before-Write**: Mandatory telemetry check <500ms before mutation.
  - **Process Isolation**: Orchestrator runs in a separate memory space from the editors.

## 🔍 V. Integrity & Forensics
- **Threat**: Silent corruption or history rewriting.
- **Mitigation**:
  - **Hash Supremacy**: Binary equivalence checks for all cross-engine transfers.
  - **Append-Only WAL**: Forensic journal linked to Intent IDs.

## 🧠 VI. AI-Specific Defense
- **Threat**: Autonomy expansion or hallucinated compliance.
- **Mitigation**:
  - **Epistemic Refusal**: Formal "Unknowable" state for unverifiable requests.
  - **Tone Constraints**: Clinical, non-manipulative output policy.

---
*Copyright (C) 2026 B-A-M-N*
