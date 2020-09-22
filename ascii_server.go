package main

import (
    "encoding/base64"
    "fmt"
    "log"
    "net"
    "bytes"
    "strings"
    "io"
    "io/ioutil"
    "net/http"
    "bufio"
)

// Main function to listen for network input.
func main() {
    l, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }

    defer l.Close()

    for {
        // Wait for a connection
        conn, err := l.Accept()
        if  err != nil {
            log.Fatal(err)
        }

        go handleConnection(conn)
    }
}


// Handles a request.  Meant to be ran independently via a go routine.
func handleConnection(c net.Conn) {
//    var filepath = "../heisei.png"
    var extension = "png"

    fmt.Println("----------")
    fmt.Println(c)
    fmt.Println("----------")
    // Open the file into a stream.
    //f, err := os.Open(filepath)
    //check(err)

    //fmt.Println("Opened File")

    // Retrieves the image from the file
    var buf bytes.Buffer

    _ , err := io.Copy(&buf, c)

    request, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(buf.String())))
    body, err := ioutil.ReadAll(request.Body)
    check(err)

    fmt.Println("Read String")
    fmt.Println(string(body))
    encoded := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(body)))
    fmt.Println(encoded)
    img, err := retrieveImage(extension, encoded)
    check(err)

    fmt.Println("Interpreted as an image")

    // Prints your ascii image.
    bi := []byte(asciify(img))
    c.Write(bi)
    c.Close()
}
