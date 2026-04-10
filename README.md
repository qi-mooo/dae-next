# dae-next

This repository packages `dae` and the custom next-ui WebUI into a wrapper binary named `dae-next`.

## Layout

- `dae/`: git submodule with the minimal backend/controller changes needed by the WebUI
- `ui/`: standalone WebUI assets and usage notes

## Development

The `dae` submodule is intentionally kept smaller than the original all-in-one branch.
The `dae-next` binary embeds the WebUI and falls back to it automatically when no local `ui/` directory is found.
If a local `ui/` directory or `DAE_WEBUI_DIR` is present, those files take priority for development and overrides.

When running `dae` from this workspace, the controller can discover `../ui` through the existing WebUI search logic.

## Build

Build the wrapped binary from the superproject root:

```bash
make dae-next
```

This produces `build/dae-next`.

To build a slimmer binary without embedded UI assets, add the `noembedui` tag and provide `ui/` from disk at runtime.
