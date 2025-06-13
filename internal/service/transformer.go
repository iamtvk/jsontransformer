package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsonata "github.com/blues/jsonata-go"
	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/iamtvk/jsontransformer/internal/repository"
)

type TransformerService struct {
	cache      *CacheLayer
	validator  *Validator
	repository repository.ScriptRepository
}

func (s *TransformerService) Transform(ctx context.Context, req *models.TransformerRequest) (models.TransformerResponse, error) {
	startTime := time.Now()
	// TODO: validate input

	// TODO: get transformation script from cache or db

	// TODO: ask Get /compile jsonata expression, add to cache

	// TODO: execute transform

	// TODO: validate output

	// TODO: add metadata

	// TODO: return response

}

func (s *TransformerService) getScript(ctx context.Context, identifier string) (*models.TransformationScript, error) {
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

func (s *TransformerService) getCompiledExpression(script string) (*jsonata.Expr, bool, error) {
	hash := sha256.Sum256([]byte(script)) // NOTE: hashing is cheaper than compiling
	scriptHash := hex.EncodeToString(hash[:])
	if expr, found := s.cache.GetCompiledExpression(scriptHash); found {
		return expr.(*jsonata.Expr), found, nil
	}
	expr, err := jsonata.Compile(script)
	if err != nil {
		return nil, false, err
	}
	s.cache.SetCompiledExpression(scriptHash, expr)
	return expr, true, nil
}

func (s *TransformerService) executeTransform(ctx context.Context, expr *jsonata.Expr, data *json.RawMessage) (json.RawMessage, error) {
	done := make(chan struct {
		result *json.RawMessage
		err    error
	}, 1)
	go func() {
		result, err := executeJSONataTransform(expr, data)
		done <- struct {
			result *json.RawMessage
			err    error
		}{result, err}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-done:
		return *res.result, res.err
	}
}

func executeJSONataTransform(expr *jsonata.Expr, data *json.RawMessage) (*json.RawMessage, error) {
	output, err := expr.Eval(data)
	if err != nil {
		return nil, err
	}
	rawJson, ok := output.(*json.RawMessage)
	if !ok {
		return nil, errors.New("")
	}
	return rawJson, nil
}

func (s *TransformerService) CreateScript(ctx context.Context, script *models.TransformationScript) error {
	if err := s.validator.ValidateScript(script.Script); err != nil {
		return fmt.Errorf("script validation failed: %w", err)
	}
	return s.repository.Create(ctx, script)
}

func (s *TransformerService) UpdateScript(ctx context.Context, script *models.TransformationScript) error {
	if err := s.validator.ValidateScript(script.Script); err != nil {
		return fmt.Errorf("script validation failed: %w", err)
	}
	return s.repository.Update(ctx, script)
}

func (s *TransformerService)
