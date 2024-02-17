CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS person (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  surname VARCHAR(80) NOT NULL,
  patronymic VARCHAR(80),
  age INT NOT NULL,
  gender gender_type
);

CREATE TABLE IF NOT EXISTS person_nationality (
  id SERIAL PRIMARY KEY,
  person_id INT NOT NULL,
  nationalize VARCHAR(80) NOT NULL,
  probability FLOAT NOT NULL,
  CONSTRAINT fk_person_nationality_person_id FOREIGN KEY (person_id) REFERENCES person(id)
);