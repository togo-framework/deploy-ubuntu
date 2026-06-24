<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/deploy-ubuntu</h1>
  <p><a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" /></a> <a href="https://pkg.go.dev/github.com/togo-framework/deploy-ubuntu"><img src="https://pkg.go.dev/badge/github.com/togo-framework/deploy-ubuntu.svg" /></a> <img src="https://img.shields.io/badge/license-MIT-blue" /></p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/deploy        # the base
togo install togo-framework/deploy-ubuntu
```
<!-- /togo-header -->

Ubuntu VPS deploy driver — builds a linux/amd64 binary, rsyncs it over SSH, writes a systemd unit, and (re)starts the service.

Configure in `togo.yaml`:

```yaml
deploy:
  provider: ubuntu
  host: 1.2.3.4
  user: root
  domain: app.example.com
```

<!-- togo-sponsors -->
---
<div align="center"><h3>Premium sponsors</h3><p><a href="https://id8media.com"><strong>ID8 Media</strong></a> · <a href="https://one-studio.co"><strong>One Studio</strong></a></p><p><sub><a href="https://github.com/sponsors/fadymondy">Become a sponsor</a>.</sub></p></div>
<!-- /togo-sponsors -->
