package http

import (
	"advertising/internal/database"
	"advertising/internal/errors"
	"advertising/internal/types"
	"encoding/json"
	"net/http"
)

// CreateAdv - метод создания нового объявления
// Аргументы:
//	name: название объявления
//	description: описание объявления
//	links: список ссылок на фотографии (первая переденная будет главной)
//	price: цена за товар в объявлении
// Возвращаемые значения:
//	id: идентификатор созданного объявления
//	status_code: код результата (200 в случае успеха)
func CreateAdv(w http.ResponseWriter, r *http.Request) {
	var (
		err          error
		adv          types.CreateRequest
		responseJSON []byte
		response     *types.CreateResponse
	)

	// Проверить метод
	if r.Method != "POST" {
		responseJSON = errors.ErrorType(http.StatusMethodNotAllowed, "Method not allowed: use POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	// Проверить валидность JSON'а
	if err = json.NewDecoder(r.Body).Decode(&adv); err != nil {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Invalid JSON format")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Проверить корректность введенных параметров
	if adv.Name == "" || adv.Description == "" || adv.Price <= 0 || len(adv.Links) == 0 || len(adv.Links) > 3 {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Wrong or missed parameters")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Создать новое объявление
	if response, err = database.CreateAdv(adv.Name, adv.Description, adv.Links, adv.Price); err != nil {
		responseJSON = errors.ErrorType(http.StatusInternalServerError, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	return
}

// GetOneAdv - метод получения конкретного объявления по его ID
// Аргументы:
//	id: уникальный идентификатор объявления.
//	fields: опциональное поле. Если оно указано [fields: true],
//		то возвращаются так же дополнительные поля: описание и все ссылки на фото.
//		Если поле fields не указано или значение false, то вышеуказанные поля не возвращаются.
// Возвращаемые значения:
//	name: название объявления
//	price: цена за объявление
//	mainlink: ссылка на главное фото
//	[description]: описание объявления
//	[alllinks]: ссылки на все фото
//	status_code: код результата (200 в случае успеха)
func GetOneAdv(w http.ResponseWriter, r *http.Request) {
	var (
		err          error
		responseJSON []byte
		response     *types.GetOneResponse
		adv          types.GetOneRequest
	)

	// Проверить метод
	if r.Method != "POST" {
		responseJSON = errors.ErrorType(http.StatusMethodNotAllowed, "Method not allowed: use POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	// Проверить валидность JSON'а
	if err = json.NewDecoder(r.Body).Decode(&adv); err != nil {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Invalid JSON format")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Проверить валидность введенных параметров
	if adv.ID <= 0 {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Wrong or missed parameters")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Достать информацию о конкретном объявлении по его ID
	if response, err = database.GetOneAdv(adv.ID, adv.Fields); err != nil {
		responseJSON = errors.ErrorType(http.StatusInternalServerError, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	return
}

// GetAllAdv - метод возвращает все объявления
// Аргументы:
//	page: пагинация. int значение, которое указывает на начало пагинации.
//		В качестве значение выступает ID объявления. Т.е., если page = 3,
//		то будут выведены объявления с ID от 3 до 12 включительно.
//	sort: сортировка. Варианты сортировки:
//		- price_asc (по возрастанию цены)
//		- price_desc (по убыванию цены)
//		- date_asc (по дате добавления)
//		- date_desc (по дате добавления в обратном порядке)
// Возвращаемые значения:
//	advertisements: список из
//		- link (ссылка на главное фото)
//		- price (цена за объявление)
//		- name (название объявления)
//	status_code: код результата (200 в случае успеха)
func GetAllAdv(w http.ResponseWriter, r *http.Request) {
	var (
		err          error
		responseJSON []byte
		response     *types.GetAllResponse
		adv          types.GetAllRequest
	)

	// Проверить метод
	if r.Method != "POST" {
		responseJSON = errors.ErrorType(http.StatusMethodNotAllowed, "Method not allowed: use POST")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(responseJSON)

		return
	}

	// Проверить валидность JSON'а
	if err = json.NewDecoder(r.Body).Decode(&adv); err != nil {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Invalid JSON format")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// Проверить валидность поля "page"
	if adv.Page <= 0 {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Wrong or missed parameters")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	sortOptions := map[string]string{
		"price_desc": "price DESC",
		"price_asc":  "price ASC",
		"date_desc":  "created_at DESC",
		"date_asc":   "created_at ASC",
	}

	// Проверить валидность поля "sort"
	if _, ok := sortOptions[adv.Sort]; !ok || adv.Sort == "" {
		responseJSON = errors.ErrorType(http.StatusBadRequest, "Wrong or missed parameters")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseJSON)

		return
	}

	// string сортировки в виде SQL
	sortOption := sortOptions[adv.Sort]

	// Достать информацию о все объявлениях
	if response, err = database.GetAllAdv(adv.Page, sortOption); err != nil {
		responseJSON = errors.ErrorType(http.StatusInternalServerError, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseJSON)

		return
	}

	responseJSON, _ = json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	return

}

// NotFound вызывается, если путь не существует
func NotFound(w http.ResponseWriter, r *http.Request) {
	responseJSON := errors.ErrorType(http.StatusNotFound, "Not found")

	w.WriteHeader(http.StatusNotFound)
	w.Write(responseJSON)

	return
}
