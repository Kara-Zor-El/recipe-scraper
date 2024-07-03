package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"log/slog"
	"recipe-scraper/parser"
	"recipe-scraper/shared"
)

var recipes []shared.Recipe

func isBlackListedExtension(link string) bool {
	for _, ext := range config.BlackListedExtensions {
		if strings.HasSuffix(link, ext) {
			return true
		}
	}
	return false
}

func isRootLink(link string, rootLink string) bool {
	return link == rootLink+"/" || link == rootLink
}

func isValidLink(link string, rootLink string) bool {
	return strings.HasPrefix(link, rootLink)
}

func handleLink(e *colly.HTMLElement, c *colly.Collector, wg *sync.WaitGroup, recipeChannel chan shared.Recipe, visitedLinks *map[string]struct{}, visitedLinksMu *sync.Mutex, parser parser.RecipeParser) {
	defer wg.Done()

	link := e.Attr("href")

	visitedLinksMu.Lock()
	defer visitedLinksMu.Unlock()

	// dont use if contains the string comment
	if strings.Contains(link, "comment") {
		return
	}

	if !isValidLink(link, parser.RootLink()) || isBlackListedExtension(link) || isRootLink(link, parser.RootLink()) || isVisitedLink(link, *visitedLinks) {
		return
	}

	(*visitedLinks)[link] = struct{}{}
	slog.Info("Link found: ", link)
	c.Visit(e.Request.AbsoluteURL(link))
}

func isVisitedLink(link string, visitedLinks map[string]struct{}) bool {
	_, found := visitedLinks[link]
	return found
}

func writeRecipesToJSON() {
	if err := os.MkdirAll("recipes", os.ModePerm); err != nil {
		slog.Error("Error creating directory:", err)
		return
	}

	for _, recipe := range recipes {
		sanitizedTitle := sanitizeFileName(recipe.Title)
		fileName := fmt.Sprintf("recipes/%s.json", sanitizedTitle)
		file, err := os.Create(fileName)
		if err != nil {
			slog.Error("Error creating file:", err)
			continue
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(recipe); err != nil {
			slog.Error("Error encoding JSON to file:", err)
		}
	}
}

func writeRecipeToJSON(recipe shared.Recipe) {
	sanitizedTitle := sanitizeFileName(recipe.Title)
	fileName := fmt.Sprintf("recipes/%s.json", sanitizedTitle)
	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(recipe); err != nil {
		slog.Error("Error encoding JSON to file:", err)
	}
}

func createRecipeDirectory() {
	if err := os.MkdirAll("recipes", os.ModePerm); err != nil {
		slog.Error("Error creating directory:", err)
		return
	}
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return reg.ReplaceAllString(name, "")
}
