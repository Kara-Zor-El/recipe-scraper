package main

import (
    "sync"

    "github.com/gocolly/colly"
    "log/slog"
    "recipe-scraper/parser"
    "recipe-scraper/shared"
)

type SiteConfig struct {
    Domain string
    Parser parser.RecipeParser
}

var siteConfigs = []SiteConfig{
    {
        Domain: "wanderingstarfarmhouse.com",
        Parser: &parser.WanderingStarFarmhouseParser{},
    },
		{
				Domain: "pinchofyum.com",
				Parser: &parser.PinchOfYumParser{},
		},
}

func main() {
    var wg sync.WaitGroup
    recipeChannel := make(chan shared.Recipe, 100)
    visitedLinks := make(map[string]struct{})
    visitedLinksMu := sync.Mutex{}

    for _, siteConfig := range siteConfigs {
        parser := siteConfig.Parser
        c := colly.NewCollector(colly.AllowedDomains(config.AllowedDomains...))

        c.OnRequest(func(r *colly.Request) {
            slog.Info("Visiting: ", r.URL.String())
        })

        c.OnHTML("a[href]", func(e *colly.HTMLElement) {
            wg.Add(1)
            go handleLink(e, c, &wg, recipeChannel, &visitedLinks, &visitedLinksMu, parser)
        })

        c.OnHTML(parser.RecipeSelector(), func(e *colly.HTMLElement) {
            parser.HandleRecipe(e, recipeChannel)
        })

        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Visit(parser.RootLink())
        }()
    }

    go func() {
        wg.Wait()
        close(recipeChannel)
    }()

    for recipe := range recipeChannel {
        recipes = append(recipes, recipe)
    }

    writeRecipesToJSON()
}
