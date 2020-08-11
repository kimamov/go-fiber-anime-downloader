package main

import (
	"fmt"

	"github.com/callicoder/packer/api"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
)

// AnimeEpisodeLink is there to display links to episodes
type AnimeEpisodeLink struct {
	EpisodeID int
	LinkID    string
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})

	app.Get("/api/video", func(c *fiber.Ctx) {
		c.Send("succes")
		data, err := api.GetData()
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	app.Get("/test", func(c *fiber.Ctx) {
		c.Send("succes")
		data, err := api.GetData()
		if err != nil {
			c.Send("failed to get video")
		} else {
			fmt.Println(data)
			_ = c.Render("search", fiber.Map{
				"Title":         "hey there! ;)",
				"AnimeEpisodes": data,
			}, "layouts/main")
		}
	})

	//app.Static("/", "./public")
	app.Get("/*", func(c *fiber.Ctx) {
		_ = c.Render("index", fiber.Map{
			"Title": "hey there! ;)",
		}, "layouts/main")
	})

	app.Listen(3000)
}
