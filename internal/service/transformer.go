package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/iamtvk/jsontransformer/internal/repository"
)

type TranformerService struct {
	cache      *CacheLayer
	validator  *Validator
	repository repository.ScriptRepository
}

func (s *TranformerService) Transform(ctx context.Context, req *models.TranformerRequest) (models.TransformerResponse, error) {
	startTime := time.Now()
	// TODO: validate input

	// TODO: get transformation script from cache or db

	// TODO: ask Get /compile jsonata expression

	// TODO: execute transform

	// TODO: validate output

	// TODO: add metadata

	// TODO: return response

}

func (s *TranformerService) getScript(ctx context.Context, identifier string) (*models.TransformationScript, error) {
	if script, found := s.cache.GetScript(identifier); found {
		return script, nil
	}
	script, err := s.repository.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}

	s.cache.SetScript(identifier, script)
	return script, nil
}

func (s *TranformerService) getCompiledExpression()
