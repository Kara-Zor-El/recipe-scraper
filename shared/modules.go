package shared

type Config struct {
	BlackListedExtensions []string
	AllowedDomains        []string
}
type Recipe struct {
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Link         string   `json:"link"`
	PrepTime     string   `json:"prep_time"`
	CookingTime  string   `json:"cooking_time"`
	Categories   []string `json:"categories"`
	Instructions []string `json:"instructions"`
	TotalTime    string   `json:"total_time"`
}
