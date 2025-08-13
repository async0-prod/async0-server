-- +goose Up
-- +goose StatementBegin
CREATE TYPE submission_status AS ENUM ('AC', 'WA', 'TLE', 'MLE', 'RE', 'CE', 'PE', 'PENDING', 'RUNNING');

CREATE TABLE IF NOT EXISTS submissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,

  code TEXT NOT NULL,
  status submission_status NOT NULL,

  runtime INTEGER,
  memory_used INTEGER,

  total_testcases INTEGER DEFAULT 0,
  passed_testcases INTEGER DEFAULT 0,
  failed_testcases INTEGER DEFAULT 0,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Fixed index names to match table name
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_submissions_user_accepted ON submissions(user_id, problem_id, status, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_submissions_user_accepted;
DROP INDEX IF EXISTS idx_submissions_created_at;
DROP INDEX IF EXISTS idx_submissions_status;
DROP INDEX IF EXISTS idx_submissions_problem_id;
DROP INDEX IF EXISTS idx_submissions_user_id;

DROP TABLE IF EXISTS submissions;
DROP TYPE IF EXISTS submission_status;
-- +goose StatementEnd