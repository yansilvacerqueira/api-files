CREATE TABLE files (
  id SERIAL,
  owner_id INT,
  folder_id INT,
  name VARCHAR(60) NOT NULL,
  type VARCHAR(50) NOT NULL,
  path VARCHAR(250) NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp,
  updated_at TIMESTAMP NOT NULL,
  deleted BOOL NOT NULL DEFAULT false,
  PRIMARY KEY(id),
  CONSTRAINT fk_users FOREIGN KEY(owner_id) REFERENCES users(id),
  CONSTRAINT fk_folders FOREIGN KEY(folder_id) REFERENCES folders(id)
)