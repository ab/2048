package main

import (
    "github.com/nsf/termbox-go"
    "fmt"
    "unicode/utf8"
)

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer fmt.Println("Goodbye!")
    defer termbox.Close()

    fmt.Println("Entered raw keyboard mode. Press any key to see it printed.\r")
    fmt.Println("Use ctrl+c or ctrl+q to exit.\r")

    // termbox.SetInputMode(termbox.InputEsc)

    for {
        fmt.Println("")
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            if ev.Key == termbox.KeyCtrlQ || ev.Key == termbox.KeyCtrlC {
                fmt.Println("Exiting...")
                return
            }
            fmt.Printf("EventKey: %#v\n\r", ev.Key)
            var buf [utf8.UTFMax]byte
            n := utf8.EncodeRune(buf[:], ev.Ch)
            fmt.Printf("Rune[%d]: %s\n\r", n, buf[:n])
        case termbox.EventResize:
            fmt.Println("EventResize")
        case termbox.EventMouse:
            fmt.Println("EventMouse")
        case termbox.EventError:
            panic(ev.Err)
        }
    }
}
