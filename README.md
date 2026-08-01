# JustSSH

Open. Select. Connect.

A minimal SSH launcher for the terminal. No dashboard, no borders, no
boxes, just a list of servers and a prompt that gets out of your way the
moment you've connected.

```
JustSSH                                                    5/5

❯ Production                                 root@192.168.1.15
  Database                                   postgres@10.0.0.5
  Hetzner                                         root@hetzner
  Oracle                                         ubuntu@oracle
  Home                                       ihnat@192.168.0.2

↵ Connect    / Search    a Add    e Edit    d Delete    q Quit
```

## Install

```sh
go install github.com/luynrs/justssh/cmd/jssh@latest
```

## Usage

Run `jssh`, pick a server, hit enter and that's all

| Key | Action |
|-----|--------|
| `↵` | Connect |
| `/` | Search |
| `a` | Add a server |
| `e` | Edit the selected server |
| `d` | Delete the selected server |
| `q` / `ctrl+c` | Quit |

JustSSH shells out to the system `ssh` binary instead of speaking the SSH
protocol itself, so your `ssh-agent`, `ProxyJump`, `ControlMaster`,
certificates, and hardware keys all keep working exactly as configured

## Configuration

Servers live in `~/.config/justssh/config.yaml`:

```yaml
servers:
  - name: Production
    host: 192.168.1.15
    user: root
    port: 22
    key: ~/.ssh/id_ed25519
```

On first run, if this file doesn't exist yet, JustSSH seeds it by importing
literal `Host` entries from `~/.ssh/config`. Wildcard and `Match` blocks are
left alone, since the system `ssh` binary already understands them and
JustSSH doesn't need to

You can override the config path with `JUSTSSH_CONFIG`.

Colors come from your terminal's own ANSI palette (cyan/gray/red) instead of
a fixed hex, so changing your terminal theme changes JustSSH too
