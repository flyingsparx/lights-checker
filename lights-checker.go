package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "time"
  "strings"
  "os"
  "gopkg.in/mailgun/mailgun-go.v1"
  "github.com/jasonlvhit/gocron"
)

type Config struct {
  Mailgun struct {
    Domain string `json:"domain"`
    ApiKey string `json:"apiKey"`
    PublicKey string `json:"publicKey"`
    From string `json:"from"`
    To string `json:"to"`
  } `json:"mailgun"`
  HueBridge struct {
    Address string `json:"address"`
    Username string `json:"username"`
  } `json:"hueBridge"`
}

var config Config
var alreadySent bool
var mg mailgun.Mailgun

func loadConfiguration(file string) Config {
  var c Config
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
    panic("Configuration not found")
  }
  jsonParser := json.NewDecoder(configFile)
  err2 := jsonParser.Decode(&c)
  if err2 != nil {
    panic("Unable to parse config file")
  }
  return c 
}

func sendNotification(lights []string) {
  body := fmt.Sprintf("Hi there,\n\nIt is late and the following lights are turned on:\n\n%s", strings.Join(lights, ",\n"))
  mg.Send(mailgun.NewMessage(config.Mailgun.From, "Lights on", body, config.Mailgun.To))
}

func checkBulbs() {
  var hour = time.Now().Hour()

  if hour > 20 || hour < 6 {
    type bulbState struct {
      On bool `json:"on"`
    }
    type bulb struct {
      State bulbState `json:"state"`
      Name string `json:"name"`
    }
    var response = make(map[int]bulb)

    resp, err := http.Get("http://" + config.HueBridge.Address + "/api/" + config.HueBridge.Username + "/lights")
    if err == nil {
      defer resp.Body.Close()
      body, err2 := ioutil.ReadAll(resp.Body)
      onBulbs := make([]string, 0)
      if err2 == nil {
        json.Unmarshal(body, &response)
        for _, value := range response {
          if value.State.On {
            onBulbs = append(onBulbs, value.Name)
          }
        }
      }
      if len(onBulbs) > 0 && alreadySent == false {
        sendNotification(onBulbs)
        alreadySent = true
      }
      if len(onBulbs) == 0 {
        alreadySent = false
      }
    }
  }
}

func main() {
  config = loadConfiguration("config.json")
  mg = mailgun.NewMailgun(config.Mailgun.Domain, config.Mailgun.ApiKey, config.Mailgun.PublicKey)
  alreadySent = false
  gocron.Every(10).Seconds().Do(checkBulbs)
  <- gocron.Start()
}
