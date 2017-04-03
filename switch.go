package hkhass

import (
	"github.com/brutella/hc/accessory"
	"strings"
)

type HassSwitch struct {
	switchAccessory *accessory.Switch
	EntityID        string
}

func NewHassSwitch(state Hkstate, client *HAClient) *HassSwitch {
	identifier := state.EntityID
	entityName := state.Attributes.FriendlyName
	if state.Attributes.homebridgeName != "" {
		entityName = state.Attributes.homebridgeName
	}

	if entityName == "" {
		s := strings.Split(identifier, ".")
		entityName = s[len(s)-1]
	}

	acc := accessory.NewSwitch(accessory.Info{
		SerialNumber: identifier,
		Name:         entityName,
	})

	hswitch := &HassSwitch{
		EntityID:        identifier,
		switchAccessory: acc,
	}

	acc.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			client.CallService("switch", "turn_on", ServiceData{EntityId: state.EntityID})
		} else {
			client.CallService("switch", "turn_off", ServiceData{EntityId: state.EntityID})
		}
	})

	// Set Current State
	acc.Switch.On.SetValue(isOn(state.State))

	// Subscribe to Topic
	client.emitter.Subscribe(identifier, hswitch.onChange)

	return hswitch
}

func (s *HassSwitch) onChange(newState string, oldState string) {
	s.switchAccessory.Switch.On.SetValue(isOn(newState))
}

func (s *HassSwitch) GetHCAccessory() *accessory.Accessory {
	return s.switchAccessory.Accessory
}

func isOn(s string) bool {
	if s == "on" {
		return true
	} else {
		return false
	}
}
