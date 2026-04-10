# dae-next

This repository packages `dae` and the custom next-ui WebUI into a wrapper binary named `dae-next`.

完整构建、部署和排障说明见 [GUIDE.md](GUIDE.md)。

## Layout

- `dae/`: git submodule with the minimal backend/controller changes needed by the WebUI
- `ui/`: standalone WebUI assets and usage notes
- `.github/workflows/build.yml`: CI build for embedded binaries

## Development

The `dae` submodule is intentionally kept smaller than the original all-in-one branch.
The `dae-next` binary embeds the WebUI and falls back to it automatically when no local `ui/` directory is found.
If a local `ui/` directory or `DAE_WEBUI_DIR` is present, those files take priority for development and overrides.

The controller shell now lives in this repository. `dae` only keeps the runtime/provider side needed by `dae-next`.

## Quick Build

```bash
GOWORK=off make -C dae ebpf CLANG=clang-15 STRIP=llvm-strip-15
make dae-next
```

This produces `build/dae-next`.
