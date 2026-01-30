// VibeSync: Zero-Trust Unity ↔ Blender Orchestrator
// Copyright (C) 2026 B-A-M-N
//
// This project is distributed under a DUAL-LICENSING MODEL:
// 1. Open-Source Path: GNU Affero General Public License v3
// 2. Commercial Path: "Work-or-Pay" Model
//
// See the LICENSE file in the project root for the full terms and conditions
// of both licensing paths.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

package main

import "time"

type EngineState string

const (
	StateStopped  EngineState = "STOPPED"
	StateStarting EngineState = "STARTING"
	StateRunning  EngineState = "RUNNING"
	StatePanic    EngineState = "PANIC"
	StateHumanReq EngineState = "HUMAN_INTERVENTION_REQUIRED"
	StateDesync   EngineState = "DESYNC"
)

type VibeBaseModel struct {
	Generation  int    `json:"generation"`
	SessionID   string `json:"session_id"`
	MonotonicID int64  `json:"monotonic_id"`
	IntentID    string `json:"intent_id,omitempty"`
}

type IntentEnvelope struct {
	InstructionHash string    `json:"instruction_hash"`
	PlanHash        string    `json:"plan_hash"`
	Rationale       string    `json:"rationale"`
	Confidence      float64   `json:"confidence"`
	Scope           []string  `json:"scope"`
	BudgetMS        int       `json:"budget_ms"`
	Signature       string    `json:"signature"`
}

type EventLevel string

const (
	LevelInfo  EventLevel = "INFO"
	LevelDebug EventLevel = "DEBUG"
	LevelWarn  EventLevel = "WARN"
	LevelError EventLevel = "ERROR"
)

type VibeEvent struct {
	Type      string                 `json:"type"`
	Level     EventLevel             `json:"level"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
	NextStep  string                 `json:"next_step,omitempty"`
	IntentID  string                 `json:"intent_id,omitempty"`
}

type SetEngineStateArgs struct {
	Target string `json:"target"`
	State  string `json:"state"`
}

type HandshakeInitArgs struct {
	Target  string `json:"target"`
	Version string `json:"version"`
}

type HandshakeResponse struct {
	Status        string   `json:"status"`
	EngineVersion string   `json:"engine_version"`
	Capabilities  []string `json:"capabilities"`
}

type ReadStateArgs struct {
	Target string `json:"target"`
}

type VerifyStateArgs struct {
	Target       string `json:"target"`
	ExpectedHash string `json:"expected_hash"`
}

type SubmitIntentArgs struct {
	Envelope IntentEnvelope `json:"envelope"`
}

type AtomicOpArgs struct {
	IntentID string `json:"intent_id"`
	Reason   string `json:"reason,omitempty"`
}

type MutateArgs struct {
	IntentID string                 `json:"intent_id"`
	OpSpec   map[string]interface{} `json:"op_spec"`
}

type SyncAssetAtomicArgs struct {
	AssetPath string `json:"asset_path"`
}

type MultiplexCallArgs struct {
	SensorID string                 `json:"sensor_id"`
	Target   string                 `json:"target"`
	Endpoint string                 `json:"endpoint"`
	Payload  map[string]interface{} `json:"payload"`
}

type SyncTransformArgs struct {
	ObjectID string    `json:"object_id"`
	Position []float64 `json:"position"`
	Rotation []float64 `json:"rotation"`
	Scale    []float64 `json:"scale"`
}

type SyncCameraArgs struct {
	Source string `json:"source"`
}

type SyncSelectionArgs struct {
	Source string   `json:"source"`
	IDs    []string `json:"ids"`
}

type SyncMaterialArgs struct {
	ObjectID string                 `json:"object_id"`
	Props    map[string]interface{} `json:"properties"`
}

type LockObjectArgs struct {
	Target   string `json:"target"`
	ObjectID string `json:"object_id"`
	Locked   bool   `json:"locked"`
}

type MapVibeIDsArgs struct {
	UnityGUID   string `json:"unity_guid"`
	BlenderName string `json:"blender_name"`
}

type RevokeIDArgs struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

type ForensicReplayArgs struct {
	TargetMonotonicID int64 `json:"target_monotonic_id"`
}