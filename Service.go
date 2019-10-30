package main

import (
	"reflect"
)

//Target systemd target
type Target string

const (
	//NetworkTarget networktarget
	NetworkTarget Target = "network.target"
	//MultiuserTarget multiuser target
	MultiuserTarget Target = "multi-user.target"
	//SocketTarget socket target
	SocketTarget Target = "socket.target"
)

//ServiceType type of service
type ServiceType string

const (
	//Simple simple service
	Simple ServiceType = "simple"
	//Notify tells the systemd if its initialzed
	Notify ServiceType = "notify"
	//Forking keep active if a fork is running but the parent has exited
	Forking ServiceType = "forking"
	//Dbus a dbus service
	Dbus ServiceType = "dbus"
	//Oneshot wait until the start action has finished until it consideres to be active
	Oneshot ServiceType = "oneshot"
	//Exec similar to simple
	Exec ServiceType = "exec"
)

//Restart when the service should be restarted
type Restart string

const (
	//No don't restart
	No Restart = "no"
	//Always restart always
	Always Restart = "always"
	//OnSuccess restart only on success (exitcode=0 or on SIGHUP, SIGINT, SIGTERM or on SIGPIPE)
	OnSuccess Restart = "on-success"
	//OnFailure restart only on failure (exitcode != 0)
	OnFailure Restart = "on-failure"
	//OnAbnormal restart if the service was terminated by a signal, or an operation timed out
	OnAbnormal Restart = "on-abnormal"
	//OnAbort restart if the service was terminated by an non clean exit signal
	OnAbort Restart = "on-abort"
	//OnWatchdog restart if the watchdog timed out
	OnWatchdog = "on-watchdog"
)

//SystemdBool a bool (true=yes/false=no)
type SystemdBool string

const (
	//True true
	True SystemdBool = "yes"
	//False false
	False SystemdBool = "no"
)

//Service service
type Service struct {
	Name    string
	Unit    Unit
	Service SService
	Install Install
}

//Unit [Unit] in .service file
type Unit struct {
	Description         string `name:"Description"`
	Documentation       string `name:"Documentation"`
	Before              Target `name:"Before"`
	After               Target `name:"After"`
	Wants               Target `name:"Wants"`
	ConditionPathExisis string `name:"ConditionPathExists"`
	Conflicts           string `name:"Conflicts"`
}

//SService [Service] in .service file
type SService struct {
	Type                     ServiceType `name:"Type"`
	ExecStartPre             string      `name:"ExecStartPre"`
	ExecStart                string      `name:"ExecStart"`
	ExecReload               string      `name:"ExecReload"`
	ExecStop                 string      `name:"ExecStop"`
	RestartSec               string      `name:"RestartSec"`
	User                     string      `name:"User"`
	Group                    string      `name:"Group"`
	Restart                  Restart     `name:"Restart"`
	TimeoutStartSec          int         `name:"TimeoutStartSec"`
	TimeoutStopSec           int         `name:"TimeoutStopSec"`
	SuccessExitStatus        string      `name:"SuccessExitStatus"`
	RestartPreventExitStatus string      `name:"RestartPreventExitStatus"`
	PIDFile                  string      `name:"PIDFile"`
	WorkingDirectory         string      `name:"WorkingDirectory"`
	RootDirectory            string      `name:"RootDirectory"`
	LogsDirectory            string      `name:"LogsDirectory"`
	KillMode                 string      `name:"KillMode"`
	ConditionPathExists      string      `name:"ConditionPathExists"`
	RemainAfterExit          SystemdBool `name:"RemainAfterExit"`
}

//Install [Install] in .service file
type Install struct {
	WantedBy Target `name:"WantedBy"`
	Alias    string `name:"Alias"`
}

//NewDefaultService creates a new default service
func NewDefaultService(name, description, execStart string) *Service {
	return &Service{
		Name: name,
		Unit: Unit{
			Description: description,
			After:       NetworkTarget,
		},
		Service: SService{
			Type:      Simple,
			ExecStart: execStart,
		},
		Install: Install{
			WantedBy: MultiuserTarget,
		},
	}
}

//NewService creates a new service
func NewService(unit Unit, service SService, install Install) *Service {
	return &Service{
		Unit:    unit,
		Service: service,
		Install: install,
	}
}

//Stop stop
func (service *Service) Stop() {

}

//Start starts a service
func (service *Service) Start() {

}

//Enable a service
func (service *Service) Enable() {

}

//Save saves a service to a .service file
func (service *Service) Save() {
	unit := service.Unit
	sservice := service.Service
	install := service.Install
	final := ""
	var part interface{}
	for i := 0; i < 3; i++ {
		if i == 0 {
			part = &unit
		} else if i == 1 {
			part = &sservice
		} else if i == 2 {
			part = &install
		}
		if i == 0 {
			final += "[Unit]\n"
		} else if i == 1 {
			final += "\n[Service]\n"
		} else if i == 2 {
			final += "\n[Install]\n"
		}
		v := reflect.ValueOf(part).Elem()
		for index := 0; index < v.NumField(); index++ {
			value := v.Field(index)
			fieldKey := v.Type().Field(index).Tag.Get("name")

			if len(value.String()) > 0 {
				final += fieldKey + "=" + value.String() + "\n"
			}
		}

	}
}
