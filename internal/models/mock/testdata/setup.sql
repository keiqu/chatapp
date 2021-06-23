CREATE TABLE test_users
(
    id              serial PRIMARY KEY,
    username        varchar(255) NOT NULL UNIQUE,
    email           varchar(255) NOT NULL UNIQUE,
    hashed_password char(60)     NOT NULL,
    created         timestamptz default (now() at time zone 'utc')
);

CREATE TABLE test_messages
(
    id      serial PRIMARY KEY,
    user_id integer   NOT NULL REFERENCES test_users (id),
    text    text      NOT NULL,
    created timestamp NOT NULL
);

CREATE INDEX idx_test_messages_created ON test_messages (created);
CREATE INDEX idx_test_users_email ON test_users (email);

INSERT INTO test_users(id, username, email, hashed_password, created)
VALUES (1, 'George', 'geor@example.com', '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG', '2021-06-12 15:02:15+00'),
       (2, 'Mary', 'maryme@example.com', '$2a$12$NuTjWXm3KKntReFwy1VHyufctoaHEwTy2eS206TNfkGfr6HzGJSWG', '2020-06-13 12:13:25+00');

INSERT INTO test_messages(user_id, text, created)
VALUES (1, 'Hello World', '2021-06-14 15:01:32'),
       (2, 'Very loooooooooooooooooooooooooooooooooooooooooooong message!', '2021-06-14 15:02:15');
