package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"tinygo.org/x/bluetooth"
)

var (
	adapter    = bluetooth.DefaultAdapter
	hrSvcUUID  = bluetooth.ServiceUUIDHeartRate
	hrCharUUID = bluetooth.CharacteristicUUIDHeartRateMeasurement
)

func bleScanAndConnect(ctx context.Context, prefix string) (*bluetooth.Device, error) {
	if err := adapter.Enable(); err != nil {
		return nil, err
	}

	ch := make(chan bluetooth.ScanResult, 1)
	err := adapter.Scan(func(a *bluetooth.Adapter, res bluetooth.ScanResult) {
		log.Printf("Found device: %s (RSSI: %d, Name: %s)\n", res.Address.String(), res.RSSI, res.LocalName())
		if strings.HasPrefix(res.LocalName(), prefix) {
			a.StopScan()
			ch <- res
		}
	})
	if err != nil {
		return nil, err
	}

	select {
	case res := <-ch:
		dev, err := adapter.Connect(res.Address, bluetooth.ConnectionParams{})
		if err != nil {
			return nil, err
		}
		return &dev, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func enableHeartRateNotifications(dev *bluetooth.Device, hrChan chan<- uint8) error {
	svcs, err := dev.DiscoverServices([]bluetooth.UUID{hrSvcUUID})
	if err != nil {
		return err
	}

	chars, err := svcs[0].DiscoverCharacteristics([]bluetooth.UUID{hrCharUUID})
	if err != nil {
		return err
	}

	char := chars[0]
	return char.EnableNotifications(func(data []byte) {
		hr := uint8(data[1])
		hrChan <- hr
	})
}

func ble(ctx context.Context, hrChan chan<- uint8) {
	dev, err := bleScanAndConnect(ctx, "WHOOP")
	if err != nil {
		log.Println("Failed to connect:", err)
		return
	}
	log.Println("Connected to device:", dev.Address.String())

	err = enableHeartRateNotifications(dev, hrChan)
	if err != nil {
		log.Println("Failed to enable notifications:", err)
		return
	}

	<-ctx.Done()
}

func websocketHandler(hrChan <-chan uint8) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			log.Println("WebSocket Accept error:", err)
			return
		}
		defer conn.Close(websocket.StatusInternalError, "Internal error")

		ctx := context.Background()

		for hr := range hrChan {
			err := wsjson.Write(ctx, conn, hr)
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}

		conn.Close(websocket.StatusNormalClosure, "Done")
	}
}

func main() {
	hrChan := make(chan uint8)
	go func() {
		ctx := context.Background()
		ble(ctx, hrChan)
	}()

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", websocketHandler(hrChan))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
