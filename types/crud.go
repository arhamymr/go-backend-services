package types

type Crud struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CrudDTO struct {
  Name string `json:"name"`
  Description string `json:"description"`
}
