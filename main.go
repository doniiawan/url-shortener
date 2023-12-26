package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/doniiawan/url-shortener/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var urls model.Urls = []model.Url{
		{
			ID:     1,
			Url:    "https://google.com",
			Count:  1,
			Short:  "gg",
			Expire: time.Now().AddDate(0, 0, 1),
		},
		{
			ID:     2,
			Url:    "https://facebook.com",
			Count:  0,
			Short:  "fb",
			Expire: time.Now().AddDate(0, 0, 1),
		},
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/:url", func(c echo.Context) error {
		urls = removeExpiredUrls(urls)
		short := c.Param("url")
		for _, url := range urls {
			if url.Short == short {
				url.Count++
				return c.Redirect(http.StatusTemporaryRedirect, url.Url)
			}
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "url not found"})
	})

	e.POST("/url", func(c echo.Context) error {
		urls = removeExpiredUrls(urls)
		var modelUrl model.Url
		// check payload
		payload := c.Bind(&modelUrl)
		if payload != nil {
			return c.JSON(http.StatusBadRequest, payload)
		}

		if !checkRequiredFields(modelUrl) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "url and short are required"})
		}
		// do something with payload
		isExist := checkDoubleShort(urls, modelUrl.Short)
		if isExist {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "short url already exists"})
		}

		if !checkUrlWeb(modelUrl.Url) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "url is not valid"})
		}

		id := len(urls) + 1
		modelUrl.ID = id
		urls = append(urls, modelUrl)
		return c.JSON(http.StatusOK, urls)
	})

	e.GET("/urls", func(c echo.Context) error {
		urls = removeExpiredUrls(urls)
		filter := c.QueryParam("count")

		if filter == "1" {
			return c.JSON(http.StatusOK, sortUrlsByCount(urls, 0))
		}

		return c.JSON(http.StatusOK, sortUrlsByCount(urls, 1))
	})
	e.Logger.Fatal(e.Start(":1323"))
}

// sortUrlsByCount sorts the given URLs based on their count in ascending or descending order.
//
// urls: The URLs to be sorted.
// sorted: An integer indicating the sorting order. If sorted is 1, the URLs will be sorted in ascending order.
//
//	Otherwise, the URLs will be sorted in descending order.
//
// Returns the sorted URLs.
func sortUrlsByCount(urls model.Urls, sorted int16) model.Urls {
	if sorted == 1 {
		sort.Slice(urls, func(i, j int) bool {
			return urls[i].Count < urls[j].Count
		})
	} else {
		sort.Slice(urls, func(i, j int) bool {
			return urls[i].Count > urls[j].Count
		})
	}
	return urls
}

// checkDoubleShort checks if a given short URL already exists in a list of URLs.
//
// Parameters:
// - urls: A list of model.Urls representing the URLs to be checked.
// - short: A string representing the short URL to be checked.
//
// Returns:
// - A boolean indicating whether the given short URL already exists in the list of URLs.
func checkDoubleShort(urls model.Urls, short string) bool {
	for _, url := range urls {
		if url.Short == short {
			return true
		}
	}
	return false
}

// checkRequiredFields checks if the required fields in the given Url struct are empty.
//
// It takes a single parameter:
// - urls: the Url struct to be checked.
//
// It returns a boolean value indicating whether the required fields are empty or not.
func checkRequiredFields(urls model.Url) bool {

	if urls.Url == "" || urls.Short == "" {
		return false
	}
	return true
}

// checkUrlWeb checks if a given URL is accessible on the web.
//
// It takes a string parameter `url` which represents the URL to be checked.
// It returns a boolean value indicating whether the URL is accessible or not.
func checkUrlWeb(url string) bool {
	_, err := http.Get(url)

	return err == nil
}

func removeExpiredUrls(urls model.Urls) model.Urls {
	var result model.Urls
	for _, url := range urls {
		if url.Expire.After(time.Now()) {
			result = append(result, url)
		}
	}
	return result
}
