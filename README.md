## Version Tagging System

This project uses a version tagging system for Docker images based on Git tags and commits. Here's how it works:

1. **Versioning**: 
   - Uses `git describe --tags --always --dirty` to generate version tags.
   - Format: `<latest-tag>-<commits-since-tag>-g<commit-hash>` (e.g., v1.0.0-3-g1234567).
   - Falls back to "dev" if no Git information is available.

2. **Building**:
   - `make build`: Creates a Docker image tagged with the current version and "latest".

3. **Running**:
   - `make run` or `make run-env`: Runs the container with the current version tag.
   - Specify version: `make run VERSION=v1.0.0`

4. **Utility Commands**:
   - `make version`: Displays the current version.
   - `make clean`: Removes both versioned and "latest" tagged images.

5. **Workflow**:
   - Tag releases: `git tag v1.0.0`
   - Build: `make build`
   - Run: `make run`

This system ensures each Docker image is tied to a specific code version, facilitating easier debugging and version management.