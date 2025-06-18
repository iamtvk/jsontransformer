package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iamtvk/jsontransformer/internal/models"
)

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (p *PostgreSQLRepository) GetByIdentifier(ctx context.Context, identifier string) (models.TransformationScript, error) {
	query := `
        SELECT id, identifier, name, script, description, 
                created_at, updated_at, created_by
        FROM transformation_scripts 
        WHERE identifier = $1
    `
	var script models.TransformationScript
	err := p.db.QueryRowContext(ctx, query, identifier).Scan(
		&script.ID, &script.Identifier, &script.Name, &script.Script,
		&script.Description, &script.CreatedAt, &script.UpdatedAt, &script.CreatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.TransformationScript{}, fmt.Errorf("script not found: %s", identifier)
		}
		return models.TransformationScript{}, err
	}

	return script, nil
}

func (p *PostgreSQLRepository) Create(ctx context.Context, script models.TransformationScript) error {
	query := `INSERT INTO transformation_scripts
		(identifier, name, script, description, created_by)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at, updated_at
	`
	err := p.db.QueryRowContext(ctx, query, script.Identifier, script.Name, script.Script, script.Description, script.CreatedBy).Scan(&script.ID, &script.CreatedAt, &script.UpdatedAt)
	return err
}

func (p *PostgreSQLRepository) Update(ctx context.Context, script models.TransformationScript) error {
	query := `UPDATE transformation_scripts 
	SET script = $2
	WHERE identifier = $1`
	err := p.db.QueryRowContext(ctx, query, script.Identifier, script.Script).Err()
	return err
}
func (p *PostgreSQLRepository) Delete(ctx context.Context, identifier string) error {
	query := `DELETE FROM transformation_scripts 
	WHERE identifier = $1`
	err := p.db.QueryRowContext(ctx, query, identifier).Err()
	return err
}

func (p *PostgreSQLRepository) List(ctx context.Context) ([]models.TransformationScript, error) {
	query := `SELECT id, identifier, name, script, description, 
                created_at, updated_at, created_by
	FROM transformation_scripts`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scripts []models.TransformationScript
	for rows.Next() {
		var script models.TransformationScript
		err := rows.Scan(&script.ID, &script.Identifier, &script.Name, &script.Script, &script.Description, &script.CreatedAt, &script.UpdatedAt, &script.CreatedBy)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}
	return scripts, nil
}
