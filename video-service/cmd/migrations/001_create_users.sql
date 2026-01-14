DROP TABLE IF EXISTS users;
CREATE TABLE users (
   id BIGINT PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR(255) NOT NULL,
   lastname VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL UNIQUE,
   passwordHash VARCHAR(255) NOT NULL,
   passwordSalt VARCHAR(255) NOT NULL,
   createdAt VARCHAR(19) NOT NULL,
   status VARCHAR(32) NOT NULL,
   isAdmin BOOLEAN NOT NULL DEFAULT FALSE
);
INSERT INTO users (name, lastname, email, passwordHash, passwordSalt, createdAt, status, isAdmin)
VALUES ('Admin', 'Pidar',  'adminpidar@gmail.com', 'd2ed23e23d', '1dw2dwe23dw', NOW(), 'active', TRUE);
-- INSERT INTO users (name, lastname, email, passwordHash, passwordSalt, createdAt, status, isAdmin)
-- VALUES ('Jeck', 'Sparrow', 'jacksparrow@gmail.com', '234fe3f', '2ed2ed', NOW(), 'active', TRUE);
-- INSERT INTO users (name, lastname, email, passwordHash, passwordSalt, createdAt, status)
-- VALUES ('Retard', 'Bebrick', 'retardbebrick@gmail.com', '234fe3f', '2ed2ed', NOW(), 'active', FALSE);
