package main
import (
    "encoding/json"
    "fmt"
)
func main() {
    var v interface{}
    err := json.Unmarshal(nil, &v)
    fmt.Println(err)
}
