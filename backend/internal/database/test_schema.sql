PRAGMA foreign_keys = ON;

-- 1) Table exists
SELECT '✓ post_votes table exists' AS test_result
WHERE EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='post_votes')
UNION ALL
SELECT '✗ post_votes table missing'
WHERE NOT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='post_votes');

-- 2) Check columns and constraints by examining CREATE statement
SELECT CASE 
    WHEN sql LIKE '%user_id INTEGER NOT NULL%' THEN '✓ user_id INTEGER NOT NULL'
    ELSE '✗ user_id INTEGER NOT NULL missing'
END AS test_result
FROM sqlite_master WHERE type='table' AND name='post_votes';

SELECT CASE 
    WHEN sql LIKE '%post_id INTEGER NOT NULL%' THEN '✓ post_id INTEGER NOT NULL'
    ELSE '✗ post_id INTEGER NOT NULL missing'
END AS test_result
FROM sqlite_master WHERE type='table' AND name='post_votes';

SELECT CASE 
    WHEN sql LIKE '%vote_type INTEGER NOT NULL%' THEN '✓ vote_type INTEGER NOT NULL'
    ELSE '✗ vote_type INTEGER NOT NULL missing'
END AS test_result
FROM sqlite_master WHERE type='table' AND name='post_votes';

SELECT CASE 
    WHEN sql LIKE '%created_at DATETIME DEFAULT CURRENT_TIMESTAMP%' THEN '✓ created_at DEFAULT CURRENT_TIMESTAMP'
    ELSE '✗ created_at DEFAULT CURRENT_TIMESTAMP missing'
END AS test_result
FROM sqlite_master WHERE type='table' AND name='post_votes';

SELECT CASE 
    WHEN sql LIKE '%UNIQUE(user_id, post_id)%' THEN '✓ UNIQUE(user_id, post_id) constraint'
    ELSE '✗ UNIQUE(user_id, post_id) constraint missing'
END AS test_result
FROM sqlite_master WHERE type='table' AND name='post_votes';

-- 3) Check foreign keys
SELECT '✓ FK user_id -> users(id)' AS test_result
WHERE EXISTS (SELECT 1 FROM pragma_foreign_key_list('post_votes') WHERE "table"='users')
UNION ALL
SELECT '✗ FK user_id -> users(id) missing'
WHERE NOT EXISTS (SELECT 1 FROM pragma_foreign_key_list('post_votes') WHERE "table"='users');

SELECT '✓ FK post_id -> posts(id)' AS test_result
WHERE EXISTS (SELECT 1 FROM pragma_foreign_key_list('post_votes') WHERE "table"='posts')
UNION ALL
SELECT '✗ FK post_id -> posts(id) missing'
WHERE NOT EXISTS (SELECT 1 FROM pragma_foreign_key_list('post_votes') WHERE "table"='posts');

-- 4) Check indexes
SELECT CASE 
    WHEN EXISTS (SELECT 1 FROM sqlite_master WHERE type='index' AND name='idx_post_votes_post_id')
    THEN '✓ idx_post_votes_post_id exists'
    ELSE '✗ idx_post_votes_post_id missing'
END AS test_result;

SELECT CASE 
    WHEN EXISTS (SELECT 1 FROM sqlite_master WHERE type='index' AND name='idx_post_votes_user_id')
    THEN '✓ idx_post_votes_user_id exists'
    ELSE '✗ idx_post_votes_user_id missing'
END AS test_result;