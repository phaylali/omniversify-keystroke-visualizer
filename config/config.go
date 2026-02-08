package config

import (
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	FontFamily string
	FontSize   int
	TextColor  string
	BgColor    string
	PaddingX   int
	PaddingY   int
	DurationMs int
	Position   string
	XOffset    int
	YOffset    int
}

func Load(filename string) (*Config, error) {
	cfg := &Config{
		FontFamily: "sans-serif",
		FontSize:   24,
		TextColor:  "white",
		BgColor:    "#2E2E2E",
		PaddingX:   20,
		PaddingY:   10,
		DurationMs: 1500,
		Position:   "bottom-center",
		XOffset:    0,
		YOffset:    -150,
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return cfg, nil
	}

	iniFile, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	if iniFile.HasSection("Appearance") {
		section := iniFile.Section("Appearance")
		if section.HasKey("font_family") {
			cfg.FontFamily = section.Key("font_family").String()
		}
		if section.HasKey("font_size") {
			cfg.FontSize, _ = section.Key("font_size").Int()
		}
		if section.HasKey("text_color") {
			cfg.TextColor = section.Key("text_color").String()
		}
		if section.HasKey("bg_color") {
			cfg.BgColor = section.Key("bg_color").String()
		}
		if section.HasKey("padding_x") {
			cfg.PaddingX, _ = section.Key("padding_x").Int()
		}
		if section.HasKey("padding_y") {
			cfg.PaddingY, _ = section.Key("padding_y").Int()
		}
		if section.HasKey("duration_ms") {
			cfg.DurationMs, _ = section.Key("duration_ms").Int()
		}
	}

	if iniFile.HasSection("Position") {
		section := iniFile.Section("Position")
		if section.HasKey("position") {
			cfg.Position = section.Key("position").String()
		}
		if section.HasKey("x_offset") {
			cfg.XOffset, _ = section.Key("x_offset").Int()
		}
		if section.HasKey("y_offset") {
			cfg.YOffset, _ = section.Key("y_offset").Int()
		}
	}

	return cfg, nil
}
