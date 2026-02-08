import tkinter as tk
import threading
import configparser
import os
import datetime
import evdev
from evdev import ecodes
import ctypes

def log(msg):
    with open("py_debug.log", "a") as f:
        f.write(f"{datetime.datetime.now()} {msg}\n")

def load_font(font_family):
    """Checks if a font family is available to Tkinter and logs alternatives."""
    from tkinter import font
    available_fonts = font.families()
    log(f"Checking for font: '{font_family}'")
    if font_family in available_fonts:
        log(f"Font family '{font_family}' is LOADED and available.")
        return True
    else:
        log(f"WARNING: Font family '{font_family}' NOT FOUND.")
        # Find similar names
        kenney_fonts = [f for f in available_fonts if 'Kenney' in f]
        if kenney_fonts:
            log(f"Found similar font names: {kenney_fonts}")
            log(f"SUGGESTION: Try setting font_family to '{kenney_fonts[0]}' in config.ini or code.")
        else:
            log("No fonts containing 'Kenney' were found in Tkinter families.")
        return False

log("PY_LOG: Script started using evdev version")

# --- Configuration Class ---
class Config:
    def __init__(self, filename='config.ini'):
        self.config = configparser.ConfigParser()
        if os.path.exists(filename):
            self.config.read(filename)
        
        self.font_family = self.config.get('Appearance', 'font_family', fallback='Kenney Input Keyboard & Mouse')
        self.font_size = self.config.getint('Appearance', 'font_size', fallback=24)
        self.text_color = self.config.get('Appearance', 'text_color', fallback='white')
        self.bg_color = self.config.get('Appearance', 'bg_color', fallback='#2E2E2E')
        self.padding_x = self.config.getint('Appearance', 'padding_x', fallback=20)
        self.padding_y = self.config.getint('Appearance', 'padding_y', fallback=10)
        self.duration_ms = self.config.getint('Appearance', 'duration_ms', fallback=1500)
        
        self.position = self.config.get('Position', 'position', fallback='bottom-center')
        self.x_offset = self.config.getint('Position', 'x_offset', fallback=0)
        self.y_offset = self.config.getint('Position', 'y_offset', fallback=-150)

# --- Main Application Class ---
class KeystrokeVisualizer:
    def __init__(self, root, config):
        self.root = root
        self.config = config
        self.root.withdraw()
        self.input_queue = []
        self.queue_lock = threading.Lock()
        self.process_queue()

    def calculate_position(self, window):
        screen_width = self.root.winfo_screenwidth()
        screen_height = self.root.winfo_screenheight()
        window.update_idletasks()
        win_width = window.winfo_width()
        win_height = window.winfo_height()

        position_map = {
            'top-left': (0, 0),
            'top-center': ((screen_width // 2) - (win_width // 2), 0),
            'top-right': (screen_width - win_width, 0),
            'center': ((screen_width // 2) - (win_width // 2), (screen_height // 2) - (win_height // 2)),
            'bottom-left': (0, screen_height - win_height),
            'bottom-center': ((screen_width // 2) - (win_width // 2), screen_height - win_height),
            'bottom-right': (screen_width - win_width, screen_height - win_height),
        }

        base_x, base_y = position_map.get(self.config.position, position_map['bottom-center'])
        return base_x + self.config.x_offset, base_y + self.config.y_offset

    def display_input(self, input_char):
        log(f"Displaying input: {input_char}")
        input_window = tk.Toplevel(self.root)
        input_window.wm_attributes("-topmost", True)
        input_window.overrideredirect(True)
        input_window.config(bg=self.config.bg_color)

        input_label = tk.Label(
            input_window,
            text=input_char,
            font=(self.config.font_family, self.config.font_size, "bold"),
            bg=self.config.bg_color,
            fg=self.config.text_color,
            padx=self.config.padding_x,
            pady=self.config.padding_y
        )
        input_label.pack()

        x_pos, y_pos = self.calculate_position(input_window)
        input_window.geometry(f"+{x_pos}+{y_pos}")
        input_window.after(self.config.duration_ms, input_window.destroy)

    def add_to_queue(self, input_char):
        log(f"Adding to queue: {input_char}")
        with self.queue_lock:
            self.input_queue.append(input_char)

    def process_queue(self):
        with self.queue_lock:
            if self.input_queue:
                input_char = self.input_queue.pop(0)
                self.display_input(input_char)
        self.root.after(50, self.process_queue)

# --- Evdev Listener Logic ---
def listen_device(device_path, app):
    try:
        device = evdev.InputDevice(device_path)
        log(f"Listening to: {device.name} ({device.path})")
        
        special_keys = {
            'KEY_0': '\ue001',
            'KEY_1': '\ue003',
            'KEY_2': '\ue005',
            'KEY_3': '\ue007',
            'KEY_4': '\ue009',
            'KEY_5': '\ue00b',
            'KEY_6': '\ue00d',
            'KEY_7': '\ue00f',
            'KEY_8': '\ue011',
            'KEY_9': '\ue013',
            'KEY_A': '\ue015',
            'KEY_APOSTROPHE': '\ue01b',
            'KEY_B': '\ue036',
            'KEY_BACKSLASH': '\ue0c1',
            'KEY_BACKSPACE': '\ue038',
            'KEY_C': '\ue046',
            'KEY_CAPSLOCK': '\ue048',
            'KEY_COMMA': '\ue050',
            'KEY_D': '\ue056',
            'KEY_DELETE': '\ue058',
            'KEY_DOT': '\ue0a9',
            'KEY_DOWN': '\ue01d',
            'KEY_E': '\ue05a',
            'KEY_END': '\ue05c',
            'KEY_ENTER': '\ue05e',
            'KEY_EQUAL': '\ue060',
            'KEY_ESC': '\ue062',
            'KEY_F': '\ue066',
            'KEY_F1': '\ue067',
            'KEY_F10': '\ue068',
            'KEY_F11': '\ue06a',
            'KEY_F12': '\ue06c',
            'KEY_F2': '\ue06f',
            'KEY_F3': '\ue071',
            'KEY_F4': '\ue073',
            'KEY_F5': '\ue075',
            'KEY_F6': '\ue077',
            'KEY_F7': '\ue079',
            'KEY_F8': '\ue07b',
            'KEY_F9': '\ue07d',
            'KEY_G': '\ue082',
            'KEY_GRAVE': '\ue0d1',
            'KEY_H': '\ue084',
            'KEY_HOME': '\ue086',
            'KEY_I': '\ue088',
            'KEY_INSERT': '\ue08a',
            'KEY_J': '\ue08c',
            'KEY_K': '\ue08e',
            'KEY_KPASTERISK': '\ue034',
            'KEY_KPENTER': '\ue09a',
            'KEY_KPMINUS': '\ue094',
            'KEY_KPPLUS': '\ue0ab',
            'KEY_KPSLASH': '\ue0c3',
            'KEY_L': '\ue090',
            'KEY_LEFT': '\ue01f',
            'KEY_LEFTALT': '\ue017',
            'KEY_LEFTBRACE': '\ue044',
            'KEY_LEFTCTRL': '\ue054',
            'KEY_LEFTMETA': '\ue0d9',
            'KEY_LEFTSHIFT': '\ue0bd',
            'KEY_M': '\ue092',
            'KEY_MINUS': '\ue094',
            'KEY_N': '\ue096',
            'KEY_NUMLOCK': '\ue098',
            'KEY_O': '\ue09e',
            'KEY_P': '\ue0a3',
            'KEY_PAGEDOWN': '\ue0a5',
            'KEY_PAGEUP': '\ue0a7',
            'KEY_PRINT': '\ue0ad',
            'KEY_Q': '\ue0af',
            'KEY_R': '\ue0b5',
            'KEY_RIGHT': '\ue021',
            'KEY_RIGHTALT': '\ue017',
            'KEY_RIGHTBRACE': '\ue03e',
            'KEY_RIGHTCTRL': '\ue054',
            'KEY_RIGHTMETA': '\ue0d9',
            'KEY_RIGHTSHIFT': '\ue0bd',
            'KEY_S': '\ue0b9',
            'KEY_SEMICOLON': '\ue0bb',
            'KEY_SLASH': '\ue0c3',
            'KEY_SPACE': '\ue0c5',
            'KEY_T': '\ue0c9',
            'KEY_TAB': '\ue0cb',
            'KEY_U': '\ue0d3',
            'KEY_UP': '\ue023',
            'KEY_V': '\ue0d5',
            'KEY_W': '\ue0d7',
            'KEY_X': '\ue0db',
            'KEY_Y': '\ue0dd',
            'KEY_Z': '\ue0df',
        }

        for event in device.read_loop():
            if event.type == ecodes.EV_KEY:
                key_event = evdev.categorize(event)
                if key_event.keystate == key_event.key_down:
                    key_code = key_event.keycode
                    if isinstance(key_code, list):
                        key_code = key_code[0]
                    
                    log(f"Raw key from {device.path}: {key_code}")
                    
                    display = ""
                    if key_code in special_keys:
                        display = special_keys[key_code]
                    elif key_code.startswith('KEY_'):
                        val = key_code.replace('KEY_', '')
                        if len(val) == 1:
                            display = val
                        else:
                            display = val.capitalize()
                    
                    if display:
                        app.add_to_queue(display)
    except Exception as e:
        log(f"Error listening to {device_path}: {e}")

def start_listeners(app):
    found_any = False
    for path in evdev.list_devices():
        try:
            dev = evdev.InputDevice(path)
            # Listen to anything that looks like a keyboard
            capabilities = dev.capabilities()
            if ecodes.EV_KEY in capabilities:
                # Basic check for a keyboard: has keys like 'A'
                if ecodes.KEY_A in capabilities[ecodes.EV_KEY]:
                    log(f"Found keyboard candidate: {dev.name} at {dev.path}")
                    thread = threading.Thread(target=listen_device, args=(path, app), daemon=True)
                    thread.start()
                    found_any = True
        except Exception as e:
            log(f"Error checking device {path}: {e}")
    
    if not found_any:
        log("No keyboard devices found!")

if __name__ == "__main__":
    config = Config('config.ini')
    root_window = tk.Tk()
    load_font("Kenney Input Keyboard & Mouse")
    app = KeystrokeVisualizer(root_window, config)
    
    start_listeners(app)

    root_window.after(2000, lambda: app.add_to_queue("Ready!"))
    root_window.mainloop()