# Security Policy

VibeSync is a **Zero-Trust Orchestrator**. We take security and state integrity extremely seriously.

## Supported Versions
Only the latest version of VibeSync (currently v0.4 "Crowbar") is supported for security updates.

## Reporting a Vulnerability
If you discover a security vulnerability (e.g., a way to bypass the semantic firewall or an unauthorized state mutation), please **do not open a public issue.**

Instead, please send a detailed report to: **[Your Contact Info / Security Email]**.

### What to include:
- A description of the vulnerability.
- A proof-of-concept (PoC) script if possible.
- The version of VibeSync and the adapters (Unity/Blender) you are using.

We will acknowledge your report within 48 hours and work with you to resolve the issue before a public disclosure.

## Security Mandates
All contributors must adhere to the rules in `AI_ENGINEERING_CONSTRAINTS.md`, which include:
- No dynamic code execution (`exec`, `eval`).
- Mandatory payload auditing.
- Token rotation for every session.
