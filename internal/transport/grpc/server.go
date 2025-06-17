package grpc

import (
	"context"
	"log"
	"time"

	pb "github.com/iamtvk/jsontransformer/api/proto/transformerPb"
	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/iamtvk/jsontransformer/internal/service"
)

type Server struct {
	pb.UnimplementedTransformerServiceServer
	transformerService *service.TransformerService
}

func NewServer(service *service.TransformerService) *Server {
	return &Server{
		transformerService: service,
	}
}

func (s *Server) Transform(ctx context.Context, req *pb.TransformRequest) (*pb.TransformResponse, error) {
	log.Printf("gRPC Transform request: identifier=%s", req.GetScriptIdentifier())
	response, err := s.transformerService.Transform(ctx, &models.TransformerRequest{
		Data:             req.GetData(),
		ScriptIdentifier: req.GetScriptIdentifier(),
		Timeout:          time.Duration(req.GetTimeoutSeconds()),
	})
	metadata := pb.TransformMetadata{
		CacheHit:         response.MetaData.CacheHit,
		ExecutionTimeMs:  response.MetaData.ExecutionTime.Microseconds(),
		ScriptIdentifier: response.MetaData.ScriptIdentifier,
		ScriptVersion:    int32(response.MetaData.ScriptVersion),
		InputSize:        int32(response.MetaData.InputSize),
		OutputSize:       int32(response.MetaData.OutputSize),
	}
	if err != nil {
		log.Printf("transform unsuccessfull,error: %v", err.Error())
		return nil, err
	}
	log.Printf("	transform successfull, identifier=%s, timetook=%dms, cache_hit=%v",
		response.MetaData.ScriptIdentifier, metadata.ExecutionTimeMs, metadata.CacheHit)
	return &pb.TransformResponse{
		Result:   response.Result,
		Metadata: &metadata,
		Error:    response.Error.Message,
	}, nil
}

func (s *Server) CreateScript(ctx context.Context, req *pb.CreateScriptRequest) (*pb.CreateScriptResponse, error) {
	log.Printf("gRPC CreateScript request: identifier=%s", req.GetScriptIdentifier())
	err := s.transformerService.CreateScript(ctx, models.TransformationScript{
		Script:      req.GetScript(),
		Identifier:  req.GetScriptIdentifier(),
		Description: req.GetDescription(),
		CreatedBy:   req.GetCreatedBy(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateScriptResponse{
		Success: true,
		Error:   err.Error(),
	}, nil
}
