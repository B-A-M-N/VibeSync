# Contributing to VibeSync

Thank you for your interest in improving VibeSync! Because this is a **Zero-Trust, High-Fidelity Orchestrator**, we have a very strict contribution policy.

---

## 🛡️ The "Review-First" Policy
All pull requests must undergo a **Mechanical Review** by the Orchestrator. We prioritize state integrity and security over feature speed.

### 1. **Zero-Trust Adherence**
Every new tool or modification **MUST** adhere to `AI_ENGINEERING_CONSTRAINTS.md`. Any use of dynamic execution, unauthorized network calls, or non-deterministic logic will be rejected.

### 2. **Verification Mandate**
Every mutation tool must include:
- A Class A or Class B verification loop (as defined in `ALLOWED_OPERATIONS.md`).
- A post-mutation telemetry check to prove reality matches intent.

### 3. **Workspace Hygiene**
Do not add files to the root directory. Store diagnostic or temporary artifacts in `.vibesync/`.

### 4. **ISA Mapping**
All new features must be expressed as a numbered entry in the **Bridge ISA**. We do not accept "helper functions" that bypass the tool registry.

### 5. **Operational Ergonomics**
Every new tool's event log MUST include a **NextStep** string providing clear human guidance in case of failure.

---

## 🛠️ Local Development
1.  **Orchestrator**: `cd mcp-server && go build`
2.  **Testing**: Use the `initiate_handshake` tool to verify your adapter connectivity before submitting.

---

## ⚖️ Contributor License Agreement (CLA)
By submitting a pull request, you agree to the terms in the `LICENSE` file, granting the Author a perpetual, irrevocable license to use and sublicense your contributions in both the open-source and commercial versions of the software. **You explicitly waive any requirement for the Author to track or notify you of future use of your contributions.**
