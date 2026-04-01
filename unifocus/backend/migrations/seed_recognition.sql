-- 插入测试数据：认定档案与规则
-- 请先确保 migration 002 已执行

-- 1. 插入认定档案 Profile
INSERT INTO recognition_profiles (school, college, major, entry_year, is_default)
VALUES 
('Test University', 'Computer Science College', 'Computer Science', 2023, true),
('Test University', 'Economics College', 'Finance', 2023, false);

-- 2. 插入全局默认规则 (低优先级)
INSERT INTO recognition_rules (
    profile_id, enabled, priority,
    match_title_regex, 
    output_competition_level, output_grade, points_by_award, rationale_template
) VALUES (
    NULL, true, 1000,
    '.*', -- Match All
    '省级', 'B', ARRAY[60, 40, 20], 'Global Fallback Rule'
);

-- 3. 插入 CS 学院规则 (Challenge Cup = A类)
INSERT INTO recognition_rules (
    profile_id, enabled, priority,
    match_title_regex, match_name_keywords,
    output_competition_level, output_grade, points_by_award, rationale_template
) VALUES (
    (SELECT id FROM recognition_profiles WHERE major='Computer Science' LIMIT 1),
    true, 10,
    NULL, ARRAY['Challenge Cup', '挑战杯'],
    '国家级A类', 'A', ARRAY[100, 80, 60], 'CS Major Policy: Challenge Cup is Top Tier'
);

-- 4. 插入 Economics 学院规则 (Challenge Cup = B类)
INSERT INTO recognition_rules (
    profile_id, enabled, priority,
    match_title_regex, match_name_keywords,
    output_competition_level, output_grade, points_by_award, rationale_template
) VALUES (
    (SELECT id FROM recognition_profiles WHERE major='Finance' LIMIT 1),
    true, 10,
    NULL, ARRAY['Challenge Cup', '挑战杯'],
    '国家级B类', 'B', ARRAY[80, 60, 40], 'Economics Policy: Challenge Cup is Tier 2'
);

-- 5. 插入一条测试机会
INSERT INTO opportunities (
    title, type, description, source_url, source_type, default_points_value,
    competition_level
) VALUES (
    '18th Challenge Cup National Competition', 'competition', 'Description...', 'http://example.com', 'manual', 50,
    '省级'
) ON CONFLICT DO NOTHING;
