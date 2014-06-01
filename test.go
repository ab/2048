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
    defer termbox.Close()

    // termbox.SetInputMode(termbox.InputEsc)

    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            if ev.Key == termbox.KeyCtrlQ {
                fmt.Println("Goodbye!")
                return
            }
            fmt.Printf("EventKey: %#v\n", ev.Key)
            var buf [utf8.UTFMax]byte
            n := utf8.EncodeRune(buf[:], ev.Ch)
            fmt.Printf("Rune: %s\n", buf[:n])
        case termbox.EventResize:
            fmt.Println("EventResize")
        case termbox.EventMouse:
            fmt.Println("EventMouse")
        case termbox.EventError:
            panic(ev.Err)
        }
    }
}
