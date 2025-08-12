-- +goose Up
-- +goose StatementBegin
INSERT INTO topics (name, slug)
VALUES
('Arrays', 'arrays'),
('Recursion', 'recursion'),
('Searching', 'searching'),
('Sorting', 'sorting'),
('Hashing', 'hashing'),
('Two Pointers', 'two-pointers'),
('Sliding Window', 'sliding-window'),
('Prefix Sum', 'prefix-sum'),
('Stack', 'stack'),
('Queue', 'queue'),
('Binary Search', 'binary-search'),
('Linked List', 'linked-list'),
('Doubly Linked List', 'doubly-linked-list'),
('Trees', 'trees'),
('DFS', 'dfs'),
('BFS', 'bfs'),
('Heap', 'heap'),
('Backtracking', 'backtracking'),
('Tries', 'tries'),
('Graphs', 'graphs'),
('1D DP', '1d-dp'),
('2D DP', '2d-dp'),
('Greedy', 'greedy'),
('Intervals', 'intervals'),
('Math & Geometry', 'math-geometry'),
('Bit Manipulation', 'bit-manipulation'),
('Miscellaneous', 'miscellaneous');

INSERT INTO lists (name, slug)
VALUES
('Neetcode 150', 'neetcode-150'),
('Striver A2Z DSA', 'striver-a2z'),
('Blind 75', 'blind-75');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM lists
WHERE slug IN (
  'neetcode-150',
  'striver-a2z',
  'blind-75'
);

DELETE FROM topics
WHERE slug IN (
  'arrays', 'recursion', 'searching', 'sorting', 'hashing',
  'two-pointers', 'sliding-window', 'prefix-sum', 'stack', 'queue',
  'binary-search', 'linked-list', 'doubly-linked-list', 'trees', 'dfs', 'bfs',
  'heap', 'backtracking', 'tries', 'graphs', '1d-dp', '2d-dp', 'greedy', 'intervals',
  'math-geometry', 'bit-manipulation', 'miscellaneous'
);
-- +goose StatementEnd
