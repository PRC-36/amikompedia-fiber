
-- Session table
CREATE TABLE "sessions" (
                         id UUID PRIMARY KEY,
                         user_id UUID REFERENCES "users"(uuid),
                         username VARCHAR(255) REFERENCES "users"(username),
                         refresh_token TEXT NOT NULL ,
                         user_agent VARCHAR(255) NOT NULL,
                         client_ip VARCHAR(255) NOT NULL,
                         is_blocked boolean NOT NULL DEFAULT FALSE,
                         expired_at TIMESTAMPTZ NOT NULL ,
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

);
