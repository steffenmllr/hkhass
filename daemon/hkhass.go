package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	"github.com/steffenmllr/hkhass"
	"net/http"
	"flag"
	"strings"
	"os"
	"time"
)

func createBridge() *accessory.Accessory {
	info := accessory.Info{
		Name:         "Home Assistant",
		Manufacturer: "-",
		SerialNumber: "1337",
		Model:        "E1337",
	}

	return accessory.New(info, accessory.TypeBridge)
}


func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}



func main() {

	var (
		hostArg     = flag.String("host", "127.0.0.1", "Host; default 127.0.0.1")
		portArg     = flag.Int("port", 8123, "Port; default 8123")
		pingArg     = flag.String("pin", "12344321", "Pin; default 12344321")
		passwordArg = flag.String("password", "", "Password")
		verboseArg  = flag.Bool("verbose", false, "Verbose; default false")
	)

    flag.Parse()

    if (*verboseArg == true) {
		log.Debug.Enable()
    }

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	host := *hostArg
	port := *portArg
	password := *passwordArg
	pin := *pingArg

	haClient := hkhass.NewHAClient(&httpClient, host, port, password)

	states, err := haClient.GetState()
	if err != nil {
		log.Info.Panic(err)
	}

	supportedTypes := []string{"switch"}
	accessories := []*accessory.Accessory{}

	for _, entity := range states {
		s := strings.Split(entity.EntityID, ".")
		entityType := s[0]

		// Check if it is supported
		isSupported := contains(supportedTypes, entityType)
		if isSupported == false || entity.Attributes.Hidden == true {
			continue
		}
		if entityType == "switch" {
			acc := hkhass.NewHassSwitch(entity, haClient)
			accessories = append(accessories, acc.GetHCAccessory())
		}
	}

	config := hc.Config{Pin: pin, StoragePath: "./db"}
	bridge := createBridge()
	t, err := hc.NewIPTransport(config, bridge, accessories...)
	if err != nil {
		log.Info.Panic(err)
	}

	hc.OnTermination(func() {
		t.Stop()
		os.Exit(1)
	})

	// Start Listing to events
	go hkhass.ListenToEvents(haClient, host, port, password)

	// Start Homekit
	found := len(accessories)
	log.Info.Printf("Found %d Accessories. Your Pin is: %s", found, pin)

	t.Start()

}
