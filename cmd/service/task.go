package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strings"
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

	reader := bufio.NewReader(os.Stdin)

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

					//Запуск команды
					for i := 0; i < len(conf.Commands); i++ {
						cmd := exec.Command("sh", "-c", conf.Commands[i])
						stdoutStderr, err := cmd.CombinedOutput()
						if err != nil {
							fmt.Println("finishhh")

							log.Fatal(err)
						}

						fmt.Printf("%s\n", stdoutStderr)
					}
				}

				fmt.Print("Enter 'q' to quit: ")
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					continue
				}
				text = strings.TrimSpace(text)
				if text == "q" {
					fmt.Println("Exiting program...")
					os.Exit(0)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("finishhh")

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
