package types

type UnsplashUrls struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
	SmallS3 string `json:"small_s3"`
}

type UnsplashUser struct {
	Username string `json:"username"`
}

type UnsplashResults struct {
	Urls UnsplashUrls `json:"urls"`
	User UnsplashUser `json:"user"`
}

type UnsplashResponse struct {
	Total       int32             `json:"total"`
	Total_pages int32             `json:"total_pages"`
	Results     []UnsplashResults `json:"results"`
}
