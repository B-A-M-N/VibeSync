# VibeSync: Zero-Trust Unity ↔ Blender Orchestrator
# Copyright (C) 2026 B-A-M-N
#
# This project is distributed under a DUAL-LICENSING MODEL:
# 1. Open-Source Path: GNU Affero General Public License v3
# 2. Commercial Path: "Work-or-Pay" Model
#
# See the LICENSE file in the project root for the full terms and conditions
# of both licensing paths.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

import bpy
from . import dev_helper
from . import bridge_server

class VIBE_PT_main_panel(bpy.types.Panel):
    bl_label = "VibeSync MCP"
    bl_idname = "VIBE_PT_main_panel"
    bl_space_type = 'VIEW_3D'
    bl_region_type = 'UI'
    bl_category = 'VibeSync'

    def draw(self, context):
        layout = self.layout
        
        layout.operator("vibe.sample_op", text="Run Module Action")
        layout.separator()
        layout.operator("vibe.dev_reload", text="Force Dev Reload", icon='FILE_REFRESH')

class VIBE_OT_dev_reload(bpy.types.Operator):
    bl_idname = "vibe.dev_reload"
    bl_label = "Vibe Dev Reload"
    
    def execute(self, context):
        dev_helper.dev_reload()
        return {'FINISHED'}

def register():
    bpy.utils.register_class(VIBE_PT_main_panel)
    bpy.utils.register_class(VIBE_OT_dev_reload)
    dev_helper.register_handler()
    bridge_server.register()

def unregister():
    bpy.utils.unregister_class(VIBE_PT_main_panel)
    bpy.utils.unregister_class(VIBE_OT_dev_reload)
    bridge_server.unregister()

if __name__ == "__main__":
    register()
