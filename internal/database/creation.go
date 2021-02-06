package database

import "database/sql"

// createDatabase - sql-скрипт создания новых таблиц и индексов
const createDatabase string = `
CREATE TABLE IF NOT EXISTS "advertisement" (
	"id" BIGSERIAL NOT NULL PRIMARY KEY,
	"name" VARCHAR(200) NOT NULL CHECK ("name" <> ''),
	"description" VARCHAR(1000) NOT NULL CHECK ("description" <> ''),
	"price" NUMERIC(16, 2) CHECK ("price" > 0.0) NOT NULL,
	"created_at" DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS "photos" (
	"id" BIGSERIAL NOT NULL PRIMARY KEY,
	"adv_id" INT NOT NULL,
	"link" TEXT NOT NULL CHECK ("link" <> ''),
	FOREIGN KEY (adv_id) REFERENCES advertisement(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS adv_id_idx ON advertisement (id);
CREATE INDEX IF NOT EXISTS adv_price_idx ON advertisement (price);
CREATE INDEX IF NOT EXISTS adv_date_idx ON advertisement (created_at);

CREATE INDEX IF NOT EXISTS photos_adv_id_idx ON photos (adv_id);
`

// CreateTableAndIndecies - метод создания новых таблиц и индексов,
//	вызывается в самом начале в main
func CreateTableAndIndecies() (err error) {
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(createDatabase)

	if err != nil {
		return err
	}

	return nil
}
