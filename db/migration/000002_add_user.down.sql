ALTER TABLE IF EXISTS "account" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE IF EXISTS "account" DROP CONSTRAINT IF EXISTS "account_owner_fkey";

-- 删除表之前需要先将引用的外键删除掉
DROP TABLE IF EXISTS "users";