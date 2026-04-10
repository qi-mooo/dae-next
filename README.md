# dae-next

This repository packages `dae` and the custom next-ui WebUI into a wrapper binary named `dae-next`.

## Layout

- `dae/`: git submodule with the minimal backend/controller changes needed by the WebUI
- `ui/`: standalone WebUI assets and usage notes

## Development

The `dae` submodule is intentionally kept smaller than the original all-in-one branch.
UI assets now live in the superproject instead of the `dae` repository itself.

When running `dae` from this workspace, the controller can discover `../ui` through the existing WebUI search logic.

## Build

Build the wrapped binary from the superproject root:

```bash
make dae-next
```

This produces `build/dae-next`.
