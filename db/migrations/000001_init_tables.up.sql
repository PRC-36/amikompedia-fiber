-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Main table
CREATE TABLE "users" (
                        uuid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                        email VARCHAR(255) UNIQUE NOT NULL,
                        nim VARCHAR(255) UNIQUE,
                        name VARCHAR(255) NOT NULL,
                        username VARCHAR(255) UNIQUE NOT NULL,
                        bio TEXT,
                        password VARCHAR(255) NOT NULL,
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- User Registration Table
CREATE TABLE "user_registrations" (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     email VARCHAR(255) NOT NULL,
                                     nim VARCHAR(10) NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     is_verified BOOLEAN DEFAULT FALSE,
                                     email_verified_at TIMESTAMPTZ,
                                     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- OTP Table
CREATE TABLE "otps" (
                       id SERIAL PRIMARY KEY,
                       user_rid INT REFERENCES "user_registrations"(id),
                       user_id UUID REFERENCES "users"(uuid),
                       otp_value VARCHAR(6) NOT NULL,
                       is_used BOOLEAN DEFAULT FALSE,
                       expired_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       ref_code VARCHAR(16) NOT NULL
);



-- Create Post table with user_id and ref_post_id
CREATE TABLE "posts" (
                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                        content TEXT NOT NULL,
                        user_id UUID REFERENCES "users"(uuid), -- Reference to the user table
                        ref_post_id UUID, -- Self-referencing column (nullable)
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- Image table
CREATE TABLE "images" (
                         id SERIAL PRIMARY KEY,
                         user_uuid UUID REFERENCES "users"(uuid),
                         post_id UUID REFERENCES "posts"(id),
                         image_type VARCHAR(50),
                         image_url VARCHAR(255),
                         file_path VARCHAR(255),
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Likes table
CREATE TABLE "post_likes" (
                        id SERIAL PRIMARY KEY,
                        user_id UUID REFERENCES "users"(uuid),
                        post_id UUID REFERENCES "posts"(id),
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Users Follows table
CREATE TABLE "user_follows" (
                        id SERIAL PRIMARY KEY,
                        follower_id UUID REFERENCES "users"(uuid),
                        following_id UUID REFERENCES "users"(uuid),
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Surveys table
CREATE TABLE "init_surveys" (
                        id SERIAL PRIMARY KEY,
                        user_id UUID REFERENCES "users"(uuid),
                        knows_amikompedia VARCHAR(100) NOT NULL,
                        impression_description TEXT,
                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);