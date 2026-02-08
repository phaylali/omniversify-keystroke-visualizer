package gui

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"time"

	"github.com/BurntSushi/freetype-go/freetype"
	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"omniversify-keystroke-visualizer/config"
)

type Overlay struct {
	X      *xgbutil.XUtil
	cfg    *config.Config
	events chan string
	win    *xwindow.Window
	font   *truetype.Font
}

func NewOverlay(cfg *config.Config) (*Overlay, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to X: %w", err)
	}

	// Load font
	fontData, err := os.ReadFile("font/kenney_input_keyboard_&_mouse.otf")
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}
	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	o := &Overlay{
		X:      X,
		cfg:    cfg,
		events: make(chan string, 10),
		font:   f,
	}

	// Calculate window size and position
	w, h := 300, 80 // Base size
	x, y := o.calculatePosition(w, h)

	win, err := xwindow.Generate(X)
	if err != nil {
		return nil, fmt.Errorf("LOG: failed to generate window ID: %w", err)
	}
	o.win = win

	// Create window with OverrideRedirect to remove decorations
	err = win.CreateChecked(X.RootWin(), x, y, w, h, xproto.CwOverrideRedirect, 1)
	if err != nil {
		fmt.Printf("LOG: Overlay failed to create window: %v\n", err)
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	fmt.Printf("LOG: Overlay window created successfully at %d,%d\n", x, y)

	go o.run()

	return o, nil
}

func (o *Overlay) calculatePosition(winW, winH int) (int, int) {
	setup := xproto.Setup(o.X.Conn())
	screen := setup.DefaultScreen(o.X.Conn())
	screenW := int(screen.WidthInPixels)
	screenH := int(screen.HeightInPixels)

	var x, y int
	switch o.cfg.Position {
	case "top-left":
		x, y = 0, 0
	case "top-center":
		x, y = (screenW/2)-(winW/2), 0
	case "top-right":
		x, y = screenW-winW, 0
	case "center":
		x, y = (screenW/2)-(winW/2), (screenH/2)-(winH/2)
	case "bottom-left":
		x, y = 0, screenH-winH
	case "bottom-center":
		x, y = (screenW/2)-(winW/2), screenH-winH
	case "bottom-right":
		x, y = screenW-winW, screenH-winH
	default:
		x, y = (screenW/2)-(winW/2), screenH-winH
	}

	return x + o.cfg.XOffset, y + o.cfg.YOffset
}

func (o *Overlay) run() {
	var timer *time.Timer
	for {
		select {
		case text := <-o.events:
			os.Stderr.WriteString(fmt.Sprintf("LOG: Overlay handling event: %s\n", text))
			o.win.Map()
			// Small delay to let XWayland/Compositor process the map
			time.Sleep(50 * time.Millisecond)

			o.draw(text)

			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(time.Duration(o.cfg.DurationMs)*time.Millisecond, func() {
				o.win.Unmap()
			})
		}
	}
}

func (o *Overlay) draw(text string) {
	logFile, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	log := func(s string) {
		logFile.WriteString(time.Now().Format("15:04:05") + " " + s + "\n")
	}

	if o.font == nil {
		log("Error: font is nil in draw")
		return
	}

	width, height := 400, 100
	img := xgraphics.New(o.X, image.Rect(0, 0, width, height))

	// TEST: Use solid BLUE
	bgBGRA := xgraphics.BGRA{R: 0, G: 0, B: 255, A: 255}
	img.For(func(x, y int) xgraphics.BGRA {
		return bgBGRA
	})

	// Debug: Draw a WHITE square
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			img.SetBGRA(x, y, xgraphics.BGRA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	log(fmt.Sprintf("Drawing text '%s'", text))

	// Draw text
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(o.font)
	c.SetFontSize(32)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))

	pt := freetype.Pt(60, 60)
	_, err := c.DrawString(text, pt)
	if err != nil {
		log(fmt.Sprintf("Error drawing string: %v", err))
	}

	img.XPaint(o.win.Id)
	img.XSurfaceSet(o.win.Id)

	// LOW-LEVEL DEBUG: Fill a white rectangle using X server directly
	id, _ := o.X.Conn().NewId()
	gc := xproto.Gcontext(id)
	xproto.CreateGC(o.X.Conn(), gc, xproto.Drawable(o.win.Id), xproto.GcForeground, []uint32{0xffffff})
	xproto.PolyFillRectangle(o.X.Conn(), xproto.Drawable(o.win.Id), gc, []xproto.Rectangle{{X: 10, Y: 10, Width: 100, Height: 30}})
	xproto.FreeGC(o.X.Conn(), gc)

	xproto.ClearArea(o.X.Conn(), false, o.win.Id, 0, 0, uint16(width), uint16(height))
	o.X.Sync()
	log("Draw completed")
}

func parseHexColor(s string) color.RGBA {
	c := color.RGBA{A: 255}
	if s == "white" {
		return color.RGBA{255, 255, 255, 255}
	}
	if len(s) == 7 && s[0] == '#' {
		fmt.Sscanf(s[1:], "%02x%02x%02x", &c.R, &c.G, &c.B)
	}
	return c
}

func (o *Overlay) Show(text string) {
	select {
	case o.events <- text:
	default:
	}
}

func (o *Overlay) Close() {
	if o.win != nil {
		o.win.Destroy()
	}
}

func init() {
	os.Remove("debug.log")
	runtime.LockOSThread()
}
