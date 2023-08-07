-- +migrate Up
CREATE TABLE users (
  id         integer NOT NULL GENERATED ALWAYS AS IDENTITY,
  email      text    NOT NULL,
  password   text    NOT NULL,
  isVerified boolean NOT NULL
);
CREATE TABLE verifyData (
  userId integer NOT NULL,
  code   text    NOT NULL
);

-- +migrate Down
DROP TABLE users;
DROP TABLE verifyData;
