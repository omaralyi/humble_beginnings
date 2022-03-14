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
func getGenresByTags(inputGenres []Genre, selectedTags []string) []Genre {
	resultsGenre := []Genre{}
	for _, genreToEvaluate := range inputGenres {
		allTagsFoundInGenre := true
		for _, selectedTag := range selectedTags {
			tagFoundInGenre := false
			for _, genreTag := range genreToEvaluate.genreTags {
				if selectedTag == genreTag {
					tagFoundInGenre = true
					break
				}
			}
			if !tagFoundInGenre {
				allTagsFoundInGenre = false
			}
		}
		if allTagsFoundInGenre {
			genreToEvaluate.selected = true
		} else {
			genreToEvaluate.selected = false
		}
		resultsGenre = append(resultsGenre, genreToEvaluate)
	}
	return resultsGenre
}
func getUniqueTagsByGenres(inputGenres []Genre) []string {
	allTags := []string{}
	for _, genre := range inputGenres {
		if genre.selected {
			allTags = append(allTags, genre.genreTags...)
		}
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
func getRootGenres() []Genre {
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
func initializeStuff() ([]Genre, []string, []string) {
	rootGenres := getRootGenres()
	availableTags := getUniqueTagsByGenres(rootGenres)
	selectedTags := []string{}

	return rootGenres, availableTags, selectedTags
}
func addSelectedTag(genres []Genre, selectedTags []string, tagToAdd string) ([]Genre, []string, []string) {
	selectedTags = append(selectedTags, tagToAdd)
	genres = getGenresByTags(genres, selectedTags)
	availableTags := getUniqueTagsByGenres(genres)
	return genres, availableTags, selectedTags
}
func removeSelectedTag(genres []Genre, selectedTags []string, tagToRemove string) ([]Genre, []string, []string) {
	for i, tagToCompare := range selectedTags {
		if tagToCompare == tagToRemove {
			selectedTags[i] = selectedTags[len(selectedTags)-1]
			selectedTags = selectedTags[:len(selectedTags)-1]
		}
	}
	genres = getGenresByTags(genres, selectedTags)
	availableTags := getUniqueTagsByGenres(genres)
	return genres, availableTags, selectedTags
}
func printTags(tags []string) {
	for _, tag := range tags {
		fmt.Print(tag, " ")
	}
	fmt.Println()
}
func printGenres(genres []Genre) {
	for _, genre := range genres {
		if genre.selected {
			fmt.Println(genre.name, ":", genre.url)
		}
	}
}

func main() {
	genres, availableTags, selectedTags := initializeStuff()
	fmt.Println("Welcome to music genre selector!")
	printTags(availableTags)

	for {
		var userInputTag string
		fmt.Println("Would you like to (add) or (remove) tags to include in your filter?")
		fmt.Scanln(&userInputTag)

		switch userInputTag {
		case "add":
			var userInputTag string
			fmt.Println("What tag would you like to add?")
			fmt.Scanln(&userInputTag)
			genres, availableTags, selectedTags = addSelectedTag(genres, selectedTags, userInputTag)
			printGenres(genres)
			fmt.Println("Your current selection:")
			printTags(selectedTags)
			fmt.Println("\nYour available options to narrow it down:")
			printTags(availableTags)
		case "remove":
			var userInputTag string
			fmt.Println("What tag would you like to remove?")
			fmt.Scanln(&userInputTag)
			genres, availableTags, selectedTags = removeSelectedTag(genres, selectedTags, userInputTag)
			printGenres(genres)
			fmt.Println("Your current selection:")
			printTags(selectedTags)
			fmt.Println("\nYour available options to narrow it down:")
			printTags(availableTags)
		}

	}
}
