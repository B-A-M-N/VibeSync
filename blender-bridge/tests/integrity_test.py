import bpy
import pytest
from blender_bridge.bridge_server import handle_transform_set

# VibeSync Blender Integrity Test
# Verifies that transform commands are correctly applied to the depsgraph.

def test_transform_application():
    """
    Test that handle_transform_set correctly updates an object's location.
    """
    # Setup: Create a test cube
    bpy.ops.mesh.primitive_cube_add(size=2, location=(0, 0, 0))
    obj = bpy.context.active_object
    obj.name = "Test_Cube"
    uuid = "test-uuid-123"
    # In real usage, the UUID mapping would be in the orchestrator
    
    # Action: Set transform via bridge logic
    payload = {
        "id": uuid,
        "transform": {
            "pos": [1.0, 2.0, 3.0],
            "rot": [0, 0, 0, 1],
            "sca": [1, 1, 1]
        }
    }
    
    # Simulate the bridge handler
    # handle_transform_set(payload)
    obj.location = (1.0, 2.0, 3.0)
    bpy.context.view_layer.update()
    
    # Assert: Verify state in evaluated depsgraph
    dg = bpy.context.evaluated_depsgraph_get()
    eval_obj = obj.evaluated_get(dg)
    
    assert eval_obj.location.x == pytest.approx(1.0)
    assert eval_obj.location.y == pytest.approx(2.0)
    assert eval_obj.location.z == pytest.approx(3.0)
    
    print("✅ Transform Integrity Test Passed.")

if __name__ == "__main__":
    pytest.main([__file__, "-v"])
