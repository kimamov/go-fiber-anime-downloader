package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

func getJSONDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// parse json from res.body
	var result map[string]string
	json.Unmarshal(body, &result)
	// use utfBody using goquery
	//fmt.Println(result["html"])
	resultReader := strings.NewReader(result["html"])
	doc, err := goquery.NewDocumentFromReader(resultReader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func getDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Convert the designated charset HTML to utf-8 encoded HTML.
	// `charset` being one of the charsets known by the iconv package.
	utfBody, err := iconv.NewReader(res.Body, "UTF-8", "utf-8")
	if err != nil {
		return nil, err
	}

	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// AnimeEpisodeLink is there to display links to episodes
type AnimeEpisodeLink struct {
	EpisodeID int
	LinkID    string
}

// GetData tries to return a list of anime episodes from an url
func GetData() ([]AnimeEpisodeLink, error) {

	// Get the HMTML
	animeDocument, err := getDocument("https://9anime.to/watch/tower-of-god-dub.kvjr/ojo9nqz")
	if err != nil {
		return nil, err
	}
	// get the data-id from the player div
	dataID, ok := animeDocument.Find("div#player").Attr("data-id")
	if !ok {
		return nil, nil
	}
	// use data-id to a list of episodes
	animeStreamsDocument, err := getJSONDocument(fmt.Sprintf("https://9anime.to/ajax/film/servers?id=%s", dataID))
	if err != nil {
		return nil, err
	}

	// try to get the div containing the episodes hosted on streamtape
	//var episodesMap map[int]string
	//episodesMap = make(map[int]string)
	var episodesList []AnimeEpisodeLink
	container := animeStreamsDocument.Find(`div[data-id="40"]`)
	container.Find("a").Each(func(i int, s *goquery.Selection) {
		LinkID, ok := s.Attr("data-id")
		if ok {
			episodesList = append(episodesList, AnimeEpisodeLink{
				EpisodeID: i,
				LinkID:    strings.TrimSpace(LinkID),
			})
		}
	})

	return episodesList, nil
}

// Anime is an anime we craped from 9anime search
type Anime struct {
	Title         string
	ThumbnailPath string
	URL           string
}

// FindAnime tries to return an array of animes for a certain title string
func FindAnime(title string) (string, error) {
	document, err := getDocument(fmt.Sprintf("https://9anime.to/search?keyword=%s", title))
	if err != nil {
		return "", err
	}
	animeContainer := document.Find(`div.film-list`).First()
	if animeContainer != nil {
		html, err := goquery.OuterHtml(animeContainer)
		if err != nil {
			return "", nil
		}
		return html, nil
	}
	return "<h1>nothing found :(</h1>", nil
}
