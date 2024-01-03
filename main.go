package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"github.com/gosuri/uilive"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func parsePTYOutput(input string) string {
	writer := uilive.New()

	var buffer bytes.Buffer

	writer.Start()

	writer.Out = &buffer

	fmt.Fprintf(writer, input)

	writer.Stop()

	out := buffer.String()

	return out
}

func main() {
	// Initialization of required SDL2 subsystems
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow("SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	font, err := ttf.OpenFont("/home/matt/.fonts/GoMonoNerdFont-Regular.ttf", 22)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	// Create a new pty and start a shell
	c := exec.Command("/bin/zsh")
	p, err := pty.Start(c)
	if err != nil {
		panic(err)
	}
	defer c.Process.Kill()

	// Write to the pty and read the output
	p.Write([]byte("ls\r"))
	time.Sleep(1 * time.Second)
	b := make([]byte, 64)
	_, err = p.Read(b) // TODO: This is blocking
	if err != nil {
		panic(err)
	}

	out := parsePTYOutput(string(b))
	fmt.Printf(out)

	// Render the ptys output to a surface and then to a texture
	surface, err := font.RenderUTF8Solid(out, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	_, _, font_width, font_height, err := texture.Query()
	if err != nil {
		panic(err)
	}

	font_rect := sdl.Rect{X: 0, Y: 0, W: font_width, H: font_height}

	err = renderer.Clear()
	if err != nil {
		panic(err)
	}

	// Copy the texture to the renderer and present it
	err = renderer.Copy(texture, nil, &font_rect)
	if err != nil {
		panic(err)
	}

	renderer.Present()

	sdl.Delay(5000)
}
