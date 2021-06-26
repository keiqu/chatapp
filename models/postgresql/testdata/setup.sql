CREATE TABLE users
(
    username        varchar(50) PRIMARY KEY,
    email           varchar(255) UNIQUE       NOT NULL,
    hashed_password char(60)                  NOT NULL,
    created         timestamptz default now() NOT NULL
);

CREATE TABLE messages
(
    id       serial PRIMARY KEY,
    username varchar(50) REFERENCES users (username) NOT NULL,
    text     text                                    NOT NULL,
    created  timestamptz default now()               NOT NULL
);

CREATE INDEX idx_test_messages_created ON messages (created);

INSERT INTO users(username, email, hashed_password, created)
VALUES ('George',
        'geor@example.com',
        '$2a$12$6vzjkqafxBK8nFtvT83.ZuYKMCVAOa..lQDjySLQ6UIUo3m.2j.um',
        '2021-06-12 15:00:00+0000');

INSERT INTO messages(username, text, created)
VALUES ('George',
        'Hello World',
        '2021-06-13 15:00:00+0000'),
       ('George',
        'Very loooooooooooooooooooooooooooooooooooooooooooong message!',
        '2021-06-13 15:00:00+0000');

