# dae-next

This repository packages `dae` and the custom next-ui WebUI into a wrapper binary named `dae-next`.

## Layout

- `dae/`: git submodule with the minimal backend/controller changes needed by the WebUI
- `ui/`: standalone WebUI assets and usage notes
- `.github/workflows/build.yml`: CI build for embedded and noembed binaries

## Development

The `dae` submodule is intentionally kept smaller than the original all-in-one branch.
The `dae-next` binary embeds the WebUI and falls back to it automatically when no local `ui/` directory is found.
If a local `ui/` directory or `DAE_WEBUI_DIR` is present, those files take priority for development and overrides.

The controller shell now lives in this repository. `dae` only keeps the runtime/provider side needed by `dae-next`.

## Build

`dae-next` depends on `dae`'s generated eBPF Go bindings. A clean checkout does not contain those generated files, so local builds must generate them first.

Typical Linux CI sequence:

```bash
GOWORK=off make -C dae ebpf CLANG=clang-15 STRIP=llvm-strip-15
make dae-next
```

This produces `build/dae-next`.

For local macOS development, a Homebrew LLVM clang also works:

```bash
GOWORK=off make -C dae ebpf CLANG=/usr/local/opt/llvm/bin/clang STRIP=llvm-strip TARGET=bpfel
make dae-next
```

To build a slimmer binary without embedded UI assets, add the `noembedui` tag and provide `ui/` from disk at runtime:

```bash
GOWORK=off GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -tags noembedui -trimpath -o build/dae-next-noembed ./cmd/dae-next
```

## GitHub Actions

The build workflow installs LLVM, runs `GOWORK=off make -C dae ebpf ...`, and then publishes four artifacts:

- `dae-next-linux-amd64`
- `dae-next-linux-amd64-noembed`
- `dae-next-linux-arm64`
- `dae-next-linux-arm64-noembed`

The default artifacts embed the WebUI.
The `-noembed` artifacts require `ui/` to be provided from disk at runtime.
