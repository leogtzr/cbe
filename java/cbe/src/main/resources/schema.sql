DROP TABLE IF EXISTS interaction;
CREATE TABLE interaction (
  id int(11) NOT NULL AUTO_INCREMENT,
  comment varchar(1000) DEFAULT NULL,
  date date DEFAULT NULL,
  person_id int(11) NOT NULL,
  PRIMARY KEY (id,person_id),
  -- KEY fk_interaction_persona1_idx (person_id),
  -- CONSTRAINT fk_interaction_persona1 FOREIGN KEY (person_id) REFERENCES person (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

DROP TABLE IF EXISTS person;
CREATE TABLE person (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(45) NOT NULL,
  type int(11) NOT NULL,
  everydays int(11) NOT NULL,
  PRIMARY KEY (id,name),
  --KEY fk_persona_person_type_idx (type),
  --CONSTRAINT fk_persona_person_type FOREIGN KEY (type) REFERENCES person_type (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE person_type (
  id int(11) NOT NULL AUTO_INCREMENT,
  type varchar(45) DEFAULT NULL,
  PRIMARY KEY (id)
); 

