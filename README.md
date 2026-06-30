# gunmap

A terminal-native network mapper. No browser, no Electron, just a fast TUI for scanning and exploring networks

Built it because every other network mapping tool either wants a GUI, a web dashboard, or fifteen flags memorized before it does anything useful. gunmap is meant to be the thing you reach for over SSH, on a server with no display, when you just want to know what's on the network right now.

![demo](demo.gif)


## What it does

- Walks you through nmap's flags page by page (host discovery, scan techniques, port selection, timing, firewall evasion, output, and so on) instead of you memorizing or googling them
- Builds the actual nmap command live as you toggle options, so you always see exactly what's about to run
- Lets you edit flag values inline where a flag needs one
- Runs the scan for you once you hit `r` and give it a target
- Fully keyboard driven no mouse, no menus you have to hunt through

## Install

```bash
git clone https://github.com/0xStxrless/gunmap.git
cd gunmap
./install.sh
```

The installer:

1. Builds the binary with `go install`
2. Writes a default config to `~/.config/gunmap.toml` (skipped if one already exists)
3. Symlinks the binary into `/usr/local/bin` so `sudo gunmap` resolves correctly — raw socket scans need root, and `sudo` doesn't see `$HOME/go/bin` by default

If you'd rather do it by hand:

```bash
go install .
sudo ln -s "$(go env GOPATH)/bin/gunmap" /usr/local/bin/gunmap
```

## Usage

```bash
sudo gunmap
```

Move through the sidebar pages (host discovery, scan techniques, port selection, timing, and so on) and pick the nmap flags you want — gunmap builds the actual command live at the bottom of the screen as you toggle options. Once it looks right, hit `r`, type in your target, and it runs the scan for you.

| Key       | Action                  |
|-----------|--------------------------|
| `↑ / ↓`   | move between flags        |
| `← / →`   | switch page                |
| `space`   | toggle flag                |
| `e`       | edit a flag's value    |
| `r`       | run (prompts for target)   |
| `q`       | quit                       |

## Config

UI colors, the cursor/selection symbols, and a couple of other display options live in `~/.config/gunmap.toml`. It's plain TOML, safe to hand-edit. If you break it, delete the file and rerun `install.sh` to regenerate the defaults.

```toml
cyan  = "#00C8E8"  # primary accent
green = "#00E87A"  # selected row
amber = "#E8A200"  # warnings
red   = "#E84040"  # errors
```

## Requirements

- Go 1.21+
- root/sudo (raw socket access for scanning)
- a terminal with truecolor support if you want the UI to look right

## Contributing

Issues and PRs are welcome. If you're adding a feature, open an issue first so we can talk through it before you write code nobody asked for.

## License

MIT
