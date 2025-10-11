-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS solutions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  title VARCHAR(100) NOT NULL,
  hint TEXT NOT NULL,
  description TEXT,
  code TEXT NOT NULL,
  code_explanation TEXT,
  notes TEXT,
  time_complexity VARCHAR(50) NOT NULL,
  space_complexity VARCHAR(50) NOT NULL,
  difficulty_level VARCHAR(20) DEFAULT 'MEDIUM',
  display_order INTEGER DEFAULT 0,
  author VARCHAR(100),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_solutions_problem_id ON solutions(problem_id);
CREATE INDEX IF NOT EXISTS idx_solutions_is_active ON solutions(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_solutions_is_active;
DROP INDEX IF EXISTS idx_solutions_problem_id;

DROP TABLE IF EXISTS solutions;
-- +goose StatementEnd
