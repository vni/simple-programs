package main
import (
    "fmt"
    "os"
    "io"
)

var base64Table []byte = []byte{
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
    'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
    'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
    'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
    'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
    'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
    'w', 'x', 'y', 'z', '0', '1', '2', '3',
    '4', '5', '6', '7', '8', '9', '+', '/',
}

func processChunk(buf []byte, out *os.File) {
    bytesWritten := 0
    outbuf := make([]byte, 4)

    fmt.Printf("got input buf of %v bytes\n", len(buf))

    for pos := 0; pos < len(buf); pos += 3 {
        left := len(buf) - pos
        switch left {
        case 1:
            a := buf[pos+0]
            outbuf[0] = base64Table[byte((a >> 2) & 0x3F)]
            outbuf[1] = base64Table[byte((a & 0x03) << 4)]
            outbuf[2] = base64Table[byte('=')]
            outbuf[3] = base64Table[byte('=')]
        case 2:
            a, b := buf[pos+0], buf[pos+1]
            outbuf[0] = base64Table[byte((a >> 2) & 0x3F)]
            outbuf[1] = base64Table[byte((((a & 0x03) << 4) | ((b & 0xF0) >> 4)) & 0x3F)]
            outbuf[2] = base64Table[byte(b & 0x0F)]
            outbuf[3] = base64Table[byte('=')]
        default:
            a, b, c := buf[pos+0], buf[pos+1], buf[pos+2]
            outbuf[0] = base64Table[byte((a >> 2) & 0x3F)]
            outbuf[1] = base64Table[byte((((a & 0x3) << 4)  | ((b & 0xF0) >> 4)) & 0x3F)]
            outbuf[2] = base64Table[byte((((b & 0x0F) << 2) | ((c & 0xC0) >> 6)) & 0x3F)]
            outbuf[3] = base64Table[byte(c & 0x3F)]
        }

        out.Write(outbuf)
        bytesWritten += 4
        if (bytesWritten % 64) == 0 {
            out.Write([]byte("\r\n"))
        }
    }

    if ((len(buf) != cap(buf)) && ((bytesWritten % 64) != 0)) {
        out.Write([]byte("\r\n"))
    }

    fmt.Printf("Processed %v bytes chunk.\n", bytesWritten);
}

func processFile(in *os.File, out *os.File) {
    buf := make([]byte, 600) // 600 to be a multiple of 3
    n, err := in.Read(buf)
    for ; err == nil; n, err = in.Read(buf) {
        processChunk(buf[:n], out)
    }
    if (err != io.EOF) {
        fmt.Printf("*ERROR* during reading the input file: %v\n", err)
        os.Exit(3)
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Printf("Usage: %v inputFileName outputFileName\n", os.Args[0])
        os.Exit(1)
    }

    inputFileName := os.Args[1]
    outputFileName := os.Args[2]

    in, err := os.Open(inputFileName)
    if err != nil {
        fmt.Printf("Failed to open input file '%v' due to: %v\n", inputFileName, err)
        os.Exit(1)
    }
    defer in.Close()

    out, err := os.OpenFile(outputFileName, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
    if err != nil {
        fmt.Printf("Failed to open output file '%v' due to: %v\n", outputFileName, err)
        os.Exit(2)
    }
    defer out.Close()

    fmt.Printf("inputFileName: %v\noutputFileName: %v\n", inputFileName, outputFileName)
    processFile(in, out)
}
