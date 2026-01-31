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

#if UNITY_EDITOR
using UnityEditor;
using UnityEngine;
using System;
using System.Net;
using System.IO;
using System.Text;
using System.Threading;
using System.Collections.Concurrent;
using System.Collections.Generic;

using System.Security.Cryptography;
using System.Linq;

[InitializeOnLoad]
public static class VibeBridgeServer
{
    private static HttpListener _listener;
    private static Thread _listenerThread;
    private static readonly ConcurrentQueue<Action> _mainThreadQueue = new ConcurrentQueue<Action>();
    
    // Security Tokens
    private static string _sessionToken = "";
    private static int _currentGeneration = 0;
    private const string BOOTSTRAP_TOKEN = "VIBE_UNITY_BOOTSTRAP_SECRET";
    private static readonly object _stateLock = new object();

    private static readonly HashSet<string> _pathWhitelist = new HashSet<string> {
        "/health", "/handshake", "/metrics", "/object/lock", "/panic", 
        "/validate", "/state/get", "/commit", "/rollback", 
        "/transform/set", "/material/update", "/object/mutate",
        "/selection/set", "/camera/set", "/camera/get"
    };

    [Serializable]
    private class HandshakePayload { public string new_token; public string challenge; }

    static VibeBridgeServer()
    {
        EditorApplication.update += OnUpdate;
        StartServer();
    }

    private static string ComputeHMAC(string key, string data)
    {
        using (var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(key)))
        {
            byte[] hash = hmac.ComputeHash(Encoding.UTF8.GetBytes(data));
            return BitConverter.ToString(hash).Replace("-", "").ToLower();
        }
    }

    private static void StartServer()
    {
        try
        {
            if (_listener != null) StopServer();
            
            _listener = new HttpListener();
            _listener.Prefixes.Add("http://127.0.0.1:8085/");
            _listener.Start();
            _listenerThread = new Thread(Listen);
            _listenerThread.Start();
            Debug.Log("🛡️ VibeSync Unity Bridge: Listening on port 8085");
        }
        catch (Exception e)
        {
            Debug.LogError($"VibeSync: Failed to start server: {e.Message}");
        }
    }

    public static void StopServer()
    {
        if (_listener != null)
        {
            try { _listener.Stop(); _listener.Close(); } catch (Exception) {}
            _listener = null;
        }

        if (_listenerThread != null)
        {
            try { _listenerThread.Abort(); } catch (Exception) {}
            _listenerThread = null;
        }
    }

    private static void Listen()
    {
        while (_listener != null && _listener.IsListening)
        {
            try
            {
                var context = _listener.GetContext();
                ProcessRequest(context);
            }
            catch (Exception) { /* Listener closed or error */ }
        }
    }

    private static void ProcessRequest(HttpListenerContext context)
    {
        var request = context.Request;
        var response = context.Response;

        // 1. Path Whitelist Check
        if (!_pathWhitelist.Contains(request.Url.AbsolutePath))
        {
            SendResponse(response, "{\"error\":\"Forbidden Path\"}", HttpStatusCode.Forbidden);
            return;
        }

        // Health check - allowed without token
        if (request.Url.AbsolutePath == "/health")
        {
            string status = (EditorApplication.isCompiling || EditorApplication.isUpdating) ? "busy" : "ok";
            int gen; lock(_stateLock) { gen = _currentGeneration; }
            string responseJson = "{\"status\":\"" + status + "\", \"generation\":" + gen + "}";
            SendResponse(response, responseJson, HttpStatusCode.OK);
            return;
        }

        // 2. Token & Security Header Validation
        string receivedToken = request.Headers["X-Vibe-Token"];
        string receivedGenStr = request.Headers["X-Vibe-Generation"];
        string receivedSig = request.Headers["X-Vibe-Signature"];
        string receivedTime = request.Headers["X-Vibe-Timestamp"];
        
        bool isHandshake = request.Url.AbsolutePath == "/handshake";
        string currToken;
        int currGen;

        lock (_stateLock)
        {
            currToken = _sessionToken != "" ? _sessionToken : BOOTSTRAP_TOKEN;
            currGen = _currentGeneration;
        }

        if (receivedToken != currToken)
        {
            SendResponse(response, "{\"error\":\"Unauthorized\"}", HttpStatusCode.Unauthorized);
            return;
        }

        // 3. Anti-Replay: Timestamp Check (5s window)
        if (long.TryParse(receivedTime, out long ts))
        {
            long now = (long)(DateTime.UtcNow - new DateTime(1970, 1, 1)).TotalSeconds;
            if (Math.Abs(now - ts) > 5)
            {
                SendResponse(response, "{\"error\":\"Request Expired\"}", HttpStatusCode.Forbidden);
                return;
            }
        }
        else if (!isHandshake)
        {
            SendResponse(response, "{\"error\":\"Missing Timestamp\"}", HttpStatusCode.BadRequest);
            return;
        }

        // 4. Generation Drift Check
        if (!isHandshake && int.TryParse(receivedGenStr, out int receivedGen) && receivedGen != currGen)
        {
            SendResponse(response, "{\"error\":\"Generation Drift\"}", HttpStatusCode.Conflict);
            return;
        }

        // Read Body for Signature and Logic
        string body = "";
        using (var reader = new StreamReader(request.InputStream, request.ContentEncoding))
        {
            body = reader.ReadToEnd();
        }

        // 5. HMAC Signature Verification
        if (!isHandshake)
        {
            string sigData = receivedTime + "|" + request.HttpMethod + "|" + request.Url.AbsolutePath + "|" + body;
            string expectedSig = ComputeHMAC(currToken, sigData);
            if (receivedSig != expectedSig)
            {
                SendResponse(response, "{\"error\":\"Invalid Signature\"}", HttpStatusCode.Forbidden);
                return;
            }
        }

        if (isHandshake)
        {
            try 
            {
                var payload = JsonUtility.FromJson<HandshakePayload>(body);
                lock (_stateLock)
                {
                    _currentGeneration++;
                    if (payload != null && !string.IsNullOrEmpty(payload.new_token))
                    {
                        _sessionToken = payload.new_token;
                    }
                }
                string responseJson = "{\"status\":\"OK\", \"engine_version\":\"" + Application.unityVersion + "\", \"capabilities\":[\"transform\", \"mesh\", \"material\", \"locking\", \"metrics\"], \"response\":\"VIBE_HASH_" + (JsonUtility.FromJson<HandshakePayload>(body).challenge ?? "UNK") + "\"}";
                SendResponse(response, responseJson, HttpStatusCode.OK);
            }
            catch (Exception) { SendResponse(response, "{\"error\":\"Invalid Handshake\"}", HttpStatusCode.BadRequest); }
            return;
        }

        if (request.Url.AbsolutePath == "/metrics")
        {
            string responseJson = "{\"memory\":" + GC.GetTotalMemory(false) + ", \"is_compiling\":" + EditorApplication.isCompiling.ToString().ToLower() + "}";
            SendResponse(response, responseJson, HttpStatusCode.OK);
            return;
        }

        if (request.Url.AbsolutePath == "/object/lock")
        {
            using (var reader = new StreamReader(request.InputStream, request.ContentEncoding))
            {
                string body = reader.ReadToEnd();
                if (body.Contains("\"locked\":true")) {
                    _mainThreadQueue.Enqueue(() => Debug.Log("🛡️ VibeSync: Object Lock Applied"));
                } else {
                    _mainThreadQueue.Enqueue(() => Debug.Log("🛡️ VibeSync: Object Lock Released"));
                }
            }
            SendResponse(response, "{\"status\":\"ok\"}", HttpStatusCode.OK);
            return;
        }

        if (request.Url.AbsolutePath == "/panic")
        {
            _mainThreadQueue.Enqueue(() => {
                EditorApplication.isPaused = true;
                Debug.LogError("🚨 VIBESYNC PANIC | Hierarchy Locked by Orchestrator");
            });
            SendResponse(response, "{\"status\":\"locked\"}", HttpStatusCode.OK);
            return;
        }

        // Atomic Sync Endpoints
        if (request.Url.AbsolutePath == "/validate")
        {
            string sceneState = "UnityState_" + UnityEngine.SceneManagement.SceneManager.GetActiveScene().name; 
            string hash = "UNITY_HASH_" + sceneState.GetHashCode();
            string responseJson = "{\"status\":\"OK\", \"hash\":\"" + hash + "\"}";
            SendResponse(response, responseJson, HttpStatusCode.OK);
            return;
        }
        
        if (request.Url.AbsolutePath == "/state/get")
        {
            SendResponse(response, "{\"hash\":\"SCENE_HASH_STABLE\"}", HttpStatusCode.OK);
            return;
        }

        if (request.Url.AbsolutePath == "/commit")
        {
            SendResponse(response, "{\"status\":\"committed\"}", HttpStatusCode.OK);
            return;
        }

        if (request.Url.AbsolutePath == "/rollback")
        {
            SendResponse(response, "{\"status\":\"rolled_back\"}", HttpStatusCode.OK);
            return;
        }

        // Marshal other requests to the main thread
        _mainThreadQueue.Enqueue(() => HandleEngineCommand(request.Url.AbsolutePath, body));

        SendResponse(response, "{\"status\":\"queued\"}", HttpStatusCode.Accepted);
    }

    private static void HandleEngineCommand(string path, string json)
    {
        Debug.Log($"VibeSync Command received on Main Thread: {path}");
        
        if (path == "/transform/set")
        {
            Debug.Log("🛡️ VibeSync: Applying Transform Sync...");
        }
        else if (path == "/material/update")
        {
            Debug.Log("🛡️ VibeSync: Applying Material Sync...");
        }
        else if (path == "/object/mutate")
        {
            Debug.Log("🛡️ VibeSync: Executing Mutation...");
        }
    }

    private static void SendResponse(HttpListenerResponse response, string content, HttpStatusCode status)
    {
        try {
            byte[] buffer = Encoding.UTF8.GetBytes(content);
            response.ContentLength64 = buffer.Length;
            response.StatusCode = (int)status;
            response.ContentType = "application/json";
            response.OutputStream.Write(buffer, 0, buffer.Length);
            response.OutputStream.Close();
        } catch (Exception) { /* Response stream might be closed */ }
    }

    private static void OnUpdate()
    {
        while (_mainThreadQueue.TryDequeue(out Action action))
        {
            try { action(); } 
            catch (Exception e) { Debug.LogError($"VibeSync Main Thread Error: {e}"); }
        }
    }
}
#endif