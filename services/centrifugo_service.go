package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CentrifugoService struct {
	apiURL string
	apiKey string
	client *http.Client
}

type CentrifugoPublishRequest struct {
	Method string                `json:"method"`
	Params CentrifugoPublishData `json:"params"`
}

type CentrifugoPublishData struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type CoefficientUpdateMessage struct {
	Type           string    `json:"type"`
	MarketID       uint      `json:"market_id"`
	EventID        uint      `json:"event_id"`
	OldCoefficient float64   `json:"old_coefficient"`
	NewCoefficient float64   `json:"new_coefficient"`
	Timestamp      time.Time `json:"timestamp"`
	MarketName     string    `json:"market_name,omitempty"`
	EventName      string    `json:"event_name,omitempty"`
}

func NewCentrifugoService(apiURL, apiKey string) *CentrifugoService {
	return &CentrifugoService{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CentrifugoService) PublishCoefficientUpdate(eventID, marketID uint, oldCoeff, newCoeff float64) error {
	message := CoefficientUpdateMessage{
		Type:           "coefficient_update",
		MarketID:       marketID,
		EventID:        eventID,
		OldCoefficient: oldCoeff,
		NewCoefficient: newCoeff,
		Timestamp:      time.Now(),
	}

	channel := "odds_updates"

	if err := c.publishToChannel(channel, message); err != nil {
		log.Printf("❌ Failed to publish to channel %s: %v", channel, err)
		return fmt.Errorf("failed to publish to channel %s: %v", channel, err)
	} else {
		log.Printf("✅ Published coefficient update to channel %s: Market %d (%.2f → %.2f)",
			channel, marketID, oldCoeff, newCoeff)
	}

	return nil
}

func (c *CentrifugoService) publishToChannel(channel string, data interface{}) error {
	request := CentrifugoPublishRequest{
		Method: "publish",
		Params: CentrifugoPublishData{
			Channel: channel,
			Data:    data,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.apiURL+"/api", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *CentrifugoService) GetChannelInfo(channel string) error {
	request := map[string]interface{}{
		"method": "info",
		"params": map[string]interface{}{
			"channel": channel,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.apiURL+"/api", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	log.Printf("📊 Channel %s info: %+v", channel, result)
	return nil
}
