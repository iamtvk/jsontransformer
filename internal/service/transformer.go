package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"

	. "github.com/iamtvk/jsontransformer/internal/config"
	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/iamtvk/jsontransformer/internal/repository"
	jsonata "github.com/xiatechs/jsonata-go"
)

type TransformerService struct {
	cache      *CacheLayer
	validator  *Validator
	repository repository.ScriptRepository
	config     *Config
}

func NewTransformerService(repo repository.ScriptRepository, config *Config, cache *CacheLayer) *TransformerService {
	return &TransformerService{
		cache:      cache,
		config:     config,
		repository: repo,
	}
}

func (s *TransformerService) Transform(ctx context.Context, req *models.TransformerRequest) (*models.TransformerResponse, error) {
	startTime := time.Now()
	script, err := s.getScript(ctx, req.ScriptIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get script by identifier:%v, err:%v", req.ScriptIdentifier, err)
	}
	// TODO: validate script & input

	expr, cacheHit, err := s.getCompiledExpression(script.Script)
	if err != nil {
		return nil, fmt.Errorf("failed to compile script to expr, err: %v", err)
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = s.config.DefaultTransformTimeout
	}
	transformCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	result, err := s.executeTransform(transformCtx, expr, req.Data)
	if err != nil {
		return nil, fmt.Errorf("error tranforming : %v", err)
	}
	executionTime := time.Since(startTime)
	// TODO: validate output
	// TODO: add metadata
	return &models.TransformerResponse{
		Result: result,
		MetaData: models.TransformerMetadata{
			ScriptIdentifier: script.Identifier,
			ExecutionTime:    executionTime,
			InputSize:        len(req.Data),
			OutputSize:       len(result),
			CacheHit:         cacheHit,
		},
	}, nil

}

func (s *TransformerService) getScript(ctx context.Context, identifier string) (models.TransformationScript, error) {
	if script, found := s.cache.GetScript(identifier); found {
		return script, nil
	}
	script, err := s.repository.GetByIdentifier(ctx, identifier)
	if err != nil {
		return models.TransformationScript{}, err
	}

	s.cache.SetScript(identifier, script)
	return script, nil
}

func (s *TransformerService) getCompiledExpression(script string) (*jsonata.Expr, bool, error) {
	hash := sha256.Sum256([]byte(script)) // NOTE: if hashing is cheaper than compiling
	scriptHash := hex.EncodeToString(hash[:])
	if expr, found := s.cache.GetCompiledExpression(scriptHash); found {
		return expr.(*jsonata.Expr), found, nil
	}
	expr := jsonata.MustCompile(script)
	// if err != nil {
	// 	log.Println("error compiling:", err.Error())
	// 	return nil, false, err
	// }
	s.cache.SetCompiledExpression(scriptHash, expr)
	return expr, false, nil
}

func (s *TransformerService) executeTransform(ctx context.Context, expr *jsonata.Expr, data json.RawMessage) (json.RawMessage, error) {
	done := make(chan struct {
		result json.RawMessage
		err    error
	}, 1)
	go func() {
		result, err := executeJSONataTransform(expr, data)
		done <- struct {
			result json.RawMessage
			err    error
		}{result, err}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-done:
		return res.result, res.err
	}
}

func executeJSONataTransform(expr *jsonata.Expr, data json.RawMessage) (json.RawMessage, error) {
	var inData any
	if err := json.Unmarshal([]byte(data), &inData); err != nil {
		log.Println("raw data: ", string(data))
		log.Println("failed to unmarshal input:", err)
		return nil, err
	}
	output, err := expr.Eval(inData)
	if err != nil {
		log.Println("executejsontranfrm error: ", err.Error())
		return nil, err
	}
	// output = output.(map[string]any)
	// log.Printf("output : %v", output)
	// log.Printf("expr in tranform : %s", expr.String())
	bytes, err := json.Marshal(output)
	raw := json.RawMessage(bytes)
	if err != nil {
		log.Println("executejsonat error rawjson conversion")
		return nil, err
	}
	return raw, nil
}

func (s *TransformerService) CreateScript(ctx context.Context, script models.TransformationScript) error {
	// if err := s.validator.ValidateScript(script.Script); err != nil {
	// 	return fmt.Errorf("script validation failed: %w", err)
	// }
	return s.repository.Create(ctx, script)
}

func (s *TransformerService) UpdateScript(ctx context.Context, script models.TransformationScript) error {
	// if err := s.validator.ValidateScript(script.Script); err != nil {
	// 	return fmt.Errorf("script validation failed: %w", err)
	// }
	return s.repository.Update(ctx, script)
}
