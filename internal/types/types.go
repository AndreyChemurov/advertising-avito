package types

import "encoding/json"

// CreateRequest - данные, которые запрашиваютя при создании объявления
type CreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Links       []string `json:"links"`
	Price       float64  `json:"price"`
}

// GetAllRequest - данные, которые запрашиваются при выборке одного объявления
type GetAllRequest struct {
	Page int    `json:"page"`
	Sort string `json:"sort"`
}

// GetOneRequest - данные, которые запрашиваются при выборке всех объявлений
type GetOneRequest struct {
	ID     int  `json:"id"`
	Fields bool `json:"fields"`
}

// CreateResponse - данные, котоорые возвращаются после создания объявления
type CreateResponse struct {
	ID int
}

// GetOneResponse - данные, котоорые возвращаются после выборки одного объявления
type GetOneResponse struct {
	Name        string
	Price       float64
	MainLink    string
	Description string
	AllLinks    []string
}

// GetAllMapType - тип данных для ответа в GetAllResponse
type GetAllMapType map[string]interface{}

// GetAllResponse - данные, котоорые возвращаются после выборки всех объявлений
type GetAllResponse struct {
	Advertisements []GetAllMapType
}

var responseOK map[string]string = map[string]string{
	"status_code":    "200",
	"status_message": "OK",
}

// ResponseOK - json 200'го статуса
var ResponseOK, _ = json.Marshal(responseOK)
