package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

func (g *Gateway) FetchServices() {
	log.Println("Fetching services from Eureka...")

	resp, err := g.Client.R().
		SetHeader("Accept", "application/json").
		Get(g.EurekaURL + "/eureka/apps")

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Println("Failed to fetch services:", err)
		return
	}

	var eurekaResp EurekaResponse
	err = json.Unmarshal(resp.Body(), &eurekaResp)
	if err != nil {
		log.Println("Error parsing Eureka response:", err)
		return
	}

	for _, app := range eurekaResp.Applications.Applications {
		if len(app.Instances) > 0 {
			instance := app.Instances[0]
			service := EurekaService{
				Name: app.Name,
				IP:   instance.IPAddr,
				Port: instance.Port.Value,
			}
			g.Services.Store(app.Name, service)
		}
	}
}
