package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"

	// "math"
	"strconv"
	"strings"
	"time"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"go.bug.st/serial"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

type MyPort struct {
	p serial.Port
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
	// log.Println("Received from serial:", line)
	return line
}

func (p *MyPort) read_float64() float64 {
	// buf := bufio.NewReader(p.p)
	// line, err := buf.ReadString('\n')
	// if err != nil {
	// 	log.Println(err)
	// }
	line := p.read_str()
	line = line[:len(line)-1]
	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		log.Println("Error parsing float:", err)
	}
	return value
}

// func (p *MyPort) read(buf_size int) []byte {
// 	buf := make([]byte, buf_size)
// 	bytes_read, err := p.p.Read(buf)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Received from serial:", buf[:bytes_read])
// 	return buf[:bytes_read]
// }

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
	p.send_str(string("m" + strconv.Itoa(num)))
	res := make([]float64, 10)
	// for j := range res {
	// 	res[j] = make([]float64, 100)
	// }
	time.Sleep(time.Duration(num+1) * time.Second)
	p.send_str("d")
	for i := 0; i < 10; i++ {
		a := p.read_float64()
		log.Println("Recieved from serial:", a)
		res[i] = a
		time.Sleep(5 * time.Millisecond)
	}
	// p.send_str("d1")
	// for i := 0; i < 100; i++ {
	// buf := p.read(128)
	// res[1][i] = float64(int(buf[0])*256 + int(buf[1]))
	// 	res[1][i] = p.read_float64()
	// }
	// p.send_str("d2")
	// for i := 0; i < 100; i++ {
	// buf := p.read(128)
	// res[2][i] = float64(int(buf[0])*256 + int(buf[1]))
	// 	res[2][i] = p.read_float64()
	// }
	// p.send_str("t")
	// buf := p.read(128)
	// t := (int(buf[0])*256 + int(buf[1])) / 1000
	x := linspace(0, float64(res[len(res)-1]), len(res))
	plt("Time of Life measurement", "t", "sigma", x, res[:])
}

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	var drop_fill string
	fmt.Println("Found ports:")
	for _, port := range ports {
		fmt.Println(port)
	}
	drop_fill = strings.Join(ports, ";")
	switch {
	case len(ports) == 0:
		log.Panicln("No serial ports found.")
	case err != nil:
		log.Fatal(err)
	}
	var my_port string
	var drop_active int32
	var p MyPort
	m := &serial.Mode{
		BaudRate: 9600,
	}

	rl.InitWindow(screenWidth, screenHeight, "Life time measurement")
	textBox := rl.NewRectangle(100, 300, 300, 60)
	textBoxActive := false
	drop_mode := false
	ser_btn_clck := false
	ser_btn_box := rl.NewRectangle(450, 40, 60, 60)
	ser_select := rl.NewRectangle(100, 40, 300, 60)
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

		if rg.DropdownBox(ser_select, drop_fill, &drop_active, drop_mode) {
			drop_mode = !drop_mode
		}

		if ser_btn_clck {
			my_port = ports[drop_active]
			// _, err = fmt.Scan(&my_port)
			s, err := serial.Open(my_port, m)
			p = MyPort{s}
			time.Sleep(time.Second)
			if err != nil {
				log.Fatal(err)
			}
			p.p.SetReadTimeout(5 * time.Second)
			log.Printf("Port set to %s successfully\n", my_port)
			ser_btn_clck = !ser_btn_clck
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			if text[0] == 'm' {
				num, _ := strconv.Atoi(string(text[1:]))
				p.measure(num)
			} else {
				p.send_str(string(text))
			}
			p.read_str()
			text = text[:0]
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.DrawRectangleLinesEx(textBox, 2, rl.Black)
		rl.DrawText("Put your command hear:", 50, 250, 26, rl.Black)
		rl.DrawText(string(text), 110, 320, 26, rl.Magenta)
		ser_btn_clck = rg.Button(ser_btn_box, "Set")
		rl.EndDrawing()

		// rl.CloseWindow()
	}
}
