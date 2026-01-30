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
	"context"
	"testing"
	"github.com/google/uuid"
)

func TestIntentConfidenceGate(t *testing.T) {
	// Initialize some state
	intents = make(map[string]IntentEnvelope)
	
	// 1. Test High Confidence
	highID := uuid.New().String()
	intents[highID] = IntentEnvelope{
		Rationale: "Testing high confidence",
		Confidence: 0.9,
	}
	
	_, _, err := validate_intent(context.Background(), nil, struct{ID string `json:"intent_id"`}{ID: highID})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestNumericalSafety(t *testing.T) {
	// Test auditPayload for NaN/Inf
	// Using a string payload that contains "nan" to trigger the auditPayload check
	badData := map[string]interface{}{
		"val": "this is nan",
	}
	
	err := auditPayload(badData)
	if err == nil {
		t.Error("Expected error for NaN payload, got nil")
	}
}

func TestSecurityAudit(t *testing.T) {
	// Test auditPayload for blocked commands
	evilData := map[string]interface{}{
		"cmd": "os.system('rm -rf /')",
	}
	
	err := auditPayload(evilData)
	if err == nil {
		t.Error("Expected error for malicious payload, got nil")
	}
}
