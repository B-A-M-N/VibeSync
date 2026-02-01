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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	UnityPort      = 8087
	BlenderPort    = 22005
	MaxWalSize     = 10 * 1024 * 1024
)

// Global State
var (
	unityPort      = 8087
	engines = map[string]*EngineData{
		"unity":   {Token: "5715493b", State: StateStopped},
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

	lastWalHash string
	walMu       sync.Mutex

	requestCounts = make(map[string]int)
	rateMu        sync.Mutex

	// Lock Table
	lockTable = make(map[string]*VibeLock)
	lockMu    sync.RWMutex

	// Intent Coalescing
	intentBuffer = make(map[string]*WalEntry)
	bufferMu     sync.Mutex

	// Visual Thought Markers
	ActivityFile  = "metadata/bridge_activity.txt"
	DiscoveryFile = "/home/bamn/ALCOM/Projects/BAMN-EXTO/metadata/vibe_status.json"
	SettingsFile  = "/home/bamn/ALCOM/Projects/BAMN-EXTO/metadata/vibe_settings.json"
	AuditFile     = "/home/bamn/ALCOM/Projects/BAMN-EXTO/logs/vibe_audit.jsonl"

	// Second-Order Invariant State
	idempotencyMap = make(map[string]string) // key -> hash
	idempMu        sync.Mutex
	sessionEntropy = 0
	maxEntropy     = 100 // Default budget
	entropyMu      sync.Mutex
	schemaVersion  = "bridge.v0.4.0"

	// Trust Tiers & Performance Mode
	performanceMode = false // Toggle for high-frequency data
	trustTier       = 0     // 0: Normal, 1: Trusted (Low Latency)

	drivers = map[string][]string{
		"vision_mcp":    {"render/capture", "material/get", "light/get"},
		"selection_mcp": {"selection/set", "hierarchy/get", "camera/frame"},
	}

	hardenedClient = &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				LocalAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1")},
			}).DialContext,
			Proxy: nil,
			MaxIdleConns: 100,
			IdleConnTimeout: 90 * time.Second,
		},
	}
)

type VibeTransaction struct {
	ID        string    `json:"id"`
	IntentID  string    `json:"intent_id"`
	StartTime time.Time `json:"start_time"`
	Status    string    `json:"status"`
}

type EngineData struct {
	Token         string      `json:"token"`
	State         EngineState `json:"state"`
	Generation    int         `json:"generation"`
	Version       string      `json:"version"`
	VersionHash   string      `json:"version_hash"`
	TrustExpiry   time.Time   `json:"trust_expiry"`
	TrustScore    int         `json:"trust_score"`
	MutationCount int         `json:"mutation_count"`
	LastMutation  time.Time   `json:"last_mutation"`
}

func discoverSettings() int {
	data, err := os.ReadFile(SettingsFile)
	if err != nil {
		return 8087 // Fallback
	}
	var settings struct {
		Ports struct {
			Control int `json:"control"`
		} `json:"ports"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return 8087
	}
	return settings.Ports.Control
}

func discoverLatestHash() (string, int64) {
	f, err := os.Open(AuditFile)
	if err != nil {
		return "", 0
	}
	defer f.Close()

	var lastLine string
	var lineCount int64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lastLine = scanner.Text()
		lineCount++
	}

	var entry struct {
		EntryHash string `json:"entryHash"`
	}
	json.Unmarshal([]byte(lastLine), &entry)
	return entry.EntryHash, lineCount
}

func discoverUnityToken() string {
	data, err := os.ReadFile(DiscoveryFile)
	if err != nil {
		return "5715493b" // Fallback
	}
	var status struct {
		State string `json:"state"`
		Nonce string `json:"nonce"`
	}
	if err := json.Unmarshal(data, &status); err != nil {
		return "5715493b"
	}
	return status.Nonce
}

func isPerformanceOp(endpoint string) bool {
	ops := []string{"transform/set", "camera/set", "playback/control", "metrics"}
	for _, op := range ops {
		if strings.Contains(endpoint, op) {
			return true
		}
	}
	return false
}

func init() {
	if _, err := os.Stat(PersistenceDir); os.IsNotExist(err) { os.Mkdir(PersistenceDir, 0755) }
	sandbox := filepath.Join(PersistenceDir, "tmp")
	if _, err := os.Stat(sandbox); os.IsNotExist(err) { os.Mkdir(sandbox, 0755) }
	loadState()

	// 1. Token & Port Discovery
	token := discoverUnityToken()
	port := discoverSettings()

	// 2. Hash & Tick Discovery
	hash, tick := discoverLatestHash()

	stateMu.Lock()
	unityPort = port
	if e, ok := engines["unity"]; ok {
		e.Token = token
		lastWalHash = hash
		monotonicID = tick
		log.Printf("🛡️ VibeSync Discovery: Port=%d | Token=%s | Hash=%s | Tick=%d", port, token, hash, tick)
	}
	stateMu.Unlock()

	go startSyncLoop()
	go startHeartbeatWatcher()
	go startCoalescingLoop()
}

func startCoalescingLoop() {
	ticker := time.NewTicker(250 * time.Millisecond)
	for range ticker.C {
		flushIntentBuffer()
	}
}

func flushIntentBuffer() {
	bufferMu.Lock()
	count := len(intentBuffer)
	if count == 0 {
		bufferMu.Unlock()
		return
	}
	log.Printf("🌊 VibeSync Batching: Finalizing %d speculative intents", count)
	
	// In a full implementation, we'd trigger a combined verification hash here
	// and promote WAL entries from PROVISIONAL to FINAL.
	
	intentBuffer = make(map[string]*WalEntry)
	bufferMu.Unlock()
}

func bufferSpeculativeIntent(uuid string, op string, data interface{}) {
	bufferMu.Lock()
	defer bufferMu.Unlock()
	
	// Create a virtual WalEntry for the buffer
	intentBuffer[uuid] = &WalEntry{
		IntentID: uint64(monotonicID),
		Engine:   "orchestrator",
		Actor:    "ai",
		Scope:    WalScope{UUIDs: []string{uuid}, Class: ClassCosmetic},
		Phase:    PhaseProvisional,
	}
}

func checkHumanLock(uuid string) error {
	lockMu.RLock()
	defer lockMu.RUnlock()
	if lock, ok := lockTable[uuid]; ok {
		if lock.Type == LockHumanActive && time.Now().Before(lock.ExpiresAt) {
			return fmt.Errorf("WAIT_HUMAN_LOCK: UUID %s is under active human manipulation", uuid)
		}
	}
	return nil
}

func startSyncLoop() {
	// Simple polling loop as a fallback for inotify
	// In a real scenario, we'd use fsnotify
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		token := discoverUnityToken()
		hash, tick := discoverLatestHash()
		
		stateMu.Lock()
		if e, ok := engines["unity"]; ok {
			if e.Token != token {
				log.Printf("🔄 VibeSync: Token Rotation Detected -> %s", token)
				e.Token = token
			}
			if lastWalHash != hash {
				log.Printf("🔄 VibeSync: Chain of Trust Updated -> %s", hash)
				lastWalHash = hash
				monotonicID = tick
			}
		}
		stateMu.Unlock()
	}
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

func checkInvariants(idempKey, targetHash string) error {
	entropyMu.Lock()
	if sessionEntropy >= maxEntropy {
		entropyMu.Unlock()
		return fmt.Errorf("INVARIANT_VIOLATION: Entropy budget exhausted")
	}
	sessionEntropy++
	entropyMu.Unlock()

	if idempKey != "" {
		idempMu.Lock()
		defer idempMu.Unlock()
		prevHash, exists := idempotencyMap[idempKey]
		if exists && prevHash != targetHash {
			return fmt.Errorf("INVARIANT_VIOLATION: Idempotency breach (Key exists with different hash)")
		}
		idempotencyMap[idempKey] = targetHash
	}
	return nil
}

func sanitizeForTarget(target string, data interface{}) interface{} {
	if data == nil { return nil }
	jsonBytes, _ := json.Marshal(data)
	payload := string(jsonBytes)

	if target == "blender" {
		payload = strings.ReplaceAll(payload, "GameObject", "Object")
		payload = strings.ReplaceAll(payload, "Prefab", "Template")
		payload = strings.ReplaceAll(payload, "MonoBehaviour", "Script")
	} else if target == "unity" {
		payload = strings.ReplaceAll(payload, "bpy.data", "EngineData")
		payload = strings.ReplaceAll(payload, "Collection", "Folder")
		payload = strings.ReplaceAll(payload, "DataBlock", "Asset")
	}

	re := regexp.MustCompile(`(0x[0-9a-fA-F]+|ptr:[0-9]+|InstanceID:[0-9]+)`)
	payload = re.ReplaceAllString(payload, "[REDACTED_HANDLE]")

	var sanitized interface{}
	json.Unmarshal([]byte(payload), &sanitized)
	return sanitized
}

func auditPayload(data interface{}) error {
	if data == nil { return nil }
	jsonBytes, _ := json.Marshal(data)
	payload := strings.ToLower(string(jsonBytes))

	blocked := []string{"os.system", "exec(", "eval(", "rm -rf", "reflection", "process.start", "import ", "__import__", "powershell", "cmd.exe", "/bin/sh", "/bin/bash"}
	for _, b := range blocked { if strings.Contains(payload, b) { return fmt.Errorf("SECURITY_VIOLATION: %s", b) } }
	if strings.Contains(payload, "nan") || strings.Contains(payload, "inf") { return fmt.Errorf("NUMERICAL_INSTABILITY: NaN/Inf detected") }
	return nil
}

func decayTrust(target string, amount int, reason string) {
	stateMu.Lock(); defer stateMu.Unlock()
	e, ok := engines[target]; if !ok { return }
	e.TrustScore -= amount
	if e.TrustScore < 0 { e.TrustScore = 0 }
	if e.TrustScore < 20 && e.State != StatePanic {
		e.State = StateQuarantine
		dispatchVibeEvent(LevelWarn, "quarantine_triggered", "", "COOL_OFF", map[string]interface{}{"target": target, "reason": reason})
	}
}

func getForensicReport() map[string]interface{} {
	report := make(map[string]interface{})
	f, err := os.Open(WalFile)
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			if len(lines) > 3 { lines = lines[1:] }
		}
		report["last_wal_entries"] = lines
	}
	stateMu.RLock()
	health := make(map[string]string)
	for name, data := range engines { health[name] = fmt.Sprintf("State: %s | Gen: %d", data.State, data.Generation) }
	stateMu.RUnlock()
	report["engine_status"], report["system_time"] = health, time.Now().Format(time.RFC3339)
	return report
}

func wrapForensicResult(data interface{}) *mcp.CallToolResult {
	report := getForensicReport()
	entropyMu.Lock()
	report["entropy_stats"] = map[string]int{"limit": maxEntropy, "used": sessionEntropy}
	entropyMu.Unlock()
	wrapped := map[string]interface{}{
		"result":          data,
		"forensic_report": report,
		"schema_version":  schemaVersion,
	}
	jsonStr, _ := json.MarshalIndent(wrapped, "", "  ")
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(jsonStr)}}}
}

func sendToEngine(target, endpoint, method string, data interface{}) (map[string]interface{}, error) {
	stateMu.RLock(); engine, ok := engines[target]; stateMu.RUnlock()
	if !ok { return nil, fmt.Errorf("unknown target") }

	log.Printf("📡 DEBUG | sendToEngine: %s %s/%s", method, target, endpoint)

	if method == "POST" && !strings.Contains(endpoint, "handshake") {
		stateMu.Lock()
		now := time.Now()
		if now.Sub(engine.LastMutation) < 200*time.Millisecond { engine.MutationCount++ } else { engine.MutationCount = 1 }
		engine.LastMutation = now
		if engine.MutationCount > 5 { stateMu.Unlock(); decayTrust(target, 5, "EXCESSIVE_MUTATION_RATE"); return nil, fmt.Errorf("RATE_LIMITED_ADAPTIVE") }
		stateMu.Unlock()
	}

	if isRateLimited(target) { return nil, fmt.Errorf("RATE_LIMIT") }
	if err := auditPayload(data); err != nil { dispatchVibeEvent(LevelError, "security_intercept", "", "PANIC", map[string]interface{}{"error": err.Error()}); decayTrust(target, 20, "AUDIT_VIOLATION"); return nil, err }
	
	endpoint = strings.TrimPrefix(endpoint, "/")
	data = sanitizeForTarget(target, data)

	var lastErr error
	for i := 0; i < 3; i++ {
		res, err := attemptSend(target, endpoint, method, data)
		if err == nil {
			if res != nil && res["error"] == "Engine Busy: Compiling or Updating" { time.Sleep(2 * time.Second); continue }
			if method == "POST" && !strings.Contains(endpoint, "handshake") {
				go func() { ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second); defer cancel(); verifyEngineState(ctx, target, endpoint) }()
			}
			return res, nil
		}
		lastErr = err
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond)
	}
	return nil, fmt.Errorf("ENGINE_ERROR | %v", lastErr)
}

func computeSignature(token, timestamp, method, path, body string) string {
	h := hmac.New(sha256.New, []byte(token))
	h.Write([]byte(timestamp + "|" + method + "|" + path + "|" + body))
	return hex.EncodeToString(h.Sum(nil))
}

func attemptSend(target, endpoint, method string, data interface{}) (map[string]interface{}, error) {
	stateMu.RLock(); engine, ok := engines[target]; stateMu.RUnlock()
	if !ok || engine.State == StatePanic || engine.State == StateHumanReq { return nil, fmt.Errorf("LOCKED") }
	if engine.State == StateQuarantine && method != "GET" && !strings.Contains(endpoint, "health") { return nil, fmt.Errorf("QUARANTINE_READ_ONLY") }
	if time.Now().After(engine.TrustExpiry) && engine.State == StateRunning { return nil, fmt.Errorf("EXPIRED") }

	isPerf := isPerformanceOp(endpoint)

	port := unityPort; if target == "blender" { port = BlenderPort }
	url := fmt.Sprintf("http://127.0.0.1:%d/%s", port, endpoint)
	log.Printf("📡 DEBUG | attemptSend: %s %s (Token: %s)", method, url, engine.Token)
	mid := nextMonotonicID()
	tid := ""
	if m, ok := data.(map[string]interface{}); ok {
		m["generation"], m["session_id"], m["monotonic_id"] = engine.Generation, currentSessionID, mid
		txMu.Lock(); if activeTransaction != nil { tid = activeTransaction.ID; m["tid"], m["parent_id"] = tid, tid }; txMu.Unlock()
	}
	
	jsonBody, _ := json.Marshal(data); bodyStr := string(jsonBody); if data == nil { bodyStr = "" }
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	signature := ""
	if isPerf && trustTier > 0 {
		// Performance optimization: skip signature for high-trust performance ops
	} else {
		signature = computeSignature(engine.Token, timestamp, method, "/"+endpoint, bodyStr)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("X-Vibe-Token", engine.Token); 
	if signature != "" {
		req.Header.Set("X-Vibe-Signature", signature); 
	}
	req.Header.Set("X-Vibe-Timestamp", timestamp); req.Header.Set("X-Vibe-Session", currentSessionID); req.Header.Set("X-Vibe-Generation", fmt.Sprintf("%d", engine.Generation))
	if tid != "" { req.Header.Set("X-Vibe-Transaction", tid) }
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := hardenedClient.Do(req); if err != nil { return nil, err }; defer resp.Body.Close()
	var res map[string]interface{}; if err := json.NewDecoder(resp.Body).Decode(&res); err != nil { if resp.StatusCode >= 400 { return nil, fmt.Errorf("HTTP %d", resp.StatusCode) }; return nil, err }
	journalOperation(map[string]interface{}{"type": "engine_call", "target": target, "endpoint": endpoint, "mid": mid})
	return res, nil
}

func verifyEngineState(ctx context.Context, target, endpoint string) {
	log.Printf("🔍 REFEREE | Verifying %s after %s", target, endpoint)
	type result struct { res map[string]interface{}; err error }
	done := make(chan result, 1)
	go func() { res, err := sendToEngine(target, "state/get", "GET", nil); done <- result{res, err} }()
	select { case <-ctx.Done(): log.Printf("🚨 VERIFICATION TIMEOUT | %s", target); case r := <-done: if r.err != nil { log.Printf("🚨 VERIFICATION FAILURE | %s: %v", target, r.err) } else { log.Printf("✅ VERIFIED | %s: %v", target, r.res["hash"]) } }
}

func startHeartbeatWatcher() {
	ticker := time.NewTicker(5 * time.Second); portMap := map[string]int{"unity": unityPort, "blender": BlenderPort}
	for range ticker.C {
		targets := make(map[string]struct{ Port, Generation int })
		stateMu.RLock()
		for name, e := range engines { 
			if e.State == StateRunning { 
				if port, known := portMap[name]; known { 
					targets[name] = struct{ Port, Generation int }{port, e.Generation} 
				} 
			} 
		}
		stateMu.RUnlock()
		
		if len(targets) == 0 { continue }
		var wg sync.WaitGroup; var panicMu sync.Mutex; panicRequired := false
		for name, target := range targets {
			wg.Add(1); go func(n string, t struct{ Port, Generation int }) {
				defer wg.Done(); client := &http.Client{Timeout: 2 * time.Second}
				endpoint := "health"; if n == "unity" { endpoint = "engine/heartbeat" }
				resp, err := client.Get(fmt.Sprintf("http://localhost:%d/%s", t.Port, endpoint))
				engineFailed := false
				if err != nil || resp.StatusCode != 200 { engineFailed = true }
			if engineFailed { stateMu.Lock(); if e, ok := engines[n]; ok { e.State = StatePanic }; stateMu.Unlock(); panicMu.Lock(); panicRequired = true; panicMu.Unlock() } else { stateMu.Lock(); if e, ok := engines[n]; ok { e.TrustExpiry = time.Now().Add(60 * time.Minute) }; stateMu.Unlock(); resp.Body.Close() }
			}(name, target)
		}
		wg.Wait(); if panicRequired { 
			for name, e := range engines { 
				if name == "blender" && e.State == StateStopped { continue } // Ignore offline Blender
				go sendToEngine(name, "panic", "POST", map[string]interface{}{"reason": "HEARTBEAT_TIMEOUT"}) 
			} 
		}
	}
}

// --- TOOLS ---

func verifyAdapterIntegrity(target string) (string, error) {
	path := "unity-bridge/Editor/VibeBridgeServer.cs"; if target == "blender" { path = "blender-bridge/bridge_server.py" }
	data, err := os.ReadFile("../" + path); if err != nil { data, err = os.ReadFile(path); if err != nil { return "", err } }
	h := sha256.New(); h.Write(data); return hex.EncodeToString(h.Sum(nil)), nil
}

func stabilize_and_start(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	updateBridgeActivity("KERNEL: STABILIZING_ENVIRONMENT")
	cmd := "python3 scripts/preflight.py"
	out, err := execCommand(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("PREFLIGHT_CRITICAL_FAILURE: %v", err)
	}
	var report map[string]interface{}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		return nil, nil, fmt.Errorf("PREFLIGHT_PARSE_ERROR")
	}
	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult(report), nil, nil
}

func get_bridge_pulse(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	stateMu.RLock()
	uState := engines["unity"].State
	bState := engines["blender"].State
	stateMu.RUnlock()

	entropyMu.Lock()
	eUsed := sessionEntropy
	entropyMu.Unlock()

	walMu.Lock()
	hash := lastWalHash
	walMu.Unlock()

	if len(hash) > 8 {
		hash = hash[:8]
	}

	pulse := fmt.Sprintf("[KERNEL: READY | UNITY: %s | BLENDER: %s | ENTROPY: %d/%d | WAL: %s]",
		uState, bState, eUsed, maxEntropy, hash)
	
	return wrapForensicResult(pulse), nil, nil
}

func verify_identity_parity(ctx context.Context, req *mcp.CallToolRequest, args struct {
	IDs []string `json:"ids"`
}) (*mcp.CallToolResult, any, error) {
	updateBridgeActivity("KERNEL: VERIFYING_IDENTITY_PARITY")
	
	results := make(map[string]map[string]string)
	
	uRes, _ := sendToEngine("unity", "object/exists", "POST", map[string]interface{}{"ids": args.IDs})
	bRes, _ := sendToEngine("blender", "object/exists", "POST", map[string]interface{}{"ids": args.IDs})
	
	results["unity"] = make(map[string]string)
	if uRes != nil && uRes["exists"] != nil {
		for k, v := range uRes["exists"].(map[string]interface{}) {
			results["unity"][k] = fmt.Sprintf("%v", v)
		}
	}

	results["blender"] = make(map[string]string)
	if bRes != nil && bRes["exists"] != nil {
		for k, v := range bRes["exists"].(map[string]interface{}) {
			results["blender"][k] = fmt.Sprintf("%v", v)
		}
	}

	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult(results), nil, nil
}

func handshake_init(ctx context.Context, req *mcp.CallToolRequest, args HandshakeInitArgs) (*mcp.CallToolResult, any, error) {
	updateBridgeActivity(fmt.Sprintf("KERNEL: HANDSHAKE_%s", strings.ToUpper(args.Target)))
	stateMu.Lock(); engines[args.Target].State, engines[args.Target].Generation = StateStarting, engines[args.Target].Generation+1; newToken, chal := uuid.New().String(), uuid.New().String(); stateMu.Unlock()
	
	endpoint := "handshake"
	if args.Target == "unity" {
		endpoint = "status" // Real VibeBridge uses status for health/handshake
	}

	res, err := sendToEngine(args.Target, endpoint, "POST", map[string]interface{}{"version": args.Version, "new_token": newToken, "challenge": chal})
	
	// Real VibeBridge returns {"status":"ok"} for /status
	if err == nil && args.Target == "unity" && res["status"] == "ok" {
		stateMu.Lock(); e := engines[args.Target]; e.State, e.TrustExpiry = StateRunning, time.Now().Add(60*time.Minute); stateMu.Unlock()
		dispatchVibeEvent(LevelInfo, "handshake_complete", "", "READY", map[string]interface{}{"target": args.Target}); saveState()
		updateBridgeActivity("KERNEL: READY")
		return wrapForensicResult("OK"), nil, nil
	}

	if err != nil || (args.Target != "unity" && res["response"] != "VIBE_HASH_"+chal) { return nil, nil, fmt.Errorf("AUTH_FAILED") }
	stateMu.Lock(); e := engines[args.Target]; e.Token, e.Version, e.State, e.TrustExpiry = newToken, fmt.Sprintf("%v", res["engine_version"]), StateRunning, time.Now().Add(60*time.Minute); stateMu.Unlock()
	dispatchVibeEvent(LevelInfo, "handshake_complete", "", "READY", map[string]interface{}{"target": args.Target}); saveState()
	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult("OK"), nil, nil
}

func read_engine_state(ctx context.Context, req *mcp.CallToolRequest, args ReadStateArgs) (*mcp.CallToolResult, any, error) {
	endpoint := "state/get"; if args.Target == "unity" { endpoint = "scene/state" }
	res, _ := sendToEngine(args.Target, endpoint, "GET", nil); return wrapForensicResult(res), nil, nil
}

func verify_engine_state(ctx context.Context, req *mcp.CallToolRequest, args VerifyStateArgs) (*mcp.CallToolResult, any, error) {
	res, _ := sendToEngine(args.Target, "state/get", "GET", nil); if fmt.Sprintf("%v", res["hash"]) == args.ExpectedHash { return wrapForensicResult("VERIFIED"), nil, nil }; return nil, nil, fmt.Errorf("DRIFT_DETECTED")
}

func submit_intent(ctx context.Context, req *mcp.CallToolRequest, args SubmitIntentArgs) (*mcp.CallToolResult, any, error) {
	if args.Envelope.Rationale == "" || args.Envelope.Provenance == "" { return nil, nil, fmt.Errorf("TECHNICAL_RATIONALE_REQUIRED") }
	id := uuid.New().String(); txMu.Lock(); intents[id] = args.Envelope; txMu.Unlock()
	return wrapForensicResult(id), nil, nil
}

func validate_intent(ctx context.Context, req *mcp.CallToolRequest, args struct{ID string `json:"intent_id"`}) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); intent, ok := intents[args.ID]; txMu.Unlock(); if !ok { return nil, nil, fmt.Errorf("UNKNOWN_INTENT") }
	if intent.Confidence < 0.8 { stateMu.Lock(); for n := range engines { engines[n].State = StateHumanReq }; stateMu.Unlock(); return wrapForensicResult("HUMAN_INTERVENTION_REQUIRED"), nil, nil }
	return wrapForensicResult("ALLOW"), nil, nil
}

func human_approve_intent(ctx context.Context, req *mcp.CallToolRequest, args struct{ID string `json:"intent_id"`}) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); for n := range engines { if engines[n].State == StateHumanReq { engines[n].State = StateRunning } }; stateMu.Unlock(); return wrapForensicResult("APPROVED"), nil, nil
}

func begin_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args AtomicOpArgs) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); defer txMu.Unlock(); tx := &VibeTransaction{ID: uuid.New().String(), IntentID: args.IntentID, StartTime: time.Now(), Status: "OPEN"}; transactions[args.IntentID], activeTransaction = tx, tx
	return wrapForensicResult("TX_OPEN"), nil, nil
}

func commit_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args CommitAtomicOpArgs) (*mcp.CallToolResult, any, error) {
	updateBridgeActivity("KERNEL: AUDITING_INTEGRITY")
	
	// Mechanical Review (Iron Box Protocol)
	_, err := execCommand("python3 ../security_gate.py")
	if err != nil {
		updateBridgeActivity("KERNEL: AUDIT_FAILED")
		return nil, nil, fmt.Errorf("MECHANICAL_AUDIT_FAILED: Security violation detected during transaction. Check security_gate.py output.")
	}

	txMu.Lock(); defer txMu.Unlock(); if args.ProofOfWork == "" { return nil, nil, fmt.Errorf("INVARIANT_VIOLATION: ProofOfWork Required") }
	delete(transactions, args.IntentID); activeTransaction = nil; 
	
	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult("COMMITTED"), nil, nil
}

func abort_atomic_operation(ctx context.Context, req *mcp.CallToolRequest, args AtomicOpArgs) (*mcp.CallToolResult, any, error) {
	txMu.Lock(); defer txMu.Unlock(); delete(transactions, args.IntentID); activeTransaction = nil; return wrapForensicResult("ABORTED"), nil, nil
}

func emit_diag_bundle(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	p := filepath.Join(PersistenceDir, "diag.zip"); f, _ := os.Create(p); w := zip.NewWriter(f); for _, fn := range []string{WalFile, EventFile, StateFile} { df, _ := os.Open(fn); zw, _ := w.Create(filepath.Base(fn)); io.Copy(zw, df); df.Close() }; w.Close(); f.Close()
	return wrapForensicResult(p), nil, nil
}

func lock_object(ctx context.Context, req *mcp.CallToolRequest, args LockObjectArgs) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine(args.Target, "object/lock", "POST", map[string]interface{}{"id": args.ObjectID, "locked": args.Locked}); if err != nil { return nil, nil, err }; return wrapForensicResult(res), nil, nil
}

func get_metrics(ctx context.Context, req *mcp.CallToolRequest, args struct{Target string `json:"target"`}) (*mcp.CallToolResult, any, error) {
	res, err := sendToEngine(args.Target, "metrics", "GET", nil); if err != nil { return nil, nil, err }; return wrapForensicResult(res), nil, nil
}

func sync_material(ctx context.Context, req *mcp.CallToolRequest, args SyncMaterialArgs) (*mcp.CallToolResult, any, error) {
	if err := checkHumanLock(args.ObjectID); err != nil { return nil, nil, err }
	data := map[string]interface{}{"id": args.ObjectID, "properties": args.Props}; journalOperation(map[string]interface{}{"type": "intent", "op": "sync_material", "id": args.ObjectID}); sendToEngine("unity", "material/update", "POST", data); sendToEngine("blender", "material/update", "POST", data)
	return wrapForensicResult("OK"), nil, nil
}

func sync_transform(ctx context.Context, req *mcp.CallToolRequest, args SyncTransformArgs) (*mcp.CallToolResult, any, error) {
	if err := checkHumanLock(args.ObjectID); err != nil { return nil, nil, err }
	for _, v := range append(append(args.Position, args.Rotation...), args.Scale...) { if math.IsNaN(v) || math.IsInf(v, 0) { return nil, nil, fmt.Errorf("NUMERICAL_INSTABILITY") } }
	
	data := map[string]interface{}{"id": args.ObjectID, "transform": map[string]interface{}{"pos": args.Position, "rot": args.Rotation, "sca": args.Scale}}
	
	// Mechanical Floor: Buffer the intent for coalescing
	bufferSpeculativeIntent(args.ObjectID, "sync_transform", data)
	
	// Speculative Execution: Background send to engines
	go sendToEngine("unity", "transform/set", "POST", data)
	go sendToEngine("blender", "transform/set", "POST", data)
	
	journalOperation(map[string]interface{}{
		"type": "intent", 
		"op": "sync_transform", 
		"id": args.ObjectID, 
		"phase": PhaseProvisional,
		"class": ClassCosmetic,
	})
	
	return wrapForensicResult("PROVISIONAL_OK"), nil, nil
}

func sync_camera(ctx context.Context, req *mcp.CallToolRequest, args SyncCameraArgs) (*mcp.CallToolResult, any, error) {
	t := "unity"; if args.Source == "unity" { t = "blender" }; res, err := sendToEngine(args.Source, "camera/get", "GET", nil); if err != nil { return nil, nil, err }; sendToEngine(t, "camera/set", "POST", res)
	return wrapForensicResult("OK"), nil, nil
}

func sync_selection(ctx context.Context, req *mcp.CallToolRequest, args SyncSelectionArgs) (*mcp.CallToolResult, any, error) {
	t := "unity"; if args.Source == "unity" { t = "blender" }; sendToEngine(t, "selection/set", "POST", map[string]interface{}{"ids": args.IDs})
	return wrapForensicResult("OK"), nil, nil
}

func sync_asset_atomic(ctx context.Context, req *mcp.CallToolRequest, args SyncAssetAtomicArgs) (*mcp.CallToolResult, any, error) {
	updateBridgeActivity("KERNEL: SYNCING_ASSET_ATOMIC")
	pre, _ := sendToEngine("blender", "preflight/run", "POST", map[string]interface{}{"path": args.AssetPath}); ex, _ := sendToEngine("blender", "export", "POST", map[string]interface{}{"path": args.AssetPath}); sendToEngine("unity", "import", "POST", map[string]interface{}{"path": args.AssetPath, "meta": ex["meta"], "mode": "sandbox"}); val, _ := sendToEngine("unity", "validate", "POST", map[string]interface{}{"path": args.AssetPath})
	if fmt.Sprintf("%v", pre["hash"]) != fmt.Sprintf("%v", val["hash"]) { sendToEngine("unity", "rollback", "POST", map[string]interface{}{"path": args.AssetPath}); stateMu.Lock(); for n := range engines { engines[n].State = StateDesync }; stateMu.Unlock(); updateBridgeActivity("KERNEL: DESYNC"); return nil, nil, fmt.Errorf("HASH_MISMATCH") }
	sendToEngine("unity", "commit", "POST", map[string]interface{}{"path": args.AssetPath})
	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult("SYNCED"), nil, nil
}

func get_operation_journal(ctx context.Context, req *mcp.CallToolRequest, args struct{Limit int `json:"limit"`}) (*mcp.CallToolResult, any, error) {
	f, _ := os.Open(WalFile); s := bufio.NewScanner(f); var out []string; for s.Scan() { out = append(out, s.Text()) }; if len(out) > args.Limit && args.Limit > 0 { out = out[len(out)-args.Limit:] }
	return wrapForensicResult(strings.Join(out, "\n")), nil, nil
}

func control_playback(ctx context.Context, req *mcp.CallToolRequest, args struct { Action string; Time float64 }) (*mcp.CallToolResult, any, error) {
	d := map[string]interface{}{"action": args.Action, "time": args.Time}; sendToEngine("unity", "playback/control", "POST", d); sendToEngine("blender", "playback/control", "POST", d)
	return wrapForensicResult("OK"), nil, nil
}

func global_id_map_resolve(ctx context.Context, req *mcp.CallToolRequest, args MapVibeIDsArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); globalIDMap[args.UnityGUID], globalIDMap[args.BlenderName] = args.BlenderName, args.UnityGUID; stateMu.Unlock(); return wrapForensicResult("RESOLVED"), nil, nil
}

func vibe_multiplex(ctx context.Context, req *mcp.CallToolRequest, args MultiplexCallArgs) (*mcp.CallToolResult, any, error) {
	allowed, ok := drivers[args.SensorID]; if !ok { return nil, nil, fmt.Errorf("DRIVER_UNREGISTERED") }; isOk := false; for _, ep := range allowed { if ep == args.Endpoint { isOk = true; break } }; if !isOk { return nil, nil, fmt.Errorf("DENIED") }
	res, err := sendToEngine(args.Target, args.Endpoint, "POST", args.Payload); if err != nil { return nil, nil, err }; return wrapForensicResult(res), nil, nil
}

func set_engine_state(ctx context.Context, req *mcp.CallToolRequest, args SetEngineStateArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); if e, ok := engines[args.Target]; ok { e.State = EngineState(args.State) }; stateMu.Unlock(); saveState(); return wrapForensicResult("OK"), nil, nil
}

func revoke_id(ctx context.Context, req *mcp.CallToolRequest, args RevokeIDArgs) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); revocationList[args.ID] = args.Reason; stateMu.Unlock(); return wrapForensicResult("REVOKED"), nil, nil
}

func epistemic_refusal(ctx context.Context, req *mcp.CallToolRequest, args struct{R string `json:"reason"`}) (*mcp.CallToolResult, any, error) {
	return wrapForensicResult("UNKNOWABLE: " + args.R), nil, nil
}

func decommission_bridge(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	stateMu.Lock(); for n := range engines { engines[n].State = StatePanic }; stateMu.Unlock(); return wrapForensicResult("DECOMMISSIONED"), nil, nil
}

func reconstruct_state(ctx context.Context, req *mcp.CallToolRequest, args ForensicReplayArgs) (*mcp.CallToolResult, any, error) {
	f, _ := os.Open(EventFile); s := bufio.NewScanner(f); var out []string; for s.Scan() { var e VibeEvent; json.Unmarshal(s.Bytes(), &e); out = append(out, e.Type) }
	return wrapForensicResult(strings.Join(out, " -> ")), nil, nil
}

func invoke_specialist(ctx context.Context, req *mcp.CallToolRequest, args InvokeSpecialistArgs) (*mcp.CallToolResult, any, error) {
	dispatchVibeEvent(LevelInfo, "specialist_invoked", args.IntentID, "DELEGATE", map[string]interface{}{"specialist": args.SpecialistID, "hash": args.CurrentHash})
	p := fmt.Sprintf("STATELESS SPECIALIST INVOCATION\nTarget: %s\nCurrent Hash: %s\nIntent: %v", args.SpecialistID, args.CurrentHash, args.TargetIntent)
	return wrapForensicResult(p), nil, nil
}

func apply_lock(ctx context.Context, req *mcp.CallToolRequest, args ApplyLockArgs) (*mcp.CallToolResult, any, error) {
	lockMu.Lock()
	defer lockMu.Unlock()
	lockTable[args.UUID] = &VibeLock{
		UUID:      args.UUID,
		Type:      args.LockType,
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Second),
	}
	updateBridgeActivity(fmt.Sprintf("KERNEL: LOCK_%s", args.UUID))
	return wrapForensicResult("LOCKED"), nil, nil
}

func release_lock(ctx context.Context, req *mcp.CallToolRequest, args ReleaseLockArgs) (*mcp.CallToolResult, any, error) {
	lockMu.Lock()
	defer lockMu.Unlock()
	delete(lockTable, args.UUID)
	updateBridgeActivity("KERNEL: READY")
	return wrapForensicResult("RELEASED"), nil, nil
}

func get_bridge_heartbeat(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	stateMu.RLock(); defer stateMu.RUnlock()
	res := BridgeHeartbeat{
		BridgePID:             os.Getpid(),
		UptimeSec:             int(time.Since(startTime).Seconds()),
		EpochID:               monotonicID,
		OrchestratorConnected: true,
		UnityConnected:        engines["unity"].State == StateRunning,
		BlenderConnected:      engines["blender"].State == StateRunning,
		LastTickHash:          lastWalHash,
		ExpectedIntervalMS:    5000,
		LastSeenMS:            500,
	}
	return wrapForensicResult(res), nil, nil
}

func get_bridge_handshake_state(ctx context.Context, req *mcp.CallToolRequest, args struct{ AssetID string `json:"asset_id"` }) (*mcp.CallToolResult, any, error) {
	u, _ := sendToEngine("unity", "state/get", "GET", nil); b, _ := sendToEngine("blender", "state/get", "GET", nil); m := fmt.Sprintf("%v", u["hash"]) == fmt.Sprintf("%v", b["hash"])
	res := BridgeHandshakeState{AssetID: args.AssetID, BlenderExportHash: fmt.Sprintf("%v", b["hash"]), UnityImportHash: fmt.Sprintf("%v", u["hash"]), HashMatch: m, LastVerified: time.Now().Format(time.RFC3339)}
	return wrapForensicResult(res), nil, nil
}

func get_bridge_wal_state(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	res := BridgeWalState{WalHead: monotonicID, WalHash: lastWalHash, LastCommittedOp: "UNKNOWN", PendingOps: 0, RollbackAvailable: true, Reversible: true}
	return wrapForensicResult(res), nil, nil
}

func get_bridge_commit_requirements(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	u, _ := sendToEngine("unity", "state/get", "GET", nil); b, _ := sendToEngine("blender", "state/get", "GET", nil); h := map[string]string{"wal": lastWalHash, "blender": fmt.Sprintf("%v", b["hash"]), "unity": fmt.Sprintf("%v", u["hash"])}
	res := BridgeCommitRequirements{RequiredHashes: h, RationaleRequired: true, CommitAllowed: h["blender"] == h["unity"]}
	return wrapForensicResult(res), nil, nil
}

func execute_governed_mutation(ctx context.Context, req *mcp.CallToolRequest, args MutateArgs) (*mcp.CallToolResult, any, error) {
	if err := auditPayload(args.OpSpec); err != nil { return nil, nil, err }
	t, _ := args.OpSpec["target"].(string); e, _ := args.OpSpec["endpoint"].(string)
	targetHash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", args.OpSpec["payload"]))))
	if err := checkInvariants(args.IdempotencyKey, targetHash); err != nil { return nil, nil, err }
	res, err := sendToEngine(t, e, "POST", args.OpSpec["payload"]); if err != nil { return nil, nil, err }
	time.Sleep(200 * time.Millisecond); v, _ := sendToEngine(t, "state/get", "GET", nil)
	r := map[string]interface{}{"engine_response": res, "verified_hash": "FAIL"}; if v != nil { r["verified_hash"] = v["hash"] }
	return wrapForensicResult(r), nil, nil
}

// --- STATE PERSISTENCE ---

func loadState() {
	stateMu.Lock(); defer stateMu.Unlock(); data, _ := os.ReadFile(StateFile); var s struct { Engines map[string]*EngineData; IDMap map[string]string; Credits int }; json.Unmarshal(data, &s); creditBalance, globalIDMap = s.Credits, s.IDMap
	for k, v := range s.Engines { if e, ok := engines[k]; ok { e.State, e.Token, e.Version = v.State, v.Token, v.Version } }
}

func saveState() {
	stateMu.RLock(); s := struct { Engines map[string]*EngineData `json:"engines"`; IDMap map[string]string `json:"id_map"`; Credits int `json:"credits"` }{Engines: engines, IDMap: globalIDMap, Credits: creditBalance}; stateMu.RUnlock(); data, _ := json.MarshalIndent(s, "", "  "); os.WriteFile(StateFile, data, 0644)
}

func updateBridgeActivity(activity string) {
	os.MkdirAll("metadata", 0755)
	os.WriteFile(ActivityFile, []byte(activity), 0644)
}

func execCommand(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func journalOperation(op map[string]interface{}) {
	walMu.Lock(); defer walMu.Unlock(); if activeTransaction != nil { op["tid"] = activeTransaction.ID }
	op["prev_hash"] = lastWalHash; data, _ := json.Marshal(op); h := sha256.New(); h.Write(data); lastWalHash = hex.EncodeToString(h.Sum(nil)); op["hash"] = lastWalHash
	final, _ := json.Marshal(op); if info, err := os.Stat(WalFile); err == nil && info.Size() > MaxWalSize { os.Rename(WalFile, WalFile+".old") }; f, _ := os.OpenFile(WalFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); defer f.Close(); f.Write(final); f.Write([]byte("\n"))
}

func startTransactionGC() {

	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {

		txMu.Lock()

		now := time.Now()

		for id, tx := range transactions {

			if now.Sub(tx.StartTime) > 60*time.Second {

				log.Printf("🚨 VibeSync: Transaction Timeout (%s) - Auto-Rolling Back", tx.ID)

				delete(transactions, id)

				if activeTransaction == tx {

					activeTransaction = nil

				}

				// Force engine rollback

				for name := range engines {

					go sendToEngine(name, "rollback", "POST", map[string]interface{}{"reason": "TX_TIMEOUT"})

				}

			}

		}

		txMu.Unlock()

	}

}



func main() {

	server := mcp.NewServer(&mcp.Implementation{Name: "VibeSync", Version: "v0.4.0"}, &mcp.ServerOptions{})

	mcp.AddTool(server, &mcp.Tool{Name: "stabilize_and_start", Description: "Self-Healing Bootstrap"}, stabilize_and_start)

	mcp.AddTool(server, &mcp.Tool{Name: "get_bridge_pulse", Description: "Bridge Pulse"}, get_bridge_pulse)

	mcp.AddTool(server, &mcp.Tool{Name: "verify_identity_parity", Description: "UUID Parity"}, verify_identity_parity)

	mcp.AddTool(server, &mcp.Tool{Name: "handshake_init", Description: "ISA 1"}, handshake_init)

	mcp.AddTool(server, &mcp.Tool{Name: "read_engine_state", Description: "ISA 2"}, read_engine_state)

	mcp.AddTool(server, &mcp.Tool{Name: "verify_engine_state", Description: "ISA 3"}, verify_engine_state)

	mcp.AddTool(server, &mcp.Tool{Name: "submit_intent", Description: "ISA 4"}, submit_intent)

	mcp.AddTool(server, &mcp.Tool{Name: "validate_intent", Description: "ISA 5"}, validate_intent)

	mcp.AddTool(server, &mcp.Tool{Name: "human_approve_intent", Description: "ISA 5b"}, human_approve_intent)

	mcp.AddTool(server, &mcp.Tool{Name: "begin_atomic_operation", Description: "ISA 6"}, begin_atomic_operation)

	mcp.AddTool(server, &mcp.Tool{Name: "commit_atomic_operation", Description: "ISA 7"}, commit_atomic_operation)

	mcp.AddTool(server, &mcp.Tool{Name: "abort_atomic_operation", Description: "ISA 8"}, abort_atomic_operation)

	mcp.AddTool(server, &mcp.Tool{Name: "emit_diag_bundle", Description: "ISA 10"}, emit_diag_bundle)

	mcp.AddTool(server, &mcp.Tool{Name: "lock_object", Description: "Locking"}, lock_object)

	mcp.AddTool(server, &mcp.Tool{Name: "apply_lock", Description: "Human Lock"}, apply_lock)

	mcp.AddTool(server, &mcp.Tool{Name: "release_lock", Description: "Release Lock"}, release_lock)

	mcp.AddTool(server, &mcp.Tool{Name: "get_metrics", Description: "Metrics"}, get_metrics)

	mcp.AddTool(server, &mcp.Tool{Name: "sync_transform", Description: "ISA 21*"}, sync_transform)

	mcp.AddTool(server, &mcp.Tool{Name: "sync_material", Description: "ISA 22"}, sync_material)

	mcp.AddTool(server, &mcp.Tool{Name: "sync_camera", Description: "ISA 23"}, sync_camera)

	mcp.AddTool(server, &mcp.Tool{Name: "sync_selection", Description: "ISA 24"}, sync_selection)

	mcp.AddTool(server, &mcp.Tool{Name: "sync_asset_atomic", Description: "ISA 26"}, sync_asset_atomic)

	mcp.AddTool(server, &mcp.Tool{Name: "get_operation_journal", Description: "WAL"}, get_operation_journal)

	mcp.AddTool(server, &mcp.Tool{Name: "control_playback", Description: "Timeline"}, control_playback)

	mcp.AddTool(server, &mcp.Tool{Name: "global_id_map_resolve", Description: "ISA 27"}, global_id_map_resolve)

	mcp.AddTool(server, &mcp.Tool{Name: "vibe_multiplex", Description: "Multiplex"}, vibe_multiplex)

	mcp.AddTool(server, &mcp.Tool{Name: "set_engine_state", Description: "State"}, set_engine_state)

	mcp.AddTool(server, &mcp.Tool{Name: "revoke_id", Description: "Revoke"}, revoke_id)

	mcp.AddTool(server, &mcp.Tool{Name: "epistemic_refusal", Description: "ISA 29"}, epistemic_refusal)

	mcp.AddTool(server, &mcp.Tool{Name: "decommission_bridge", Description: "ISA 32"}, decommission_bridge)

	mcp.AddTool(server, &mcp.Tool{Name: "reconstruct_state", Description: "Forensic"}, reconstruct_state)

	mcp.AddTool(server, &mcp.Tool{Name: "invoke_specialist", Description: "Delegate"}, invoke_specialist)

	mcp.AddTool(server, &mcp.Tool{Name: "get_bridge_heartbeat", Description: "Bridge Liveness"}, get_bridge_heartbeat)

	mcp.AddTool(server, &mcp.Tool{Name: "get_bridge_handshake_state", Description: "Reality Align"}, get_bridge_handshake_state)

	mcp.AddTool(server, &mcp.Tool{Name: "get_bridge_wal_state", Description: "WAL State"}, get_bridge_wal_state)

	mcp.AddTool(server, &mcp.Tool{Name: "get_bridge_commit_requirements", Description: "Semantic Gate"}, get_bridge_commit_requirements)

	mcp.AddTool(server, &mcp.Tool{Name: "execute_governed_mutation", Description: "Gov Mutate"}, execute_governed_mutation)



	// Start Background Services (Flow Amplifiers)

	go startControlPlane()

	go startTransactionGC()



	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {


				log.Printf("MCP Server stopped: %v", err)
			}
			
			log.Printf("🛡️ Orchestrator entering background mode (Control Plane active)")
			select {} // Keep alive for Control Plane
		}
		
		func startControlPlane() {
		mux := http.NewServeMux()
			mux.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
				stateMu.RLock()
				uState := engines["unity"].State
				bState := engines["blender"].State
				token := engines["unity"].Token
				stateMu.RUnlock()
		
				walMu.Lock()
				hash := lastWalHash
				walMu.Unlock()
		
				entropyMu.Lock()
				eUsed := sessionEntropy
				entropyMu.Unlock()
		
				status := map[string]interface{}{
					"kernel":   "READY",
					"unity":    map[string]interface{}{"state": uState, "port": unityPort, "token": token},
					"blender":  map[string]interface{}{"state": bState, "port": BlenderPort},
					"wal_hash": hash,
					"entropy":  fmt.Sprintf("%d/%d", eUsed, maxEntropy),
					"uptime":   time.Since(startTime).String(),
				}
		
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(status)
			})
		
			mux.HandleFunc("/recover", func(w http.ResponseWriter, r *http.Request) {
				stateMu.Lock()
				for n := range engines {
					engines[n].State = StateStopped
				}
				stateMu.Unlock()
				log.Printf("🛡️ VibeSync: Manual Recovery Triggered - All engines reset to STOPPED")
				w.Write([]byte("RECOVERY_INITIATED"))
			})
					mux.HandleFunc("/activity", func(w http.ResponseWriter, r *http.Request) {
	
				data, _ := os.ReadFile(ActivityFile)
	
				w.Header().Set("Content-Type", "text/plain")
	
				w.Write(data)
	
			})
	
		
	
			mux.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
	
				var call struct {
	
					Name      string          `json:"name"`
	
					Arguments json.RawMessage `json:"arguments"`
	
				}
	
				if err := json.NewDecoder(r.Body).Decode(&call); err != nil {
	
					http.Error(w, err.Error(), 400)
	
					return
	
				}
	
		
	
				// Simple dispatcher for critical tools
	
				// Note: In a full proxy, we would use reflection or the server's own registry
	
				var res interface{}
	
				var err error
	
		
	
				switch call.Name {
	
				case "handshake_init":
	
					var args HandshakeInitArgs
	
					json.Unmarshal(call.Arguments, &args)
	
					res, _, err = handshake_init(context.Background(), nil, args)
	
				case "get_bridge_pulse":
	
					res, _, err = get_bridge_pulse(context.Background(), nil, struct{}{})
	
				case "stabilize_and_start":
	
					res, _, err = stabilize_and_start(context.Background(), nil, struct{}{})
	
				default:
	
					http.Error(w, "Tool not supported via proxy yet", 404)
	
					return
	
				}
	
		
	
				if err != nil {
	
					http.Error(w, err.Error(), 500)
	
					return
	
				}
	
		
	
				w.Header().Set("Content-Type", "application/json")
	
				json.NewEncoder(w).Encode(res)
	
			})
	
		
	
			log.Printf("📡 VibeSync Control Plane: Listening on http://localhost:8080")
	
		
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("🚨 Control Plane Error: %v", err)
		}
	}
	