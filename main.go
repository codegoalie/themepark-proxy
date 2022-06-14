package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	waitTimeURLTemplate = "https://api.themeparks.wiki/preview/parks/%s/waittime"
)

type Attraction struct {
	ID         string    `json:"id"`
	WaitTime   int       `json:"waitTime"`
	Status     string    `json:"status"`
	Active     bool      `json:"active"`
	LastUpdate time.Time `json:"lastUpdate"`
	Name       string    `json:"name"`
	FastPass   bool      `json:"fastPass"`
	Meta       struct {
		Type        string  `json:"type"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
		EntityID    string  `json:"entityId"`
		SingleRider bool    `json:"singleRider"`
		ReturnTime  struct {
			State       string      `json:"state"`
			ReturnEnd   interface{} `json:"returnEnd"`
			ReturnStart string      `json:"returnStart"`
		} `json:"returnTime"`
	} `json:"meta"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/:parkID/waittimes", waitTimeProxy)

	app.Listen(":3000")
}

func waitTimeProxy(c *fiber.Ctx) error {
	waitTimes, err := fetchWaitTimes(c.Params("parkID"))
	if err != nil {
		err = fmt.Errorf("failed to fetch wait times: %w", err)
		return err
	}

	waitTimeMap := map[string]Attraction{}
	for _, attraction := range waitTimes {
		waitTimeMap[attraction.ID] = attraction
	}

	return c.JSON(waitTimeMap)
}

func fetchWaitTimes(parkID string) ([]Attraction, error) {
	attractions := []Attraction{}

	httpClient := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * 5,
			}).Dial,
			TLSHandshakeTimeout: time.Second * 5,
		},
	}

	resp, err := httpClient.Get(fmt.Sprintf(waitTimeURLTemplate, parkID))
	if err != nil {
		err = fmt.Errorf("failed to issue wait time GET: %w", err)
		return attractions, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to GET wait times (%d): %s", resp.StatusCode, string(body))
		return attractions, err
	}

	err = json.Unmarshal(body, &attractions)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal wait time attractions: %w", err)
		return attractions, err
	}

	return attractions, nil
}
