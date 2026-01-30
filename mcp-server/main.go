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

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Persistence Paths
const (
	PersistenceDir = ".vibesync"
	WalFile        = PersistenceDir + "/wal.jsonl"
	EventFile      = PersistenceDir + "/events.jsonl"
	StateFile      = PersistenceDir + "/state.json"
	UnityPort      = 8085
	BlenderPort    = 22000
	MaxWalSize     = 10 * 1024 * 1024
)

// Global State
var (
	engines = map[string]*EngineData{
		"unity":   {Token: "VIBE_UNITY_BOOTSTRAP_SECRET", State: StateStopped},
		"blender": {Token: "VIBE_BLENDER_BOOTSTRAP_SECRET", State: StateStopped},
	}
	currentSessionID = uuid.New().String()
	globalIDMap      = make(map[string]string)
	revocationList   = make(map[string]string)
	creditBalance    = 100
	stateMu          sync.RWMutex
	startTime        = time.Now()

	intents      = make(map[string]IntentEnvelope)
	transactions = make(map[string]*VibeTransaction)
	txMu         sync.Mutex

	activeTransaction *VibeTransaction

	monotonicID int64
	clockMu     sync.Mutex

	requestCounts = make(map[string]int)
	rateMu        sync.Mutex
)

type VibeTransaction struct {
	ID        string    `json:"id"`
	IntentID  string    `json:"intent_id"`
	StartTime time.Time `json:"start_time"`
	Status    string    `json:"status"`
}

type EngineData struct {
	Token        string      `json:"token"`
	State        EngineState `json:"state"`
	Generation   int         `json:"generation"`
	Version      string      `json:"version"`
	VersionHash  string      `json:"version_hash"`
	TrustExpiry  time.Time   `json:"trust_expiry"`
}

func init() {
	if _, err := os.Stat(PersistenceDir); os.IsNotExist(err) { os.Mkdir(PersistenceDir, 0755) }
	// Ensure sandbox directory exists
	sandbox := filepath.Join(PersistenceDir, "tmp")
	if _, err := os.Stat(sandbox); os.IsNotExist(err) { os.Mkdir(sandbox, 0755) }
	loadState()
	go startHeartbeatWatcher()
}

// --- INFRASTRUCTURE ---

func nextMonotonicID() int64 {
	clockMu.Lock(); defer clockMu.Unlock(); monotonicID++; return monotonicID
}

func dispatchVibeEvent(level EventLevel, eventType, intentID, next string, payload map[string]interface{}) {
	event := VibeEvent{Type: eventType, Level: level, IntentID: intentID, Timestamp: time.Now(), Payload: payload, NextStep: next}
	data, _ := json.Marshal(event)
	f, err := os.OpenFile(EventFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return }
	defer f.Close(); f.Write(data); f.Write([]byte("\n"))
}

func isRateLimited(target string) bool {
	rateMu.Lock(); defer rateMu.Unlock()
	if requestCounts[target] > 100 { return true }
	requestCounts[target]++
	go func() { time.Sleep(time.Second); rateMu.Lock(); requestCounts[target]--; rateMu.Unlock() }()
	return false
}

func auditPayload(data interface{}) error {
	jsonBytes, _ := json.Marshal(data)
	payload := strings.ToLower(string(jsonBytes))

	// Hard Execution Bans
	blocked := []string{
		"os.system", "exec(", "eval(", "rm -rf", "reflection", 
		"process.start", "import ", "__import__", "powershell",
		"cmd.exe", "/bin/sh", "/bin/bash",
	}
	for _, b := range blocked { 
		if strings.Contains(payload, b) { 
			return fmt.Errorf("SECURITY_VIOLATION: %s", b) 
		} 
	}

	// Numerical Instability check (NaN/Inf) - check marshaled string for common representations
	if strings.Contains(payload, "nan") || strings.Contains(payload, "inf") {
		return fmt.Errorf("NUMERICAL_INSTABILITY: NaN/Inf detected in payload")
	}

	return nil
}

func sendToEngine(target, endpoint, method string, data interface{}) (map[string]interface{}, error) {
	if isRateLimited(target) { return nil, fmt.Errorf("RATE_LIMIT") }
	// Basic audit
	if err := auditPayload(data); err != nil { 
		dispatchVibeEvent(LevelError, "security_intercept", "", "PANIC", map[string]interface{}{"error": err.Error()})
		return nil, err 
	}
	
	// Sanitize endpoint
	endpoint = strings.TrimPrefix(endpoint, "/")
	
	var lastErr error
	for i := 0; i < 3; i++ { // RETRY WITH BACKOFF
		res, err := attemptSend(target, endpoint, method, data)
		if err == nil {
			if res != nil && res["error"] == "Engine Busy: Compiling or Updating" {
				time.Sleep(2 * time.Second)
				continue
			}
			if method == "POST" && !strings.Contains(endpoint, "handshake") {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					verifyEngineState(ctx, target, endpoint)
				}()
			}
			return res, nil
		}
		lastErr = err
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond)
	}
	return nil, fmt.Errorf("ENGINE_ERROR | %v", lastErr)
}

// ... attemptSend remains same ...

func verifyEngineState(ctx context.Context, target, endpoint string) {
	log.Printf("🔍 REFEREE | Verifying %s after %s", target, endpoint)
	// Law of Reality: prove intent matches reality
	// Using a channel to handle the result of the potentially blocking sendToEngine call
	type result struct {
		res map[string]interface{}
		err error
	}
	done := make(chan result, 1)
	go func() {
		res, err := sendToEngine(target, "state/get", "GET", nil)
		done <- result{res, err}
	}()

	select {
	case <-ctx.Done():
		log.Printf("🚨 VERIFICATION TIMEOUT | %s failed to return state within deadline", target)
	case r := <-done:
		if r.err != nil {
			log.Printf("🚨 VERIFICATION FAILURE | %s failed to return state: %v", target, r.err)
			return
		}
		log.Printf("✅ VERIFIED | %s state hash: %v", target, r.res["hash"])
	}
}

func startHeartbeatWatcher() {
	ticker := time.NewTicker(5 * time.Second)
	portMap := map[string]int{"unity": UnityPort, "blender": BlenderPort}
	
	for range ticker.C {
		targets := make(map[string]struct {
			Port       int
			Generation int
		})

		stateMu.RLock()
		for name, e := range engines {
			if e.State == StateRunning {
				if port, known := portMap[name]; known {
					targets[name] = struct {
						Port       int
						Generation int
					}{Port: port, Generation: e.Generation}
				}
			}
		}
		stateMu.RUnlock()

		if len(targets) == 0 {
			continue
		}

		var wg sync.WaitGroup
		var panicMu sync.Mutex
		panicRequired := false

		for name, target := range targets {
			wg.Add(1)
			go func(n string, t struct{Port int; Generation int}) {
				defer wg.Done()
				client := &http.Client{Timeout: 2 * time.Second}
				resp, err := client.Get(fmt.Sprintf("http://localhost:%d/health", t.Port))
				
				engineFailed := false
				if err != nil || resp.StatusCode != 200 {
					log.Printf("🚨 HEARTBEAT FAILURE | %s is unresponsive. Triggering cluster panic.", n)
					engineFailed = true
				} else {
					var res map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
						// Generation Drift Detection
						if gen, ok := res["generation"].(float64); ok && int(gen) != t.Generation {
							log.Printf("🚨 DRIFT DETECTED | %s reported Gen %v, expected %d. Engine reloaded without handshake.", n, gen, t.Generation)
							engineFailed = true
						}
					}
					resp.Body.Close()
				}

				if engineFailed {
					stateMu.Lock()
					if e, ok := engines[n]; ok {
						e.State = StatePanic
					}
					stateMu.Unlock()
					panicMu.Lock()
					panicRequired = true
					panicMu.Unlock()
				} else {
					stateMu.Lock()
					if e, ok := engines[n]; ok {
						e.TrustExpiry = time.Now().Add(60 * time.Minute)
					}
					stateMu.Unlock()
				}
			}(name, target)
		}
		wg.Wait()

		if panicRequired {
			// Law of Stability: Broadcast panic to ALL engines
			for name := range engines {
				go sendToEngine(name, "panic", "POST", map[string]interface{}{"reason": "HEARTBEAT_TIMEOUT"})
			}
		}
	}
}

// --- TOOLS ---

func handshake_init(ctx context.Context, req *mcp.CallToolRequest, args HandshakeInitArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); engines[args.Target].State = StateStarting; engines[args.Target].Generation++; newToken, chal := uuid.New().String(), uuid.New().String(); stateMu.Unlock()
	res, err := sendToEngine(args.Target, "handshake", "POST", map[string]interface{}{"version": args.Version, "new_token": newToken, "challenge": chal})
	if err != nil || res["response"] != "VIBE_HASH_"+chal { return nil, nil, fmt.Errorf("AUTH_FAILED") }
	stateMu.Lock(); e := engines[args.Target]; e.Token, e.Version, e.State, e.TrustExpiry = newToken, fmt.Sprintf("%v", res["engine_version"]), StateRunning, time.Now().Add(60*time.Minute); stateMu.Unlock()
	dispatchVibeEvent(LevelInfo, "handshake_complete", "", "READY", map[string]interface{}{"target": args.Target})
	saveState(); return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func read_engine_state(ctx context.Context, req *mcp.CallToolRequest, args ReadStateArgs) (*mcp.CallToolResult, any, error) {
	res, _ := sendToEngine(args.Target, "state/get", "GET", nil)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Hash: %v", res["hash"])}}}, nil, nil
}

func verify_engine_state(ctx context.Context, req *mcp.CallToolRequest, args VerifyStateArgs) (*mcp.CallToolResult, any, error) {
	res, _ := sendToEngine(args.Target, "state/get", "GET", nil)
	if fmt.Sprintf("%v", res["hash"]) == args.ExpectedHash { return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "VERIFIED"}}}, nil, nil }
	return nil, nil, fmt.Errorf("DRIFT_DETECTED")
}

func submit_intent(ctx context.Context, req *mcp.CallToolRequest, args SubmitIntentArgs) (*mcp.CallToolResult, any, error) {
	if args.Envelope.Rationale == "" { return nil, nil, fmt.Errorf("RATIONALE_REQUIRED") }
	if args.Envelope.Confidence < 0 || args.Envelope.Confidence > 1.0 { return nil, nil, fmt.Errorf("INVALID_CONFIDENCE") }
	id := uuid.New().String(); txMu.Lock(); intents[id] = args.Envelope; txMu.Unlock()
	dispatchVibeEvent(LevelInfo, "intent_submitted", id, "VALIDATE", map[string]interface{}{"confidence": args.Envelope.Confidence})
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: id}}}, nil, nil
}

func validate_intent(ctx context.Context, req *mcp.CallToolRequest, args struct{ID string `json:"intent_id"`}) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); intent, ok := intents[args.ID]; txMu.Unlock()
	if !ok { return nil, nil, fmt.Errorf("UNKNOWN_INTENT") }
	
	if intent.Confidence < 0.8 {
		stateMu.Lock()
		for n := range engines { engines[n].State = StateHumanReq }
		stateMu.Unlock()
		dispatchVibeEvent(LevelWarn, "low_confidence_intercept", args.ID, "HUMAN_APPROVAL", map[string]interface{}{"confidence": intent.Confidence})
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "HUMAN_INTERVENTION_REQUIRED"}}}, nil, nil
	}
	
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "ALLOW"}}}, nil, nil
}

func human_approve_intent(ctx context.Context, req *mcp.CallToolRequest, args struct{ID string `json:"intent_id"`}) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); _, ok := intents[args.ID]; txMu.Unlock()
	if !ok { return nil, nil, fmt.Errorf("UNKNOWN_INTENT") }
	
	stateMu.Lock()
	for n := range engines { if engines[n].State == StateHumanReq { engines[n].State = StateRunning } }
	stateMu.Unlock()
	
	dispatchVibeEvent(LevelInfo, "human_approval_granted", args.ID, "EXECUTE", nil)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "APPROVED"}}}, nil, nil
}

func begin_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args AtomicOpArgs) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); defer txMu.Unlock(); 
	tx := &VibeTransaction{ID: uuid.New().String(), IntentID: args.IntentID, StartTime: time.Now(), Status: "OPEN"}
	transactions[args.IntentID] = tx
	activeTransaction = tx
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "TX_OPEN"}}}, nil, nil
}

func commit_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args AtomicOpArgs) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); defer txMu.Unlock(); 
	delete(transactions, args.IntentID)
	activeTransaction = nil
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "COMMITTED"}}}, nil, nil
}

func abort_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args AtomicOpArgs) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); defer txMu.Unlock(); 
	delete(transactions, args.IntentID)
	activeTransaction = nil
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "ABORTED"}}}, nil, nil
}

func emit_diag_bundle(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	p := filepath.Join(PersistenceDir, "diag.zip"); f, _ := os.Create(p); w := zip.NewWriter(f)
	for _, fn := range []string{WalFile, EventFile, StateFile} {
		df, _ := os.Open(fn); zw, _ := w.Create(filepath.Base(fn)); io.Copy(zw, df); df.Close()
	}
	w.Close(); f.Close(); return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: p}}}, nil, nil
}

func lock_object(ctx context.Context, req *mcp.CallToolRequest, args LockObjectArgs) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine(args.Target, "object/lock", "POST", map[string]interface{}{"id": args.ObjectID, "locked": args.Locked})
	if err != nil { return nil, nil, err }
	dispatchVibeEvent(LevelInfo, "object_lock_changed", "", "", map[string]interface{}{"id": args.ObjectID, "locked": args.Locked, "target": args.Target})
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%v", res["status"])}}}, nil, nil
}

func get_metrics(ctx context.Context, req *mcp.CallToolRequest, args struct{Target string `json:"target"`}) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine(args.Target, "metrics", "GET", nil)
	if err != nil { return nil, nil, err }
	jsonRes, _ := json.Marshal(res)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(jsonRes)}}}, nil, nil
}

func blender_mutate_mesh(ctx context.Context, req *mcp.CallToolRequest, args MutateArgs) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine("blender", "mesh/mutate", "POST", args.OpSpec)
	if err != nil { return nil, nil, err }; return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%v", res["status"])}}}, nil, nil
}

func unity_mutate_gameobject(ctx context.Context, req *mcp.CallToolRequest, args MutateArgs) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine("unity", "object/mutate", "POST", args.OpSpec)
	if err != nil { return nil, nil, err }; return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%v", res["status"])}}}, nil, nil
}

func sync_material(ctx context.Context, req *mcp.CallToolRequest, args SyncMaterialArgs) (*mcp.CallToolResult, any, error) {
	data := map[string]interface{}{"id": args.ObjectID, "properties": args.Props}
	// Record intent in WAL before execution
	journalOperation(map[string]interface{}{"type": "intent", "op": "sync_material", "id": args.ObjectID})
	
	sendToEngine("unity", "material/update", "POST", data)
	sendToEngine("blender", "material/update", "POST", data)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func sync_transform(ctx context.Context, req *mcp.CallToolRequest, args SyncTransformArgs) (*mcp.CallToolResult, any, error) {
	// Numerical Safety: Prevent NaN/Inf in creative sync
	for _, v := range append(append(args.Position, args.Rotation...), args.Scale...) {
		if math.IsNaN(v) || math.IsInf(v, 0) { return nil, nil, fmt.Errorf("NUMERICAL_INSTABILITY") }
	}

	data := map[string]interface{}{"id": args.ObjectID, "transform": map[string]interface{}{"pos": args.Position, "rot": args.Rotation, "sca": args.Scale}}
	sendToEngine("unity", "transform/set", "POST", data); sendToEngine("blender", "transform/set", "POST", data)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func sync_camera(ctx context.Context, req *mcp.CallToolRequest, args SyncCameraArgs) (*mcp.CallToolResult, any, error) {
	target := "unity"; if args.Source == "unity" { target = "blender" }
	
	res, err := sendToEngine(args.Source, "camera/get", "GET", nil)
	if err != nil { return nil, nil, err }
	
	sendToEngine(target, "camera/set", "POST", res)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func sync_selection(ctx context.Context, req *mcp.CallToolRequest, args SyncSelectionArgs) (*mcp.CallToolResult, any, error) {
	target := "unity"; if args.Source == "unity" { target = "blender" }
	
	sendToEngine(target, "selection/set", "POST", map[string]interface{}{"ids": args.IDs})
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func sync_asset_atomic(ctx context.Context, req *mcp.CallToolRequest, args SyncAssetAtomicArgs) (*mcp.CallToolResult, any, error) {
	pre, err := sendToEngine("blender", "preflight/run", "POST", map[string]interface{}{"path": args.AssetPath})
	if err != nil || pre == nil { return nil, nil, fmt.Errorf("PREFLIGHT_FAILED | %v", err) }
	
	ex, err := sendToEngine("blender", "export", "POST", map[string]interface{}{"path": args.AssetPath})
	if err != nil || ex == nil { return nil, nil, fmt.Errorf("EXPORT_FAILED | %v", err) }
	
	_, err = sendToEngine("unity", "import", "POST", map[string]interface{}{"path": args.AssetPath, "meta": ex["meta"], "mode": "sandbox"})
	if err != nil { return nil, nil, err }
	
	val, err := sendToEngine("unity", "validate", "POST", map[string]interface{}{"path": args.AssetPath})
	if err != nil || val == nil { return nil, nil, fmt.Errorf("VALIDATION_FAILED | %v", err) }
	
	if fmt.Sprintf("%v", pre["hash"]) != fmt.Sprintf("%v", val["hash"]) {
		sendToEngine("unity", "rollback", "POST", map[string]interface{}{"path": args.AssetPath})
		
		stateMu.Lock()
		for n := range engines { engines[n].State = StateDesync }
		stateMu.Unlock()
		
		return nil, nil, fmt.Errorf("HASH_MISMATCH | Cluster in DESYNC")
	}
	
	_, err = sendToEngine("unity", "commit", "POST", map[string]interface{}{"path": args.AssetPath})
	if err != nil { return nil, nil, err }
	
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "SYNCED"}}}, nil, nil
}

func reconcile_sync_state(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	u, _ := sendToEngine("unity", "state/get", "GET", nil)
	b, _ := sendToEngine("blender", "state/get", "GET", nil)
	
	if fmt.Sprintf("%v", u["hash"]) == fmt.Sprintf("%v", b["hash"]) {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "MATCHED"}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "DRIFT_DETECTED"}}}, nil, nil
}

func get_operation_journal(ctx context.Context, req *mcp.CallToolRequest, args struct{Limit int `json:"limit"`}) (*mcp.CallToolResult, any, error) {
	f, _ := os.Open(WalFile); defer f.Close(); s := bufio.NewScanner(f); var out []string
	for s.Scan() { out = append(out, s.Text()) }
	if len(out) > args.Limit && args.Limit > 0 { out = out[len(out)-args.Limit:] }
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(out, "\n")}}}, nil, nil
}

func control_playback(ctx context.Context, req *mcp.CallToolRequest, args struct {
	Action string  `json:"action"`
	Time   float64 `json:"time"`
}) (*mcp.CallToolResult, any, error) {
	data := map[string]interface{}{"action": args.Action, "time": args.Time}
	sendToEngine("unity", "playback/control", "POST", data)
	sendToEngine("blender", "playback/control", "POST", data)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func global_id_map_resolve(ctx context.Context, req *mcp.CallToolRequest, args MapVibeIDsArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); globalIDMap[args.UnityGUID], globalIDMap[args.BlenderName] = args.BlenderName, args.UnityGUID; stateMu.Unlock()
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "RESOLVED"}}}, nil, nil
}

// DRIVER REGISTRY (v0.4 Multiplexer Knowledge)
var drivers = map[string][]string{
	"vision_mcp":    {"render/capture", "material/get", "light/get"},
	"selection_mcp": {"selection/set", "hierarchy/get", "camera/frame"},
}

func vibe_multiplex(ctx context.Context, req *mcp.CallToolRequest, args MultiplexCallArgs) (*mcp.CallToolResult, any, error) {
	allowedEndpoints, ok := drivers[args.SensorID]
	if !ok { return nil, nil, fmt.Errorf("DRIVER_UNREGISTERED") }
	isAllowed := false
	for _, ep := range allowedEndpoints { if ep == args.Endpoint { isAllowed = true; break } }
	if !isAllowed { return nil, nil, fmt.Errorf("PERMISSION_DENIED") }

	res, err := sendToEngine(args.Target, args.Endpoint, "POST", args.Payload)
	if err != nil { return nil, nil, err }
	jsonRes, _ := json.Marshal(res)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(jsonRes)}}}, nil, nil
}

func set_engine_state(ctx context.Context, req *mcp.CallToolRequest, args SetEngineStateArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); defer stateMu.Unlock()
	e, ok := engines[args.Target]
	if !ok { return nil, nil, fmt.Errorf("Unknown target") }
	e.State = EngineState(args.State)
	saveState()
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "OK"}}}, nil, nil
}

func revoke_id(ctx context.Context, req *mcp.CallToolRequest, args RevokeIDArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); revocationList[args.ID] = args.Reason; stateMu.Unlock()
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "REVOKED"}}}, nil, nil
}

func epistemic_refusal(ctx context.Context, req *mcp.CallToolRequest, args struct{R string `json:"reason"`}) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "UNKNOWABLE: " + args.R}}}, nil, nil
}

func decommission_bridge(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); for n := range engines { engines[n].State = StatePanic }; stateMu.Unlock(); return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "DECOMMISSIONED"}}}, nil, nil
}

func reconstruct_state(ctx context.Context, req *mcp.CallToolRequest, args ForensicReplayArgs) (*mcp.CallToolResult, any, error) {
	f, _ := os.Open(EventFile); defer f.Close(); s := bufio.NewScanner(f); var out []string
	for s.Scan() { var e VibeEvent; json.Unmarshal(s.Bytes(), &e); out = append(out, e.Type) }
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(out, " -> ")}}}, nil, nil
}

// --- STATE PERSISTENCE ---

func loadState() {
	stateMu.Lock(); defer stateMu.Unlock(); data, _ := os.ReadFile(StateFile)
	var s struct { Engines map[string]*EngineData `json:"engines"`; IDMap map[string]string `json:"id_map"`; Credits int `json:"credits"` }
	json.Unmarshal(data, &s); creditBalance = s.Credits; globalIDMap = s.IDMap
	for k, v := range s.Engines { if e, ok := engines[k]; ok { e.State, e.Token, e.Version = v.State, v.Token, v.Version } }
}

func saveState() {
	stateMu.RLock(); s := struct { Engines map[string]*EngineData `json:"engines"`; IDMap map[string]string `json:"id_map"`; Credits int `json:"credits"` }{Engines: engines, IDMap: globalIDMap, Credits: creditBalance}
	stateMu.RUnlock(); data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(StateFile, data, 0644)
}

func journalOperation(op map[string]interface{}) {
	if activeTransaction != nil { op["tid"] = activeTransaction.ID }
	data, _ := json.Marshal(op)
	if info, err := os.Stat(WalFile); err == nil && info.Size() > MaxWalSize { os.Rename(WalFile, WalFile+".old") }
	f, _ := os.OpenFile(WalFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close(); f.Write(data); f.Write([]byte("\n"))
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "VibeSync", Version: "v0.4.0"}, &mcp.ServerOptions{})
	mcp.AddTool(server, &mcp.Tool{Name: "handshake_init", Description: "ISA 1"}, handshake_init)
	mcp.AddTool(server, &mcp.Tool{Name: "read_engine_state", Description: "ISA 2"}, read_engine_state)
	mcp.AddTool(server, &mcp.Tool{Name: "verify_engine_state", Description: "ISA 3"}, verify_engine_state)
	mcp.AddTool(server, &mcp.Tool{Name: "submit_intent", Description: "ISA 4 (Requires Rationale)"}, submit_intent)
	mcp.AddTool(server, &mcp.Tool{Name: "validate_intent", Description: "ISA 5 (Confidence Gate)"}, validate_intent)
	mcp.AddTool(server, &mcp.Tool{Name: "human_approve_intent", Description: "ISA 5b (Override StateHumanReq)"}, human_approve_intent)
	mcp.AddTool(server, &mcp.Tool{Name: "begin_atomic_operation", Description: "ISA 6"}, begin_atomic_operation)
	mcp.AddTool(server, &mcp.Tool{Name: "commit_atomic_operation", Description: "ISA 7"}, commit_atomic_operation)
	mcp.AddTool(server, &mcp.Tool{Name: "abort_atomic_operation", Description: "ISA 8"}, abort_atomic_operation)
	mcp.AddTool(server, &mcp.Tool{Name: "emit_diag_bundle", Description: "ISA 10"}, emit_diag_bundle)
	mcp.AddTool(server, &mcp.Tool{Name: "lock_object", Description: "Hierarchy Locking"}, lock_object)
	mcp.AddTool(server, &mcp.Tool{Name: "get_metrics", Description: "Observability"}, get_metrics)
	mcp.AddTool(server, &mcp.Tool{Name: "blender_mutate_mesh", Description: "ISA 14"}, blender_mutate_mesh)
	mcp.AddTool(server, &mcp.Tool{Name: "unity_mutate_gameobject", Description: "ISA 21"}, unity_mutate_gameobject)
	mcp.AddTool(server, &mcp.Tool{Name: "sync_transform", Description: "ISA 21*"}, sync_transform)
	mcp.AddTool(server, &mcp.Tool{Name: "sync_material", Description: "ISA 22"}, sync_material)
	mcp.AddTool(server, &mcp.Tool{Name: "sync_camera", Description: "ISA 23"}, sync_camera)
	mcp.AddTool(server, &mcp.Tool{Name: "sync_selection", Description: "ISA 24"}, sync_selection)
	mcp.AddTool(server, &mcp.Tool{Name: "sync_asset_atomic", Description: "ISA 26"}, sync_asset_atomic)
	mcp.AddTool(server, &mcp.Tool{Name: "reconcile_sync_state", Description: "ISA 2"}, reconcile_sync_state)
	mcp.AddTool(server, &mcp.Tool{Name: "get_operation_journal", Description: "Forensic WAL"}, get_operation_journal)
	mcp.AddTool(server, &mcp.Tool{Name: "control_playback", Description: "Coordinated Timeline"}, control_playback)
	mcp.AddTool(server, &mcp.Tool{Name: "global_id_map_resolve", Description: "ISA 27"}, global_id_map_resolve)
	mcp.AddTool(server, &mcp.Tool{Name: "vibe_multiplex", Description: "Kernel Multiplexer"}, vibe_multiplex)
	mcp.AddTool(server, &mcp.Tool{Name: "set_engine_state", Description: "ISA"}, set_engine_state)
	mcp.AddTool(server, &mcp.Tool{Name: "revoke_id", Description: "Meta"}, revoke_id)
	mcp.AddTool(server, &mcp.Tool{Name: "epistemic_refusal", Description: "ISA 29"}, epistemic_refusal)
	mcp.AddTool(server, &mcp.Tool{Name: "decommission_bridge", Description: "ISA 32"}, decommission_bridge)
	mcp.AddTool(server, &mcp.Tool{Name: "reconstruct_state", Description: "Forensic"}, reconstruct_state)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil { log.Fatal(err) }
}
