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

using UnityEngine;

public static class ModuleA
{
    public static void Init()
    {
        // This is called automatically by MCPDevAutoReload on script save
        Debug.Log("ModuleA: State Re-synchronized with VibeSync MCP Server");
        
        // Example: Cache common references or reset transient states
    }

    public static void PerformAction()
    {
        Debug.Log("ModuleA: Performing synced operation...");
    }
}
