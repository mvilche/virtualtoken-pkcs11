package main

import (
	"errors"
	"flag"
)

type Flag struct {
	Start bool
	Stop  bool
	Init  bool
}

func getFlag() (Flag, error) {

	var err error

	var f Flag
	start := flag.Bool("start", false, "INICIA VIRTUAL TOKEN")
	stop := flag.Bool("stop", false, "DETIENE VIRTUAL TOKEN")
	init := flag.Bool("init", false, "INICIALIZA VIRTUAL ETOKEN Y GENERA CERTIFICADO AUTOFIRMADO")
	flag.Parse()

	f.Start = *start
	f.Stop = *stop
	f.Init = *init

	if f.Start && f.Stop {
		err = errors.New("DEBE SELECCIONAR START, STOP o INIT")
	}

	if f.Init && f.Stop {
		err = errors.New("DEBE SELECCIONAR START, STOP o INIT")
	}

	if f.Init && f.Start {
		err = errors.New("DEBE SELECCIONAR START, STOP o INIT")
	}

	if f.Init && f.Start && f.Stop {
		err = errors.New("DEBE SELECCIONAR START, STOP o INIT")
	}

	if !f.Start && !f.Stop && !f.Init {
		err = errors.New("ERROR DEBE SELECCIONAR STOP O START")
	}

	return f, err
}
