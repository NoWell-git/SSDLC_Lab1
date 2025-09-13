package main

import (
  "bufio"
  "database/sql"
  "fmt"
  "os"
  "strings"

  "golang.org/x/term"
  "gopkg.in/yaml.v3"

  _ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
  Host     string `yaml:"host"`
  Port     int    `yaml:"port"`
  Database string `yaml:"database"`
  User     string `yaml:"user"`
  Password string `yaml:"password"`
  SSLMode  string `yaml:"sslmode"`
}

var sanitizeChars = []string{";", "'", "\"", "<", ">"}
var id = 0

func sanitizeInput(input string) string {
  sanitized := strings.TrimSpace(input)
  flag := false
  id++

  for _, ch := range sanitizeChars {
    if strings.Contains(sanitized, ch) {
      flag = true
    }
    sanitized = strings.ReplaceAll(sanitized, ch, "")
  }

  if flag {
    if id == 1 {
      fmt.Println("Логин содержит запрещенные символы")
    }
    if id == 2 {
      fmt.Println("Пароль содержит запрещенные символы")
    }
  }

  return sanitized
}

func main() {

  data, _ := os.ReadFile("config.yaml")
  var cfg Config
  yaml.Unmarshal(data, &cfg)

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Введите имя пользователя: ")
  userInput, _ := reader.ReadString('\n')
  userInput = sanitizeInput(userInput)

  fmt.Print("Введите пароль: ")
  bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
  passInput := string(bytePassword)
  fmt.Println()
  passInput = sanitizeInput(passInput)

  cfg.User = userInput
  cfg.Password = passInput

  connStr := fmt.Sprintf("host=%s port=%d database=%s user=%s password=%s sslmode=%s",
    cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.SSLMode)

  db, _ := sql.Open("pgx", connStr)
  defer db.Close()

  var version string
  err := db.QueryRow("SELECT VERSION()").Scan(&version)
  if err != nil {
    if strings.Contains(err.Error(), "connect: connection refused") {
      fmt.Println("Ошибка подключения к базе данных")
    } else {
      fmt.Println("Неверное имя пользователя или пароль")
    }
  } else {
    fmt.Println(version)
  }
}
