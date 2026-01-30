# VibeSync: Architectural Blueprint (v0.4)

This diagram illustrates the layered orchestration, governance, and verification stack.

```mermaid
graph TD
    %% External Layer
    AI[External AI / Human] -- Submit Intent --> INT[Intent Envelope]
    
    %% Orchestrator Layer
    subgraph "🛡️ Go Orchestrator (The Brain)"
        INT --> VAL[Semantic Firewall / validate_intent]
        VAL --> CONF[Confidence Scoring]
        CONF -- "Allow" --> TX[Transaction Manager]
        CONF -- "Low Confidence" --> HUM[Human-in-the-Loop Gate]
        
        TX --> REF[Referee Loop / verifyEngineState]
        TX --> WAL[Append-Only WAL]
        TX --> EB[Structured Event Bus]
        
        subgraph "Observability Stack"
            MET[Metrics Engine]
            TEL[Real-time Telemetry]
        end
    end
    
    %% Transport Layer
    REF -- "JSON/HTTP (Tokens + Monotonic IDs)" --> BUS{VibeSync Bus}
    
    %% Engine Layer
    subgraph "🎮 Unity Adapter (Limb)"
        BUS --> UL[Listener Thread]
        UL --> UQ[Main Thread Queue]
        UQ --> US[Sandbox / .vibesync/tmp]
        US --> UC[Commit/Rollback]
    end
    
    subgraph "🧊 Blender Adapter (Limb)"
        BUS --> BL[Listener Thread]
        BL --> BT[bpy.app.timers]
        BT --> BS[Sandbox / Hidden Coll]
        BS --> BC[Commit/Rollback]
    end

    %% Conflict / State
    UC -- Telemetry --> TEL
    BC -- Telemetry --> TEL
    TEL --> MET
```

## Security & Integrity Zones
1. **The Intent Zone**: Cryptographic provenance and semantic validation.
2. **The Verification Zone**: The "Referee" loop ensures reality matches intent.
3. **The Sandbox Zone**: Engine mutations are staged in isolation before being committed.
4. **The Forensic Zone**: Immuatable journals (WAL/Events) for replay and audit.

---
*Copyright (C) 2026 B-A-M-N*
