package main

import (
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Path     string   `yaml:"path"`
	Commands []string `yaml:"commands"`
}

func main() {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var conf Config

	if err := yaml.Unmarshal(f, &conf); err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%+v\n", conf)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(conf.Path)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})

}
