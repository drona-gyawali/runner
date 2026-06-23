# Runner Engine

A minimal, high-performance continuous integration (CI) pipeline execution engine written in Go.

The platform securely evaluates multi-stage workflow pipelines parsed from graph architectures concurrently using an isolated containerized infrastructure sandbox.


## Core Architecture & Design (v1)

The platform separates the orchestration layer from the execution environment to ensure secure, parallel, and predictable multi-stage workflow execution

> <img width="526" height="553" alt="image" src="https://github.com/user-attachments/assets/961e0090-eb9f-4b7e-bd00-4494a3d69267" />


## Why gVisor + Docker?

Standard container architectures share the host system's Linux kernel directly. If a developer runs a malicious pipeline script or hits an unknown dependency flaw, they can escape the container space and corrupt the core server.

* **gVisor Integration:** This platform wraps container workloads inside a custom `runsc` application kernel virtualization layer. It intercepts unprivileged host system calls, completely shielding the host kernel from potential multi-tenant execution panics or security compromises.
* **Docker Context:** Docker is used to manage dynamic image layer replication (`python`, `alpine`, etc.) and provide an ephemeral volume lifecycle, mounting active workspaces directly onto isolated paths.



## Getting Started

# Prerequisites

* **Go:** `1.22+`
* **Docker Engine** (with `gVisor` runtime configured)

## Local Configuration Setup

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/drona-gyawali/runner.git
   cd runner

2. **Compile the Core Binary:**
    ```bash
    go build -o bin/runner cmd/runner/main.go

3. **Execute a Workflow Pipeline:**

To run a pipeline workflow, pass the target configuration file path using the -config flag:
```bash
    ./bin/runner -config YOUR_FILE_PATH/runner/workflows/YOUR_CI_NAME.toml
```

## Integration Test

The engine utilizes a zero-cache integration test harness that tests the compiled binary against live project repositories under realistic runtime conditions.

The test suite runs an end-to-end project fixture validation matrix. It spins up concurrent jobs and forces a dependent bottleneck sync, and evaluates core unit tests.

**Running the Integration Suite**

Execute the following command in your terminal space:

```bash
go test -v ./tests/...
```

This engine was engineered completely from scratch out of a deep first-principles curiosity to understand the inner workings of distributed infrastructure, **Docker internals**, and **sandbox virtualization layers**. 

If you appreciate the architecture or find the implementation patterns helpful, **drop a star on the repository!** It helps keep the project visible to other systems engineers. ⭐