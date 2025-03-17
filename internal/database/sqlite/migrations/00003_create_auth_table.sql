create table e_auth (
    user_id PRIMARY KEY REFERENCES e_user(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    hashed_password BLOB NOT NULL,
    session_id TEXT,
    session_expires_at DATE,
    UNIQUE(email)
);

