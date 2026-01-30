# 🏗️ Contributing to the Iron Box

> [!IMPORTANT]
> **READ THIS BEFORE OPENING A PULL REQUEST.**
> We do not accept "vibes-based" code. This is a distributed state machine operating on destructive endpoints. Your code must be **defensive, typed, and paranoid.**

## ⚔️ The Code Philosophy: "Zero Trust"

We operate under a simple axiom: **The engines (Unity/Blender) are hostile.** They want to crash, they want to hang, and they want to corrupt user data.

Your code is the shield.

1.  **No "Happy Paths"**: Assume the socket is dead, the texture is null, and the user is spamming the button. Handle the edge case first.
2.  **Explicit > Implicit**: No magic auto-wiring. No hidden side effects. If a function mutates state, it must return a `Result` or `Error` type explicitly.
3.  **Performance is a Feature**: We operate in the 16ms (60fps) window. Allocations in hot paths are rejected.

---

## 🏛️ Language Standards

### 🐹 Go (The Orchestrator)
*The central nervous system. It must never crash.*

*   **Formatter**: `gofmt` (Standard 1.23+).
*   **Linter**: `golangci-lint` with strict settings.
*   **Error Handling**:
    *   ❌ **FORBIDDEN**: `_` (Blank identifier) for error returns. ALWAYS handle or bubble the error.
    *   ❌ **FORBIDDEN**: `panic()` in runtime code. Panics are allowed ONLY during `init()` or bootstrap.
    *   ✅ **REQUIRED**: Structured logging via `slog`.
*   **Concurrency**:
    *   Use `channels` for signaling, `mutexes` for state.
    *   No "fire and forget" goroutines without a `WaitGroup` or context cancellation.

### 🐍 Python (The Blender Adapter)
*The chaotic script environment. Must be disciplined.*

*   **Formatter**: `Black` (Line length: 88).
*   **Linter**: `Ruff` or `Flake8`.
*   **Type Hinting**: **MANDATORY**. All function signatures must have type hints (`def sync_mesh(name: str) -> bool:`).
*   **Blender Context**:
    *   Never assume `bpy.context` is correct. Always override context or verify `poll()` methods.
    *   Wrap all `bpy.ops` calls in `try/except` blocks to catch internal Blender failures.

### #️⃣ C# (The Unity Adapter)
*The simulation runtime. Watch the Garbage Collector.*

*   **Style**: Standard .NET / Unity Style (K&R braces).
*   **Memory Discipline**:
    *   ❌ **AVOID**: LINQ in `Update()` loops.
    *   ❌ **AVOID**: String concatenation in hot paths (use `StringBuilder`).
    *   ✅ **REQUIRED**: Use `using` blocks for `IDisposable` resources (Streams, WebRequests).
*   **Threading**:
    *   Unity API calls **MUST** run on the Main Thread.
    *   Use `EditorApplication.delayCall` or our internal `Dispatcher` for cross-thread marshalling.

---

## 🛡️ The Gauntlet (PR Requirements)

Before submitting a PR, verify it survives **The Gauntlet**:

1.  **The "Destruction" Test**: What happens if I force-quit Unity in the middle of your function? Does the Orchestrator recover, or does it hang forever waiting for a handshake?
2.  **The "Spam" Test**: What happens if I click the button 50 times in 1 second?
3.  **Documentation**:
    *   Public methods must have XML/Docstring comments explaining *what* they do and *what they mutate*.
    *   If you change the Protocol, you **MUST** update `ADAPTER_CONTRACT.md`.

---

## 📜 Legal & DCO

By contributing, you certify that you own the code you are submitting or have the right to contribute it under the **AGPLv3**.

**Sign-off Procedure**:
All commits must be signed off to certify the Developer Certificate of Origin (DCO).

```bash
git commit -s -m "feat: implement atomic handshake retry logic"
```

### Commit Message Convention
We follow **Conventional Commits**:
*   `feat:` New capability (e.g., "add material sync").
*   `fix:` Bug fix (e.g., "resolve nan coordinates").
*   `perf:` Optimization (e.g., "reduce allocation in update loop").
*   `docs:` Documentation only.
*   `refactor:` Code change that neither fixes a bug nor adds a feature.

---
**Welcome to the Iron Box. Build it strong.**