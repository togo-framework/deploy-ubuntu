# deploy-ubuntu — docs

**Ubuntu VPS deploy.** SSH/systemd VPS driver — build linux/amd64, rsync over SSH, write a systemd unit, restart.

## Install

```bash
togo install togo-framework/deploy-ubuntu
```

Registers on the [`deploy`](https://github.com/togo-framework/deploy) base; select it with **deploy.provider in togo.yaml (or DEPLOY_PROVIDER)**, then use **`togo deploy`**.

## Interface

`Deployer` — `Provision`/`Deploy`/`Destroy`/`Status` over a `Spec{App,Dir,BuildCmd,Host,User,Image,Region,Domain}` built from your `togo.yaml`.

## Usage & notes

Config comes from the `togo.yaml` `deploy` block: `host` (required), `user` (default `root`). Installs runtime prerequisites with apt, deploys to `/opt/<app>`, manages a systemd service. Needs SSH key access to the host.

## Example

```bash
togo deploy --provider ubuntu --dry-run   # preview the plan
togo deploy --provider ubuntu
```

## Links

- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/deploy-ubuntu)
