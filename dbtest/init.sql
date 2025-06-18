-- Transformation scripts table
CREATE TABLE transformation_scripts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) UNIQUE NOT NULL, -- e.g., "apparel21_to_shopify_v1"
    name VARCHAR(255) NOT NULL,
    description TEXT,
    script TEXT NOT NULL, -- JSONata transformation script
    -- source_schema JSONB, -- Optional: source data schema validation
    -- target_schema JSONB, -- Optional: target data schema validation
    -- version INTEGER NOT NULL DEFAULT 1,
    -- is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255)
);

-- Version history for auditing
-- CREATE TABLE transformation_script_versions (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     script_id UUID REFERENCES transformation_scripts(id),
--     version INTEGER NOT NULL,
--     script TEXT NOT NULL,
--     changes TEXT, -- Description of changes
--     created_at TIMESTAMP DEFAULT NOW(),
--     created_by VARCHAR(255),
--
--     UNIQUE(script_id, version)
-- );

-- -- Performance metrics
-- CREATE TABLE transformation_metrics (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     script_identifier VARCHAR(255) NOT NULL,
--     execution_time_ms INTEGER,
--     input_size_bytes INTEGER,
--     output_size_bytes INTEGER,
--     success BOOLEAN,
--     error_message TEXT,
--     executed_at TIMESTAMP DEFAULT NOW(),
--
--     INDEX idx_script_identifier (script_identifier),
--     INDEX idx_executed_at (executed_at)
-- );
