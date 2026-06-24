// Package ubuntu is a VPS deploy driver for togo deploy targeting Ubuntu hosts.
// It builds a linux/amd64 binary locally, rsyncs it to the host over SSH, writes a
// systemd unit, and (re)starts the service. Select with DEPLOY_PROVIDER=ubuntu.
//
// deploy-centos and deploy-debian are sibling drivers with the same flow; they
// differ only in the package manager used to install runtime prerequisites.
package ubuntu

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/togo-framework/deploy"
	"github.com/togo-framework/togo"
)

func init() {
	deploy.RegisterDriver("ubuntu", func(k *togo.Kernel) (deploy.Deployer, error) {
		return New("ubuntu", "apt-get update -y && apt-get install -y ca-certificates"), nil
	})
}

// New returns a generic SSH/systemd VPS driver. distro is the name; pkgInstall is
// the shell snippet that installs runtime prerequisites (apt/yum/dnf). Exported so
// deploy-centos / deploy-debian can reuse it.
func New(distro, pkgInstall string) deploy.Deployer { return &driver{distro: distro, pkgInstall: pkgInstall} }

type driver struct {
	distro     string
	pkgInstall string
}

func run(ctx context.Context, dir, name string, args ...string) (string, error) {
	if _, err := exec.LookPath(name); err != nil {
		return "", fmt.Errorf("deploy-%s: %q not found on PATH", "vps", name)
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	return out.String(), cmd.Run()
}

func target(spec deploy.Spec) string {
	user := spec.User
	if user == "" {
		user = "root"
	}
	return user + "@" + spec.Host
}

func (d *driver) Provision(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	if spec.Host == "" {
		return nil, fmt.Errorf("deploy-%s: host required", d.distro)
	}
	// Ensure remote prerequisites + an app dir.
	if out, err := run(ctx, "", "ssh", target(spec), d.pkgInstall+" && mkdir -p /opt/"+spec.App); err != nil {
		return nil, fmt.Errorf("provision %s: %w\n%s", spec.Host, err, out)
	}
	return &deploy.Result{Message: d.distro + " host prepared: " + spec.Host}, nil
}

func (d *driver) Deploy(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	if spec.Host == "" {
		return nil, fmt.Errorf("deploy-%s: host required", d.distro)
	}
	dir := spec.Dir
	if dir == "" {
		dir = "."
	}
	bin := spec.Binary
	if bin == "" {
		bin = "/tmp/" + spec.App + "-bin"
		env := os.Environ()
		c := exec.CommandContext(ctx, "go", "build", "-o", bin, "./...")
		c.Dir = dir
		c.Env = append(env, "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
		if out, err := c.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("build: %w\n%s", err, out)
		}
	}
	remote := "/opt/" + spec.App + "/" + spec.App
	if out, err := run(ctx, "", "rsync", "-az", bin, target(spec)+":"+remote); err != nil {
		return nil, fmt.Errorf("rsync: %w\n%s", err, out)
	}
	// systemd unit + (re)start.
	var envLines bytes.Buffer
	for k, v := range spec.Env {
		fmt.Fprintf(&envLines, "Environment=%s=%s\n", k, v)
	}
	unit := fmt.Sprintf("[Unit]\nDescription=%s (togo)\nAfter=network.target\n\n[Service]\nExecStart=%s\n%sRestart=always\nUser=%s\n\n[Install]\nWantedBy=multi-user.target\n",
		spec.App, remote, envLines.String(), firstNonEmpty(spec.User, "root"))
	script := fmt.Sprintf("chmod +x %s && cat > /etc/systemd/system/%s.service <<'UNIT'\n%sUNIT\nsystemctl daemon-reload && systemctl enable %s && systemctl restart %s",
		remote, spec.App, unit, spec.App, spec.App)
	if out, err := run(ctx, "", "ssh", target(spec), script); err != nil {
		return nil, fmt.Errorf("systemd: %w\n%s", err, out)
	}
	url := "http://" + spec.Host + ":8080"
	if spec.Domain != "" {
		url = "https://" + spec.Domain
	}
	return &deploy.Result{Message: "deployed " + spec.App + " to " + spec.Host + " (systemd)", URL: url}, nil
}

func (d *driver) Destroy(ctx context.Context, spec deploy.Spec) error {
	_, err := run(ctx, "", "ssh", target(spec),
		fmt.Sprintf("systemctl disable --now %s; rm -f /etc/systemd/system/%s.service; rm -rf /opt/%s; systemctl daemon-reload", spec.App, spec.App, spec.App))
	return err
}

func (d *driver) Status(ctx context.Context, spec deploy.Spec) (*deploy.Status, error) {
	out, err := run(ctx, "", "ssh", target(spec), "systemctl is-active "+spec.App)
	active := strings.TrimSpace(out) == "active"
	return &deploy.Status{Healthy: active && err == nil, Detail: strings.TrimSpace(out)}, nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
