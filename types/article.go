package types

type CreateArticleDTO struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Preview   string `json:"preview"`
	Thumbnail string `json:"thumbnail"`
	Slug      string `json:"slug"`
}

type Article struct {
	Uuid      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreateArticleDTO
}
