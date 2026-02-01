using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using UnityEditor;
using UnityEngine;

namespace VibeSync.Editor.Intelligence
{
    /// <summary>
    /// Exports public tool signatures to JSON for AI API Intelligence.
    /// Part of the "Hallucination Prevention" system.
    /// </summary>
    public static class ApiMetadataMapper
    {
        [MenuItem("VibeSync/Intelligence/Export API Map")]
        public static void ExportApiMap()
        {
            var metadata = new List<MethodMetadata>();
            
            // Search for classes containing "VibeTool" or bridge endpoints
            var types = AppDomain.CurrentDomain.GetAssemblies()
                .SelectMany(a => a.GetTypes())
                .Where(t => t.Name.Contains("Vibe") || t.Namespace?.Contains("VibeSync") == true);

            foreach (var type in types)
            {
                var methods = type.GetMethods(BindingFlags.Public | BindingFlags.Static | BindingFlags.Instance)
                    .Where(m => m.DeclaringType == type);

                foreach (var method in methods)
                {
                    metadata.Add(new MethodMetadata
                    {
                        ClassName = type.FullName,
                        MethodName = method.Name,
                        Parameters = method.GetParameters().Select(p => $"{p.ParameterType.Name} {p.Name}").ToList(),
                        ReturnType = method.ReturnType.Name
                    });
                }
            }

            string json = JsonUtility.ToJson(new MetadataWrapper { tools = metadata }, true);
            string path = Path.Combine(Application.dataPath, "../metadata/unity_api_map.json");
            File.WriteAllText(path, json);
            Debug.Log($"✅ VibeSync: Exported API Metadata Map to {path}");
        }

        [Serializable]
        public class MethodMetadata
        {
            public string ClassName;
            public string MethodName;
            public List<string> Parameters;
            public string ReturnType;
        }

        [Serializable]
        public class MetadataWrapper
        {
            public List<MethodMetadata> tools;
        }
    }
}
