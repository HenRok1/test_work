//package main
//
//import (
//	"fmt"
//	"github.com/fsnotify/fsnotify"
//	"github.com/jmoiron/sqlx"
//	_ "github.com/lib/pq"
//	"gopkg.in/yaml.v3"
//	"log"
//	"os"
//	"os/exec"
//	"time"
//)
//
//type Config struct {
//	WatchedPaths []struct {
//		Path     string   `yaml:"path"`
//		Commands []string `yaml:"commands"`
//	} `yaml:"watched_paths"`
//}
//
//type FileChange struct {
//	ID          int       `db:"id"`
//	FilePath    string    `db:"file_path"`
//	Method      string    `db:"method"`
//	Time_change time.Time `db:"time_change"`
//}
//
//func main() {
//	//Чтение yaml-файла
//	f, err := os.ReadFile("config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//fmt.Println("Here")
//	//Подключение к БД
//	db, err := sqlx.Connect("postgres", "user=postgres dbname=test sslmode=disable")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//	//fmt.Println("Here2")
//
//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS file_changes_2 (
//	id SERIAL PRIMARY KEY,
//	file_path TEXT NOT NULL,
//	method TEXT NOT NULL,
//	time_change TIMESTAMP WITH TIME ZONE NOT NULL
// )`)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var conf Config
//	if err := yaml.Unmarshal(f, &conf); err != nil {
//		log.Fatal(err)
//	}
//
//	//fmt.Printf("%+v\n", conf)
//
//	// Создание нового Watcher
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer watcher.Close()
//
//	// Start listening for events.
//	go func() {
//		for {
//			select {
//			case event, ok := <-watcher.Events:
//				if !ok {
//					fmt.Println("finishhh")
//					return
//				}
//				log.Println("event:", event)
//				if event.Has(fsnotify.Write) || event.Has(fsnotify.Chmod) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
//					log.Println("modified file:", event.Name)
//
//					//Запуск команды
//					for i := 0; i < len(conf.Commands); i++ {
//						cmd := exec.Command("sh", "-c", conf.Commands[i])
//						stdoutStderr, err := cmd.CombinedOutput()
//						if err != nil {
//							fmt.Println("finishhh")
//
//							log.Fatal(err)
//						}
//
//						fmt.Printf("%s\n", stdoutStderr)
//					}
//					change := &FileChange{FilePath: conf.Path[0], Method: event.String(), Time_change: time.Now().UTC()}
//					_, err = db.NamedExec(`INSERT INTO file_changes_2(file_path, method, time_change)
//												VALUES (:file_path, :method, :time_change)`, map[string]interface{}{
//						"file_path":   change.FilePath,
//						"method":      change.Method,
//						"time_change": change.Time_change.Format(time.RFC3339),
//					})
//					if err != nil {
//						log.Println(err)
//					}
//				}
//
//			case err, ok := <-watcher.Errors:
//				if !ok {
//					fmt.Println("finishhh")
//
//					return
//				}
//				log.Println("error:", err)
//			}
//		}
//	}()
//
//	// Add a path.
//	err = watcher.Add(conf.WatchedPaths)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Block main goroutine forever.
//	<-make(chan struct{})
//
//}

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	WatchedPaths []struct {
		Path     string   `yaml:"path"`
		Commands []string `yaml:"commands"`
	} `yaml:"watched_paths"`
}

type FileChange struct {
	ID          int       `db:"id"`
	FilePath    string    `db:"file_path"`
	Method      string    `db:"method"`
	Time_change time.Time `db:"time_change"`
}

func main() {
	// Чтение yaml-файла
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к БД
	db, err := sqlx.Connect("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS file_changes_2 (   
    id SERIAL PRIMARY KEY,   
    file_path TEXT NOT NULL,   
    method TEXT NOT NULL,   
    time_change TIMESTAMP WITH TIME ZONE NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	var conf Config
	if err := yaml.Unmarshal(f, &conf); err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%+v\n", conf)

	// Создание нового Watcher
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
					fmt.Println("finishhh")

					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Chmod) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					log.Println("modified file:", event.Name)

					// Запуск команды
					for _, watchedPath := range conf.WatchedPaths {
						str1 := strings.Split(watchedPath.Path, "/")
						str2 := strings.Split(event.Name, "/")

						str1 = str1[1:]
						str1 = str1[:len(str1)-1]
						str2 = str2[1:]

						//fmt.Println("watchedPath.Path = ", watchedPath.Path)
						//fmt.Println("watchedPath.Path = ", str1, "len(str1) = ", len(str1))
						//fmt.Println("event = ", str2, "len(str2) = ", len(str2))
						if len(str2)-len(str1) == 1 {
							flag := 0
							for i := range str1 {
								if str1[i] == str2[i] {
									flag = 1
								} else {
									flag = 0
								}
							}
							if flag == 1 {
								for _, command := range watchedPath.Commands {
									cmd := exec.Command("sh", "-c", command)
									stdoutStderr, err := cmd.CombinedOutput()
									if err != nil {
										log.Fatal(err)
									}

									fmt.Printf("%s\n", stdoutStderr)
								}
								change := &FileChange{FilePath: event.Name, Method: event.String(), Time_change: time.Now().UTC()}
								_, err = db.NamedExec(`INSERT INTO file_changes_2(file_path, method, time_change) VALUES (:file_path, :method, :time_change)`, map[string]interface{}{
									"file_path": change.FilePath, "method": change.Method, "time_change": change.Time_change.Format(time.RFC3339)})
								if err != nil {
									log.Println(err)
								}
							}
						}

					}

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
	for _, watchedPath := range conf.WatchedPaths {
		err = watcher.Add(watchedPath.Path)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
