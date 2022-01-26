package database

import (
	"database/sql"
	"log"
)

type dataBase struct {
	database *sql.DB
}

func New(db *sql.DB) *dataBase {
	return &dataBase{
		database: db,
	}
}

func (db dataBase) Fetch(id int) (string, error) {
	var code string
	err := db.database.QueryRow("SELECT code FROM ISO4217 WHERE id=?", id).Scan(&code)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("select from database: no code with id %d\n", id)
	case err != nil:
		log.Println(err)
		return "", err
	}
	log.Printf("select from database: id is %d, code is %s\n", id, code)
	return code, nil
}

func (db dataBase) FetchAll() (map[int]string, error) {
	rows, err := db.database.Query("SELECT * from ISO4217;")
	if err != nil {
		log.Println(err)
		return map[int]string{}, err
	}
	defer rows.Close()

	m := make(map[int]string)
	for rows.Next() {
		var id int
		var code string
		err := rows.Scan(&id, &code)
		if err != nil {
			log.Printf("select from database: error is %s\n", err)
			return map[int]string{}, err
		}
		m[id] = code
	}

	return m, nil
}
