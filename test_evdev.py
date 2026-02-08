import evdev

devices = [evdev.InputDevice(path) for path in evdev.list_devices()]
for device in devices:
    print(f"path: {device.path}, name: {device.name}, phys: {device.phys}")
