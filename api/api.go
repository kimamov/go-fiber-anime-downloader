package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

func getNineAnimeRootURL() string {
	val, ok := os.LookupEnv("NINE_ANIME_ROOT_URL")
	if ok {
		return val
	}
	return "https://www10.9anime.to/"
}

var nineAnimeRootURL string = getNineAnimeRootURL()

func getJSONResponse(url string) (map[string]string, error) {
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
	if result != nil {
		return result, nil
	}
	return nil, errors.New("failed to parse JSON")
}

func getJSONDocument(url string) (*goquery.Document, error) {
	result, err := getJSONResponse(url)
	if err != nil {
		return nil, err
	}
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
type AnimeWithEpisodes struct {
	Title    string
	Episodes []AnimeEpisodeLink
}

// GetAnimeEpisodes tries to return a list of anime episodes from an url
func GetAnimeEpisodes(animeURL string) (AnimeWithEpisodes, error) {
	anime := AnimeWithEpisodes{}
	// Get the HMTML
	animeDocument, err := getDocument(animeURL)
	if err != nil {
		return anime, err
	}
	// get the data-id from the player div
	dataID, ok := animeDocument.Find("div#player").Attr("data-id")
	if !ok || dataID == "" {
		return anime, errors.New("could not find anime data-id")
	}
	// use data-id to a list of episodes
	animeStreamsDocument, err := getJSONDocument(fmt.Sprintf("%sajax/film/servers?id=%s", nineAnimeRootURL, dataID))
	if err != nil {
		return anime, err
	}

	// try to get the div containing the episodes hosted on streamtape
	//var episodesMap map[int]string
	//episodesMap = make(map[int]string)
	//var episodesList []AnimeEpisodeLink
	container := animeStreamsDocument.Find(`div[data-id="40"]`)
	container.Find("a").Each(func(i int, s *goquery.Selection) {
		LinkID, ok := s.Attr("data-id")
		if ok {
			anime.Episodes = append(anime.Episodes, AnimeEpisodeLink{
				EpisodeID: i + 1,
				LinkID:    strings.TrimSpace(LinkID),
			})
		}
	})
	// if we found all the needed stuff successfully find the anime title too
	titleContainer := animeDocument.Find(`div.widget-title`)
	if titleContainer != nil {
		title := titleContainer.Find(`h1`).First().Text()
		if title != "" {
			anime.Title = title
		}
	}
	//fmt.Println(goquery.OuterHtml(anime))

	return anime, nil
}

// Anime is an anime we craped from 9anime search
type Anime struct {
	Title         string
	ThumbnailPath string
	URL           string
}

// FindAnime tries to return an array of animes for a certain title string
func FindAnime(title string) (string, error) {
	document, err := getDocument(fmt.Sprintf("%ssearch?keyword=%s", nineAnimeRootURL, title))
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

type AnimeStream struct {
	Link    string
	Episode string
}

func GetStream(videoID string) (AnimeStream, error) {
	// https://9anime.to/ajax/episode/info?id=550fd1cbd47a12d12729279913a9eb7040ea828c6a169fec38e2169641146d70&server=40
	stream := AnimeStream{}
	streamTapeLink, err := getJSONResponse(fmt.Sprintf("%sajax/episode/info?id=%s&server=40", nineAnimeRootURL, videoID))
	if err != nil {
		return stream, err
	}
	val, ok := streamTapeLink["target"]
	if !ok || val == "" {
		return stream, errors.New("could not find a valid link")
	}
	episode, ok := streamTapeLink["name"]
	if !ok || val == "" {
		return stream, errors.New("could not find valid episode")
	}
	// val will be the url of the streamTape iframe content
	// lets get that contet
	/* res, err := http.Get(val)
	if err != nil {
		log.Panic(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
		str, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf(string(str))
	} */
	playerDocument, err := getDocument(val)
	//fmt.Println(playerDocument.Text())
	//fmt.Println(val)
	if err != nil {
		return stream, err
	}
	videoSrc := ""
	playerDocument.Find("div").Each(func(index int, node *goquery.Selection) {
		node.RemoveAttr("hidden")
		node.RemoveAttr("style")
		//fmt.Printf(goquery.OuterHtml(node))
		videoLinkContainerID, ok := node.Attr("id")
		if ok && videoLinkContainerID == "videolink" {
			//fmt.Println(videoLinkContainerID)
			videoLink := node.Text()
			if videoLink != "" {
				videoSrc = fmt.Sprintf("https:%s", videoLink)
			}
		}
	})
	if videoSrc != "" {
		stream.Link = videoSrc
		stream.Episode = episode
		return stream, nil
	}
	return stream, errors.New("could not find video src")
}
