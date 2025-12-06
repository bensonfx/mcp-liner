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
        ldflags = f"-s -w -X main.appVersion={version}"
        build_cmd = [
            "go",
            "build",
            "-trimpath",
            "-ldflags",
            ldflags,
            "-buildmode=c-shared",
            "-o",
            output_path,
            cmd_dir,
        ]

        self.app.display_info(f"Running: {' '.join(build_cmd)}")

        try:
            subprocess.check_call(build_cmd)
        except subprocess.CalledProcessError as e:
            self.app.display_error(f"Go build failed: {e}")
            sys.exit(1)

        self.app.display_info("Go build completed successfully.")

        # self.app.display_info("DEBUG: Entering initialize hook")

        # Explicitly include the generated file in the wheel
        # This is necessary because the file might be ignored by default or not tracked
        if "force_include" not in build_data:
            build_data["force_include"] = {}

        # Map local path to path in wheel
        # Since it's inside mcp_liner/, we want it in mcp_liner dir in the wheel
        build_data["force_include"][output_path] = os.path.join("mcp_liner", lib_name)

        # Mark as not pure python and set tag explicitly to py3-none-<platform>
        # This ensures the wheel is platform-specific but python-version agnostic
        # (since we use ctypes to load the Go shared library)
        try:
            from packaging.tags import sys_tags

            # Get the most specific platform tag (first one in the list usually)
            platform_tag = next(tag.platform for tag in sys_tags())

            # Construct the final tag
            final_tag = f"py3-none-{platform_tag}"

            build_data["pure_python"] = False
            build_data["tag"] = final_tag

        except ImportError:
            self.app.display_warning(
                "Could not import packaging.tags, falling back to pure_python=False only"
            )
            build_data["pure_python"] = False
            build_data["infer_tag"] = True

        # Clean up header file
        header_path = os.path.splitext(output_path)[0] + ".h"
        if os.path.exists(header_path):
            os.remove(header_path)
