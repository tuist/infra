#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

docker_image="${DOCKER_SMOKE_IMAGE:-alpine:3.20}"
go_image="${GO_SMOKE_IMAGE:-golang:1.24}"
go_bin="/usr/local/go/bin/go"

echo "==> Checking Docker connectivity"
docker info >/dev/null

echo "==> Verifying the Docker runtime"
docker run --rm "${docker_image}" echo docker-ok

echo "==> Running controller smoke tests in ${go_image}"
docker run --rm \
  -v "${repo_root}:/workspace" \
  -w /workspace \
  "${go_image}" \
  sh -lc "${go_bin} test ./internal/controller -run 'Test(ReconciliationSmokeFlow|ProviderMachineReconcileUpdatesObservedHost)$' -v"

echo "==> Building the repository in ${go_image}"
docker run --rm \
  -v "${repo_root}:/workspace" \
  -w /workspace \
  "${go_image}" \
  sh -lc "${go_bin} build ./..."
