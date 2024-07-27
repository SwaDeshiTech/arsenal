package v1

import (
	"log"
	"net"

	"github.com/google/uuid"
)

var NodeInterface string

func GenerateUID() (string, error) {

	setInterface()

	uuid.SetClockSequence(1)
	uuid.SetNodeInterface(NodeInterface)

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}

func setInterface() error {

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces:", err)
		return err
	}

	var chosenInterface net.Interface
	for _, iface := range interfaces {
		if iface.Name == "eth0" {
			chosenInterface = iface
			break
		}
	}

	if chosenInterface.Name == "" {
		log.Println("Interface not found")
		return err
	}

	NodeInterface = chosenInterface.Name

	return nil
}
