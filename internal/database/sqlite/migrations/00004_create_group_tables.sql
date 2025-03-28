create table e_group (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    owner_id NOT NULL REFERENCES e_user(id) ON DELETE CASCADE,
    created_at DATE NOT NULL,
    hashed_password BLOB NOT NULL
);

create table e_group_member (
    group_id NOT NULL REFERENCES e_group(id) ON DELETE CASCADE,
    user_id NOT NULL REFERENCES e_user(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);
