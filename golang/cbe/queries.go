package main

func getInteractionsPerType(personType int) ([]Interaction, error) {
	interactions := []Interaction{}

	stmt, err := db.Query(`
	select concat(p.name, ' (', pt.type, ')') person,inter.date, inter.comment, inter.id
	FROM interaction inter
	inner join person p on p.id = inter.person_id
	inner join person_type pt on pt.id = p.type
	where p.type = ?`,
		personType)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var person, date, comment, id string
	for stmt.Next() {
		stmt.Scan(&person, &date, &comment, &id)
		interactions = append(interactions, Interaction{Person: person, Date: date, Comment: comment, ID: id})
	}

	return interactions, nil
}

func getInteractions(interactionID int) (Interaction, error) {
	interaction := Interaction{}

	stmt, err := db.Query(`
	select concat(p.name, ' (', pt.type, ')') person,inter.date, inter.comment, inter.id
	FROM interaction inter
	inner join person p on p.id = inter.person_id
	inner join person_type pt on pt.id = p.type
	where inter.id = ?`, interactionID)
	if err != nil {
		return Interaction{}, err
	}
	defer stmt.Close()

	for stmt.Next() {
		stmt.Scan(&interaction.Person, &interaction.Date, &interaction.Comment, &interaction.ID)
	}

	return interaction, nil
}

func personsPerType(personType int) ([]Person, error) {
	persons := []Person{}

	stmt, err := db.Query(`
	select distinct(p.name), pt.type, p.id
	FROM person p
	inner join person_type pt on pt.id = p.type
	where p.type = ?`,
		personType)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var name, pt, id string
	for stmt.Next() {
		stmt.Scan(&name, &pt, &id)
		persons = append(persons, Person{Name: name, Type: pt, ID: id})
	}

	return persons, nil
}

func personInfo(id int) (PersonInfo, error) {
	info := PersonInfo{}

	stmt, err := db.Query(`
	select p.id, p.name, pt.id, pt.type from person p
	inner join person_type pt on pt.id = p.type
	where p.id = ?`,
		id)
	if err != nil {
		return PersonInfo{}, err
	}
	defer stmt.Close()

	for stmt.Next() {
		stmt.Scan(&info.ID, &info.Name, &info.Type, &info.TypeName)
	}

	return info, nil
}

func personTypes() ([]PersonType, error) {

	personTypes := []PersonType{}

	rows, err := db.Query("SELECT id, type FROM person_type")
	if err != nil {
		return nil, err
	}

	var tp, id string
	for rows.Next() {
		rows.Scan(&id, &tp)
		personTypes = append(personTypes, PersonType{ID: id, Type: tp})
	}

	return personTypes, nil
}
