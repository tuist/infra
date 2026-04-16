# Kubebox

Kubebox is an open-source Kubernetes control plane for provisioning ephemeral full-VM sandbox environments on bare-metal machines.

It is designed for products that need fast, isolated environments on:

- Linux bare metal with `cloud-hypervisor`
- macOS bare metal with `tart`
- provider-backed capacity such as AWS and Scaleway

## ✨ What Kubebox Aims To Provide

- ⚡ Warm pools of ready bare-metal hosts
- 🧹 Ephemeral full-VM sandboxes with TTL-based cleanup
- 🔒 Dedicated and shared tenancy policies
- 🖥️ SSH and GUI access for supported images
- ☸️ A Kubernetes-native control plane that can reconcile host capacity up and down

## 🚧 Status

Kubebox is in early bootstrap. The repository currently contains the initial project docs and design direction; the control plane and host agents have not been implemented yet.

The first implementation will target:

- Linux and macOS backends from day one
- full VMs only
- provider-backed host acquisition
- single-active-VM-per-host policies where required, especially for macOS

## 🧭 How It Works

Kubebox uses Kubernetes as the control plane, not as the sandbox runtime.

At a high level:

1. operators install Kubebox into a management Kubernetes cluster
2. Kubebox reconciles host pools against capacity providers such as AWS, Scaleway, or self-managed bare metal
3. host agents running on the provisioned machines register those hosts with the control plane
4. products request sandboxes, and Kubebox places them onto eligible hosts

## 👥 Who It Is For

Kubebox is meant to be infrastructure that any product can build on, including Tuist and other teams that need reproducible sandbox environments without building their own host orchestration layer.

## 📚 Documentation

- [Architecture](./docs/architecture.md)
- [License](./LICENSE)

## 📄 License

Kubebox is available under the MIT license. See [LICENSE](./LICENSE).
