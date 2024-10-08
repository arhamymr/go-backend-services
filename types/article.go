package types

type CreateArticleDTO struct {
	Title      string       `json:"title" validate:"required"`
	Content    string       `json:"content" validate:"required"`
	Author     string       `json:"author" validate:"required"`
	Image      UnsplashUrls `json:"image" validate:"required"`
	Slug       string       `json:"slug" validate:"required"`
	Excerpt    string       `json:"excerpt" validate:"required"`
	CategoryId string       `json:"category_id" validate:"required"`
}

type Article struct {
	Uuid      string       `json:"uuid"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Title     string       `json:"title" `
	Content   string       `json:"content" `
	Author    string       `json:"author" `
	Image     UnsplashUrls `json:"thumbnail" `
	Slug      string       `json:"slug" `
	Excerpt   string       `json:"excerpt" `
	Category  string       `json:"category"`
}
