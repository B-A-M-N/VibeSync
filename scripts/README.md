# VibeSync Enterprise Infrastructure Scripts

This directory contains the hardening and validation tools for commercial-grade creation.

## 🛡️ Hardening
- `snap_commit.py`: Automated Iron Box snapshot tool.
- `log_rotate.py`: Forensic log truncation and rotation.
- `license_audit.py`: MIT/Apache 2.0 license compliance auditor.
- `pip_blender_orchestrator.py`: Isolated dependency management for Blender.

## 🔍 Validators & Sentinels
- `validators/dual_validator.py`: Cross-engine static analysis gate.
- `validators/setup_blender_intel.py`: Installs type stubs and pyright.
- `sentinels/upm_sentinel.py`: Unity package dependency auditor.
- `hardening/integrity_suite.py`: End-to-end Blender ↔ Unity verification.

## 🚀 Usage
Agents are required to run `dual_validator.py` and `snap_commit.py` as part of the Phase -3 and Phase 2 workflows. 

Pre-commit hooks are configured in `.pre-commit-config.yaml` to automate these checks.
