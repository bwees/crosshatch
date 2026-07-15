package services

import (
	"context"
	"crosshatch/internal/bambu"
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/repositories"
	"crosshatch/internal/socketio"
	"fmt"
	"log/slog"
	"sync"

	"go.uber.org/fx"
)

type PrinterService struct {
	printerRepo  *repositories.PrinterRepository
	cameraRepo   *repositories.CameraRepository
	filamentRepo *repositories.FilamentRepository
	socketio     *socketio.SocketIO

	clientsMu sync.Mutex
	clients   map[string]*bambu.BambuClient

	statusMu    sync.RWMutex
	statusCache map[string]dtos.PrinterStatus

	observersMu sync.Mutex
	observers   []StatusObserver
}

// StatusObserver is notified of every printer status transition. prev is nil
// when no prior status has been cached for the printer.
type StatusObserver func(serial string, prev, next *dtos.PrinterStatus)

func (s *PrinterService) AddObserver(o StatusObserver) {
	s.observersMu.Lock()
	s.observers = append(s.observers, o)
	s.observersMu.Unlock()
}

func (s *PrinterService) GetPrinters() ([]models.Printer, error) {
	return s.printerRepo.GetPrinters()
}

func (s *PrinterService) CreatePrinter(dto dtos.CreatePrinterDto) (*models.Printer, error) {
	printer, err := s.printerRepo.CreatePrinter(dto)
	if err != nil {
		return nil, err
	}

	s.connectClient(*printer)
	s.reconcileCameraStreams()

	return printer, nil
}

func (s *PrinterService) UpdatePrinter(serial string, dto dtos.UpdatePrinterDto) (*models.Printer, error) {
	existing, err := s.printerRepo.GetPrinterBySerial(serial)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("printer with serial %s not found", serial)
	}

	updated, err := s.printerRepo.UpdatePrinter(serial, dto)
	if err != nil {
		return nil, err
	}

	s.reconcileCameraStreams()

	return updated, nil
}

func (s *PrinterService) DeletePrinter(serial string) error {
	if err := s.printerRepo.DeletePrinter(serial); err != nil {
		return err
	}

	s.disconnectClient(serial)
	s.reconcileCameraStreams()

	return nil
}

// connectClient creates and registers a Bambu MQTT client for a printer,
// replacing any existing client for the same serial.
func (s *PrinterService) connectClient(printer models.Printer) {
	client := bambu.NewBambuClient(printer.HostIP, printer.AccessCode, printer.Serial, s.onPrinterStatusUpdate)

	s.clientsMu.Lock()
	if existing, ok := s.clients[printer.Serial]; ok {
		existing.Close()
	}
	s.clients[printer.Serial] = client
	s.clientsMu.Unlock()
}

// disconnectClient closes and removes the Bambu MQTT client for a serial.
func (s *PrinterService) disconnectClient(serial string) {
	s.clientsMu.Lock()
	client, ok := s.clients[serial]
	delete(s.clients, serial)
	s.clientsMu.Unlock()

	if ok {
		client.Close()
	}
}

// client returns the Bambu MQTT client for a serial, or an error if no client
// is registered (e.g. the printer does not exist).
func (s *PrinterService) client(serial string) (*bambu.BambuClient, error) {
	s.clientsMu.Lock()
	client, ok := s.clients[serial]
	s.clientsMu.Unlock()

	if !ok {
		return nil, fmt.Errorf("MQTT client for printer %s not found", serial)
	}
	return client, nil
}

func (s *PrinterService) StopPrint(serial string) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.StopPrint()
}

func (s *PrinterService) PausePrint(serial string) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.PausePrint()
}

func (s *PrinterService) ResumePrint(serial string) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.ResumePrint()
}

func (s *PrinterService) SetLight(serial string, on bool) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.SetLight(on)
}

func (s *PrinterService) UnloadMaterial(serial string, amsID int) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.UnloadMaterial(amsID)
}

func (s *PrinterService) StartDrying(serial string, amsID int, dto dtos.StartDryingDto) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.StartDrying(amsID, dto.Temperature, dto.Duration, dto.CoolingTemp, dto.Filament, dto.RotateTray)
}

func (s *PrinterService) StopDrying(serial string, amsID int) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.StopDrying(amsID)
}

func (s *PrinterService) SetPrintSpeed(serial string, level int) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.SetPrintSpeed(level)
}

var fanNodes = map[string]int{
	"part":    bambu.FanPart,
	"aux":     bambu.FanAux,
	"chamber": bambu.FanChamber,
}

func (s *PrinterService) SetFan(serial string, fan string, speed int) error {
	node, ok := fanNodes[fan]
	if !ok {
		return fmt.Errorf("unknown fan %q", fan)
	}
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.SetFanSpeed(node, speed)
}

func (s *PrinterService) SetFilament(serial string, dto dtos.SetFilamentDto) error {
	client, err := s.client(serial)
	if err != nil {
		return err
	}
	return client.SetFilament(dto.AmsID, dto.TrayID, dto.TrayInfoIdx, dto.TrayColor, dto.TrayType, dto.NozzleTempMin, dto.NozzleTempMax)
}

// printerStatusPayload flattens the printer status alongside its serial, so the
// emitted event matches the `{ serial, ...status }` shape the clients expect.
type printerStatusPayload struct {
	Serial string `json:"serial"`
	dtos.PrinterStatus
}

// onPrinterStatusUpdate is handed to each Bambu client; it projects the merged
// MQTT print state into a status DTO and broadcasts it.
func (s *PrinterService) onPrinterStatusUpdate(serial string, state *dtos.BambuPrintState) {
	status := dtos.StatusFromMQTT(state)

	if err := s.filamentRepo.CreateMissing(filamentsFromStatus(status)); err != nil {
		slog.Error("failed to record filaments", "serial", serial, "error", err)
	}

	s.statusMu.RLock()
	var prev *dtos.PrinterStatus
	if cached, ok := s.statusCache[serial]; ok {
		prevCopy := cached
		prev = &prevCopy
	}
	s.statusMu.RUnlock()

	s.BroadcastStatus(serial, status)
	s.notifyObservers(serial, prev, &status)
}

func (s *PrinterService) notifyObservers(serial string, prev, next *dtos.PrinterStatus) {
	s.observersMu.Lock()
	observers := s.observers
	s.observersMu.Unlock()

	for _, o := range observers {
		o(serial, prev, next)
	}
}

func (s *PrinterService) BroadcastStatus(serial string, status dtos.PrinterStatus) {
	s.statusMu.Lock()
	s.statusCache[serial] = status
	s.statusMu.Unlock()

	s.socketio.Emit("printer.status", printerStatusPayload{Serial: serial, PrinterStatus: status})
}

// GetStatus returns the latest known status for a printer, or an error if no
// status has been received yet (e.g. the printer is offline or unknown).
func (s *PrinterService) GetStatus(serial string) (*dtos.PrinterStatus, error) {
	s.statusMu.RLock()
	status, ok := s.statusCache[serial]
	s.statusMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no status available for printer %s", serial)
	}
	return &status, nil
}

func (s *PrinterService) replayStatus(emit socketio.EmitFunc) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	for serial, status := range s.statusCache {
		emit("printer.status", printerStatusPayload{Serial: serial, PrinterStatus: status})
	}
}

func (s *PrinterService) connectPrinterClients() {
	printers, err := s.GetPrinters()
	if err != nil {
		slog.Error("failed to fetch printers", "error", err)
		return
	}

	for _, printer := range printers {
		s.connectClient(printer)
	}
}

func containsPrinter(printers []models.Printer, serial string) bool {
	for _, printer := range printers {
		if printer.Serial == serial {
			return true
		}
	}
	return false
}

func (s *PrinterService) reconcileCameraStreams() {
	slog.Info("reconciling camera streams with printers")

	printers, err := s.GetPrinters()
	if err != nil {
		slog.Error("failed to fetch printers for stream reconciliation", "error", err)
		return
	}

	streams, err := s.cameraRepo.GetStreams()
	if err != nil {
		slog.Error("failed to fetch camera streams for reconciliation", "error", err)
		return
	}

	for _, printer := range printers {
		stream, exists := streams[printer.Serial]
		if !exists || len(stream.Producers) == 0 || stream.Producers[0].URL != printer.CameraURL() {
			if exists {
				if err := s.cameraRepo.UpdateStream(printer.Serial, printer.CameraURL()); err != nil {
					slog.Error("failed to update camera stream", "serial", printer.Serial, "error", err)
				}
			} else {
				if err := s.cameraRepo.AddStream(printer.Serial, printer.CameraURL()); err != nil {
					slog.Error("failed to add camera stream", "serial", printer.Serial, "error", err)
				}
			}
		}
	}

	for serial := range streams {
		if !containsPrinter(printers, serial) {
			if err := s.cameraRepo.DeleteStream(serial); err != nil {
				slog.Error("failed to delete camera stream", "serial", serial, "error", err)
			}
		}
	}
}

func NewPrinterService(lc fx.Lifecycle, printerRepo *repositories.PrinterRepository, cameraRepo *repositories.CameraRepository, filamentRepo *repositories.FilamentRepository, socket *socketio.SocketIO) *PrinterService {
	svc := &PrinterService{
		printerRepo:  printerRepo,
		cameraRepo:   cameraRepo,
		filamentRepo: filamentRepo,
		socketio:     socket,
		clients:      make(map[string]*bambu.BambuClient),
		statusCache:  make(map[string]dtos.PrinterStatus),
	}

	// Send each newly connected client the latest known status for every printer.
	socket.OnConnect(svc.replayStatus)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go svc.connectPrinterClients()
			go svc.reconcileCameraStreams()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			svc.clientsMu.Lock()
			defer svc.clientsMu.Unlock()
			for _, client := range svc.clients {
				client.Close()
			}
			return nil
		},
	})

	return svc
}
