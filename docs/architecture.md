# Architecture

This document captures the current technical design for Infra.

## Goals

- Provision ephemeral full-VM sandboxes on bare-metal machines
- Support Linux and macOS under one control plane
- Use warm pools to reduce startup latency
- Support dedicated and shared placement where the VM runtime allows it
- Reconcile host acquisition and release through Kubernetes controllers
- Keep product integration simple through a public API rather than direct CRD usage

## Non-Goals for v1

- Container-native sandboxes as a first-class primitive
- A single machine reproducing every real VM runtime
- Hiding provider constraints behind a fake generic model

## High-Level Architecture

```text
Product / CI / CLI
        |
        v
+------------------------+
| Infra API              |
| - authn/authz          |
| - sandbox lifecycle    |
| - access information   |
+------------------------+
        |
        v
+-----------------------------------------------+
| Kubernetes Control Plane                      |
| - sandbox controller                          |
| - placement controller                        |
| - warm-pool controller                        |
| - provider controllers                        |
| - garbage collection                          |
+-----------------------------------------------+
        |
        v
+-----------------------------------------------+
| Host Agents                                   |
| - linux -> cloud-hypervisor                   |
| - macOS -> tart                               |
| - local -> simulated runtime                  |
+-----------------------------------------------+
        |
        v
+-----------------------------------------------+
| Bare-Metal Hosts                              |
| - self-managed                                |
| - AWS-backed                                  |
| - Scaleway-backed                             |
+-----------------------------------------------+
```

## Core Resources

### `HostPool`

Defines desired capacity and policy for a pool of interchangeable hosts.

Typical policy includes:

- VM runtime and architecture
- provider type and region or zone
- allowed tenancy modes
- warm-pool targets
- desired warm images
- scale-down behavior

`HostPool` is the main scaling input.

### `ProviderMachine`

Represents one external request for a bare-metal machine from a provider such as AWS, Scaleway, or Metal3.

Suggested lifecycle:

`Pending -> Provisioning -> Bootstrapping -> Registered -> Releasing -> Released -> Failed`

### `Host`

Represents one registered machine that is healthy and schedulable.

A `Host` reports:

- VM runtime support
- OS and architecture
- allocatable capacity
- image-cache summary counts
- agent health

Only healthy `Host` resources count as usable capacity.

### `HostImage`

Represents one desired image cache entry for a specific host and VM runtime.

Suggested lifecycle:

`Pending -> Pulling -> Ready -> Failed`

`HostImage` is the first-class state that lets Infra reconcile image warming
asynchronously instead of assuming images are already present on a host.

### `HostLease`

Represents an exclusive claim on a host.

This is required for dedicated placement and for policies such as `maxActiveVMsPerHost = 1`.

### `SandboxClass`

Defines a reusable sandbox shape:

- VM runtime
- CPU, memory, and disk defaults
- tenancy policy
- access policy
- TTL

### `Sandbox`

Represents one requested ephemeral full VM.

Suggested lifecycle:

`Pending -> Scheduled -> Creating -> Ready -> Draining -> Deleted`

## Reconciliation Model

Host acquisition and release should be handled by Kubernetes reconciliation.

The intended flow is:

1. an operator creates or updates a `HostPool`
2. the `HostPool` controller computes the required warm capacity
3. Infra creates or deletes `ProviderMachine` resources
4. provider controllers reconcile those resources against AWS, Scaleway, or Metal3
5. provisioned machines boot and start the Infra host service
6. the host service registers a `Host`
7. the image cache controller creates `HostImage` resources for the pool's desired warm images
8. the host service pulls those images asynchronously and reports `HostImage` state
9. the placement controller binds `Sandbox` resources only to eligible `Host`s whose required `HostImage` is `Ready`
10. for dedicated placement, Infra creates a `HostLease`
11. on teardown, the lease is released and the host either returns to the pool or is drained for scale-down

## Host Agents

Host services are binaries developed as part of Infra and installed on each bare-metal machine.

They are responsible for:

- registering the host with the control plane
- heartbeating health and capacity
- pulling and caching images
- reporting per-image cache state back through `HostImage` resources
- creating, starting, stopping, and deleting VMs
- publishing access details
- cleaning the host after teardown

Deployment model:

- Linux: service managed by `systemd`
- macOS: service managed by `launchd`
- local development: same agent in simulated mode

## VM Runtimes

### Linux

Linux hosts use `cloud-hypervisor`.

Shared placement is expected to be supported first on Linux.

### macOS

macOS hosts use `tart`.

macOS should be modeled as dedicated single-VM host capacity. In practice, the design assumes one active VM per host.

Note: Infra can remain MIT-licensed while the macOS runtime depends on Tart, which has its own separate licensing terms. Operators need to evaluate that dependency independently.

## Capacity Providers

The current design assumes provider-backed capacity for both Linux and macOS.

Planned provider integrations:

- `metal3`
- `aws-ec2-bare-metal`
- `aws-ec2-mac`
- `scaleway-elastic-metal`
- `scaleway-apple-silicon`

Provider-specific constraints should be handled inside reconciliation logic rather than ignored.

## Access Model

Access is part of the control plane contract.

Supported access modes:

- SSH
- GUI

For macOS, both SSH and GUI should be supported.

For Linux, SSH is mandatory and GUI is optional per sandbox class or image.

## Local Development

The first implementation should be reproducible locally at the control-plane level.

That means:

- a local Kubernetes cluster such as `kind`
- the real Infra controllers
- the real API server
- a `local` host runtime for simulated machines
- a reconciled image-cache loop that can move `HostImage` resources through `Pending`, `Pulling`, and `Ready`

VM-runtime-specific smoke tests can then be run where the host hardware allows it:

- Linux machine for real `cloud-hypervisor`
- Apple Silicon macOS machine for real `tart`

## Recommended Implementation Order

### Phase 0

- define API types and CRDs
- scaffold the controller manager
- scaffold the API server
- build the `local` runtime
- implement the basic sandbox lifecycle

### Phase 1

- implement the Linux host service
- implement the macOS host service
- add placement and lease handling
- add SSH and GUI access publication

### Phase 2

- integrate provider controllers
- add warm-pool autoscaling
- add host draining and replacement

## Readiness

The project is ready to start implementation.

The design is specific enough to begin Phase 0:

- the control-plane boundary is clear
- the core resources are identified
- the host role is defined
- the provider-backed scaling model is defined

The remaining open questions are implementation details, not blockers for starting the skeleton.
