CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_us_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    "Mr Test",
    "test@example.com",
    "$2a$12$aECsJPM7awGa8zFtBUxwsep6VFin1Fu0DvMGLM3WJoX3f5WiXXG1a",
    "2026-10-31 10:10:10"
);
