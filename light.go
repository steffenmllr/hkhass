package hkhass

import (
	"github.com/brutella/hc/accessory"
	"strings"
)

type HassLight struct {
	lightAccessory *accessory.Lightbulb
}

func NewHassLight(state Hkstate, client *HAClient) *HassLight {
	identifier := state.EntityID
	entityName := state.Attributes.FriendlyName
	if state.Attributes.homebridgeName != "" {
		entityName = state.Attributes.homebridgeName
	}

	if entityName == "" {
		s := strings.Split(identifier, ".")
		entityName = s[len(s)-1]
	}

	acc := accessory.NewLightbulb(accessory.Info{
		SerialNumber: identifier,
		Manufacturer: "Home Assistant",
		Model:        "Light",
		Name:         entityName,
	})

	hslight := &HassLight{
		lightAccessory: acc,
	}

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			client.CallService("light", "turn_on", ServiceData{EntityId: state.EntityID})
		} else {
			client.CallService("light", "turn_off", ServiceData{EntityId: state.EntityID})
		}
	})

	// Set Current State
	acc.Lightbulb.On.SetValue(isOn(state.State))

	// Subscribe to Topic
	client.emitter.Subscribe(identifier, hslight.onChange)

	return hslight
}

func (s *HassLight) onChange(newState string, oldState string) {
	s.lightAccessory.Lightbulb.On.SetValue(isOn(newState))
}

func (s *HassLight) GetHCAccessory() *accessory.Accessory {
	return s.lightAccessory.Accessory
}
