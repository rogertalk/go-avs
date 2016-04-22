Alexa Voice Service for Go
==========================

A simple package for communicating with Amazon’s HTTP/2 API for AVS.

Requires Go 1.6 or later.


Example
-------

```go
package main

import (
  "fmt"
  "io/ioutil"
  "os"

  "github.com/rogertalk/go-avs"
)

func main() {
  request := avs.NewRequest("YOUR ACCESS TOKEN")
  request.Event = avs.NewRecognize("abc123", "abc123dialog")
  request.Audio, _ = os.Open("./speech.wav")
  response, err := avs.DefaultClient.Do(request)
  if err != nil {
    fmt.Printf("Failed to call AVS: %v\n", err)
    return
  }
  // Depending on the request, there may be more (or less) directives.
  switch d := response.Directives[0].Typed().(type) {
  case *avs.Speak:
    data := response.Content[d.ContentId()]
    ioutil.WriteFile("./response.mp3", data, 0666)
    fmt.Println("Wrote Alexa’s reply to response.mp3")
  }
}
```