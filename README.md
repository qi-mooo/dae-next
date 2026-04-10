# dae-next-ui

This repository keeps the custom WebUI and workspace-level integration around `dae`.

## Layout

- `dae/`: git submodule with the minimal backend/controller changes needed by the WebUI
- `ui/`: standalone WebUI assets and usage notes

## Development

The `dae` submodule is intentionally kept smaller than the original all-in-one branch.
UI assets now live in the superproject instead of the `dae` repository itself.

When running `dae` from this workspace, the controller can discover `../ui` through the existing WebUI search logic.
