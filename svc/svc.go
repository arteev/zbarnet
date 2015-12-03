package svc

import (
	"github.com/arteev/zbarnet/logger"
	"github.com/kardianos/service"
)

type program struct {
	exit        chan struct{}
	internalrun func()
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info.Println("Running in terminal.")
	} else {
		logger.Info.Println("Running under service manager.")
	}
	p.exit = make(chan struct{})
	go p.Run()
	return nil
}

func (p *program) Run() {
	p.internalrun()
}
func (p *program) Stop(s service.Service) error {
	logger.Info.Println("service stopping!")
	close(p.exit)
	return nil
}

//New create instance of *service.Service using config,
//internalrun it a function which will be runned
func New(config *service.Config, internalrun func()) (service.Service, error) {
	prg := &program{}
	prg.internalrun = internalrun
	serviceInctance, err := service.New(prg, config)
	if err != nil {
		return nil, err
	}
	return serviceInctance, nil
}
