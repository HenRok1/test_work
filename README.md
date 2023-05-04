# Test task for job

# Information
Консольное приложение, для отслеживания изменений в директории на языке Go. \
В разработке использовался Golang 1.18, база данных Postgres. \
Были использованы такие библиотеки:
```go
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
```

## Особенности:
Конфигурация считывается с config.yaml файла, в котором прописывается путь, команды и папка для логов. \
### Пример
```yaml
watched_paths:
  - path: "/Users/project/test_work/cmd/"
    commands:
      - go build -o ../../build/bin/app1 service/main.go
      - ../../build/bin/app1
    log_file: /tmp/log1.out

  - path: "/Users/project/test_work/"
    commands:
      - go build -o ../../build/bin/app2 service/main2.go
      - ../../build/bin/app2
    log_file: /tmp/log2.out

  - path: "/Users/project/test_work/cmd/service/"
    commands:
      - echo hello >> ../build/Nail.txt
    log_file: /tmp/log3.out
```
### Подключение к базе данных - подключение к серверу postgres:
```go
func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
```

## Вопросы которым возникали
- Логи переводятся в файл только для изменений директорий
- Во время обработки нескольких путей, появился вопрос, как обрабатывать изменения в конкретной директории и отрбатывать нужные команды(решение: проеврка длинны пути изменений и длинны event, и если оно отличается на 1, то значит изменения произошли в нужной директории, однако возникла проблема с лишними пробелами, из-за чего длина была некорректной, решил при помощи слайсов)
```go
str1 := strings.Split(watchedPath.Path, "/")
str2 := strings.Split(event.Name, "/")

str1 = str1[1:]
str1 = str1[:len(str1)-1]
str2 = str2[1:]

if len(str2)-len(str1) == 1 {
```

# Installation:
1) Нужно перейти в директорию cmd
2) Необходимо прописать в config.yaml файле путь, команды и папку для логирования
3) Нужно поднять сервер postgres и прописать в user и db.name
```go
db, err := sqlx.Connect("postgres", "user=postgres dbname=test sslmode=disable")
```
4) При успешном подключении выведитется окно об этом
5) Выполнить команду:
```go
go build -o main_app main.go
```
6) Запустить:
```./main_app```