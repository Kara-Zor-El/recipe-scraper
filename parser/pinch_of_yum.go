
package parser

import (
    "github.com/gocolly/colly"
    "recipe-scraper/shared"
		"strings"
)

type PinchOfYumParser struct{}

func (p *PinchOfYumParser) ParseRecipe(e *colly.HTMLElement) shared.Recipe {
	// title | h2.tasty-recipes-title
	// ingredients | span inside div.tasty-recipes-ingredients
	// totalTime | span.tasty-recipes-total-time
	// categories | capital words of the title split by space
	// instructions | div.tasty-recipes-instructions
	// link | current page
	title := e.ChildText("h2.tasty-recipes-title")
	ingredients := []string{}
	e.ForEach("li[data-tr-ingredient-checkbox]", func(_ int, elem *colly.HTMLElement) {
        ingredient := elem.Text
        ingredients = append(ingredients, ingredient)
    })
	totalTime := e.ChildText("span.tasty-recipes-total-time")
	categories := strings.Split(title, " ")
	var instructions []string
	e.ForEach("div.tasty-recipes-instructions ol li", func(_ int, elem *colly.HTMLElement) {
		// append text and urls for images
		instructions = append(instructions, elem.Text)
	})

	return shared.Recipe{
		Title:        title,
		Ingredients:  ingredients,
		TotalTime:    totalTime,
		Categories:   categories,
		Instructions: instructions,
		Link:         e.Request.URL.String(),
	}
}

func (p *PinchOfYumParser) RootLink() string {
    return "https://pinchofyum.com"
}

func (p *PinchOfYumParser) RecipeSelector() string {
    return ".tasty-recipes-has-image"
}

func (p *PinchOfYumParser) HandleRecipe(e *colly.HTMLElement, recipeChannel chan shared.Recipe) {
    recipe := p.ParseRecipe(e)
		recipeChannel <- recipe
}
