package parser

import (
    "github.com/gocolly/colly"
    "recipe-scraper/shared"
)

type RecipeParser interface {
    ParseRecipe(e *colly.HTMLElement) shared.Recipe
    RootLink() string
    RecipeSelector() string
    HandleRecipe(e *colly.HTMLElement, recipeChannel chan shared.Recipe)
}
