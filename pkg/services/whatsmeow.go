package services

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsmeowService struct {
	store          *sqlstore.Container
	clientsPointer map[string]*whatsmeow.Client
}

func NewWhatsmeowService(store *sqlstore.Container, clientsPointer map[string]*whatsmeow.Client) *WhatsmeowService {
	return &WhatsmeowService{
		store:          store,
		clientsPointer: clientsPointer,
	}
}

func (s *WhatsmeowService) eventHandler(instanceID string) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Connected:
			fmt.Printf("Instance %s connected\n", instanceID)

		case *events.Disconnected:
		// Connection lost — most likely a restart or a logout from the phone. You usually want to reconnect here.

		case *events.LoggedOut:
			fmt.Printf("Instance %s logged out (reason: %v)\n", instanceID, v.Reason)

			delete(s.clientsPointer, instanceID)

		case *events.Message:
			// Some message received — you could do something with it here
			fmt.Printf("Instance %s received a message\n", instanceID)

		}
	}
}

func (s *WhatsmeowService) StartAllConnections() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	deviceStore, err := s.store.GetAllDevices(ctx)
	if err != nil {
		fmt.Println("Error fetching devices from store:", err)
		return
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	for _, device := range deviceStore {
		// Skip unpaired devices — they have no JID yet
		if device.ID == nil {
			continue
		}

		jid := device.ID.String()

		client := whatsmeow.NewClient(device, clientLog)
		client.AddEventHandler(s.eventHandler(jid))

		if err := client.Connect(); err != nil {
			fmt.Printf("Error reconnecting instance %s (JID %s): %v\n", jid, jid, err)
			continue
		}

		s.clientsPointer[jid] = client

		fmt.Printf("Reconnected instance %s\n", jid)
	}
}

func (s *WhatsmeowService) StartNewConnection(id string) error {
	if s.clientsPointer[id] != nil {
		return nil
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	device := s.store.NewDevice()
	client := whatsmeow.NewClient(device, clientLog)

	qrChan, _ := client.GetQRChannel(context.Background())

	if err := client.Connect(); err != nil {
		return err
	}

	go func() {
		const maxQRCodes = 10
		qrCount := 0

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				qrCount++

				if qrCount > maxQRCodes {
					fmt.Printf("Instance %s exceeded %d QR codes — deleting\n", id, maxQRCodes)
					client.Disconnect()
					delete(s.clientsPointer, id)
					return
				}

			case "success":
				// Pairing done — save the JID so we can reconnect on restart
				if client.Store.ID != nil {
					// Successfully paired, save the JID to the database
				}

			default:
				fmt.Printf("QR event for instance %s: %s\n", id, evt.Event)
			}
		}
	}()

	s.clientsPointer[id] = client
	return nil
}

type ClientInfo struct {
	ID          string
	Name        string
	IsConnected bool
}

func (s *WhatsmeowService) GetClients() map[string]*ClientInfo {
	return s.iterateClients()
}

func (s *WhatsmeowService) iterateClients() map[string]*ClientInfo {
	result := make(map[string]*ClientInfo)

	for key, client := range s.clientsPointer {
		result[key] = &ClientInfo{
			ID:          key,
			Name:        client.Store.PushName,
			IsConnected: client.Store.ID != nil || client.IsConnected(),
		}
	}

	return result
}
