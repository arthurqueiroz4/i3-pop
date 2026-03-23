# i3-pop

Small daemon for i3 workspace history.

It gives you browser-like navigation:
- `back` (previous workspace)
- `front` (next workspace)

## How it works

- The daemon subscribes to i3 workspace focus events.
- Every time you change workspace, it stores history in two stacks (`back` and `front`).
- It also runs a small TCP server on `127.0.0.1:43223`.
- When you send `back`, it moves to the previous workspace.
- When you send `front`, it moves to the next workspace.
- Your i3 keybindings just send those commands to the daemon with `nc`.

## Quick install

Requirements:
- Linux with i3
- Go installed
- `systemd --user`
- `nc` (netcat)

From the project folder:

```bash
./script/install_daemon.sh
```

This will:
- build `i3-pop` to `~/.local/bin/i3-pop`
- install `systemd/i3-pop.service` to `~/.config/systemd/user/`
- enable and start the daemon

Check if it is running:

```bash
systemctl --user status i3-pop.service
```

## i3 keybindings

Add this to `~/.config/i3/config`:

```i3
bindsym $mod+Tab exec --no-startup-id sh -c "printf back | nc 127.0.0.1 43223"
bindsym $mod+Shift+Tab exec --no-startup-id sh -c "printf front | nc 127.0.0.1 43223"
```

Reload i3:

```bash
i3-msg reload
```

## Logs

```bash
journalctl --user -u i3-pop.service -f
```

## Remove

```bash
./script/uninstall_daemon.sh
```
