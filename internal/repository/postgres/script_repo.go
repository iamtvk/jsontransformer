package postgres

import (
	"context"
	"database/sql"

	"github.com/iamtvk/jsontransformer/internal/models"
)

type ScriptRepository interface {
	GetByIdentifier(ctx context.Context, identifier string) (*models.TransformationScript, error)
	Create(ctx context.Context, script *models.TransformationScript) error
	Update(ctx context.Context, script *models.TransformationScript) error
	Delete(ctx context.Context, identifier string) error
	List(ctx context.Context) []*models.TransformationScript
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

// TODO: implement methods
func (p *PostgreSQLRepository) GetByIdentifier(ctx context.Context, Identifier string) (*models.TranformerRequest, error) {
	return nil, nil
}

func (p *PostgreSQLRepository) Create(ctx context.Context, script *models.TranformerRequest) error {
	return nil
}
func (p *PostgreSQLRepository) Update(ctx context.Context, script *models.TranformerRequest) error {
	return nil
}

func (p *PostgreSQLRepository) Delete(ctx context.Context, identifier string) error {
	return nil
}

func (p *PostgreSQLRepository) List(ctx context.Context) []*models.TranformerRequest {
	return nil
}
