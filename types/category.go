package types

type CreateCategoryDTO struct {
	Name string `json:"name"`
}

type Category struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
