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

        # On Linux, repair the shared library with auditwheel for manylinux compatibility
        if sys.platform.startswith("linux"):
            self.app.display_info("Detected Linux platform, running auditwheel repair...")
            self._repair_wheel_linux(output_path, output_dir, lib_name)

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

    def _repair_wheel_linux(self, lib_path, output_dir, lib_name):
        """
        Use auditwheel to repair the shared library for manylinux compatibility.
        This ensures the wheel can be uploaded to PyPI and works across different Linux distributions.
        """
        try:
            # Check if auditwheel is available
            result = subprocess.run(["auditwheel", "--version"], capture_output=True, text=True, check=False)

            if result.returncode != 0:
                self.app.display_warning("auditwheel not found. Skipping wheel repair. "
                                         "Install auditwheel for manylinux compatibility: pip install auditwheel")
                return

            # Create a temporary directory for repaired libraries
            repair_dir = os.path.join(output_dir, ".auditwheel_temp")
            os.makedirs(repair_dir, exist_ok=True)

            # Run auditwheel repair
            repair_cmd = [
                "auditwheel",
                "repair",
                lib_path,
                "-w",
                repair_dir,
                "--plat",
                "manylinux2014_x86_64"  # Adjust based on architecture
            ]

            self.app.display_info(f"Running: {' '.join(repair_cmd)}")

            result = subprocess.run(repair_cmd, capture_output=True, text=True, check=False)

            if result.returncode == 0:
                self.app.display_info("auditwheel repair completed successfully.")

                # Find the repaired library
                repaired_files = os.listdir(repair_dir)
                if repaired_files:
                    repaired_lib = os.path.join(repair_dir, repaired_files[0])
                    # Replace the original library with the repaired one
                    if os.path.exists(repaired_lib):
                        os.replace(repaired_lib, lib_path)
                        self.app.display_info(f"Replaced {lib_path} with repaired version.")

                # Clean up temporary directory
                import shutil
                shutil.rmtree(repair_dir, ignore_errors=True)
            else:
                self.app.display_warning(f"auditwheel repair failed: {result.stderr}\n"
                                         "Continuing with unrepaired library. This may affect PyPI compatibility.")

        except Exception as e:
            self.app.display_warning(f"Error during auditwheel repair: {e}\n"
                                     "Continuing with unrepaired library.")
