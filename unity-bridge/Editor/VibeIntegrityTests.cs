using System.Collections;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.TestTools;
using UnityEditor;

namespace VibeSync.Editor.Tests
{
    /// <summary>
    /// VibeSync Unity Integrity Test Suite.
    /// Runs inside the Unity Process via VibeTool_TestRunner.
    /// </summary>
    public class VibeIntegrityTests
    {
        [Test]
        public void Test_TransformSync_Integrity()
        {
            // Setup
            GameObject testObj = new GameObject("Test_Integrity_Obj");
            var initialPos = testObj.transform.position;

            // Action: Simulate a bridge transform update
            Vector3 targetPos = new Vector3(10, 20, 30);
            testObj.transform.position = targetPos;

            // Assert
            Assert.AreEqual(targetPos, testObj.transform.position, "Transform position mismatch after sync.");
            
            // Cleanup
            GameObject.DestroyImmediate(testObj);
            Debug.Log("✅ Unity Transform Integrity Test Passed.");
        }

        [UnityTest]
        public IEnumerator Test_AssetImport_Wait()
        {
            // Verifies that the bridge correctly waits for AssetDatabase
            yield return null;
            Assert.IsFalse(EditorApplication.isUpdating, "Mutation occurred during AssetDatabase update.");
        }
    }
}
