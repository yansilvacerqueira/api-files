CREATE TABLE users (
  id SERIAL,
  full_name VARCHAR(60) NOT NULL,
  email VARCHAR(60) NOT NULL UNIQUE,
  password VARCHAR(200) NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL,
  last_login TIMESTAMP DEFAULT current_timestamp,
  deleted BOOL NOT NULL DEFAULT false,
  PRIMARY KEY(id)
)