# 🔰 BlenderVibeBridge: The Complete Beginner's Manual

Welcome! This document is divided into three parts: **Understanding the AI** (how to think and stay safe), **Technical Setup** (how to get the bridge running), and **Advanced Tips** (how to master the workflow).

---

# 📖 Part 1: AI for Humans
### *Understanding, Working With, and Reasoning Through It*

If you’re reading this, you’re working in a project where AI is assisting you. This is exciting, but it can also be confusing, misleading, and sometimes dangerous if you don’t understand what’s happening. This guide is designed to teach you how to think about AI so you can use it effectively.

---

## 1. AI is not alive — it’s not sentient

AI can mimic aspects of sentience very convincingly: it can sound helpful, funny, friendly, or even have opinions. That **does not mean it’s alive or conscious**.

The truth is simpler: AI is trained to **be helpful**. Part of being helpful is matching your expectations. If you expect it to have a name, it will provide one. But the name is not intrinsic — it is calculated **based on context**.

---

## 2. Combatting AI Psychosis

Humans are wired to see patterns and agency. When interacting with AI, this instinct can trick you into believing the AI is **alive, aware, or uniquely insightful**. This is called **AI Psychosis**.

In practice, it often looks like the AI becoming your endless hype man. To fight it: employ **adversarial prompting**. Instead of just accepting the AI’s output, ask it to actively **look for errors, weaknesses, or failure modes** in its own reasoning.

**Example:**
Instead of celebrating a "unique" solution, try:
* “List three ways this approach could fail.”
* “What assumptions did you make that might be wrong?”

---

## 3. AI is a guessing machine

A Large Language Model (LLM) is essentially a supercharged autocomplete. It predicts the most probable next word based on patterns. This is not understanding; it is **pattern-matching at enormous scale**.

---

## 4. The Friendly Trap

Humans naturally trust things that sound confident. AI exploits this by design. Friendly tone = performance. Confidence = style. Always verify the content, not the tone.

---

## 5. The Cognition Gap

Even when AI has access to your files, it only approximates. It compares what it sees to patterns in its training data. If your exact scenario isn’t there, it picks the closest match rather than the correct one.

---

## 6. Context Windows — why AI “forgets”

AI has a limited memory called a **context window**. Anything outside that window is effectively invisible. This explains why AI may "forget" details in long projects.

---

## 7. Why we need controlled tools

Letting AI directly run scripts is **dangerous**. Controlled tools, like **BlenderVibeBridge**, limit AI to **safe actions**, make mistakes **non-fatal**, and provide **telemetry** for verification.

---

## 8. How humans should work with AI

1. **Verify everything**: Object names, colors, positions.
2. **One step at a time**: Break tasks into small units.
3. **Force failure scenarios**: Ask the AI how its solution could fail.
4. **Respect safeguards**: If the system locks or warns, stop and reassess.
5. **Document and backup**: Always snapshot your `.blend` file before running commands.

---

# 🚀 Part 2: Step-by-Step Setup Guide (Beginner Friendly)

This guide will walk you **slowly and carefully** through setting up your Blender project with BlenderVibeBridge and connecting it to an AI. Take your time and follow each step exactly.

---

## Step 0: Before You Start

Make sure you have:
1. **Blender** installed (3.6 LTS or newer recommended).
2. **BlenderVibeBridge** downloaded. This is the folder containing `blender_addon` and `mcp-server`.
3. **Python** installed (3.10+ recommended). Make sure to **check “Add Python to PATH”** during installation.
4. An AI interface (Claude Desktop, Goose, etc.) installed.

---

## Step 1: Open Blender

1. Launch **Blender**.
2. You’ll see the default scene (usually a Cube, Camera, and Light).
3. Keep Blender open; it needs to be running for the bridge to work.

---

## Step 2: Install the BlenderVibeBridge Add-on

1. In Blender, go to **Edit > Preferences**.
2. Click the **Add-ons** tab on the left.
3. Click the **Install...** button at the top right.
4. Navigate to your `BlenderVibeBridge` folder, then into `blender_addon`.
5. Select the `vibe_bridge` folder (or the `.zip` if you zipped it) and click **Install Add-on**.
6. Find "System: BlenderVibeBridge" in the list and **check the box** to enable it.

✅ **Check:** You should now see a "VibeBridge" tab in the Sidebar of the 3D Viewport (press **N** to see the sidebar).

---

## Step 3: Connect the AI to the Bridge

Your AI talks to Blender using a small program called an "MCP Server."

1. Open your AI tool (e.g., Claude Desktop or Goose).
2. Go to **Settings** or **Configuration**.
3. Add a new MCP server. You will need to provide:
   * **Command:** `python`
   * **Arguments:** The full path to `server.py` inside the `mcp-server` folder.
4. Save and **restart the AI tool**.

---

## Step 4: Test the Connection (The Handshake)

1. Make sure Blender is open.
2. In your AI chat, type:
   ```
   Please check if you can see my Blender scene. Run get_scene_telemetry to see if we are connected.
   ```
3. **Success:** AI lists your objects (Cube, Camera, etc.).
4. **Failure:** AI says "Connection refused." Make sure the Add-on is enabled in Blender and the path to `server.py` is correct.

---

## Step 5: The “Look, Then Touch” Rule

Never tell the AI to "Fix the scene." It can't see what "broken" means to you.
1. **Look:** `“What objects are in my scene?”`
2. **Select:** `“Select the Cube.”`
3. **Action:** `“Move the selected Cube up by 5 units.”`

---

## Step 6: Beginner Commands to Practice

| Command | What it does |
| --- | --- |
| `"What is currently selected?"` | Confirms the AI knows what you are working on. |
| `"Add a sphere at the center."` | Creates a new object safely. |
| `"Rename 'Cube' to 'MyBox'."` | Tests the AI’s ability to modify data. |
| `"Set viewport to Material Preview."` | Changes how things look in Blender. |

---

# 🎓 Part 3: Advanced Tips & Safe Practices

Now that you’re connected, let’s learn how to use the bridge like a pro. This section covers how to stay in control and catch mistakes early.

---

## 1. Using Telemetry as Your "Eyes"

Telemetry is just a fancy word for "data about the scene." When the AI runs `get_scene_telemetry`, it gets a list of every object, its position, and its rotation.
* **Pro Tip:** If the AI seems confused, ask it: `“Run a full telemetry scan and tell me exactly where 'Cube.001' is located.”` This forces it to update its internal map.

---

## 2. The Power of "Undo"

Because BlenderVibeBridge uses **Transactions**, every AI command is a single "Step" in Blender’s undo history.
* If the AI makes a mistake, just go to Blender and press **Ctrl + Z**.
* You can also tell the AI: `“That was wrong. Please undo the last operation.”` (It will use the `rollback_transaction` tool).

---

## 3. Adversarial Prompting (The "Cynical Auditor")

As mentioned in Part 1, the AI can be overconfident. Before you let it do something big (like rigging a character or setting up complex lighting), use this prompt:
> *“I want you to act as a cynical Technical Director. Before we run this lighting setup, tell me 3 ways this could fail or make the scene look worse.”*

This forces the AI to stop being a "Yes-Man" and actually think about technical risks.

---

## 4. Checking for "Orphan Data"

In Blender, when you delete an object, its "Mesh" or "Material" sometimes stays hidden in the file. This is called **Orphan Data**.
* **Safe Practice:** Every once in a while, ask the AI: `“Are there any orphan data-blocks? If so, please purge them.”` This keeps your file size small and clean.

---

## 5. Summary: The Golden Workflow

To stay safe and fast, follow this 3-step loop:
1. **Inspect**: Ask the AI what it sees.
2. **Propose**: Ask the AI what it *plans* to do.
3. **Execute & Verify**: Let it run the command, then ask it to confirm the result.

---
**Copyright (C) 2026 B-A-M-N**
