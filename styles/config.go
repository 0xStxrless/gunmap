package styles

import (
	"image/color"
	"os"
	"os/user"
	"path/filepath"

	"charm.land/lipgloss/v2"
	"github.com/BurntSushi/toml"
)

// configColors mirrors the keys in ~/.config/gnmap.conf. Every field is a
// string hex color (or, for Selected/Cursor, a literal glyph) and is
// optional — any field left blank in the file keeps the built-in default.
type configColors struct {
	BgBase  string `toml:"bgBase"`
	BgPanel string `toml:"bgPanel"`
	BgCard  string `toml:"bgCard"`

	BorderSubtle string `toml:"borderSubtle"`
	BorderActive string `toml:"borderActive"`

	Cyan      string `toml:"cyan"`
	CyanDim   string `toml:"cyanDim"`
	CyanGhost string `toml:"cyanGhost"`

	Green string `toml:"green"`
	Amber string `toml:"amber"`
	Red   string `toml:"red"`

	TextPrimary string `toml:"textPrimary"`
	TextMuted   string `toml:"textMuted"`
	TextGhost   string `toml:"textGhost"`
	TextInverse string `toml:"textInverse"`

	CrumbActive string `toml:"crumbActive"`
	CrumbDim    string `toml:"crumbDim"`

	SidebarBgActive string `toml:"sidebarBgActive"`
	SidebarRule     string `toml:"sidebarRule"`

	Selected string `toml:"selected"`
	Cursor   string `toml:"cursor"`
}

// configPath returns ~/.config/gnmap.conf, or "" if the home dir can't be
// resolved.
func configPath() string {
	// When run via `sudo`, os.UserHomeDir() resolves to /root instead of
	// the invoking user's home. Prefer SUDO_USER's home in that case so
	// the config file is found regardless of how the binary is launched.
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if u, err := user.Lookup(sudoUser); err == nil {
			return filepath.Join(u.HomeDir, ".config", "gunmap.toml")
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "gunmap.toml")
}

// loadConfig reads the user's config file if present. A missing file is not
// an error — it just means "use defaults" — but a malformed file is
// reported so the user can fix it.
func loadConfig() (configColors, error) {
	var cfg configColors

	path := configPath()
	if path == "" {
		return cfg, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// applyOverrides replaces any non-empty field in cfg onto the package-level
// color/glyph vars, leaving the hardcoded defaults in place otherwise.
func applyOverrides(cfg configColors) {
	set := func(dst *color.Color, v string) {
		if v != "" {
			*dst = lipgloss.Color(v)
		}
	}

	set(&bgBase, cfg.BgBase)
	set(&bgPanel, cfg.BgPanel)
	set(&bgCard, cfg.BgCard)

	set(&borderSubtle, cfg.BorderSubtle)
	set(&borderActive, cfg.BorderActive)

	set(&cyan, cfg.Cyan)
	set(&cyanDim, cfg.CyanDim)
	set(&cyanGhost, cfg.CyanGhost)

	set(&green, cfg.Green)
	set(&amber, cfg.Amber)
	set(&red, cfg.Red)

	set(&textPrimary, cfg.TextPrimary)
	set(&textMuted, cfg.TextMuted)
	set(&textGhost, cfg.TextGhost)
	set(&textInverse, cfg.TextInverse)

	set(&crumbActive, cfg.CrumbActive)
	set(&crumbDim, cfg.CrumbDim)

	set(&sidebarBgActive, cfg.SidebarBgActive)
	set(&sidebarRule, cfg.SidebarRule)

	if cfg.Selected != "" {
		selectedGlyph = cfg.Selected
	}
	if cfg.Cursor != "" {
		cursorGlyph = cfg.Cursor
	}
}

func init() {
	cfg, err := loadConfig()
	if err != nil {
		// Don't crash the TUI over a bad config file — fall back to
		// defaults and let the user know via stderr.
		os.Stderr.WriteString("gnmap: warning: failed to load " + configPath() + ": " + err.Error() + "\n")
		return
	}
	applyOverrides(cfg)
}
