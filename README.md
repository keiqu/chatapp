# ChatApp

ChatApp is a simple chat that uses WebSockets to communicate with users and stores all data in PostgreSQL.

## Usage

Build with:

```shell
go build cmd/web/*.go
```

You can configure application via command-line options:

```shell
-addr string
    Address that will be used by the server. (default ":4000")
-dsn string
    PostgreSQL connection URI. (default "postgresql://web:pass@localhost/chatapp")
-secret string
    Secret for the session manager. (default "946IpCV9y5Vlur8YvODJEhaOY8m9J1E4")
```

Database schemas:

```postgresql
CREATE TABLE users
(
    id              serial       PRIMARY KEY,
    username        varchar(255) UNIQUE NOT NULL,
    email           varchar(255) UNIQUE NOT NULL,
    hashed_password char(60)     NOT NULL,
    created         timestamptz  default now() NOT NULL
);

CREATE TABLE messages
(
    id      serial      PRIMARY KEY,
    user_id integer     REFERENCES users (id) NOT NULL,
    text    text        NOT NULL,
    created timestamptz default now() NOT NULL
);
```