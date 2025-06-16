package repository

import (
	"context"

	"github.com/iamtvk/jsontransformer/internal/models"
)

type ScriptRepository interface {
	GetByIdentifier(ctx context.Context, identifier string) (models.TransformationScript, error)
	Create(ctx context.Context, script models.TransformationScript) error
	Update(ctx context.Context, script models.TransformationScript) error
	Delete(ctx context.Context, identifier string) error
	List(ctx context.Context) ([]models.TransformationScript, error)
}
