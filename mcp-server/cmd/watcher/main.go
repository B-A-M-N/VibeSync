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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	log.Println("🚀 Starting VibeSync Go Watcher...")
	
	// Initial Build
	rebuild()
	
	// Start the server process
	cmd := startServer()
	
	lastMod := lastModified()
	
	for {
		time.Sleep(1 * time.Second)
		currentMod := lastModified()
		
		if currentMod.After(lastMod) {
			log.Println("🔄 Change detected. Restarting...")
			
			// Kill existing process
			if cmd != nil && cmd.Process != nil {
				cmd.Process.Kill()
				cmd.Wait()
			}
		
			rebuild()
			cmd = startServer()
			lastMod = currentMod
		}
	}
}

func lastModified() time.Time {
	var latest time.Time
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if filepath.Ext(path) == ".go" && info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		return nil
	})
	return latest
}

func rebuild() {
	log.Println("🔨 Rebuilding vibe-mcp-server...")
	// Adjusted paths for the new location cmd/watcher/
	cmd := exec.Command("go", "build", "-o", "../../vibe-mcp-server", "../../main.go", "../../contract.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ Build Failed: %v\n%s", err, out)
		return
	}
	log.Println("✅ Build Successful")
}

func startServer() *exec.Cmd {
	// Adjusted path to the binary
	cmd := exec.Command("../../vibe-mcp-server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Start()
	if err != nil {
		log.Printf("❌ Failed to start server: %v", err)
		return nil
	}
	log.Println("🛰️ Server Running")
	return cmd
}
