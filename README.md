# Omniversify Keystroke Visualizer

A beautiful, gamified keystroke visualizer for Linux and Windows. It displays your keystrokes in real-time using high-quality icons or fonts, making it perfect for screen recording, coding demonstrations, or gaming.

> [!IMPORTANT]
> This project currently features a **Python (evdev)** implementation which is highly recommended for Wayland users (GNOME/KDE) as it bypasses security restrictions that block other visualizers.

## Features

- **Gamified Aesthetics**: Uses the **Kenney Input** font for beautiful, controller-style key icons.
- **Wayland Support**: Built-in compatibility for modern Linux desktops via hardware-level event reading.
- **Customizable**: Control colors, font size, position, and duration via `config.ini`.
- **Lightweight**: Minimal performance overhead.

## Installation & Setup (Linux)

The Linux version uses `evdev` to read keystrokes directly from your hardware.

### Requirements

- [uv](https://docs.astral.sh/uv/) (Python package manager)
- Sudo privileges (required to read hardware input devices on Wayland)

### Font Installation (Required for Icons)

To see the gamified icons, you must install the font to your system:

```bash
mkdir -p ~/.local/share/fonts
cp ./font/* ~/.local/share/fonts
fc-cache -f
```

### Running the Visualizer

```bash
sudo env "PATH=$PATH" uv run --with evdev python3 visualizer.py
```

## Configuration

Edit `config.ini` to customize the appearance:

```ini
[Appearance]
font_family = Kenney Input Keyboard & Mouse
font_size = 80
text_color = white
bg_color = #2E2E2E
duration_ms = 1500

[Position]
position = bottom-right
x_offset = -20
y_offset = -150
```

### Position Options

- `top-left`, `top-center`, `top-right`
- `center`
- `bottom-left`, `bottom-center`, `bottom-right`

## License

MIT
