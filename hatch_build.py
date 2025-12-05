import os
import subprocess
import sys
from hatchling.builders.hooks.plugin.interface import BuildHookInterface


class CustomBuildHook(BuildHookInterface):

    def initialize(self, version, build_data):
        self.app.display_info("Starting Go build process...")

        # Define input and output paths
        cmd_dir = "./cmd/mcp-liner"
        output_dir = "mcp_liner"

        # Determine shared library extension based on platform
        if sys.platform.startswith("win"):
            lib_name = "_mcp_liner.dll"
        elif sys.platform.startswith("darwin"):
            lib_name = "_mcp_liner.so"
        else:
            lib_name = "_mcp_liner.so"

        output_path = os.path.join(output_dir, lib_name)

        # Ensure output directory exists
        os.makedirs(output_dir, exist_ok=True)

        # Build command
        build_cmd = ["go", "build", "-buildmode=c-shared", "-o", output_path, cmd_dir]

        self.app.display_info(f"Running: {' '.join(build_cmd)}")

        try:
            subprocess.check_call(build_cmd)
        except subprocess.CalledProcessError as e:
            self.app.display_error(f"Go build failed: {e}")
            sys.exit(1)

        self.app.display_info("Go build completed successfully.")

        # Explicitly include the generated file in the wheel
        # This is necessary because the file might be ignored by default or not tracked
        if 'force_include' not in build_data:
            build_data['force_include'] = {}

        # Map local path to path in wheel
        # Since it's inside mcp_liner/, we want it in mcp_liner dir in the wheel
        build_data['force_include'][output_path] = os.path.join("mcp_liner", lib_name)

        # Clean up header file
        header_path = os.path.splitext(output_path)[0] + ".h"
        if os.path.exists(header_path):
            os.remove(header_path)
