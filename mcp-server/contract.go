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

type Actor string

const (
	ActorHuman Actor = "human"
	ActorAI    Actor = "ai"
	ActorSystem Actor = "system"
)

const (
	StateStopped  EngineState = "STOPPED"
	StateStarting EngineState = "STARTING"
	StateRunning  EngineState = "RUNNING"
	StatePanic    EngineState = "PANIC"
	StateHumanReq EngineState = "HUMAN_INTERVENTION_REQUIRED"
	StateDesync   EngineState = "DESYNC"
	StateQuarantine EngineState = "QUARANTINE"
)

type VibeOpcode uint8

const (
	OpTransform   VibeOpcode = 0x03
	OpModifier    VibeOpcode = 0x04
	OpNode        VibeOpcode = 0x05
	OpMaterial    VibeOpcode = 0x09
	OpBake        VibeOpcode = 0x0A
	OpIO          VibeOpcode = 0x0B
	OpSystem      VibeOpcode = 0x0F
	OpAudit       VibeOpcode = 0x10
)

type IntentType string

const (
	IntentOptimize   IntentType = "OPTIMIZE"
	IntentRig        IntentType = "RIG"
	IntentLight      IntentType = "LIGHT"
	IntentAnimate    IntentType = "ANIMATE"
	IntentSceneSetup IntentType = "SCENE_SETUP"
	IntentGeneral    IntentType = "GENERAL"
)

type VibeUnitSettings struct {
	System      string  `json:"system"` // Metric, Imperial
	ScaleLength float64 `json:"scale_length"`
}

type VibeBaseModel struct {
	Generation    int    `json:"generation"`
	SessionID     string `json:"session_id"`
	MonotonicID   int64  `json:"monotonic_id"`
	MonotonicTick int64  `json:"monotonic_tick"`
	IntentID      string `json:"intent_id,omitempty"`
	SchemaVersion string `json:"schema_version"`
}

type IntentEnvelope struct {
	InstructionHash string            `json:"instruction_hash"`
	PlanHash          string            `json:"plan_hash"`
	Rationale         string            `json:"rationale"`
	Confidence        float64           `json:"confidence"`
	Scope             []string          `json:"scope"`
	Capabilities      []string          `json:"capabilities"`
	Provenance        string            `json:"provenance"`
	BudgetMS          int               `json:"budget_ms"`
	Signature         string            `json:"signature"`
	BasedOnHashes     map[string]string `json:"based_on_hashes"`
	Intent            IntentType        `json:"intent"`
	Opcode            VibeOpcode        `json:"opcode,omitempty"`
	DryRun            bool              `json:"dry_run,omitempty"`
}

type VibeState string

const (
	StateKnown   VibeState = "KNOWN"
	StateUnknown VibeState = "UNKNOWN"
	StateInvalid VibeState = "INVALID"
)

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
	IntentID       string                 `json:"intent_id"`
	IdempotencyKey string                 `json:"idempotency_key"`
	OpSpec         map[string]interface{} `json:"op_spec"`
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
	ExpectedIntervalMS    int    `json:"expected_interval_ms"`
	LastSeenMS            int    `json:"last_seen_ms"`
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
	Reversible        bool   `json:"reversible"`
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

type IntentClass string

const (
	ClassCosmetic    IntentClass = "cosmetic"
	ClassStructural  IntentClass = "structural"
	ClassDestructive IntentClass = "destructive"
)

type WalPhase string

const (
	PhasePending     WalPhase = "PENDING"
	PhaseAttempted   WalPhase = "ATTEMPTED"
	PhaseFailed      WalPhase = "FAILED"
	PhaseHalted      WalPhase = "HALTED"
	PhaseTerminal    WalPhase = "TERMINAL"
	PhaseProvisional WalPhase = "PROVISIONAL"
	PhaseFinal       WalPhase = "FINAL"
	PhaseRolledBack  WalPhase = "ROLLED_BACK"
	PhaseQuarantined WalPhase = "QUARANTINED"
	PhaseWaitHuman   WalPhase = "WAIT_HUMAN_LOCK"
)

type FailureClass string

const (
	FailSyntax            FailureClass = "SyntaxError"
	FailDependency        FailureClass = "DependencyMissing"
	FailNamespace         FailureClass = "NamespaceCollision"
	FailAssetMismatch     FailureClass = "AssetMismatch"
	FailInvalidState      FailureClass = "InvalidState"
	FailToolUnavailable   FailureClass = "ToolUnavailable"
	FailPolicyViolation   FailureClass = "PolicyViolationRisk"
	FailUnknown           FailureClass = "Unknown"
)

type RolePermissions struct {
	CanExecute  bool `json:"can_execute"`
	CanRetry    bool `json:"can_retry"`
	MaxRetries  int  `json:"max_retries,omitempty"`
	CanEscalate bool `json:"can_escalate,omitempty"`
	CanFreeze   bool `json:"can_freeze,omitempty"`
}

type PermissionsMask map[string]RolePermissions

type WalEntry struct {
	IntentID   uint64           `json:"intent_id"`
	ParentHash string           `json:"parent_hash"`
	EntryHash  string           `json:"entry_hash"`
	Timestamp  int64            `json:"timestamp"`
	Engine     string           `json:"engine"`
	Actor      string           `json:"actor"`
	Scope      WalScope         `json:"scope"`
	Phase      WalPhase         `json:"phase"`
	Verify     WalVerify        `json:"verification"`
	Rollback   WalRoll          `json:"rollback"`
	Conflict   ConflictMetadata `json:"conflict,omitempty"`
	FailureSig string           `json:"failure_signature,omitempty"`
	FailureClass FailureClass   `json:"failure_class,omitempty"`
	RetryCount int              `json:"retry_count"`
	EscalationLevel int         `json:"escalation_level"`
	Permissions PermissionsMask `json:"permissions_mask,omitempty"`
	SystemHealth string         `json:"system_health"` // SAFE | QUARANTINED
}

type FailureSignature struct {
	Engine     string `json:"engine"`
	Opcode     uint8  `json:"opcode"`
	Target     string `json:"target_uuid"`
	ErrorCode  string `json:"error_code"`
	Location   string `json:"location"` // File:Line
}

type ConflictMetadata struct {
	Type           string `json:"type"`
	Resolution     string `json:"resolution"`
	WinnerIntentID uint64 `json:"winner_intent_id,omitempty"`
	Reason         string `json:"reason"`
}

type ObjectKind string

const (
	KindPrefabDef      ObjectKind = "PREFAB_DEF"
	KindPrefabInstance ObjectKind = "PREFAB_INSTANCE"
	KindObject         ObjectKind = "OBJECT"
)

type ObjectIdentity struct {
	UUID        string     `json:"uuid"`
	Kind        ObjectKind `json:"kind"`
	PrefabDepth int        `json:"prefab_depth"`
}

type LockType string

const (
	LockHumanActive LockType = "HUMAN_ACTIVE"
	LockAISpeculative LockType = "AI_SPECULATIVE"
	LockPerimeter LockType = "PERIMETER_LOCK"
)

type VibeLock struct {
	UUID      string    `json:"uuid"`
	Type      LockType  `json:"type"`
	Actor     Actor     `json:"actor"`
	Timestamp time.Time `json:"timestamp"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ApplyLockArgs struct {
	UUID     string   `json:"uuid"`
	LockType LockType `json:"lock_type"`
}

type ReleaseLockArgs struct {
	UUID string `json:"uuid"`
}

type WalScope struct {
	UUIDs        []string    `json:"uuids"`
	ClosureUUIDs []string    `json:"closure_uuids,omitempty"`
	Class        IntentClass `json:"intent_class"`
}

type WalVerify struct {
	ExpectedHash string  `json:"expected_hash"`
	ObservedHash string  `json:"observed_hash,omitempty"`
	Epsilon      float64 `json:"epsilon"`
	VerifiedAt   int64   `json:"verified_at,omitempty"`
}

type WalRoll struct {
	UndoToken   string `json:"undo_token"`
	SnapshotRef string `json:"snapshot_ref,omitempty"`
}

type WorkOrder struct {
	ID          string     `json:"id"`
	MonotonicID int64      `json:"monotonic_id"`
	Intent      IntentType `json:"intent"`
	Opcode      VibeOpcode `json:"opcode"`
	UUID        string     `json:"uuid"`
	Description string     `json:"description"`
	Context     map[string]interface{} `json:"context"` // Snapshot of relevant state
}

type WorkResult struct {
	ID           string `json:"id"`
	WorkOrderID  string `json:"work_order_id"`
	Status       string `json:"status"` // SUCCESS | FAILURE | BUSY
	Hash         string `json:"hash"`
	Error        string `json:"error,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
}

type VibeSitrep struct {
	Timestamp     time.Time              `json:"timestamp"`
	Pulse         string                 `json:"pulse"`
	EngineStatus  map[string]string      `json:"engine_status"`
	AffordanceMap map[string][]string    `json:"affordance_map"` // UUID -> list of valid Opcodes/Actions
	GlobalPerimeter bool                 `json:"global_perimeter_locked"`
}

type StrategicPlan struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Steps       []PlanStep   `json:"steps"`
	Rationale   string       `json:"rationale"`
	TotalBudget int          `json:"total_budget_ms"`
}

type PlanStep struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Intent      IntentType `json:"intent"`
	Opcode      VibeOpcode `json:"opcode"`
	UUID        string     `json:"uuid"`
}

type MutationIntegrityArgs struct {
	UUID     string     `json:"uuid"`
	Opcode   VibeOpcode `json:"opcode"`
	Expected map[string]interface{} `json:"expected_state"`
}

type ForensicSnapshotArgs struct {
	Reason string `json:"reason"`
	Tid    string `json:"transaction_id,omitempty"`
}

type EntropyBudget struct {
	Limit int `json:"limit"`
	Used  int `json:"used"`
}
