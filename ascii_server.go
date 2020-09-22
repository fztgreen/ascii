package main

import (
    "fmt"
    "log"
    "net"
    "os"
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
    var filepath = "../heisei.png"
    var extension = "png"

    // Open the file into a stream.
    f, err := os.Open(filepath)
    check(err)

    fmt.Println("Opened File")

    // Retrieves the image from the file
    img, err := retrieveImage(extension, f)
    check(err)

    fmt.Println("Interpreted as an image")

    // Prints your ascii image.
    b := []byte(asciify(img))
    c.Write(b)
    c.Close()
}
