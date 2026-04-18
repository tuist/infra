# Infra

Infra is an open-source Kubernetes control plane for provisioning ephemeral full-VM sandbox environments on bare-metal machines.

It is designed for products that need fast, isolated environments on:

- Linux bare metal with `cloud-hypervisor`
- macOS bare metal with `tart`
- provider-backed capacity such as AWS and Scaleway

## ✨ What Infra Aims To Provide

- ⚡ Warm pools of ready bare-metal hosts
- 🧹 Ephemeral full-VM sandboxes with TTL-based cleanup
- 🔒 Dedicated and shared tenancy policies
- 🖥️ SSH and GUI access for supported images
- ☸️ A Kubernetes-native control plane that can reconcile host capacity up and down

## 🚧 Status

Infra is in early bootstrap. The repository currently contains the initial project docs and design direction; the control plane and host services have not been implemented yet.

The first implementation will target:

- Linux and macOS backends from day one
- full VMs only
- provider-backed host acquisition
- single-active-VM-per-host policies where required, especially for macOS

## 🧭 How It Works

Infra uses Kubernetes as the control plane, not as the sandbox runtime.

At a high level:

1. operators install Infra into a management Kubernetes cluster
2. Infra reconciles host pools against capacity providers such as AWS, Scaleway, or self-managed bare metal
3. host services running on the provisioned machines register those hosts with the control plane
4. products request sandboxes, and Infra places them onto eligible hosts

## 👥 Who It Is For

Infra is meant to model Tuist's infrastructure and provide a base any product can build on when it needs reproducible sandbox environments without building its own host orchestration layer.

## 📚 Documentation

- [Architecture](./docs/architecture.md)
- [License](./LICENSE)

## 📄 License

Infra is available under the MIT license. See [LICENSE](./LICENSE).
