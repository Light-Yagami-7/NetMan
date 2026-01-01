# NetMan

NetMan is a lightweight Linux Wi-Fi manager GUI built using **Go** and **Fyne**.  
It provides a simple graphical interface for listing available Wi-Fi networks,
connecting to them, and disconnecting — powered by **NetworkManager (`nmcli`)**.

This project is intended as a minimal, no-bloat alternative for users who prefer
simple tools over full desktop network managers.

---

## Features

- List available Wi-Fi networks
- Connect to WPA/WPA2 secured networks
- Disconnect from the active network
- Display currently connected SSID
- Minimal and lightweight GUI
- Linux-only (NetworkManager based)

---

## Screenshots

![NetMan Screenshot](assets/netman.png)

---

## Requirements

- Linux system
- NetworkManager installed
- `nmcli` available in PATH
- Go 1.20 or newer (for building from source)
- Wayland or X11  
  *(Tested on Hyprland)*

---

## Installation

### Build from source

```bash
git clone https://github.com/Light-Yagami-7/NetMan.git
cd NetMan
go build -o netman

Optional: install system-wide

mv netman ~/.local/bin/

Make sure ~/.local/bin is in your $PATH.
Usage

Run the application:

netman

    Select a Wi-Fi network from the list

    Click Connect and enter the password if required

    Use Disconnect to disconnect from the current network

    Connection status is displayed at the bottom of the window

Notes

    NetMan relies entirely on nmcli; NetworkManager must be running

    This is an experimental project — expect rough edges
