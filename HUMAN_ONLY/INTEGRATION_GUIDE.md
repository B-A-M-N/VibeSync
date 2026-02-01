# Getting Started: VibeSync Engine Adapters

This guide explains how to connect your creative engines to the Go Orchestrator.

---

## 🎮 Unity Integration
1.  **Copy Files**: Move the `unity-bridge/` folder into your Unity Project's `Assets/` directory.
2.  **Verify Server**: Unity will automatically compile and start the `VibeBridgeServer` on port `8085`.
3.  **Check Console**: You should see `🛡️ VibeSync Unity Bridge: Listening on port 8085`.
4.  **Handshake**: Use the MCP tool `initiate_handshake(target="unity")` to establish trust.

---

## 🧊 Blender Integration
1.  **Install Addon**:
    *   Zip the `blender-bridge/` directory.
    *   In Blender, go to `Edit -> Preferences -> Add-ons -> Install`.
    *   Select the zip and enable "VibeSync MCP".
2.  **Verify Server**: Check the Blender console (Toggle System Console) for `🛡️ VibeSync Blender Bridge: Listening on port 22000`.
3.  **UI Panel**: Look for the "VibeSync" tab in the 3D Viewport sidebar (N-panel).
4.  **Handshake**: Use the MCP tool `initiate_handshake(target="blender")` to establish trust.

---

## 🧪 The First Sync
1.  Launch the Go Orchestrator: `cd mcp-server && go run watcher.go`.
2.  In your AI interface, call `initiate_handshake` for both engines.
3.  Execute `sync_transform` on a Cube to verify the link.
4.  Execute `sync_asset_atomic` to test the full validated pipeline.

---
*Copyright (C) 2026 B-A-M-N*
