-- +goose Up
-- +goose StatementBegin
CREATE TYPE difficulty_enum AS ENUM ('EASY', 'MEDIUM', 'HARD', 'NA');

CREATE TABLE IF NOT EXISTS topics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(60) NOT NULL,
  slug VARCHAR(60) NOT NULL UNIQUE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  display_order INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(50) NOT NULL,
  slug VARCHAR(120) NOT NULL UNIQUE,
  link VARCHAR(255),
  author VARCHAR(100),
  total_problems INTEGER DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  display_order INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS problems (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(120) NOT NULL UNIQUE,
  link VARCHAR(255),
  problem_number INTEGER UNIQUE,
  difficulty difficulty_enum NOT NULL DEFAULT 'NA',
  starter_code JSONB NOT NULL,
  solution_code JSONB,
  time_limit INTEGER NOT NULL DEFAULT 2000,
  memory_limit INTEGER NOT NULL DEFAULT 256,
  acceptance_rate DECIMAL(5,2),
  total_submissions INTEGER DEFAULT 0,
  successful_submissions INTEGER DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS problem_topics (
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (problem_id, topic_id)
);

CREATE TABLE IF NOT EXISTS list_problems (
  list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  position INTEGER NOT NULL,
  is_required BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (list_id, problem_id),
  CONSTRAINT unique_position_per_list UNIQUE (list_id, position)
);

CREATE INDEX idx_problem_topics_problem_id ON problem_topics(problem_id);
CREATE INDEX idx_problem_topics_topic_id ON problem_topics(topic_id);
CREATE INDEX idx_list_problems_list_id ON list_problems(list_id);
CREATE INDEX idx_list_problems_problem_id ON list_problems(problem_id);
CREATE INDEX idx_list_problems_position ON list_problems(list_id, position);

CREATE INDEX idx_topics_slug ON topics(slug);
CREATE INDEX idx_topics_active_order ON topics(is_active, display_order);

CREATE INDEX idx_lists_slug ON lists(slug);

CREATE INDEX idx_problems_slug ON problems(slug);
CREATE INDEX idx_problems_number ON problems(problem_number);
CREATE INDEX idx_problems_difficulty ON problems(difficulty);
CREATE INDEX idx_problems_active ON problems(is_active);
CREATE INDEX idx_problems_acceptance_rate ON problems(acceptance_rate DESC);

CREATE INDEX idx_problem_topics_problem_id ON problem_topics(problem_id);
CREATE INDEX idx_problem_topics_topic_id ON problem_topics(topic_id);

CREATE INDEX idx_list_problems_list_id ON list_problems(list_id);
CREATE INDEX idx_list_problems_problem_id ON list_problems(problem_id);
CREATE INDEX idx_list_problems_position ON list_problems(list_id, position);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ language 'plpgsql';

CREATE TRIGGER update_topics_updated_at BEFORE UPDATE ON topics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lists_updated_at BEFORE UPDATE ON lists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_problems_updated_at BEFORE UPDATE ON problems
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_list_problems_updated_at BEFORE UPDATE ON list_problems
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_list_problems_updated_at ON list_problems;
DROP TRIGGER IF EXISTS update_problems_updated_at ON problems;
DROP TRIGGER IF EXISTS update_lists_updated_at ON lists;
DROP TRIGGER IF EXISTS update_topics_updated_at ON topics;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_list_problems_position;
DROP INDEX IF EXISTS idx_list_problems_problem_id;
DROP INDEX IF EXISTS idx_list_problems_list_id;
DROP INDEX IF EXISTS idx_problem_topics_topic_id;
DROP INDEX IF EXISTS idx_problem_topics_problem_id;
DROP INDEX IF EXISTS idx_problems_acceptance_rate;
DROP INDEX IF EXISTS idx_problems_active;
DROP INDEX IF EXISTS idx_problems_difficulty;
DROP INDEX IF EXISTS idx_problems_number;
DROP INDEX IF EXISTS idx_problems_slug;
DROP INDEX IF EXISTS idx_lists_slug;
DROP INDEX IF EXISTS idx_topics_active_order;
DROP INDEX IF EXISTS idx_topics_slug;

DROP TABLE IF EXISTS list_problems;
DROP TABLE IF EXISTS problem_topics;

DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS lists;
DROP TABLE IF EXISTS topics;

DROP TYPE IF EXISTS difficulty_enum;
-- +goose StatementEnd


-- Arrays
-- Recursion
-- Searching
-- Sorting
-- Hashing
-- Two Pointers
-- Sliding Window
-- Prefix Sum
-- Stack
-- Queue
-- Binary Search
-- Linked List
-- Doubly Linked List
-- Trees
-- DFS
-- BFS
-- Heap
-- Backtracking
-- Tries
-- Graphs
-- 1D DP
-- 2D DP
-- Greedy
-- Intervals
-- Math & Geometry
-- Bit Manipulation
-- Miscalleneous