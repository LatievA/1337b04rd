CREATE TABLE user_sessions {
    id SERIAL PRIMARY KEY,
    session_token TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '1 week'
};

CREATE TABLE posts {
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES user_sessions(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    archived_at TIMESTAMP DEFAULT NOW() + INTERVAL '15 minute',
    is_archived BOOLEAN DEFAULT FALSE
};

CREATE TABLE comments {
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES user_sessions(id) ON DELETE CASCADE,
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    parent_comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
};