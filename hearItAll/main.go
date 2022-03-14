package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Genre struct {
	name      string
	url       string
	genreTags []string
	selected  bool
}

func trimLastChars(s string) string {
	return s[:len(s)-3]
}

func selectGenresByTag(inputGenres []Genre, searchTag string) []Genre {
	resultsGenre := []Genre{}
	for _, genreToEvaluate := range inputGenres {
		for _, tagToEvaluate := range genreToEvaluate.genreTags {
			if tagToEvaluate == searchTag {
				genreToEvaluate.selected = true
				break
			}
			genreToEvaluate.selected = false
		}
		resultsGenre = append(resultsGenre, genreToEvaluate)
	}
	return resultsGenre
}

func getUniqueTagsByGenres(inputGenres []Genre) []string {
	allTags := []string{}
	for _, genre := range inputGenres {
		allTags = append(allTags, genre.genreTags...)
	}
	return removeDuplicateStr(allTags)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func checkError(err error) {
	if err != nil {
		fmt.Print("error")
	}
}

func getGenres() []Genre {
	genres := []Genre{}
	url := "https://everynoise.com"
	response, error := http.Get(url)
	checkError(error)
	defer response.Body.Close()

	if response.StatusCode > 400 {
		fmt.Println("Status code:", response.StatusCode)
	}
	doc, error := goquery.NewDocumentFromReader(response.Body)
	checkError(error)

	doc.Find("div.canvas").Find("div.genre.scanme").Each(func(index int, item *goquery.Selection) {
		name := trimLastChars(item.Text())
		url, exists := item.Attr("preview_url")
		checkError(error)
		if exists {
			tagsToAdd := strings.Split(name, " ")
			genreTags := []string{}
			for _, tag := range tagsToAdd {
				genreTags = append(genreTags, strings.TrimSpace(tag))
			}
			genres = append(genres, Genre{name, url, genreTags, true})
		}
	})
	return genres
}

func main() {
	genres := getGenres()
	//rootTags := getUniqueTagsByGenres(rootGenres)
	fmt.Println("Discover music genres as you please! What do you want to look for? Try something random")
	var userSearchTerm string
	// Taking input from user
	fmt.Scanln(&userSearchTerm)
	genres = selectGenresByTag(genres, userSearchTerm)
	//availableTags := getUniqueTagsByGenres(searchedGenres)
	fmt.Println("Available Genres to look through:")
	for i, genre := range genres {
		if genre.selected {
			fmt.Println("Genre #", i, ":", genre.name, "URL:", genre.url, genre.selected)
		}
	}
}
