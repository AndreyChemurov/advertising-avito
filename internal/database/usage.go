package database

import (
	"advertising/internal/types"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" //
)

const (
	host     string = "db"
	port     string = "5432"
	user     string = "postgres"
	password string = "postgres"
	dbname   string = "postgres"
)

var (
	psqlinfo string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db *sql.DB
)

func openDB() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", psqlinfo)

	if err != nil {
		return db, err
	}

	return db, nil
}

// CreateAdv - метод создания объявления в БД
func CreateAdv(name string, description string, links []string, price float64) (response *types.CreateResponse, err error) {
	if db, err = openDB(); err != nil {
		return response, err
	}

	defer db.Close()

	var id int

	// Начало транзакции
	tx, err := db.Begin()

	if err != nil {
		return response, err
	}

	// Откат будет игнориться, если был коммит
	defer tx.Rollback()

	// Подготовка к добавлению в таблицу "advertisement"
	stmt, err := tx.Prepare("INSERT INTO advertisement VALUES (DEFAULT, $1, $2, $3) RETURNING id;")

	if err != nil {
		return response, err
	}

	defer stmt.Close()

	// Добавление в таблицу "advertisement" с возвращением ID
	if err = stmt.QueryRow(name, description, price).Scan(&id); err != nil {
		return response, err
	}

	// Подготовка к добавлению в таблицу "photos"
	stmt, err = tx.Prepare("INSERT INTO photos VALUES (DEFAULT, $1, $2);")

	if err != nil {
		return response, err
	}

	// Добавление в таблицу "photos"
	for _, link := range links {
		_, err = stmt.Exec(id, link)

		if err != nil {
			_ = tx.Rollback()
			return response, err
		}
	}

	// Коммит
	if err = tx.Commit(); err != nil {
		return response, err
	}

	// Выделить память для структуры ответа
	response = new(types.CreateResponse)

	// Формирование response
	response.ID = id

	return response, nil
}

// GetOneAdv - метод получения одного объявления из БД
func GetOneAdv(id int, fields bool) (response *types.GetOneResponse, err error) {
	if db, err = openDB(); err != nil {
		return response, err
	}

	defer db.Close()

	var (
		stmt   *sql.Stmt
		rows   *sql.Rows
		exists bool
	)

	// Начало транзакции
	tx, err := db.Begin()

	if err != nil {
		return response, err
	}

	// Откат будет игнориться, если был коммит
	defer tx.Rollback()

	// Проверить, существовует ли объявление с конкретным ID
	stmt, err = tx.Prepare("SELECT EXISTS(SELECT * FROM advertisement WHERE id=$1);")

	if err != nil {
		return response, err
	}

	if err = stmt.QueryRow(id).Scan(&exists); err != nil {
		// Если не удается счесть в переменную...
		return response, err
	} else if !exists {
		// ... или если объявления с таким ID не существует
		return response, errors.New("Advertisement with such ID does not exist")
	}

	// Подготовка к селекту
	if fields {
		// Если указано поле 'fields', то нужно вывести дополнительные поля описание и все ссылки на фото
		stmt, err = tx.Prepare("SELECT name, link, price, description FROM advertisement INNER JOIN photos ON (advertisement.id=$1 and adv_id=$1);")
	} else {
		// Иначе только название, цену и ссылку на главное фото
		stmt, err = tx.Prepare("SELECT name, link, price FROM advertisement INNER JOIN photos ON (advertisement.id=$1 and adv_id=$1) LIMIT 1;")
	}

	if err != nil {
		return response, err
	}

	// Выделить память для структуры ответа
	response = new(types.GetOneResponse)

	// Выборка объявления
	if fields {
		// Если вариант с доп. атрибутами, то в выборке будет несколько рядов, поэтому используется stmt.Query
		if rows, err = stmt.Query(id); err != nil {
			return response, err
		}

		defer rows.Close()

		for rows.Next() {
			var (
				name        string
				description string
				link        string
				price       float64
			)

			if err = rows.Scan(&name, &link, &price, &description); err != nil {
				_ = tx.Rollback()
				return response, err
			}

			// Формирование response
			response.AllLinks = append(response.AllLinks, link)

			response.MainLink = response.AllLinks[0]
			response.Name = name
			response.Description = description
			response.Price = price
		}

		if err = rows.Err(); err != nil {
			return response, err
		}
	} else {
		// Иначе, так как будет один ряд, - stmt.QueryRow
		var (
			name  string
			link  string
			price float64
		)

		if err = stmt.QueryRow(id).Scan(&name, &link, &price); err != nil {
			return response, err
		}

		// Формирование response
		response.MainLink = link
		response.Name = name
		response.Price = price
	}

	// Коммит
	if err = tx.Commit(); err != nil {
		return response, err
	}

	return response, nil
}

// GetAllAdv - метод получение всех объявлний в соответствии с page
func GetAllAdv(page int, sort string) (response *types.GetAllResponse, err error) {
	if db, err = openDB(); err != nil {
		return response, err
	}

	defer db.Close()

	var rows *sql.Rows

	// Начало транзакции
	tx, err := db.Begin()

	if err != nil {
		return response, err
	}

	// Откат будет игнориться, если был коммит
	defer tx.Rollback()

	// Подготовка к селекту
	SQLString := fmt.Sprintf(`
	SELECT name, link, price FROM (
		SELECT DISTINCT ON (p.adv_id) name, link, price, created_at 
		FROM advertisement a INNER JOIN photos p on (a.id=p.adv_id) 
		WHERE a.id BETWEEN %d AND %d + 9) subquery
	ORDER BY %s;
	`, page, page, sort)

	stmt, err := tx.Prepare(SQLString)

	if err != nil {
		return response, err
	}

	// Получение из выборки
	if rows, err = stmt.Query(); err != nil {
		return response, err
	}

	defer rows.Close()

	// Выделить память для структуры ответа
	response = new(types.GetAllResponse)

	for rows.Next() {
		var (
			name  string
			link  string
			price float64

			m types.GetAllMapType
		)

		if err = rows.Scan(&name, &link, &price); err != nil {
			_ = tx.Rollback()
			return response, err
		}

		m = types.GetAllMapType{
			"name":  name,
			"link":  link,
			"price": price,
		}

		// Формирование response
		response.Advertisements = append(response.Advertisements, m)
	}

	if err = rows.Err(); err != nil {
		return response, err
	}

	// Коммит
	if err = tx.Commit(); err != nil {
		return response, err
	}

	return response, nil
}
