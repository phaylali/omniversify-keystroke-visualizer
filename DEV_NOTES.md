# Developer Notes - Omniversify Keystroke Visualizer

This project evolved from a Go-based X11 visualizer to a Python-based hardware listener to handle Wayland security restrictions.

## 1. Technical Strategy: Why Python & Evdev?

- **The Problem**: On Wayland (GNOME/KDE), security policies prevent standard applications from capturing keystrokes globally. The Go-based X11 implementation resulted in "black boxes" or "Permission Denied" errors.
- **The Solution**: Direct hardware access via `/dev/input/event*`. Since Linux treats everything as a file, we can read raw keyboard events if we have root permissions.
- **Library**: `evdev` (Python) is used instead of `pynput` as it provides true global hardware access that bypasses the Wayland compositor's filters.

## 2. Dependencies & Runtime

- **Manager**: `uv` is used for high-speed, isolated dependency management.
- **Command**: MUST be run with `sudo` and preserved path to work on Wayland:
  ```bash
  sudo env "PATH=$PATH" uv run --with evdev python3 visualizer.py
  ```

## 3. UI Implementation

- **Framework**: Tkinter (Python's built-in GUI).
- **Overlay**: Uses `Toplevel` windows with `overrideredirect(True)` and `-topmost` attributes.
- **Font**: Uses **Kenney Input Keyboard & Mouse** for a gamified look.
  - _Setup_: Run `mkdir -p ~/.local/share/fonts && cp ./font/* ~/.local/share/fonts && fc-cache -f` to register the font with the OS so Tkinter can find it.

## 4. Architecture

- **Multi-threaded**: A background thread is spawned for every keyboard device found in `evdev.list_devices()`.
- **Event Queue**: Keystrokes are sent to a thread-safe queue.
- **Main Loop**: Tkinter processes the queue every 50ms and spawns transient overlay windows.

## 5. Key Mappings

- The `listen_device` function in `visualizer.py` contains a dictionary mapping hardware `KEY_X` codes to Kenney Font PUA glyphs (e.g., `\ue0c5` for Space).

## 6. Maintenance Checklist for AI/Devs

- [ ] If no keys work, check `/dev/input/` permissions or ensure `sudo` is used.
- [ ] If font looks generic, the Kenney font isn't registered with the OS.
- [ ] Debug logs are written to `py_debug.log`.
