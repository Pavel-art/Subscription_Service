CREATE TABLE subscriptions (
                               id UUID PRIMARY KEY,
                               service_name TEXT NOT NULL,
                               price BIGINT NOT NULL,
                               user_id UUID NOT NULL,
                               start_date TIMESTAMP NOT NULL,
                               end_date TIMESTAMP NULL,
                               created_at TIMESTAMP NOT NULL DEFAULT now(),
                               updated_at TIMESTAMP NOT NULL DEFAULT now()
);
