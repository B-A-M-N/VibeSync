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

class VIBE_OT_sample_op(bpy.types.Operator):
    bl_idname = "vibe.sample_op"
    bl_label = "Vibe Sample Op"
    
    def execute(self, context):
        print("Sample Op Executed")
        return {'FINISHED'}

classes = [VIBE_OT_sample_op]

def init():
    print("Module A Initialized")
