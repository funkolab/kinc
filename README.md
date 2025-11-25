# kinc (Kubernetes in Containers)

A tool for creating and managing local Kubernetes clusters using Apple container 'nodes'.

## Prerequisites

- macOS with Apple Silicon (M1/M2/M3 or later)
- Apple containers version 0.6 or later
- Go 1.25 or later (for building from source)


## Installation

```bash
go install github.com/funkolab/kinc@latest
```

## Usage

```bash
# Create a cluster
kinc create cluster

# Get clusters
kinc get clusters

# Delete a cluster
kinc delete cluster
```

## License

Apache License 2.0
