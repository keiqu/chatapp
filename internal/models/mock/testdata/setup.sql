CREATE TABLE test_messages
(
    id      serial primary key,
    text    text      not null,
    created timestamp not null
);

CREATE INDEX idx_test_messages_created ON test_messages (created);

INSERT INTO test_messages(text, created)
VALUES ('Hello World', '2021-06-14 15:01:32'),
       ('Very loooooooooooooooooooooooooooooooooooooooooooong message!', '2021-06-14 15:02:15');
