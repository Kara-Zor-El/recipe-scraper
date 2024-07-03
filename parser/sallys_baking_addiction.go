package parser

import (
	"github.com/gocolly/colly"
	"recipe-scraper/shared"
	"strings"
)

type SallysBakingAddictionParser struct{}

func (p *SallysBakingAddictionParser) ParseRecipe(e *colly.HTMLElement) shared.Recipe {
	// title | h2.tasty-recipes-title
	// ingredients | span inside div.tasty-recipes-ingredients
	// prepTime | span.tasty-recipes-prep-time
	// cookTime | span.tasty-recipes-cook-time
	// totalTime | span.tasty-recipes-total
	// categories | capital words of the title split by space
	// instructions | div.tasty-recipes-instructions
	// link | current page
	title := e.ChildText("h2.tasty-recipes-title")
	ingredients := []string{}
	e.ForEach("li[data-tr-ingredient-checkbox]", func(_ int, elem *colly.HTMLElement) {
		ingredient := elem.Text
		ingredients = append(ingredients, ingredient)
	})
	prepTime := e.ChildText("span.tasty-recipes-prep-time")
	cookTime := e.ChildText("span.tasty-recipes-cook-time")
	totalTime := e.ChildText("span.tasty-recipes-total")
	categories := strings.Split(title, " ")
	var instructions []string
	e.ForEach("div.tasty-recipes-instructions ol li", func(_ int, elem *colly.HTMLElement) {
		// append text and urls for images
		instructions = append(instructions, elem.Text)
	})

	return shared.Recipe{
		Title:        title,
		Ingredients:  ingredients,
		PrepTime:     prepTime,
		CookingTime:     cookTime,
		TotalTime:    totalTime,
		Categories:   categories,
		Instructions: instructions,
		Link:         e.Request.URL.String(),
	}
}

func (p *SallysBakingAddictionParser) RootLink() string {
	return "https://sallysbakingaddiction.com"
}

func (p *SallysBakingAddictionParser) RecipeSelector() string {
	return ".tasty-recipes-has-image"
}

func (p *SallysBakingAddictionParser) HandleRecipe(e *colly.HTMLElement, recipeChannel chan shared.Recipe) {
	recipe := p.ParseRecipe(e)
	recipeChannel <- recipe
}
