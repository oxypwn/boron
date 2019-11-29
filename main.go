package main

import (
	//"fmt"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
	"path/filepath"
	"github.com/scizorman/go-scpi"
)



func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<-c
		os.Exit(1)
	}()
	os.Setenv("BORON_COFIG", filepath.Join(os.Getenv("HOME"), ".boron.yaml"))
}


func main() {
	var deviceip,config string
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "device-ip",
			Usage:   "V4 address of SCPI device `IP`",
			EnvVars: []string{"BORON_DEVICE_IP"},
			Destination: &deviceip,
		}),
		&cli.StringFlag{
			Name:    "config",
			Value:   filepath.Join(os.Getenv("HOME"), ".boron.yaml"),
			EnvVars: []string{"BORON_CONFIG"},
			Destination: &config,
		},
	}
	app := &cli.App{
		Name:  "boron",
		Usage: "Control SCPI capable devices and forward data to Apache Kafka.",
		Commands: []*cli.Command{
			{
				Name:  "headless",
				Usage: "Forward data from device to kafka.",
				Action: func(c *cli.Context) error {
					headlessLogg(deviceip)
					return nil
				},
			},
		},
	}
	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))
	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func headlessLogg(deviceip string) func(cli *cli.Context) error {

	var seconds time.Duration
	proto := "tcp"
	t,err := scpi.NewClient(proto, deviceip, seconds)
	if err != nil {
				log.Fatal(err)
	}

	log.Println(t)
	return nil
}
