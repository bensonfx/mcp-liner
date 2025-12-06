"""
Basic tests for mcp-liner package.

Tests the Python package structure and shared library loading.
"""

import sys
from pathlib import Path

import pytest


class TestPackageStructure:
    """Test basic package structure and imports."""

    def test_package_exists(self):
        """Test that mcp_liner package can be imported."""
        # Note: This test may fail if the package is not installed
        # It will work in CI when the package is built and installed
        try:
            import mcp_liner  # noqa: F401
        except ImportError:
            pytest.skip("Package not installed, skipping import test")

    def test_shared_library_exists(self):
        """Test that the shared library file exists in the package directory."""
        # Get the mcp_liner package location
        try:
            import mcp_liner

            package_dir = Path(mcp_liner.__file__).parent
        except (ImportError, AttributeError):
            pytest.skip("Package not installed, skipping library test")

        # Determine expected library name based on platform
        if sys.platform.startswith("win"):
            lib_name = "_mcp_liner.dll"
        else:
            lib_name = "_mcp_liner.so"

        lib_path = package_dir / lib_name

        # Check if the library file exists
        assert lib_path.exists(), (f"Shared library not found at {lib_path}. "
                                   "The package may not be properly built.")

    def test_shared_library_loadable(self):
        """Test that the shared library can be loaded."""
        try:
            import mcp_liner

            package_dir = Path(mcp_liner.__file__).parent
        except (ImportError, AttributeError):
            pytest.skip("Package not installed, skipping library load test")

        # Determine library path
        if sys.platform.startswith("win"):
            lib_name = "_mcp_liner.dll"
        else:
            lib_name = "_mcp_liner.so"

        lib_path = package_dir / lib_name

        if not lib_path.exists():
            pytest.skip("Shared library not found")

        # Try to load the library using ctypes
        import ctypes

        try:
            if sys.platform.startswith("win"):
                lib = ctypes.CDLL(str(lib_path))
            else:
                lib = ctypes.CDLL(str(lib_path), mode=ctypes.RTLD_GLOBAL)

            # If we get here, the library loaded successfully
            assert lib is not None
        except OSError as e:
            pytest.fail(f"Failed to load shared library: {e}")


class TestBuildArtifacts:
    """Test build artifacts and configuration."""

    def test_version_file_or_git(self):
        """Test that version can be determined from git or _version.py."""
        # This test ensures that versioning works
        project_root = Path(__file__).parent.parent

        # Check if we're in a git repo or if version file exists
        git_dir = project_root / ".git"
        version_file = project_root / "mcp_liner" / "_version.py"

        assert (git_dir.exists() or version_file.exists()), "Neither .git directory nor _version.py found"

    def test_pyproject_toml_exists(self):
        """Test that pyproject.toml exists."""
        project_root = Path(__file__).parent.parent
        pyproject_path = project_root / "pyproject.toml"

        assert pyproject_path.exists(), "pyproject.toml not found"

    def test_hatch_build_hook_exists(self):
        """Test that hatch_build.py exists."""
        project_root = Path(__file__).parent.parent
        hatch_build_path = project_root / "hatch_build.py"

        assert hatch_build_path.exists(), "hatch_build.py not found"


class TestGoComponents:
    """Test Go-related components."""

    def test_go_mod_exists(self):
        """Test that go.mod exists."""
        project_root = Path(__file__).parent.parent
        go_mod_path = project_root / "go.mod"

        assert go_mod_path.exists(), "go.mod not found"

    def test_go_source_exists(self):
        """Test that Go source directories exist."""
        project_root = Path(__file__).parent.parent

        # Check for common Go source directories
        cmd_dir = project_root / "cmd"
        internal_dir = project_root / "internal"

        assert cmd_dir.exists() or internal_dir.exists(), ("No Go source directories found")
