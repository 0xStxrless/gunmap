package styles

import (
	"image/color"
	"os"
	"os/user"
	"path/filepath"

	"charm.land/lipgloss/v2"
	"github.com/BurntSushi/toml"
)

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

func configPath() string {
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
		// Don't crash the TUI over a bad config file fall back to
		// defaults and let the user know via stderr.
		os.Stderr.WriteString("gnmap: warning: failed to load " + configPath() + ": " + err.Error() + "\n")
		return
	}
	applyOverrides(cfg)
}
