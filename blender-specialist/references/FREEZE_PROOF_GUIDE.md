# ❄️ Ultimate VibeBridge Freeze-Proof Guide (v0.4)

**Goal:** Ensure Unity, Blender, and the VibeSync bridge never freeze, hang, or deadlock during live-sync operations.

---

## 🏗️ 1. Freeze-Proof Architecture: The Cheat Sheet

| Freeze Source | Primary Cause | Professional Fix | Example Code |
| :--- | :--- | :--- | :--- |
| **Heartbeat Loops** | Synchronous wait for response. | Background thread + Timeout. | `cts.CancelAfter(1000)` |
| **Network I/O** | Blocking socket reads/writes. | Non-blocking / Async sockets. | `await networkClient.Send()` |
| **Asset Transfer** | Synchronous heavy file copy. | Async transfer + Atomic Swap. | `File.Move(tmp, live)` |
| **Engine Launch** | `Process.Start()` on main thread. | `Task.Run()` / `Thread` spawn. | `Task.Run(() => Process.Start())` |
| **Main Thread Access**| Thread-unsafe API calls. | Concurrent Queues + Dispatcher. | `queue.Enqueue(() => { ... })` |
| **Asset DB Lock** | Unity/Blender file locks. | Lock checking + Atomic Rename. | `if (!IsFileLocked(path))` |
| **Sync Loops** | Version drift / Conflict. | Strict version/hash checking. | `if (src.hash != dst.hash)` |
| **Console Flood** | High-frequency logging. | Rate-limited logging hooks. | `if (time > lastLog + 1.0)` |

---

## 📜 2. Detailed Technical Procedures

### A. Heartbeat & Failure Recovery
Synchronous heartbeats are the #1 cause of bridge hangs.
- **Fix**: Use `async/await` with a `CancellationToken` linked to a timeout.
```csharp
// Unity async heartbeat pattern
async Task SendHeartbeatAsync(CancellationToken token) {
    try {
        using var cts = CancellationTokenSource.CreateLinkedTokenSource(token);
        cts.CancelAfter(1000); // 1-second timeout
        await networkClient.Ping(cts.Token);
    } catch (OperationCanceledException) {
        Debug.LogWarning("🛡️ VibeSync: Heartbeat timed out. Triggering stability protocol.");
    } catch (Exception e) {
        Debug.LogError($"🛡️ VibeSync: Heartbeat failed: {e}");
    }
}
```

### B. Blender Main-Thread Integrity (Modal Pattern)
Never use `time.sleep()` in Blender.
- **Fix**: Use `bpy.app.timers` to schedule non-blocking calls.
```python
def heartbeat_modal():
    try:
        # Perform non-blocking logic (check queue, etc.)
        return 1.0  # Reschedule in 1 second
    except Exception as e:
        print(f"🚨 VibeSync Heartbeat Error: {e}")
        return 2.0  # Retry with backoff
```

### C. Atomic Asset Swap
Overwriting a live asset while Unity/Blender is reading it causes deadlocks.
- **Fix**: Write to a temporary file in `.vibesync/tmp`, verify the hash, then perform an atomic rename.

---

## 🛡️ 3. Bridge Freeze-Proof Principles (Mandatory)

1.  **Never block the main thread** (Unity or Blender).
2.  **Background all I/O**: Network, Disk, and Process spawning must be off-thread.
3.  **Strict Timeouts**: No external call should wait longer than 1000ms.
4.  **Queue-Based Mutation**: Only the Main Thread modifies objects; all other threads must use thread-safe queues.
5.  **Rate-Limit Logging**: Max 1 log/second for repetitive status updates to prevent console freeze.
6.  **Watchdog Monitoring**: If an engine doesn't respond to 3 consecutive heartbeats, trigger a `PANIC` and lock hierarchies.
7.  **Resource Cleanup**: Always use `using` blocks or `try/finally` to release sockets and file handles.

---
**VibeSync: Engineered for Stability.**
*Copyright (C) 2026 B-A-M-N*