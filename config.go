package main

type Config struct {
	BlackListedExtensions []string
	AllowedDomains        []string
}

var config = Config{
	BlackListedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	AllowedDomains:        []string{"wanderingstarfarmhouse.com", "pinchofyum.com"},
}
