import re

def parse_map(filepath):
    mapping = {}
    with open(filepath, 'r') as f:
        for line in f:
            match = re.match(r'(.+): U\+([0-9A-F]+)', line)
            if match:
                name, hex_val = match.groups()
                mapping[name] = chr(int(hex_val, 16))
    return mapping

def generate_evdev_mapping(font_map):
    evdev_to_kenney = {
        # Letters
        **{f'KEY_{c.upper()}': font_map[f'keyboard_{c}'] for c in 'abcdefghijklmnopqrstuvwxyz' if f'keyboard_{c}' in font_map},
        # Numbers
        **{f'KEY_{i}': font_map[f'keyboard_{i}'] for i in range(10) if f'keyboard_{i}' in font_map},
        # Special Keys
        'KEY_SPACE': font_map.get('keyboard_space'),
        'KEY_ENTER': font_map.get('keyboard_enter'),
        'KEY_BACKSPACE': font_map.get('keyboard_backspace'),
        'KEY_LEFTSHIFT': font_map.get('keyboard_shift'),
        'KEY_RIGHTSHIFT': font_map.get('keyboard_shift'),
        'KEY_LEFTCTRL': font_map.get('keyboard_ctrl'),
        'KEY_RIGHTCTRL': font_map.get('keyboard_ctrl'),
        'KEY_LEFTALT': font_map.get('keyboard_alt'),
        'KEY_RIGHTALT': font_map.get('keyboard_alt'),
        'KEY_TAB': font_map.get('keyboard_tab'),
        'KEY_LEFTMETA': font_map.get('keyboard_win'),
        'KEY_RIGHTMETA': font_map.get('keyboard_win'),
        'KEY_ESC': font_map.get('keyboard_escape'),
        'KEY_CAPSLOCK': font_map.get('keyboard_capslock'),
        'KEY_UP': font_map.get('keyboard_arrow_up'),
        'KEY_DOWN': font_map.get('keyboard_arrow_down'),
        'KEY_LEFT': font_map.get('keyboard_arrow_left'),
        'KEY_RIGHT': font_map.get('keyboard_arrow_right'),
        'KEY_APOSTROPHE': font_map.get('keyboard_apostrophe'),
        'KEY_COMMA': font_map.get('keyboard_comma'),
        'KEY_DOT': font_map.get('keyboard_period'),
        'KEY_SLASH': font_map.get('keyboard_slash_forward'),
        'KEY_BACKSLASH': font_map.get('keyboard_slash_back'),
        'KEY_SEMICOLON': font_map.get('keyboard_semicolon'),
        'KEY_EQUAL': font_map.get('keyboard_equals'),
        'KEY_MINUS': font_map.get('keyboard_minus'),
        'KEY_LEFTBRACE': font_map.get('keyboard_bracket_open'),
        'KEY_RIGHTBRACE': font_map.get('keyboard_bracket_close'),
        'KEY_GRAVE': font_map.get('keyboard_tilde'),
        'KEY_PRINT': font_map.get('keyboard_printscreen'),
        'KEY_DELETE': font_map.get('keyboard_delete'),
        'KEY_PAGEUP': font_map.get('keyboard_page_up'),
        'KEY_PAGEDOWN': font_map.get('keyboard_page_down'),
        'KEY_HOME': font_map.get('keyboard_home'),
        'KEY_END': font_map.get('keyboard_end'),
        'KEY_INSERT': font_map.get('keyboard_insert'),
        'KEY_NUMLOCK': font_map.get('keyboard_numlock'),
        'KEY_KPSLASH': font_map.get('keyboard_slash_forward'),
        'KEY_KPASTERISK': font_map.get('keyboard_asterisk'),
        'KEY_KPMINUS': font_map.get('keyboard_minus'),
        'KEY_KPPLUS': font_map.get('keyboard_plus'),
        'KEY_KPENTER': font_map.get('keyboard_numpad_enter'),
        # Function Keys
        **{f'KEY_F{i}': font_map.get(f'keyboard_f{i}') for i in range(1, 13) if f'keyboard_f{i}' in font_map},
    }
    return {k: v for k, v in evdev_to_kenney.items() if v is not None}

if __name__ == "__main__":
    font_map = parse_map('font/kenney_input_keyboard_&_mouse_map.txt')
    evdev_map = generate_evdev_mapping(font_map)
    print("special_keys = {")
    for k, v in sorted(evdev_map.items()):
        print(f"    '{k}': {repr(v)},")
    print("}")
