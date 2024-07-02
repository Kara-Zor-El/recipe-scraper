package parser

import (
    "strings"

    "github.com/gocolly/colly"
    "recipe-scraper/shared"
)

type WanderingStarFarmhouseParser struct{}

func (p *WanderingStarFarmhouseParser) ParseRecipe(e *colly.HTMLElement) shared.Recipe {
    title := e.ChildText(".recipe-overview h2")
    ingredients := []string{}
    e.ForEach(".recipe-ingredients li span:not(.bold)", func(_ int, elem *colly.HTMLElement) {
        ingredients = append(ingredients, elem.Text)
    })
    prepTime := e.ChildText(".recipe-overview .meta-row .recipe-meta-item:contains('Prep')")
    cookTime := e.ChildText(".recipe-overview .meta-row .recipe-meta-item:contains('Cook')")
    categories := append(strings.Split(title, " "), e.ChildText(".meta-row .recipe-meta-item .fa-folder"))

    var instructions []string
    e.ForEach(".recipe-method p", func(_ int, elem *colly.HTMLElement) {
        instructions = append(instructions, elem.Text)
    })

    return shared.Recipe{
        Title:        title,
        Ingredients:  ingredients,
        PrepTime:     prepTime,
        CookingTime:  cookTime,
        Categories:   categories,
        Instructions: instructions,
        Link:         e.Request.URL.String(),
    }
}

func (p *WanderingStarFarmhouseParser) RootLink() string {
    return "https://wanderingstarfarmhouse.com"
}

func (p *WanderingStarFarmhouseParser) RecipeSelector() string {
    return "#printthis"
}

func (p *WanderingStarFarmhouseParser) HandleRecipe(e *colly.HTMLElement, recipeChannel chan shared.Recipe) {
    recipe := p.ParseRecipe(e)
		 recipeChannel <- recipe
}
