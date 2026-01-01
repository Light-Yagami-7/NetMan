package main

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Get available Wi-Fi networks
	cmd := exec.Command("nmcli", "-t", "-f", "SSID", "device", "wifi", "list")
	outCmd, _ := cmd.Output()
	uCmd := string(outCmd)
	trimCmd := strings.Split(uCmd, "\n")

	// Filter out empty SSIDs
	var networks []string
	for _, ssid := range trimCmd {
		if strings.TrimSpace(ssid) != "" {
			networks = append(networks, ssid)
		}
	}

	a := app.New()
	w := a.NewWindow("Network Manager")

	selected := -1
	list := widget.NewList(
		func() int {
			return len(networks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(networks[i])
		},
	)
	list.OnSelected = func(id int) {
		selected = id
		fmt.Println("SSID =", networks[selected])
	}

	// Scrollable list for networks
	scrollList := container.NewVScroll(list)
	scrollList.SetMinSize(fyne.NewSize(300, 200))

	status := widget.NewLabel("...")
	connectedLabel := widget.NewLabel("Connected Network:")

	bashActiveSSID := exec.Command(
		"bash", "-c",
		"nmcli -t -f active,ssid dev wifi | grep '^yes' | cut -d: -f2",
	)
	activeSSID, _ := bashActiveSSID.Output()
	trimActSSID := strings.TrimSpace(string(activeSSID))
	connectedSSID := widget.NewLabel(trimActSSID) // Shows the active SSID

	var globalWifiConName string
	var interfaceName string

	connectBtn := widget.NewButton("Connect", func() {
		if selected == -1 {
			dialog.ShowInformation("Error", "Please select a network", w)
			return
		}

		interfaceOut, _ := exec.Command(
			"nmcli", "-t", "-f", "DEVICE,TYPE", "device",
		).Output()

		for _, line := range strings.Split(string(interfaceOut), "\n") {
			fields := strings.Split(line, ":")
			if len(fields) == 2 && fields[1] == "wifi" {
				interfaceName = fields[0]
				break
			}
		}

		if interfaceName == "" {
			dialog.ShowError(fmt.Errorf("No Wi-Fi interface found"), w)
			return
		}

		ssid := networks[selected]
		globalWifiConName = ssid

		tempConCmdBuild := exec.Command(
			"bash", "-c",
			fmt.Sprintf("nmcli -t -f NAME connection show | grep -Fx '%s'", ssid),
		)
		tempConCmd, _ := tempConCmdBuild.Output()
		tempCon := strings.TrimSpace(string(tempConCmd))

		if ssid == tempCon {
			exec.Command("nmcli", "con", "up", "id", ssid).Run()
			status.SetText("Connected!")
			return
		}

		password := widget.NewPasswordEntry()
		dialog.ShowForm(
			"Enter Wi-Fi Password",
			"Connect",
			"Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Password", password),
			},
			func(ok bool) {
				status.SetText("Making a Connection...")
				exec.Command(
					"nmcli", "con", "add",
					"type", "wifi",
					"ifname", interfaceName,
					"con-name", ssid,
					"ssid", ssid,
				).Run()
				exec.Command("nmcli", "con", "modify", ssid, "wifi-sec.key-mgmt", "wpa-psk").Run()
				exec.Command("nmcli", "con", "modify", ssid, "wifi-sec.psk", password.Text).Run()
				exec.Command("nmcli", "con", "up", "id", ssid).Run()

				status.SetText("Connected!!")
				status.SetText("Connected to " + ssid)
				connectedSSID.SetText(ssid) // Show SSID at the top
				fmt.Println("Connected to", ssid)
			},
			w,
		)
	})

	disconnectBtn := widget.NewButton("Disconnect", func() {
		if globalWifiConName == "" {
			dialog.ShowInformation("Error", "No active connection to disconnect", w)
			return
		}
		exec.Command("nmcli", "dev", "disconnect", interfaceName).Run()
		status.SetText("Disconnected.")
		connectedSSID.SetText("None") // Reset SSID at the top
	})

	// Add separator between Connected and Available sections
	separator := widget.NewSeparator()

	content := container.NewVBox(
		connectedLabel, // Top section
		connectedSSID,
		separator, // Horizontal line separator
		widget.NewLabel("Available Networks:"),
		scrollList,
		connectBtn,
		disconnectBtn,
		status,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 450)) // Adjust size for top section
	w.ShowAndRun()
}
