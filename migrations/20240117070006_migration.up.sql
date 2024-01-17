CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS persons (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  surname VARCHAR(80) NOT NULL,
  patronymic VARCHAR(80),
  age INT NOT NULL,
  gender gender_type,
  nationalize VARCHAR(80) NOT NULL
);