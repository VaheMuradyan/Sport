package generator

import (
	"context"
	"github.com/VaheMuradyan/Sport/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

type CoefficientGenerator struct {
	client   proto.CoefficientServiceClient
	conn     *grpc.ClientConn
	markets  []MarketConfig
	running  bool
	stopChan chan bool
}

type MarketConfig struct {
	ID             uint32
	EventID        uint32
	Name           string
	MinCoefficient float64
	MaxCoefficient float64
	Current        float64
	Volatility     float64
}

func NewCoefficientGenerator(grpcAddr string) (*CoefficientGenerator, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := proto.NewCoefficientServiceClient(conn)

	markets := []MarketConfig{
		{ID: 1, EventID: 1, Name: "Atletico Ottawa Win", MinCoefficient: 1.01, MaxCoefficient: 5.0, Current: 1.24, Volatility: 0.15},
		{ID: 2, EventID: 1, Name: "Draw", MinCoefficient: 2.0, MaxCoefficient: 8.0, Current: 4.65, Volatility: 0.20},
		{ID: 3, EventID: 1, Name: "York 9 FC Win", MinCoefficient: 5.0, MaxCoefficient: 20.0, Current: 11.1, Volatility: 0.25},
		{ID: 4, EventID: 2, Name: "Over 2.5 Goals", MinCoefficient: 1.5, MaxCoefficient: 4.0, Current: 2.1, Volatility: 0.18},
		{ID: 5, EventID: 2, Name: "Under 2.5 Goals", MinCoefficient: 1.5, MaxCoefficient: 4.0, Current: 1.8, Volatility: 0.18},
	}

	return &CoefficientGenerator{
		client:   client,
		conn:     conn,
		markets:  markets,
		stopChan: make(chan bool),
	}, nil
}

func (og *CoefficientGenerator) Start(ctx context.Context) {
	og.running = true
	ticker := time.NewTicker(time.Duration(3+rand.Intn(3)) * time.Second)
	defer ticker.Stop()

	log.Println("🎲 Coefficient generator started! Generating Coefficient every 3-5 seconds...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Coefficient Generator stopped by context")
			og.running = false
			return
		case <-og.stopChan:
			log.Println("🛑 Coefficient Generator stopped")
			og.running = false
			return
		case <-ticker.C:
			og.generateAndUpdateCoefficient()
			ticker.Reset(time.Duration(3+rand.Intn(3)) * time.Second)
		}
	}
}

func (og *CoefficientGenerator) Stop() {
	if og.running {
		og.stopChan <- true
	}
	og.conn.Close()
}

func (og *CoefficientGenerator) generateAndUpdateCoefficient() {
	market := &og.markets[rand.Intn(len(og.markets))]

	changePercent := (rand.Float64() - 0.5) * 2 * market.Volatility
	newOdds := market.Current * (1 + changePercent)

	if newOdds < market.MinCoefficient {
		newOdds = market.MinCoefficient
	}
	if newOdds > market.MaxCoefficient {
		newOdds = market.MaxCoefficient
	}

	newOdds = float64(int(newOdds*100)) / 100

	if newOdds == market.Current {
		return
	}

	oldOdds := market.Current
	market.Current = newOdds

	direction := "📈"
	if newOdds < oldOdds {
		direction = "📉"
	}

	reasons := []string{
		"Market movement",
		"Betting volume increase",
		"Team news impact",
		"Live game events",
		"Automated adjustment",
	}
	reason := reasons[rand.Intn(len(reasons))]

	log.Printf("🎯 Updating Market [%s] %s: %.2f → %.2f %s",
		market.Name, direction, oldOdds, newOdds, reason)

	req := &proto.UpdateCoefficientRequest{
		MarketId:       market.ID,
		NewCoefficient: newOdds,
		UserId:         1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := og.client.UpdateCoefficient(ctx, req)
	if err != nil {
		log.Printf("❌ Failed to update coefficient for market %d: %v", market.ID, err)
		return
	}

	if resp.Success {
		log.Printf("✅ Successfully updated market %d ocefficient: %.2f → %.2f",
			resp.MarketId, resp.OldCoefficient, resp.NewCoefficient)
	} else {
		log.Printf("❌ Failed to update market %d: %s", market.ID, resp.Message)
	}
}
