package bambu

import (
	"crosshatch/internal/dtos"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"sync/atomic"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Fan node ids as used by the M106 P<fan> gcode.
const (
	FanPart    = 1
	FanAux     = 2
	FanChamber = 3
)

// AMS drying control modes, matching BambuStudio's DevAms::DryCtrlMode.
const (
	DryModeOff    = 0
	DryModeOnTime = 1
)

// StatusUpdateHandler is invoked with the latest merged print state whenever a
// report message containing a "print" payload is received.
type StatusUpdateHandler func(serial string, state *dtos.BambuPrintState)

type BambuClient struct {
	ip        string
	accessCode string
	serial    string

	mqttClient mqtt.Client

	onStatusUpdate StatusUpdateHandler

	stateMu sync.RWMutex
	state   *dtos.BambuPrintState

	commandSequence atomic.Int64
}

func (c *BambuClient) reportTopic() string {
	return fmt.Sprintf("device/%s/report", c.serial)
}

func (c *BambuClient) commandTopic() string {
	return fmt.Sprintf("device/%s/request", c.serial)
}

// State returns a reference to the latest merged print state, or nil if no
// report has been received yet.
func (c *BambuClient) State() *dtos.BambuPrintState {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

func (c *BambuClient) onConnect(client mqtt.Client) {
	fmt.Printf("Connected to Bambu printer %s\n", c.serial)

	if token := client.Subscribe(c.reportTopic(), 0, c.onMessage); token.Wait() && token.Error() != nil {
		fmt.Printf("Failed to subscribe to %q: %v\n", c.reportTopic(), token.Error())
	}
}

func (c *BambuClient) onMessage(_ mqtt.Client, msg mqtt.Message) {
	var envelope struct {
		Print json.RawMessage `json:"print"`
	}
	if err := json.Unmarshal(msg.Payload(), &envelope); err != nil {
		fmt.Printf("Received invalid MQTT message on %q: %v\n", c.reportTopic(), err)
		return
	}

	if len(envelope.Print) == 0 {
		return
	}

	c.stateMu.Lock()
	if c.state == nil {
		c.state = &dtos.BambuPrintState{}
	}
	// Decoding onto the accumulated state merges the partial report: fields
	// absent from this message keep their previous value, while arrays are
	// replaced wholesale.
	if err := json.Unmarshal(envelope.Print, c.state); err != nil {
		c.stateMu.Unlock()
		fmt.Printf("Failed to decode print state on %q: %v\n", c.reportTopic(), err)
		return
	}
	state := c.state
	c.stateMu.Unlock()

	if c.onStatusUpdate != nil {
		c.onStatusUpdate(c.serial, state)
	}
}

func (c *BambuClient) onDisconnect(client mqtt.Client, err error) {
	fmt.Printf("Disconnected: %v\n", err)
}

func (c *BambuClient) Close() {
	c.mqttClient.Disconnect(250)
}

// sendCommand stamps the payload's "print", "info", and "system" sub-objects
// with an incrementing sequence_id and publishes it to the command topic.
func (c *BambuClient) sendCommand(payload map[string]any) error {
	for _, key := range []string{"print", "info", "system"} {
		if sub, ok := payload[key].(map[string]any); ok {
			sub["sequence_id"] = strconv.FormatInt(c.commandSequence.Add(1)-1, 10)
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	token := c.mqttClient.Publish(c.commandTopic(), 0, false, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish command: %w", token.Error())
	}
	return nil
}

// Control methods

func (c *BambuClient) StopPrint() error {
	return c.sendCommand(map[string]any{"print": map[string]any{"command": "stop"}})
}

func (c *BambuClient) PausePrint() error {
	return c.sendCommand(map[string]any{"print": map[string]any{"command": "pause"}})
}

func (c *BambuClient) ResumePrint() error {
	return c.sendCommand(map[string]any{"print": map[string]any{"command": "resume"}})
}

func (c *BambuClient) SetLight(on bool) error {
	mode := "off"
	if on {
		mode = "on"
	}

	return c.sendCommand(map[string]any{
		"system": map[string]any{
			"command":       "ledctrl",
			"led_node":      "chamber_light",
			"led_mode":      mode, // "on" | "off" | "flashing"
			"led_on_time":   500,
			"led_off_time":  500,
			"loop_times":    1,
			"interval_time": 1000,
		},
	})
}

// SetPrintSpeed selects a speed profile: 1=silent, 2=standard, 3=sport,
// 4=ludicrous. The level is sent as a string param, matching OrcaSlicer.
func (c *BambuClient) SetPrintSpeed(level int) error {
	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command": "print_speed",
			"param":   strconv.Itoa(level),
		},
	})
}

// SetFanSpeed sets a fan (FanPart/FanAux/FanChamber) to a 0-100 percentage,
// issued as an M106 gcode line with the 0-255 duty cycle Bambu expects.
func (c *BambuClient) SetFanSpeed(fan int, percent int) error {
	duty := int(math.Round(float64(percent) / 100.0 * 255.0))
	duty = min(max(duty, 0), 255)

	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command": "gcode_line",
			"param":   fmt.Sprintf("M106 P%d S%d\n", fan, duty),
		},
	})
}

// SetFilament assigns a filament preset, colour and temperature window to an
// AMS slot (ams_id/tray_id) via the ams_filament_setting command.
func (c *BambuClient) SetFilament(amsID, trayID int, trayInfoIdx, trayColor, trayType string, nozzleTempMin, nozzleTempMax int) error {
	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command":         "ams_filament_setting",
			"ams_id":          amsID,
			"tray_id":         trayID,
			"tray_info_idx":   trayInfoIdx,
			"tray_color":      trayColor,
			"nozzle_temp_min": nozzleTempMin,
			"nozzle_temp_max": nozzleTempMax,
			"tray_type":       trayType,
			"setting_id":      "",
		},
	})
}

// StartDrying begins a timed drying cycle on an AMS unit. temp is the target
// temperature in °C, duration is in hours, and coolingTemp is the temperature
// the unit cools to afterwards. Only AMS 2 Pro / AMS HT units support this.
func (c *BambuClient) StartDrying(amsID, temp, durationHours, coolingTemp int, filament string, rotateTray bool) error {
	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command":              "ams_filament_drying",
			"ams_id":               amsID,
			"mode":                 DryModeOnTime,
			"filament":             filament,
			"temp":                 temp,
			"duration":             durationHours,
			"humidity":             0,
			"rotate_tray":          rotateTray,
			"cooling_temp":         coolingTemp,
			"close_power_conflict": false,
		},
	})
}

// StopDrying ends an in-progress drying cycle on an AMS unit.
func (c *BambuClient) StopDrying(amsID int) error {
	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command":              "ams_filament_drying",
			"ams_id":               amsID,
			"mode":                 DryModeOff,
			"filament":             "",
			"temp":                 0,
			"duration":             0,
			"humidity":             0,
			"rotate_tray":          false,
			"cooling_temp":         0,
			"close_power_conflict": false,
		},
	})
}

func (c *BambuClient) UnloadMaterial(amsID int) error {
	return c.sendCommand(map[string]any{
		"print": map[string]any{
			"command":   "ams_change_filament",
			"curr_temp": 255,
			"tar_temp":  255,
			"ams_id":    amsID,
			"target":    255,
			"slot_id":   255,
		},
	})
}

func NewBambuClient(ip string, accessCode string, serial string, onStatusUpdate StatusUpdateHandler) *BambuClient {
	client := &BambuClient{
		ip:             ip,
		accessCode:      accessCode,
		serial:         serial,
		onStatusUpdate: onStatusUpdate,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtts://%s:%d", ip, 8883))
	opts.SetClientID(fmt.Sprintf("crosshatch-%s", serial))
	opts.SetUsername("bblp")
	opts.SetPassword(accessCode)
	opts.SetKeepAlive(60)
	opts.SetAutoReconnect(true)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	opts.OnConnect = client.onConnect
	opts.OnConnectionLost = client.onDisconnect

	client.mqttClient = mqtt.NewClient(opts)

	go func() {
		if token := client.mqttClient.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("Error connecting to MQTT broker: %v\n", token.Error())
		}
	}()

	return client
}
