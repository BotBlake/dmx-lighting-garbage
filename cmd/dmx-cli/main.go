package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"dmx-lighting-garbage/internal/device"
	"dmx-lighting-garbage/internal/dmx"
	serialhw "dmx-lighting-garbage/internal/hardware/serial"
)

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	fmt.Println("Scanning DMX interfaces...")

	devs, err := serialhw.ListDevices()
	if err != nil || len(devs) == 0 {
		fmt.Println("No serial devices found")
		return
	}

	fmt.Println("We have selected the following interfaces:")
	for i, d := range devs {
		fmt.Printf("[%d] %s\n", i, d.Name)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Select the one to use: ")
	line, _ := reader.ReadString('\n')
	idx, _ := strconv.Atoi(strings.TrimSpace(line))

	if idx < 0 || idx >= len(devs) {
		fmt.Println("Invalid selection")
		return
	}

	driver := serialhw.New(devs[idx].Name)
	if err := driver.Open(); err != nil {
		fmt.Println("Failed to open device:", err)
		return
	}
	defer driver.Close()

	fmt.Println("Device driver connected!")

	fmt.Print("Bind this device to universe: ")
	line, _ = reader.ReadString('\n')
	universeID, _ := strconv.Atoi(strings.TrimSpace(line))

	universes := map[int]*dmx.Universe{
		universeID: dmx.NewUniverse(universeID),
	}

	binding := dmx.Binding{
		Universe: universeID,
		Driver:   driver,
	}

	fmt.Println("Setup Completed... Loading DLG (DmxLightingGarbage)")
	ClearScreen()
	fmt.Println("Welcome to DLG (DmxLightingGarbage)")

	engine := dmx.NewEngine([]dmx.Binding{binding}, 25)
	go engine.Run(universes)

	// Load DDF profiles
	profiles, err := device.LoadProfiles("./ddf")
	if err != nil || len(profiles) == 0 {
		fmt.Println("Warning: no DDF profiles loaded")
	}

	// Device registry
	devices := device.NewRegistry()

	fmt.Println("Type `help` for commands")

	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		if handleCommand(
			strings.TrimSpace(cmd),
			universes,
			engine,
			profiles,
			devices,
			reader,
		) {
			break
		}
	}
}

func handleCommand(
	cmd string,
	universes map[int]*dmx.Universe,
	engine *dmx.Engine,
	profiles map[string]*device.DDFDevice,
	devices *device.Registry,
	reader *bufio.Reader,
) bool {

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return false
	}

	switch parts[0] {

	case "help":
		fmt.Println("Commands:")
		fmt.Println("  set <universe> <channel> <value>")
		fmt.Println("  get all")
		fmt.Println("  clear")
		fmt.Println("  device create")
		fmt.Println("  device set <device> <brightness|color> <value>")
		fmt.Println("  quit | exit")
		return false

	case "quit", "exit":
		fmt.Println("Shutting down DMX engine...")
		engine.Stop()
		return true

	/* ---------- RAW DMX ---------- */

	case "set":
		if len(parts) != 4 {
			fmt.Println("usage: set <universe> <channel> <value>")
			return false
		}

		u, _ := strconv.Atoi(parts[1])
		ch, _ := strconv.Atoi(parts[2])
		val, _ := strconv.Atoi(parts[3])

		if ch < 1 || ch > 512 || val < 0 || val > 255 {
			fmt.Println("Invalid channel or value")
			return false
		}

		if _, ok := universes[u]; !ok {
			universes[u] = dmx.NewUniverse(u)
		}

		universes[u].Frame[ch] = byte(val)
		fmt.Printf("Universe %d Channel %d = %d\n", u, ch, val)
		return false

	case "get":
		if len(parts) == 2 && parts[1] == "all" {
			for u, uni := range universes {
				fmt.Printf("Universe %d:\n", u)
				for i := 1; i <= 512; i++ {
					if uni.Frame[i] != 0 {
						fmt.Printf("  CH %d = %d\n", i, uni.Frame[i])
					}
				}
			}
		}
		return false

	case "clear":
		for _, uni := range universes {
			for i := 1; i <= 512; i++ {
				uni.Frame[i] = 0
			}
		}
		fmt.Println("All channels cleared")
		return false

	/* ---------- DEVICE ---------- */

	case "device":
		if len(parts) < 2 {
			fmt.Println("usage: device <create|set>")
			return false
		}

		switch parts[1] {

		case "create":
			if len(profiles) == 0 {
				fmt.Println("No DDF profiles available")
				return false
			}

			keys := make([]string, 0, len(profiles))
			fmt.Println("Available device profiles:")
			i := 0
			for k := range profiles {
				fmt.Printf("[%d] %s\n", i, k)
				keys = append(keys, k)
				i++
			}

			fmt.Print("Select profile: ")
			line, _ := reader.ReadString('\n')
			pidx, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil || pidx < 0 || pidx >= len(keys) {
				fmt.Println("Invalid selection")
				return false
			}

			profileKey := keys[pidx]
			profile := profiles[profileKey]

			fmt.Print("Device name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Universe: ")
			line, _ = reader.ReadString('\n')
			u, _ := strconv.Atoi(strings.TrimSpace(line))

			if _, ok := universes[u]; !ok {
				universes[u] = dmx.NewUniverse(u)
			}

			fmt.Print("DMX start address (1–512): ")
			line, _ = reader.ReadString('\n')
			addr, _ := strconv.Atoi(strings.TrimSpace(line))

			if addr < 1 || addr > 512 {
				fmt.Println("Invalid DMX address")
				return false
			}

			dev := &device.Instance{
				Name:     name,
				Profile:  profile,
				Universe: universes[u],
				Address:  addr,
			}

			devices.Add(dev)

			fmt.Printf(
				"Device '%s' created (%s @ U%d:%d)\n",
				name,
				profileKey,
				u,
				addr,
			)
			return false

		case "set":
			if len(parts) < 5 {
				fmt.Println("usage: device set <device> <brightness|color> <value>")
				return false
			}

			devName := parts[2]
			param := parts[3]
			value := parts[4]

			dev, ok := devices.Get(devName)
			if !ok {
				fmt.Println("Device not found")
				return false
			}

			switch param {
			case "brightness":
				v, err := strconv.Atoi(value)
				if err != nil || v < 0 || v > 100 {
					fmt.Println("Brightness must be 0–100")
					return false
				}
				dev.SetBrightness(v)
				fmt.Printf("Set %s brightness to %d%%\n", devName, v)

			case "color":
				dev.SetColor(value)
				fmt.Printf("Set %s color to %s\n", devName, value)

			default:
				fmt.Println("Unknown device parameter")
			}
			return false
		}

	default:
		fmt.Println("Unknown command")
		return false
	}

	return false
}
