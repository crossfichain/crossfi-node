# CrossFi Node

## Deterministic Builds

CrossFi Node supports deterministic builds to ensure that anyone can verify that the published binaries match the source code. This is critical for blockchain infrastructure where consensus-critical binaries must be reproducible and verifiable.

### What are Deterministic Builds?

Deterministic builds ensure that compiling the same source code with the same build environment produces identical binaries every time. This allows:
- Node operators to verify binaries match the source code
- Community members to audit releases independently
- Enhanced security through build transparency

### Triggering a Deterministic Build

The deterministic build workflow can be triggered manually via GitHub Actions:

1. Navigate to the [Actions tab](../../actions/workflows/deterministic-build.yml)
2. Click "Run workflow"
3. Optionally specify a version tag (e.g., `v1.0.0`), or leave empty to use the current commit hash
4. Click "Run workflow" to start the build

The workflow will build `crossfid` binaries for all supported platforms:
- `darwin-amd64` (macOS Intel)
- `darwin-arm64` (macOS Apple Silicon)
- `linux-amd64` (Linux x86_64)
- `linux-arm64` (Linux ARM64)

### Build Artifacts

After the workflow completes, you can download:
- **Binaries**: `crossfid-binaries-{version}` artifact containing all platform binaries
- **Checksums**: `checksums-{version}` artifact containing SHA256 checksums

### Verifying Checksums

To verify a downloaded binary:

```bash
# Download the binary and checksums.txt
sha256sum crossfid-linux-amd64

# Compare with the value in checksums.txt
cat checksums.txt
```

### Building Locally (Reproducible)

To build deterministic binaries locally, you need Docker installed:

```bash
# Build for Linux AMD64 (example)
docker run --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  ghcr.io/goreleaser/goreleaser-cross:v1.21 \
  bash -c "
    export CGO_ENABLED=1
    export GOOS=linux
    export GOARCH=amd64
    export CC=gcc
    export CXX=g++

    go build \
      -mod=readonly \
      -tags='cgo' \
      -trimpath \
      -buildvcs=false \
      -ldflags='-X github.com/cosmos/cosmos-sdk/version.Name=crossfi \
        -X github.com/cosmos/cosmos-sdk/version.AppName=crossfid \
        -X github.com/cosmos/cosmos-sdk/version.Version=\$(git describe --tags --always) \
        -X github.com/cosmos/cosmos-sdk/version.Commit=\$(git rev-parse --short HEAD) \
        -X github.com/cosmos/cosmos-sdk/version.BuildTags=cgo \
        -w -s' \
      -o crossfid-linux-amd64 \
      ./cmd/crossfid
  "

# Generate checksum
sha256sum crossfid-linux-amd64
```

**Note**: The checksum should match the one from the GitHub Actions workflow, proving the build is reproducible.

### Deterministic Build Flags

The following flags ensure build reproducibility:

- `-mod=readonly`: Prevents dependency modifications during build
- `-tags='cgo'`: Enables CGO for RocksDB support
- `-trimpath`: Removes absolute file paths from the binary
- `-buildvcs=false`: Excludes version control information
- `-w -s`: Strips debugging symbols (reduces binary size)

### Technical Details

- **Docker Image**: `ghcr.io/goreleaser/goreleaser-cross:v1.21`
- **Go Version**: 1.21.0
- **CGO**: Enabled (required for RocksDB and other native dependencies)
- **Cross-Compilation Toolchains**:
  - macOS Intel: `o64-clang`
  - macOS ARM: `oa64-clang`
  - Linux AMD64: `gcc`
  - Linux ARM64: `aarch64-linux-gnu-gcc`