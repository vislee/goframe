package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/vislee/goframe/srv"
	"log"
)

func gfFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("goframed", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	return flagSet
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("[ERROR] error args")
	}
	rdMgFlag := gfFlagSet()
	rdMgFlag.Parse(os.Args[1:])

	configFile := rdMgFlag.Lookup("config").Value.String()
	if configFile == "" {
		log.Fatalln("not config files")
	}

	var cfg srv.Config
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		log.Fatalf("ERROR: failed to load config file %s - %s", configFile, err.Error())
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	mainSrv := NewMainSrv(&cfg)
	err = mainSrv.Init()
	if err != nil {
		log.Fatalf("ERROR: mainSrv init error. error: %s", err.Error())
	}
	mainSrv.Main()

	<-sigCh
	mainSrv.Stop()
}
