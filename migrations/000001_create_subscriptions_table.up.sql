CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE subscriptions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               service_name TEXT NOT NULL,
                               monthly_cost INTEGER NOT NULL CHECK (monthly_cost > 0),
                               user_id UUID NOT NULL,
                               start_date DATE NOT NULL,
                               end_date DATE NULL
);

COMMENT ON TABLE subscriptions IS 'Таблица подписок пользователей';
COMMENT ON COLUMN subscriptions.service_name IS 'Название сервиса';
COMMENT ON COLUMN subscriptions.monthly_cost IS 'Месячная стоимость в рублях';
COMMENT ON COLUMN subscriptions.user_id IS 'ID пользователя (UUID)';
COMMENT ON COLUMN subscriptions.start_date IS 'Дата начала подписки';
COMMENT ON COLUMN subscriptions.end_date IS 'Дата окончания подписки (Опционально)';

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
