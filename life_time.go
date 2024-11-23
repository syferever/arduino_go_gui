package main

import (
	"bufio"
	"image/color"
	"log"

	// "math"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tarm/serial"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

type MyPort struct {
	p *serial.Port
}

func (p *MyPort) send_str(s string) {
	msg := strings.TrimRight(s, "\n") + "\n"
	_, err := p.p.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Send to serial: \"%v\" -> %v", msg[:len(msg)-1], []rune(msg))
}

func (p *MyPort) read_str() string {
	buf := bufio.NewReader(p.p)
	line, err := buf.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Received from serial:", line)
	return line
}

func (p *MyPort) read_float64() float64 {
	buf := bufio.NewReader(p.p)
	line, err := buf.ReadString('\r')
	if err != nil {
		log.Println(err)
	}
	line = line[1 : len(line)-1]
	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		log.Println("Error parsing float:", err)
	}
	return value
}

func (p *MyPort) read(buf_size int) []byte {
	buf := make([]byte, buf_size)
	bytes_read, err := p.p.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Received from serial:", buf[:bytes_read])
	return buf[:bytes_read]
}

func linspace(start, end float64, n int) []float64 {
	if n <= 0 {
		return []float64{}
	}
	if n == 1 {
		return []float64{start}
	}
	step := (end - start) / float64(n-1)
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = start + step*float64(i)
	}

	return values
}

func plt(title, x_ax, y_ax string, x, y []float64) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = x_ax
	p.Y.Label.Text = y_ax

	n := len(x)
	points := make(plotter.XYs, n)
	for i := range points {
		points[i].X = x[i]
		points[i].Y = y[i]
	}
	graph, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(graph)
	graph.Color = color.RGBA64{0, 255, 0, 0}

	if err := p.Save(6*vg.Inch, 6*vg.Inch, "measured.png"); err != nil {
		panic(err)
	}
}

func (p *MyPort) measure(num int) {
	p.send_str(string('m' + rune(num)))
	res := make([][]float64, 3)
	for j := range res {
		res[j] = make([]float64, 100)
	}
	time.Sleep(time.Duration(num) * time.Second)
	p.send_str("d0")
	for i := 0; i < 100; i++ {
		// buf := p.read(128)
		// res[0][i] = float64(int(buf[0])*256 + int(buf[1]))
		res[0][i] = p.read_float64()
		// fmt.Println(res)
	}
	p.send_str("d1")
	for i := 0; i < 100; i++ {
		// buf := p.read(128)
		// res[1][i] = float64(int(buf[0])*256 + int(buf[1]))
		res[1][i] = p.read_float64()
	}
	p.send_str("d2")
	for i := 0; i < 100; i++ {
		// buf := p.read(128)
		// res[2][i] = float64(int(buf[0])*256 + int(buf[1]))
		res[2][i] = p.read_float64()
	}
	p.send_str("t")
	buf := p.read(128)
	t := (int(buf[0])*256 + int(buf[1])) / 1000
	x := linspace(0, float64(t), 100)
	plt("Time of Life measurement", "t", "sigma", x, res[2][:])
}

func main() {
	c := &serial.Config{
		Name: "COM11",
		Baud: 9600,
	}
	c.ReadTimeout = time.Millisecond * 100
	s, err := serial.OpenPort(c)
	p := MyPort{s}
	time.Sleep(time.Second)
	if err != nil {
		log.Fatal(err)
	}

	rl.InitWindow(screenWidth, screenHeight, "Life time measurement")
	textBox := rl.NewRectangle(100, 100, 300, 60)
	textBoxActive := false
	var text []rune
	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			if rl.CheckCollisionPointRec(rl.GetMousePosition(), textBox) {
				textBoxActive = true
			}
		}

		if textBoxActive {
			key := rl.GetCharPressed()
			for key > 0 {
				if key >= 32 && key <= 125 {
					text = append(text, rune(key))
				}
				key = rl.GetCharPressed()
			}
			if rl.IsKeyPressed(rl.KeyBackspace) {
				if len(text) > 0 {
					text = text[:len(text)-1]
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			p.send_str(string(text))
			if text[0] == 'm' {
				p.measure(int(text[1]))
			}
			p.read_str()
			text = text[:0]
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.DrawRectangleLinesEx(textBox, 2, rl.Black)
		rl.DrawText("Put your command hear:", 30, 50, 26, rl.Black)
		rl.DrawText(string(text), 110, 120, 26, rl.Magenta)
		rl.EndDrawing()

		// rl.CloseWindow()
	}
}
