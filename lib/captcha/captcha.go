package captcha

import (
	"crypto/rand"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/statikFS"
	"math"
	"math/big"
	"strconv"
)

const (
	operator        = "+-*/"
	defaultLen      = 4
	defaultFontSize = 25
	defaultDpi      = 72
)

// Captcha Graphical CAPTCHA Use font default ttf format
// w image width, h image height, CodeLen number of CAPTCHA
// FontSize font size, Dpi clarity
// mode validation mode 0: normal string, 1: simple mathematical formula up to 10
type Captcha struct {
	W, H, CodeLen int
	FontSize      float64
	Dpi           int
	mode          int
}

// NewCaptcha Instantiating CAPTCHA
func NewCaptcha(w, h, CodeLen int) *Captcha {
	return &Captcha{W: w, H: h, CodeLen: CodeLen}
}

// OutPut generates a captcha image and returns the image in RGBA format
func (captcha *Captcha) OutPut() (string, *image.RGBA) {
	img := captcha.initCanvas()
	return captcha.doImage(img)
}

// RangeRand Get the random number in the interval [-m, n]
func (captcha *Captcha) RangeRand(min, max int64) int64 {
	if min > max {
		log.ErrorPrint("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))

		return min + result.Int64()
	}
}

// Random strings
func (captcha *Captcha) getRandCode() string {
	if captcha.CodeLen <= 0 {
		captcha.CodeLen = defaultLen
	}

	return tool.GetToken(captcha.CodeLen)
}

// getFormulaMixData returns a tuple with an arithmetic operation string and an array of strings containing the operands and the operator
func (captcha *Captcha) getFormulaMixData() (string, []string) {
	num1 := int(captcha.RangeRand(6, 12))
	num2 := int(captcha.RangeRand(0, 6))
	opArr := []rune(operator)
	opRand := opArr[captcha.RangeRand(0, 2)]

	strNum1 := strconv.Itoa(num1)
	strNum2 := strconv.Itoa(num2)

	var ret int
	var opRet string
	switch string(opRand) {
	case "+":
		ret = num1 + num2
		opRet = "+"
	case "-":
		ret = num1 - num2
		opRet = "-"
	case "*":
		ret = num1 * num2
		opRet = "×"
	}

	return strconv.Itoa(ret), []string{strNum1, opRet, strNum2, "=", "?"}
}

// Initialising the canvas
func (captcha *Captcha) initCanvas() *image.RGBA {
	dest := image.NewRGBA(image.Rect(0, 0, captcha.W, captcha.H))

	// Random colours
	r := uint8(255) // uint8(captcha.RangeRand(50, 250))
	g := uint8(255) // uint8(captcha.RangeRand(50, 250))
	b := uint8(255) // uint8(captcha.RangeRand(50, 250))

	// Fill background colour
	for x := 0; x < captcha.W; x++ {
		for y := 0; y < captcha.H; y++ {
			dest.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255}) //Set the transparency of the alpha image
		}
	}

	return dest
}

// doImage generates an image and the corresponding captcha code
func (captcha *Captcha) doImage(dest *image.RGBA) (string, *image.RGBA) {
	gc := draw2dimg.NewGraphicContext(dest)

	defer gc.Close()
	defer gc.FillStroke()

	captcha.setFont(gc)
	captcha.doPoint(gc)
	captcha.doLine(gc)
	captcha.doSinLine(gc)

	var codeStr string
	if captcha.mode == 1 {
		ret, formula := captcha.getFormulaMixData()
		log.DebugPrint("Challenge: %s, Answer: %s", formula, ret)
		codeStr = ret
		captcha.doFormula(gc, formula)
	} else {
		codeStr = captcha.getRandCode()
		captcha.doCode(gc, codeStr)
	}

	return codeStr, dest
}

// doCode captcha characters set to the image
func (captcha *Captcha) doCode(gc *draw2dimg.GraphicContext, code string) {
	for l := 0; l < len(code); l++ {
		y := captcha.RangeRand(int64(captcha.FontSize)-1, int64(captcha.H)+6)
		x := captcha.RangeRand(1, 20)

		// Random colours
		r := uint8(captcha.RangeRand(0, 200))
		g := uint8(captcha.RangeRand(0, 200))
		b := uint8(captcha.RangeRand(0, 200))

		gc.SetFillColor(color.RGBA{R: r, G: g, B: b, A: 255})
		gc.FillStringAt(string(code[l]), float64(x)+captcha.FontSize*float64(l), float64(int64(captcha.H)-y)+captcha.FontSize)
		gc.Stroke()
	}
}

// doFormula captcha characters set to the image
func (captcha *Captcha) doFormula(gc *draw2dimg.GraphicContext, formulaArr []string) {
	for l := 0; l < len(formulaArr); l++ {
		y := captcha.RangeRand(0, 10)
		x := captcha.RangeRand(5, 10)

		// Random colours
		r := uint8(captcha.RangeRand(10, 200))
		g := uint8(captcha.RangeRand(10, 200))
		b := uint8(captcha.RangeRand(10, 200))

		gc.SetFillColor(color.RGBA{R: r, G: g, B: b, A: 255})

		gc.FillStringAt(formulaArr[l], float64(x)+captcha.FontSize*float64(l), captcha.FontSize+float64(y))
		gc.Stroke()
	}
}

// doLine Adding interference lines
func (captcha *Captcha) doLine(gc *draw2dimg.GraphicContext) {
	// Setting up interference lines
	for n := 0; n < 5; n++ {
		// gc.SetLineWidth(float64(captcha.RangeRand(1, 2)))
		gc.SetLineWidth(1)

		// Random background colours
		r := uint8(captcha.RangeRand(0, 255))
		g := uint8(captcha.RangeRand(0, 255))
		b := uint8(captcha.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{R: r, G: g, B: b, A: 255})

		// Initialisation position
		gc.MoveTo(float64(captcha.RangeRand(0, int64(captcha.W)+10)), float64(captcha.RangeRand(0, int64(captcha.H)+5)))
		gc.LineTo(float64(captcha.RangeRand(0, int64(captcha.W)+10)), float64(captcha.RangeRand(0, int64(captcha.H)+5)))

		gc.Stroke()
	}
}

// Adding points of disturbance
func (captcha *Captcha) doPoint(gc *draw2dimg.GraphicContext) {
	for n := 0; n < 50; n++ {
		gc.SetLineWidth(float64(captcha.RangeRand(1, 3)))

		// Random colours
		r := uint8(captcha.RangeRand(0, 255))
		g := uint8(captcha.RangeRand(0, 255))
		b := uint8(captcha.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{R: r, G: g, B: b, A: 255})

		x := captcha.RangeRand(0, int64(captcha.W)+10) + 1
		y := captcha.RangeRand(0, int64(captcha.H)+5) + 1

		gc.MoveTo(float64(x), float64(y))
		gc.LineTo(float64(x+captcha.RangeRand(1, 2)), float64(y+captcha.RangeRand(1, 2)))

		gc.Stroke()
	}
}

// Adding sine interference lines
func (captcha *Captcha) doSinLine(gc *draw2dimg.GraphicContext) {
	h1 := captcha.RangeRand(-12, 12)
	h2 := captcha.RangeRand(-1, 1)
	w2 := captcha.RangeRand(5, 20)
	h3 := captcha.RangeRand(5, 10)

	h := float64(captcha.H)
	w := float64(captcha.W)

	// Random colours
	r := uint8(captcha.RangeRand(128, 255))
	g := uint8(captcha.RangeRand(128, 255))
	b := uint8(captcha.RangeRand(128, 255))

	gc.SetStrokeColor(color.RGBA{R: r, G: g, B: b, A: 255})
	gc.SetLineWidth(float64(captcha.RangeRand(2, 4)))

	var i float64
	for i = -w / 2; i < w/2; i = i + 0.1 {
		y := h/float64(h3)*math.Sin(i/float64(w2)) + h/2 + float64(h1)

		gc.LineTo(i+w/2, y)

		if h2 == 0 {
			gc.LineTo(i+w/2, y+float64(h2))
		}
	}

	gc.Stroke()
}

// SetMode Setting mode
func (captcha *Captcha) SetMode(mode int) {
	captcha.mode = mode
}

// SetFontSize Set font size
func (captcha *Captcha) SetFontSize(fontSize float64) {
	captcha.FontSize = fontSize
}

// setFont Setting the font
func (captcha *Captcha) setFont(gc *draw2dimg.GraphicContext) {
	font := statikFS.CaptchaFont

	// Set custom font information
	gc.FontCache = draw2d.NewSyncFolderFontCache("./arphic.ttf")
	gc.FontCache.Store(draw2d.FontData{Name: "Arphic Roman-Mincho Ultra JIS", Family: 0, Style: draw2d.FontStyleNormal}, font)
	gc.SetFontData(draw2d.FontData{Name: "Arphic Roman-Mincho Ultra JIS", Style: draw2d.FontStyleNormal})

	//Set DPI
	if captcha.Dpi <= 0 {
		captcha.Dpi = defaultDpi
	}
	gc.SetDPI(captcha.Dpi)

	// Set font size
	if captcha.FontSize <= 0 {
		captcha.FontSize = defaultFontSize
	}
	gc.SetFontSize(captcha.FontSize)
}
