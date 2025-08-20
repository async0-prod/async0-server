-- +goose Up
-- +goose StatementBegin

DELETE FROM testcases;
DELETE FROM problem_topics;
DELETE FROM list_problems;
DELETE FROM problems;

INSERT INTO problems (id, name, slug, link, problem_number, difficulty, starter_code, acceptance_rate)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'Two Sum', 'two-sum', 'https://leetcode.com/problems/two-sum/', 1, 'EASY', 'function twoSum(nums, target) { }', 45.0),
  ('22222222-2222-2222-2222-222222222222', 'Valid Parentheses', 'valid-parentheses', 'https://leetcode.com/problems/valid-parentheses/', 2, 'EASY', 'function isValid(s) { }', 38.0),
  ('33333333-3333-3333-3333-333333333333', 'Merge Intervals', 'merge-intervals', 'https://leetcode.com/problems/merge-intervals/', 3, 'MEDIUM', 'function merge(intervals) { }', 58.0),
  ('44444444-4444-4444-4444-444444444444', 'Best Time to Buy and Sell Stock', 'best-time-to-buy-sell-stock', 'https://leetcode.com/problems/best-time-to-buy-and-sell-stock/', 4, 'EASY', 'function maxProfit(prices) { }', 53.0),
  ('55555555-5555-5555-5555-555555555555', 'Maximum Subarray', 'maximum-subarray', 'https://leetcode.com/problems/maximum-subarray/', 5, 'MEDIUM', 'function maxSubArray(nums) { }', 49.0),
  ('66666666-6666-6666-6666-666666666666', 'Product of Array Except Self', 'product-of-array-except-self', 'https://leetcode.com/problems/product-of-array-except-self/', 6, 'MEDIUM', 'function productExceptSelf(nums) { }', 61.0),
  ('77777777-7777-7777-7777-777777777777', 'Maximum Depth of Binary Tree', 'maximum-depth-binary-tree', 'https://leetcode.com/problems/maximum-depth-of-binary-tree/', 7, 'EASY', 'function maxDepth(root) { }', 74.0),
  ('88888888-8888-8888-8888-888888888888', 'Invert Binary Tree', 'invert-binary-tree', 'https://leetcode.com/problems/invert-binary-tree/', 8, 'EASY', 'function invertTree(root) { }', 79.0),
  ('99999999-9999-9999-9999-999999999999', 'Binary Search', 'binary-search', 'https://leetcode.com/problems/binary-search/', 9, 'EASY', 'function search(nums, target) { }', 54.0),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Climbing Stairs', 'climbing-stairs', 'https://leetcode.com/problems/climbing-stairs/', 10, 'EASY', 'function climbStairs(n) { }', 65.0),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Coin Change', 'coin-change', 'https://leetcode.com/problems/coin-change/', 11, 'MEDIUM', 'function coinChange(coins, amount) { }', 37.0),
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Number of Islands', 'number-of-islands', 'https://leetcode.com/problems/number-of-islands/', 12, 'MEDIUM', 'function numIslands(grid) { }', 52.0);

INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 1, TRUE FROM lists l, problems p WHERE l.slug = 'neetcode-150' AND p.slug = 'two-sum';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 2, TRUE FROM lists l, problems p WHERE l.slug = 'neetcode-150' AND p.slug = 'valid-parentheses';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 3, TRUE FROM lists l, problems p WHERE l.slug = 'neetcode-150' AND p.slug = 'merge-intervals';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 4, TRUE FROM lists l, problems p WHERE l.slug = 'neetcode-150' AND p.slug = 'maximum-subarray';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 5, TRUE FROM lists l, problems p WHERE l.slug = 'neetcode-150' AND p.slug = 'coin-change';

INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 1, TRUE FROM lists l, problems p WHERE l.slug = 'blind-75' AND p.slug = 'two-sum';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 2, TRUE FROM lists l, problems p WHERE l.slug = 'blind-75' AND p.slug = 'best-time-to-buy-sell-stock';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 3, TRUE FROM lists l, problems p WHERE l.slug = 'blind-75' AND p.slug = 'product-of-array-except-self';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 4, TRUE FROM lists l, problems p WHERE l.slug = 'blind-75' AND p.slug = 'number-of-islands';

INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 1, TRUE FROM lists l, problems p WHERE l.slug = 'striver-a2z' AND p.slug = 'binary-search';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 2, TRUE FROM lists l, problems p WHERE l.slug = 'striver-a2z' AND p.slug = 'climbing-stairs';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 3, TRUE FROM lists l, problems p WHERE l.slug = 'striver-a2z' AND p.slug = 'invert-binary-tree';
INSERT INTO list_problems (list_id, problem_id, position, is_required)
SELECT l.id, p.id, 4, TRUE FROM lists l, problems p WHERE l.slug = 'striver-a2z' AND p.slug = 'maximum-depth-binary-tree';

INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'two-sum' AND t.slug = 'arrays';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'valid-parentheses' AND t.slug = 'stack';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'merge-intervals' AND t.slug = 'intervals';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'best-time-to-buy-sell-stock' AND t.slug = 'greedy';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'maximum-subarray' AND t.slug = '1d-dp';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'product-of-array-except-self' AND t.slug = 'arrays';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'maximum-depth-binary-tree' AND t.slug = 'trees';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'invert-binary-tree' AND t.slug = 'trees';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'binary-search' AND t.slug = 'binary-search';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'climbing-stairs' AND t.slug = '1d-dp';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'coin-change' AND t.slug = '2d-dp';
INSERT INTO problem_topics (problem_id, topic_id)
SELECT p.id, t.id FROM problems p, topics t WHERE p.slug = 'number-of-islands' AND t.slug = 'dfs';

INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[2,7,11,15],9', '[0,1]', 'array,target', 1 FROM problems p WHERE p.slug = 'two-sum';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '"()[]{}"', 'true', 'string', 1 FROM problems p WHERE p.slug = 'valid-parentheses';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[[1,3],[2,6],[8,10],[15,18]]', '[[1,6],[8,10],[15,18]]', 'intervals', 1 FROM problems p WHERE p.slug = 'merge-intervals';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[7,1,5,3,6,4]', '5', 'array', 1 FROM problems p WHERE p.slug = 'best-time-to-buy-sell-stock';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[-2,1,-3,4,-1,2,1,-5,4]', '6', 'array', 1 FROM problems p WHERE p.slug = 'maximum-subarray';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[1,2,3,4]', '[24,12,8,6]', 'array', 1 FROM problems p WHERE p.slug = 'product-of-array-except-self';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[3,9,20,null,null,15,7]', '3', 'tree', 1 FROM problems p WHERE p.slug = 'maximum-depth-binary-tree';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[4,2,7,1,3,6,9]', '[4,7,2,9,6,3,1]', 'tree', 1 FROM problems p WHERE p.slug = 'invert-binary-tree';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[-1,0,3,5,9,12],9', '4', 'array,target', 1 FROM problems p WHERE p.slug = 'binary-search';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '5', '8', 'integer', 1 FROM problems p WHERE p.slug = 'climbing-stairs';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[1,2,5],11', '3', 'coins,amount', 1 FROM problems p WHERE p.slug = 'coin-change';
INSERT INTO testcases (problem_id, input, output, ui, position)
SELECT p.id, '[["1","1","0","0"],["1","1","0","0"],["0","0","1","0"],["0","0","0","1"]]', '3', 'grid', 1 FROM problems p WHERE p.slug = 'number-of-islands';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM testcases WHERE problem_id IN (SELECT id FROM problems WHERE slug IN (
  'two-sum','valid-parentheses','merge-intervals','best-time-to-buy-sell-stock','maximum-subarray',
  'product-of-array-except-self','maximum-depth-binary-tree','invert-binary-tree','binary-search',
  'climbing-stairs','coin-change','number-of-islands'
));

DELETE FROM problem_topics WHERE problem_id IN (SELECT id FROM problems WHERE slug IN (
  'two-sum','valid-parentheses','merge-intervals','best-time-to-buy-sell-stock','maximum-subarray',
  'product-of-array-except-self','maximum-depth-binary-tree','invert-binary-tree','binary-search',
  'climbing-stairs','coin-change','number-of-islands'
));

DELETE FROM list_problems WHERE problem_id IN (SELECT id FROM problems WHERE slug IN (
  'two-sum','valid-parentheses','merge-intervals','best-time-to-buy-sell-stock','maximum-subarray',
  'product-of-array-except-self','maximum-depth-binary-tree','invert-binary-tree','binary-search',
  'climbing-stairs','coin-change','number-of-islands'
));

DELETE FROM problems WHERE slug IN (
  'two-sum','valid-parentheses','merge-intervals','best-time-to-buy-sell-stock','maximum-subarray',
  'product-of-array-except-self','maximum-depth-binary-tree','invert-binary-tree','binary-search',
  'climbing-stairs','coin-change','number-of-islands'
);
-- +goose StatementEnd
