package main

import (
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

	app.Get("/api/getepisodes", func(c *fiber.Ctx) {
		url := c.Query("url")
		if url == "" {
			c.Send("please provide a valid 9anime url")
			return
		}
		data, err := api.GetData(url)
		if err != nil {
			c.Send("failed to get video")
		} else {
			c.Send(data)
		}
	})

	app.Get("/getepisodes", func(c *fiber.Ctx) {
		url := c.Query("url")
		if url == "" {
			c.Send("please provide a valid 9anime url")
			return
		}
		data, err := api.GetData(url)
		if err != nil {
			c.Send("failed to get video")
		} else {
			//fmt.Println(data)
			_ = c.Render("search", fiber.Map{
				"Title":         "hey there! ;)",
				"AnimeEpisodes": data,
			}, "layouts/main")
		}
	})

	app.Get("/search", func(c *fiber.Ctx) {
		key := c.Query("key")
		c.Set("Content-Type", "text/html")

		if key == "" {
			c.Send(`<h1>please provide a valid search key</h1>`)
		} else {
			animeList, err := api.FindAnime(key)
			if err != nil {
				c.Send("<h1>failed to find anything</h1>")
			} else {
				//fmt.Println(animeList)
				//c.Send(animeList)
				_ = c.Render("animeSearch", fiber.Map{
					"Title":         "hey there! ;)",
					"AnimeListHTML": animeList,
				}, "layouts/main")
			}
		}
	})

	app.Get("/getstream", func(c *fiber.Ctx) {
		videoID := c.Query("id")

		if videoID == "" {
			c.Send(`please provide a valid id`)
			return
		}
		stream, err := api.GetStream(videoID)
		if err != nil {
			c.Send(err)
			return
		}
		if stream == "" {
			c.Send("could not find video src")
			return
		}
		c.Render("stream", fiber.Map{
			"Title":     "stream or download your anime",
			"StreamSrc": stream,
		}, "layouts/main")
	})

	//app.Static("/", "./public")
	app.Get("/*", func(c *fiber.Ctx) {
		_ = c.Render("index", fiber.Map{
			"Title": "hey there! ;)",
		}, "layouts/main")
	})

	app.Listen(80)
}
