// Create: cmd/client/main.go
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/centrifugal/centrifuge-go"
)

type CoefficientUpdate struct {
	Type           string    `json:"type"`
	MarketID       uint      `json:"market_id"`
	EventID        uint      `json:"event_id"`
	OldCoefficient float64   `json:"old_coefficient"`
	NewCoefficient float64   `json:"new_coefficient"`
	Timestamp      time.Time `json:"timestamp"`
}

func main() {
	log.Println("🚀 Starting Centrifugo WebSocket client...")

	// Create client with anonymous connection (no token needed)
	client := centrifuge.NewProtobufClient(
		"ws://localhost:8000/connection/websocket",
		centrifuge.Config{
			// No token needed for anonymous connection
		},
	)

	// Set up connection event handlers
	client.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Printf("✅ Connected to Centrifugo!")
		log.Printf("   Client ID: %s", e.ClientID)
		//log.Printf("   Protocol: %s", e.Transport)
	})

	client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Printf("❌ Disconnected from Centrifugo")
		log.Printf("   Code: %d", e.Code)
		log.Printf("   Reason: %s", e.Reason)
	})

	client.OnError(func(e centrifuge.ErrorEvent) {
		log.Printf("❌ Connection error: %s", e.Error.Error())
	})

	// Connect to Centrifugo
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	// Wait a moment for connection to establish
	time.Sleep(1 * time.Second)

	// Subscribe to ONLY ONE channel to avoid permission issues
	channels := []string{
		"odds_updates", // Only subscribe to the global channel
	}

	var subscriptions []*centrifuge.Subscription

	// Subscribe to each channel
	for _, channelName := range channels {
		sub, err := client.NewSubscription(channelName)
		if err != nil {
			log.Printf("❌ Failed to create subscription for %s: %v", channelName, err)
			continue
		}

		// Set up subscription event handlers
		sub.OnSubscribed(func(e centrifuge.SubscribedEvent) {
			log.Printf("📡 Successfully subscribed to: %s", channelName)
		})

		sub.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
			log.Printf("📴 Unsubscribed from: %s", channelName)
		})

		sub.OnPublication(func(e centrifuge.PublicationEvent) {
			handleMessage(channelName, e.Data)
		})

		sub.OnSubscribing(func(e centrifuge.SubscribingEvent) {
			log.Printf("🔄 Subscribing to: %s", channelName)
		})

		sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
			log.Printf("❌ Subscription error for %s: %s", channelName, e.Error.Error())
		})

		// Subscribe
		err = sub.Subscribe()
		if err != nil {
			log.Printf("❌ Failed to subscribe to %s: %v", channelName, err)
			continue
		}

		subscriptions = append(subscriptions, sub)
	}

	log.Println("👂 Listening for messages... Press Ctrl+C to stop")

	// Handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-interrupt
	log.Println("🛑 Shutting down client...")

	// Unsubscribe from all channels
	for _, sub := range subscriptions {
		sub.Unsubscribe()
	}

	log.Println("✅ Client stopped")
}

func handleMessage(channel string, data []byte) {
	log.Printf("📥 [%s] Raw data: %s", channel, string(data))

	// Try to parse as coefficient update
	var coeffUpdate CoefficientUpdate
	if err := json.Unmarshal(data, &coeffUpdate); err == nil && coeffUpdate.Type == "coefficient_update" {
		printCoefficientUpdate(channel, coeffUpdate)
		return
	}

	// Try to parse as generic JSON for pretty printing
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		log.Printf("📋 [%s] Formatted data:\n%s", channel, string(prettyJSON))
		return
	}

	// Just print as string if not JSON
	log.Printf("📝 [%s] Message: %s", channel, string(data))
}

func printCoefficientUpdate(channel string, update CoefficientUpdate) {
	direction := "📈"
	if update.NewCoefficient < update.OldCoefficient {
		direction = "📉"
	}

	change := update.NewCoefficient - update.OldCoefficient
	changePercent := 0.0
	if update.OldCoefficient != 0 {
		changePercent = (change / update.OldCoefficient) * 100
	}

	log.Printf("🎯 [%s] COEFFICIENT UPDATE:", channel)
	log.Printf("   Market: %d | Event: %d", update.MarketID, update.EventID)
	log.Printf("   %.2f → %.2f %s (%.2f%% change)",
		update.OldCoefficient, update.NewCoefficient, direction, changePercent)
	log.Printf("   Time: %s", update.Timestamp.Format("15:04:05"))
}
