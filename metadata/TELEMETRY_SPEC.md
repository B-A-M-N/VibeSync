# VibeSync: Telemetry & Diagnostics Specification (v0.4)

This document defines the schema for forensic logging and performance monitoring within the cluster.

---

## 📈 1. Performance Metrics
Every heartbeat, engines report:
- `engine_load_fps`: Current rendering frame rate.
- `sync_latency_ms`: Delta between Command Reception and Main Thread execution.
- `memory_delta_kb`: Change in engine heap size since last sync.

## 📁 2. Standardized Event Types
The `events.jsonl` file records the following transitions:
- `TX_PREFLIGHT`: Start of atomic sync sequence.
- `TX_COMMIT`: Successful completion of state mutation.
- `TX_ROLLBACK`: Reversion due to hash mismatch or engine error.
- `TRUST_DEGRADE`: Monotonic trust score reduction.
- `QUARANTINE_LIFTED`: Trust score recovery (Manual only).

## 📊 3. Sync Success Auditing
The Orchestrator maintains a rolling window of:
- **Total Throughput**: Bytes synced per session.
- **Error Rate**: Ratio of `ROLLBACK` to `COMMIT` events.
- **Drift Frequency**: Number of times `state/get` hash mismatched expectation.

---
*VibeSync: Transparent Reality.*
