# beam

Type `beam` in an SSH session, and VS Code opens on your Mac, connected to that machine, in that folder.

Basically `code .` for when you're SSH'd into your homelab and VS Code lives on your laptop.

```
~/projects/my-app ❯ beam
sent to beamd
```

…and VS Code pops up on your Mac, already connected over Remote - SSH to `~/projects/my-app`.

## How it works

beam is two small programs:

- **`beam`** runs on the remote machine (your server, homelab, dev box). It grabs your username, the machine's hostname and your current directory, and sends them as JSON to `127.0.0.1:7777`.
- **`beamd`** runs on your Mac. It listens on `127.0.0.1:7777`, builds a `vscode://vscode-remote/ssh-remote+...` link from the request, and hands it to macOS with `open`.

The trick is that `127.0.0.1:7777` on the remote machine isn't really local. An SSH `RemoteForward` carries anything sent there back through your existing SSH connection to `beamd` on your Mac:

```
 remote machine                         your Mac
┌──────────────────────┐   SSH tunnel   ┌──────────────────────┐
│ beam                 │ ─────────────▶ │ beamd                │
│ → 127.0.0.1:7777     │                │ 127.0.0.1:7777       │
└──────────────────────┘                │ → open vscode://...  │
                                        └──────────────────────┘
```

No ports exposed to your network, no extra connections, and it only works while you're actually SSH'd in. Because the tunnel already exists, it's basically instant.

## Requirements

- **Your Mac:** macOS, [VS Code](https://code.visualstudio.com/) with the [Remote - SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh) extension
- **Remote machine:** Linux (or anything Go can build for)
- **To build:** Go 1.24 or newer

## Installing

### 1. Install `beamd` on your Mac

```bash
go install github.com/Alasdair-Roddick/beam/cmd/beamd@latest
```

This puts `beamd` in `~/go/bin` (or `$GOBIN` if you've set it). Make sure that's on your `PATH`, then start it:

```bash
beamd
```

Leave it running in a terminal for now. (Starting it automatically at login is on the roadmap.)

### 2. Install `beam` on the remote machine

```bash
go install github.com/Alasdair-Roddick/beam/cmd/beam@latest
```

Same deal: make sure `~/go/bin` is on your `PATH`.

**Don't want Go on your server?** Build it on any machine and copy the binary over. Go cross-compiles, so from your Mac:

```bash
git clone https://github.com/Alasdair-Roddick/beam.git
cd beam
GOOS=linux GOARCH=amd64 go build -o bin/beam ./cmd/beam
scp bin/beam your-server:~/.local/bin/
```

Use `GOARCH=arm64` for ARM machines like a Raspberry Pi. `uname -m` on the remote tells you which you need (`x86_64` means `amd64`, `aarch64` means `arm64`).

### 3. Set up your SSH config (on your Mac)

Add a `Host` entry to `~/.ssh/config` for the remote machine. Two important bits:

- **The `Host` name must match the remote machine's hostname** (whatever `hostname` prints on it). That's the name `beam` sends, and the name VS Code will use to connect.
- **`RemoteForward`** creates the tunnel `beam` talks through.

```
Host devenv
    HostName 192.168.1.50
    User you
    IdentityFile ~/.ssh/id_ed25519
    RemoteForward 7777 127.0.0.1:7777
```

Replace `devenv` with your machine's hostname and `HostName` with its IP or DNS name.

### 4. Try it

```bash
ssh devenv
cd ~/some/project
beam
```

VS Code should open on your Mac, connected to `devenv`, in `~/some/project`.

## Troubleshooting

**`yelling at an empty room rn. reconnect please :)`**
Nothing is listening on port 7777 on the remote, which means this SSH session has no tunnel. Make sure you connected with `ssh <host>` using the config entry above, so `RemoteForward` applies. Sessions opened directly on the machine (or via a different config) won't have it.

**`The remote is there (can connect to the mac) but nothing is listening`**
The tunnel is up but `beamd` isn't running on your Mac. Start it.

**`Warning: remote port forwarding failed for listen port 7777` when you SSH in**
Something on the remote already owns port 7777. Usually that's another SSH session to the same machine (including VS Code's own connection, which reads the same config). It's harmless: whichever session got the port first still forwards to your Mac. If it's something else, check with `ss -tlnp | grep 7777`.

**VS Code opens but won't connect**
`beamd` did its job; the problem is the SSH connection itself. Check that `ssh <host>` works from a terminal on your Mac, using exactly the name `beam` sent (the remote's hostname).

## Security notes

- `beamd` only listens on `127.0.0.1`, so nothing else on your network can reach it.
- On the remote side, the forwarded port is also bound to `127.0.0.1`. **However, any user logged into that remote machine can connect to it** while your session is open, and ask your Mac to open VS Code. On a personal homelab that's you; on a shared server, think twice until token auth lands (see roadmap).
- `beamd` launches `open` directly, never through a shell, so a weird path can't be turned into a command.

## Roadmap

- [ ] Reject relative paths in `beamd`
- [ ] Shared token so only your `beam` can talk to your `beamd`
- [ ] Configurable port
- [ ] Start `beamd` at login with launchd
- [ ] Optional hostname → SSH alias mapping, for when they don't match
- [ ] Unix socket forwarding, to avoid port clashes between sessions
- [ ] `beamd` support for Linux desktops

## Project layout

```
beam/
├── cmd/
│   ├── beam/        # CLI that runs on the remote machine
│   └── beamd/       # daemon that runs on your Mac
└── internal/
    └── protocol/    # request type shared by both
```

## Contributing

Issues and PRs are welcome. This started as a learning-Go project, so small, readable changes are preferred over clever ones.

1. Fork the repo and create a branch from `main` (e.g. `feat/token-auth` or `fix/path-validation`).
2. Make your change. Keep to the standard library where reasonable; beam has no dependencies and it'd be nice to keep it that way.
3. Before pushing, make sure these are clean:

   ```bash
   gofmt -l .
   go vet ./...
   go build ./...
   ```

4. Test the real flow if your change touches it: `beamd` on a Mac, `beam` on a remote over SSH with the tunnel.
5. Open a PR describing what changed and why. If it changes the request format in `internal/protocol`, mention it, since `beam` and `beamd` need to agree.

Picking something from the roadmap is a great place to start. If it's a bigger change, open an issue first so we can chat about the approach.

## License

TBD.