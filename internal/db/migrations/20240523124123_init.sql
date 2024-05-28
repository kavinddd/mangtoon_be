-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO "user" ( username, email, password, is_active)
VALUES ('admin', 'admin@admin.com', '$2a$10$O6u/djR/sK8cknrQFm3HCu5ibtgNoOTnx9QY1g6NoJMC5oeBURjEu', true);

CREATE TABLE role (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
  name VARCHAR(255) UNIQUE NOT NULL,
  description TEXT
);

INSERT INTO role (name, description)
VALUES ('admin', 'Can access to back office managing resources'),
       ('writer', 'Can create/edit of their mangas'),
       ('reader', 'Can read their purchased mangas');

CREATE TABLE user_role (
  user_id UUID NOT NULL,
  role_id UUID NOT NULL,
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE
);

INSERT INTO user_role (user_id, role_id)
VALUES (
    (SELECT id FROM "user" WHERE username = 'admin'),
    (SELECT id FROM role WHERE name = 'admin')
);


CREATE TABLE session (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  user_agent VARCHAR NOT NULL,
  client_ip VARCHAR NOT NULL,
  refresh_token VARCHAR NOT NULL,
  is_blocked boolean NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL ,
  FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd

