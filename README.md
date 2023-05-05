# Test task for job

# Инструкция
1) Внести в config.yaml путь, команды и папку для логов, как это показано в примере(в начала и в конце путей должны быть "/")
2) Запустить docker
3) Прописать: ```make run```
4) В новой командной строке прописать: ```make all```
5) Отдельным терминалом вносить изменения в директориях, которые вы прописали в config.yaml


# Информация
Консольное приложение для отслеживания изменений в директории на языке Go. \
В разработке использовался Golang 1.18, база данных PostgreSQL. \
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

## Особенности
Конфигурация считывается с config.yaml файла, в котором прописывается путь, команды и папка для логов.
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
```
### Подключение к серверу postgres:
```go
func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
```

## Вопросы которые возникали
- Логи переводятся в файл только для изменений директорий
- Во время обработки нескольких путей, появился вопрос, как обрабатывать изменения в конкретной директории и отрабатывать нужные команды(решение: проверка длинны пути изменений и длинны event соответственно, и если оно отличается на 1, то значит изменения произошли в нужной директории, однако возникла проблема с лишними пробелами, из-за чего длина была некорректной, решил при помощи слайсов)
```go
str1 := strings.Split(watchedPath.Path, "/")
str2 := strings.Split(event.Name, "/")

str1 = str1[1:]
str1 = str1[:len(str1)-1]
str2 = str2[1:]

if len(str2)-len(str1) == 1 {
```
- Для подключения к базе данных по хорошему создать переменную окружения, но в данном случае харкод
