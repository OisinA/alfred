package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Strum355/log"
)

// Service represents a command service that can be triggered through a command
type Service struct {
	// Name of the service
	Name string
	// The command that can be used to trigger the service
	CommandTrigger string
	// The URL of the service
	URL string
	// The port of the service
	Port string
	// Alive keeps track of if the service is replying to keep alive requests
	Alive bool
}

func (s *Service) Send(sc SendCommand) (string, error) {
	marshal, err := json.Marshal(sc)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", "http://"+s.URL+":"+s.Port+"/"+sc.Command, bytes.NewBuffer(marshal))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Response, nil
}

// Registry is an in-memory storage of the registered service
type Registry struct {
	Services []Service
}

func NewRegistry() Registry {
	return Registry{Services: make([]Service, 0)}
}

func FromFile(file string) (Registry, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return Registry{}, err
	}

	var reg Registry
	err = json.Unmarshal(f, &reg)
	if err != nil {
		return Registry{}, err
	}

	return reg, nil
}

func (r *Registry) ToFile(file string) error {
	marshal, err := json.Marshal(*r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, marshal, 0777)
}

func (r *Registry) PrintRegistry() {
	log.Info("Services:")
	for _, service := range r.Services {
		log.Info(fmt.Sprintf("\tService: %s, Trigger: %s, URL: %s, Port: %s, Alive: %t", service.Name, service.CommandTrigger, service.URL, service.Port, service.Alive))
	}
}

func (r *Registry) interceptRegister(message SendCommand) string {
	split := strings.Split(message.Args, " ")

	name := split[0]
	command := split[1]
	url := split[2]
	port := split[3]

	service := Service{Name: name, CommandTrigger: command, URL: url, Port: port, Alive: true}
	log.WithFields(log.Fields{
		"service": service,
	}).Info("Service registered.")

	r.Services = append(r.Services, service)

	err := r.ToFile("services.json")
	if err != nil {
		log.WithError(err).Error("Could not write services to file.")
	}

	return "Service registered."
}

func (r *Registry) interceptDeregister(message SendCommand) string {
	split := strings.Split(message.Args, " ")

	command := split[0]

	i := -1
	for x, s := range r.Services {
		if s.CommandTrigger == command {
			i = x
		}
	}

	if i == -1 {
		return "Command not found."
	}

	r.Services[i] = r.Services[len(r.Services)-1]
	r.Services = r.Services[:len(r.Services)-1]

	err := r.ToFile("services.json")
	if err != nil {
		log.WithError(err).Error("Could not write services to file.")
	}

	log.WithFields(log.Fields{
		"command": command,
	}).Info("Service deregistered.")

	return "Service deregistered."
}

func (r *Registry) interceptCommands(message SendCommand) string {
	msg := ""
	for _, s := range r.Services {
		msg += s.CommandTrigger + "\n"
	}

	return msg
}

// Send the command information by the command string
func (r *Registry) SendByCommand(command string, message SendCommand) (string, error) {
	if command == "register" {
		return r.interceptRegister(message), nil
	}
	if command == "deregister" {
		return r.interceptDeregister(message), nil
	}
	if command == "commands" {
		return r.interceptCommands(message), nil
	}
	var service *Service = nil
	for _, s := range r.Services {
		if s.CommandTrigger == command {
			service = &s
			break
		}
	}

	if service == nil {
		return "", errors.New(fmt.Sprintf("Command not found %s", command))
	}

	fmt.Println(*service)

	resp, err := service.Send(message)
	if err != nil {
		return "", err
	}

	return resp, nil
}
