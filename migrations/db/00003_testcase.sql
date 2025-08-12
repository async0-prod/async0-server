-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS testcases (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  input TEXT NOT NULL,
  output TEXT NOT NULL,
  position INTEGER NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS testcases;
-- +goose StatementEnd
