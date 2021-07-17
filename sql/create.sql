package
sql

CREATE
database snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE
snippetbox;

CREATE TABLE snippets
(
    id      INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title   VARCHAR(100) NOT NULL,
    context TEXT         NOT NULL,
    created DATETIME     NOT NULL,
    expires DATETATE     NOT NULL,
);

CREATE
INDEX idx_snippets_created on snippets(created);

INSERT INTO snippets(title, context, created, expires)
VALUES ('An old silent pond', 'frog jumps into pond', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY);

);
INSERT INTO snippets(title, context, created, expires)
VALUES ('An old silent pond', 'frog jumps into pond', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY);
);
INSERT INTO snippets(title, context, created, expires)
VALUES ('An old silent pond', 'frog jumps into pond', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY);

);