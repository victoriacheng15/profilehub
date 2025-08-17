package main

import (
    "fmt"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v3"
)


func printYAML(key string, value interface{}, indent int) {
    pad := ""
    for i := 0; i < indent; i++ {
        pad += "  "
    }
    switch v := value.(type) {
    case map[string]interface{}:
        fmt.Printf("%s%v:\n", pad, key)
        for k, val := range v {
            printYAML(k, val, indent+1)
        }
    case []interface{}:
        fmt.Printf("%s%v:\n", pad, key)
        for i, item := range v {
            printYAML(fmt.Sprintf("[%d]", i), item, indent+1)
        }
    default:
        fmt.Printf("%s%v: %v\n", pad, key, v)
    }
}

func main() {
    configFile := "src/config/config.yml"
    data, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Fatalf("Error reading config.yml: %v", err)
    }

    var config map[string]interface{}
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        log.Fatalf("Error parsing config.yml: %v", err)
    }

    for key, value := range config {
        printYAML(key, value, 0)
    }
}
