package main

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"runtime/debug"
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
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

// Handles a request.  Meant to be ran independently via a go routine.
func handleConnection(c net.Conn) {
	debug.SetPanicOnFault(true)

	// Set up handling if this routine crashes
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in handle connection error is: %v \n", r)
			c.Write([]byte("Took too long to process.  Try sending smaller sized image."))
			c.Close()
		}
	}()

	fmt.Println("copying connecting data to buffer")

	// Read request headers
	request, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		log.Println("Error reading request:", err)
		c.Close()
		return
	}

	// Get the content length from request headers
	contentLength := request.ContentLength
	if contentLength < 0 {
		// Content-Length header not provided or invalid
		log.Println("Content-Length header not provided or invalid")
		c.Close()
		return
	}

	c.SetReadDeadline(time.Now().Add(10 * time.Second))

	fmt.Println("Parsing form")
	var maxMemoryAmount int64 = 20 * 1024 * 1024
	request.ParseMultipartForm(maxMemoryAmount)

	fmt.Println("Unloading body")
	var body []*multipart.FileHeader = (*request.MultipartForm).File["image"]

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

	// Set response headers
	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/plain\r\n\r\n", len(bi))

	// Write headers
	_, err = c.Write([]byte(headers))
	if err != nil {
		log.Println("Error writing headers:", err)
		return
	}

	fmt.Println("Writing out connection")
	c.Write(bi)
	c.Close()
}
