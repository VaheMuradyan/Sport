package services

import (
	"context"
	"github.com/VaheMuradyan/Sport/proto"
	"log"
)

type GRPCCoefficientServer struct {
	proto.UnimplementedCoefficientServiceServer
	coefficientService *CoefficientService
	centrifugoService  *CentrifugoService
}

func NewGRPCCoefficientServer(coefficientService *CoefficientService, centrifugoService *CentrifugoService) *GRPCCoefficientServer {
	return &GRPCCoefficientServer{
		coefficientService: coefficientService,
		centrifugoService:  centrifugoService,
	}
}

func (s *GRPCCoefficientServer) UpdateCoefficient(ctx context.Context, req *proto.UpdateCoefficientRequest) (*proto.UpdateCoefficientResponse, error) {
	response, err := s.coefficientService.UpdateMarketCoefficient(
		uint(req.MarketId),
		req.NewCoefficient,
		uint(req.UserId),
	)
	if err != nil {
		return &proto.UpdateCoefficientResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	market, err := s.coefficientService.GetMarketWithHistory(uint(req.MarketId))
	if err != nil {
		return &proto.UpdateCoefficientResponse{
			Success: false,
			Message: "Failed to get market details",
		}, nil
	}

	err = s.centrifugoService.PublishCoefficientUpdate(
		market.EventID,
		uint(req.MarketId),
		response.OldCoefficient,
		response.NewCoefficient,
	)
	if err != nil {
		log.Println("cant publish coefficient !!!!!!!!!!!!!!!!!!!!!!!")
	}

	return &proto.UpdateCoefficientResponse{
		Success:        response.Success,
		Message:        response.Message,
		MarketId:       req.MarketId,
		OldCoefficient: response.OldCoefficient,
		NewCoefficient: response.NewCoefficient,
		UpdatedAt:      response.UpdatedAt.Unix(),
	}, nil
}
