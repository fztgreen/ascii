package main

import (
    "fmt"
    "log"
    "net"
    "bytes"
    "strings"
    "io"
    "net/http"
    "bufio"
    "mime/multipart"
    "time"
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
    fmt.Println("copying connecting data to buffer")

    //defer c.Close()
    //bitties := make([]byte, 5242880)
    //reqLen, _ := c.Read(bitties)

    //batties := bytes.NewBuffer(bitties)

    c.SetReadDeadline(time.Now().Add(5 * time.Second))

    var buf bytes.Buffer
    num , err := io.CopyN(&buf, c, 995242880)

    fmt.Println(num, "bytes read")

    fmt.Println("Converting bytes.Buffer to the http.Request")
    request, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(buf.String())))

    fmt.Println("Parsing form")
    request.ParseMultipartForm(65536)

    fmt.Println("Unloading body")
    var body []*multipart.FileHeader
    body = (*request.MultipartForm).File["image"]

    // Verify extension is valid.
    fmt.Println("Identifying Extension")
    extension, err := checkFileExtension((*body[0]).Filename)
    check(err)

    fmt.Println("Opening file")
    file, err := (*body[0]).Open()
    check(err)

    fmt.Println("Retrieving Image")
    img, err := retrieveImage(extension, file)
    check(err)

    fmt.Println("Converting to Ascii")

    // Prints your ascii image.
    ascii := asciify(img)
    fmt.Println(ascii)
    bi := []byte(ascii)

    fmt.Println("Writing out connection")
    c.Write(bi)
    c.Close()
}
