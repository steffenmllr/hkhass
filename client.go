package hkhass

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/brutella/hc/log"
	"github.com/asaskevich/EventBus"
	"net/http"
)


type Hkstate struct {
	Attributes struct {
		FriendlyName     string `json:"friendly_name"`
		homebridgeHidden bool   `json:"homebridge_hidden"`
		homebridgeName   string `json:"homebridge_name"`
		Hidden           bool   `json:"hidden"`
	} `json:"attributes"`
	EntityID string `json:"entity_id"`
	State    string `json:"state"`
}


type HAClient struct {
	sling *sling.Sling
	emitter EventBus.Bus
}

type ServiceData struct {
	EntityId string `json:"entity_id"`
}

func NewHAClient(client *http.Client, host string, port int, password string) *HAClient {
    BaseURL := fmt.Sprintf("http://%s:%d/api/", host, port)
	sle := sling.New().Client(client).Base(BaseURL)
	if password != "" {
		sle.Set("x-ha-access", password)
	}

	bus := EventBus.New()

	return &HAClient{
		sling: sle,
		emitter: bus,
	}
}

func (client *HAClient) GetState() ([]Hkstate, error) {
	items := []Hkstate{}
	_, err := client.sling.New().Get("states").ReceiveSuccess(&items)
	if len(items) == 0 {
	    log.Info.Println("Did not found any attributes please check your password")
	}

	return items, err
}

func (client *HAClient) CallService(domain string, service string, data ServiceData) error {
	url := fmt.Sprintf("services/%s/%s", domain, service)
	items := []Hkstate{}
	_, err := client.sling.New().Post(url).BodyJSON(&data).ReceiveSuccess(items)
	return err
}
