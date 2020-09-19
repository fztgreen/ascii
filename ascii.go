package main

import (
    "fmt"
    "flag"
    "github.com/nfnt/resize"
    "errors"
    "image"
    "image/png"
    "image/jpeg"
    "image/color"
    "strings"
    "os"
)

// Define the enum of valid file extensions.
const(
    e_png = "png"
    e_jpeg = "jpeg"
    e_jpg = "jpg"
)

// Handles errors.
func check(e error) {
    if e != nil {
         panic(e)
    }
}

// Retireves an image from a file.
func retrieveImage(extension string, file *os.File) (image.Image, error) {
    if (strings.Compare(extension, e_png) == 0) {
        return png.Decode(file)
    }

    if (strings.Compare(extension, e_jpeg) == 0) {
        return jpeg.Decode(file)
    }

    if (strings.Compare(extension, e_jpg) == 0) {
       return jpeg.Decode(file)
    }

    return nil, errors.New("Image could not be decoded")
}

// Validates the file extension.
func checkFileExtension(filepath string) (string, error) {
    var split = strings.Split(filepath, ".")
    var extension = strings.ToLower(split[len(split)  - 1])

    if (strings.Compare(extension, e_png) == 0) {
        return e_png, nil
    }

    if (strings.Compare(extension, e_jpeg) == 0) {
        return e_jpeg, nil
    }

    if (strings.Compare(extension, e_jpg) == 0) {
        return e_jpg, nil
    }

    return "", errors.New("Invalid file extension: " + extension)
}

// Main routine.
func main() {
    // Read filepath from args.
    filepath, err := readFilepath()
    check(err)

    fmt.Println("Using filepath:", filepath)

    // Verify extension is valid.
    extension, err := checkFileExtension(filepath)
    check(err)

    fmt.Println("Valid extension was recorded:", extension)

    // Open the file into a stream.
    f, err := os.Open(filepath)
    check(err)

    fmt.Println("Opened File")

    // Retrieves the image from the file
    img, err := retrieveImage(extension, f)
    check(err)

    fmt.Println("Interpreted as an image")

    // Prints your ascii image.
    asciify(img)
}

// Reads a file path from the command line.
func readFilepath() (string, error) {
    // Get command line Args for file location.
    var filepath = flag.String("file", "", "This is the relative file path to the image to convert.")

    flag.Parse()

    if (strings.Compare(*filepath, "") != 0) {
        return *filepath, nil
    }

    return *filepath, errors.New("No file provided!")
}

// Prints out an Ascii image of the picture.
func asciify(img image.Image) {
    levels := []string{"@", "-", "/", "|", "$"}

    res := img.Bounds().Max.Y / img.Bounds().Max.X

    img = resize.Resize((uint)(32 * res), 32, img,  resize.NearestNeighbor)

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 1 {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x += 1 {
            c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
            level := c.Y / 51 // 51 * 5 = 255
            if level == 5 {
                level--
            }
            fmt.Print(levels[level])
        }
        fmt.Print("\n")
    }
    fmt.Println("Image is now Ascii")
}
