-- 001_init_schema.down.sql
-- 回滚初始化Schema

-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_opportunities_updated_at ON opportunities;
DROP TRIGGER IF EXISTS update_user_opportunities_updated_at ON user_opportunities;
DROP TRIGGER IF EXISTS update_competition_rules_updated_at ON competition_level_rules;

-- 删除函数
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除表（按依赖顺序倒序删除）
DROP TABLE IF EXISTS nlp_tasks;
DROP TABLE IF EXISTS crawl_tasks;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS user_opportunities;
DROP TABLE IF EXISTS opportunities;
DROP TABLE IF EXISTS competition_level_rules;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;

-- 删除扩展（可选，如果其他数据库也在使用则不删除）
-- DROP EXTENSION IF EXISTS "pg_trgm";
-- DROP EXTENSION IF EXISTS "uuid-ossp";
