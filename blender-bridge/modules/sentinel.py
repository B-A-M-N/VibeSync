import bpy
import time

# Blender Depsgraph Sentinel
# Prevents mutations while the depsgraph is recalculating.

def wait_for_stable_depsgraph(timeout=5.0):
    """
    Waits until the Blender depsgraph is stable.
    """
    start_time = time.time()
    
    # Check if we are in a state where depsgraph can be evaluated
    if not hasattr(bpy.context, "evaluated_depsgraph_get"):
        return True
        
    while time.time() - start_time < timeout:
        # Trigger an update to see if things change
        bpy.context.view_layer.update()
        
        # If no updates are pending, we assume stability
        if not bpy.context.view_layer.depsgraph.is_updated:
            return True
            
        time.sleep(0.1)
        
    print("⚠️ VibeSync: Depsgraph stability timeout reached.")
    return False

def is_engine_busy():
    """
    Returns True if Blender is in a state that should block mutations.
    """
    if bpy.context.screen.is_animation_playing:
        return True
    
    # Check if any operator is running (difficult in Python, but we can check mode)
    if bpy.context.mode != 'OBJECT' and bpy.context.mode != 'EDIT_MESH':
        # Block if in modes like sculpt/paint which are heavy
        return True
        
    return False
