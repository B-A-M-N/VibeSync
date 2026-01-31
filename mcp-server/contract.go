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
	StateQuarantine EngineState = "QUARANTINE"
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
	Capabilities    []string  `json:"capabilities"`
	Provenance      string    `json:"provenance"`
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

type CommitAtomicOpArgs struct {
	IntentID    string `json:"intent_id"`
	ProofOfWork string `json:"proof_of_work"`
	Reason      string `json:"reason,omitempty"`
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

type BridgeHeartbeat struct {
	BridgePID             int    `json:"bridge_pid"`
	UptimeSec             int    `json:"uptime_sec"`
	EpochID               int64  `json:"epoch_id"`
	OrchestratorConnected bool   `json:"orchestrator_connected"`
	UnityConnected        bool   `json:"unity_connected"`
	BlenderConnected       bool   `json:"blender_connected"`
	LastTickHash          string `json:"last_tick_hash"`
}

type BridgeHandshakeState struct {
	AssetID           string `json:"asset_id"`
	BlenderExportHash string `json:"blender_export_hash"`
	UnityImportHash   string `json:"unity_import_hash"`
	HashMatch         bool   `json:"hash_match"`
	LastVerified      string `json:"last_verified"`
}

type BridgeWalState struct {
	WalHead           int64  `json:"wal_head"`
	WalHash           string `json:"wal_hash"`
	LastCommittedOp   string `json:"last_committed_op"`
	PendingOps        int    `json:"pending_ops"`
	RollbackAvailable bool   `json:"rollback_available"`
}

type BridgeTransactionState struct {
	TransactionID string   `json:"transaction_id"`
	Status        string   `json:"status"`
	LockedAssets  []string `json:"locked_assets"`
	TimeoutSec    int      `json:"timeout_sec"`
}

type BridgeDeltaState struct {
	LastDeltaID    int64    `json:"last_delta_id"`
	Source         string   `json:"source"`
	AffectedAssets []string `json:"affected_assets"`
	DeltaHash      string   `json:"delta_hash"`
	Applied        bool     `json:"applied"`
}

type BridgeCommitRequirements struct {
	RequiredHashes    map[string]string `json:"required_hashes"`
	RationaleRequired bool              `json:"rationale_required"`
	CommitAllowed     bool              `json:"commit_allowed"`
}

type TechnicalRationaleCheck struct {
	Wal     string `json:"wal"`
	Blender string `json:"blender"`
	Unity   string `json:"unity"`
	Reason  string `json:"reason"`
}

type InvokeSpecialistArgs struct {
	SpecialistID string                 `json:"specialist_id"`
	IntentID     string                 `json:"intent_id"`
	CurrentHash  string                 `json:"current_hash"`
	TargetIntent map[string]interface{} `json:"target_intent"`
}