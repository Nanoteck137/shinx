-- +goose Up
CREATE TABLE projects (
    id TEXT,

    CONSTRAINT projects_pk PRIMARY KEY(id)
);

CREATE TABLE users (
    id TEXT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY(id),
    CONSTRAINT users_username_unique UNIQUE(username)
);

CREATE TABLE project_user_links (
    id TEXT,
    user_id TEXT NOT NULL,
    project_id TEXT NOT NULL,

    CONSTRAINT project_user_links_pk PRIMARY KEY(id),

    CONSTRAINT project_user_links_user_fk FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT project_user_links_project_fk FOREIGN KEY (project_id)
        REFERENCES projects(id)
);

-- +goose Down
DROP TABLE project_user_links;
DROP TABLE users;
DROP TABLE projects;
