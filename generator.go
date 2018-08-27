package kyBlockies

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"bytes"
)

var randSeed = [4]int64{}

/*
	Generate Icon Bytes with target address
	@param address: the target eth address needed to create identicon
    @return target image bytes
 */
func CreateIconBytes(address string) []byte {
	options := NewBuildOption(address)
	data := CreateImageData(options.size)
	return createImage(data, options.color, options.bgColor, options.spotColor, 16)
}


/*
	Generate Icon png and save to target path
	@param path:  full path to save png
	@param address: the target eth address needed to create identicon
 */
//========================================================================
//
//========================================================================
func CreateIconPng(path string, address string) {
	options := NewBuildOption(address)
	data := CreateImageData(options.size)
	createPngAndSaveToPath(path, data, options.color, options.bgColor, options.spotColor, 16)
}

type IconOptions struct {
	size      int
	scale     int
	color     color.Color
	bgColor   color.Color
	spotColor color.Color
}

func NewBuildOption(seed string) IconOptions {
	seedRand(seed)
	return IconOptions{
		size:      8,
		scale:     4,
		color:     createColor(),
		bgColor:   createColor(),
		spotColor: createColor(),
	}
}

func createColor() color.Color {
	h := math.Floor(rand() * 360)
	s := (rand() * 60) + 40
	l := (rand() + rand() + rand() + rand()) * 25
	return HSL{
		H: uint32(h),
		S: uint32(s),
		L: uint32(l),
	}
}

func seedRand(seed string) {
	for i := 0; i < len(randSeed); i++ {
		randSeed[i] = 0
	}
	max := int32(1<<31 - 1)
	min := int32(-1 << 31)
	for i := range seed {
		var t1 int64
		t := int64(randSeed[i%4] << 5)
		if t > int64(max<<1) || t < int64(min<<1) {
			t1 = int64(int32(t))
		} else {
			t1 = t
		}
		t2 := t1 - int64(randSeed[i%4])
		randSeed[i%4] = t2 + int64([]rune(seed)[i])
	}
	for i := range randSeed {
		randSeed[i] = int64(int32(randSeed[i]))
	}
}

func CreateImageData(size int) []float64 {
	width := size
	height := size

	dataWidth := math.Ceil(float64(width / 2))
	mirrorWidth := width - int(dataWidth)
	data := []float64{}
	for i := 0; i < height; i++ {
		row := []float64{}
		r := []float64{}
		for j := 0; j < int(dataWidth); j++ {
			row = append(row, math.Floor(rand()*2.3))
			if j < mirrorWidth {
				r = append(r, row[j])
			}
		}
		r = reverse(r)
		row = append(row, r...)
		data = append(data, row...)
	}
	return data
}

func reverse(target []float64) []float64 {
	for i, j := 0, len(target)-1; i < j; i, j = i+1, j-1 {
		target[i], target[j] = target[j], target[i]
	}
	return target
}

func rand() float64 {
	var t int32
	t = int32(randSeed[0] ^ (randSeed[0] << 11))
	randSeed[0] = randSeed[1]
	randSeed[1] = randSeed[2]
	randSeed[2] = randSeed[3]
	randSeed[3] = randSeed[3] ^ (randSeed[3] >> 19) ^ int64(t) ^ int64(t>>8)
	t1 := math.Abs(float64(randSeed[3]))
	return float64(t1 / (1<<31 - 1))
}

func createPngAndSaveToPath(fullPath string ,data []float64, c color.Color, bgColor color.Color, spotColor color.Color, scale int) {
	width := int(math.Sqrt(float64(len(data))))
	w := width * scale
	h := width * scale
	imgFile, _ := os.Create(fullPath)
	defer imgFile.Close()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(data); i++ {
		row := int(math.Floor(float64(i / width)))
		col := i % width
		tc := spotColor
		if data[i] == 1 {
			tc = c
		}
		if data[i] > 0 {
			for x := col * scale; x < col*scale+scale; x++ {
				for y := row * scale; y < row*scale+scale; y++ {
					img.Set(x, y, tc)
				}
			}
		} else {
			for x := col * scale; x < col*scale+scale; x++ {
				for y := row * scale; y < row*scale+scale; y++ {
					img.Set(x, y, bgColor)
				}
			}
		}
	}
	png.Encode(imgFile, img)
}

func createImage(data []float64, c color.Color, bgColor color.Color, spotColor color.Color, scale int) []byte {
	width := int(math.Sqrt(float64(len(data))))
	w := width * scale
	h := width * scale
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(data); i++ {
		row := int(math.Floor(float64(i / width)))
		col := i % width
		tc := spotColor
		if data[i] == 1 {
			tc = c
		}
		if data[i] > 0 {
			for x := col * scale; x < col*scale+scale; x++ {
				for y := row * scale; y < row*scale+scale; y++ {
					img.Set(x, y, tc)
				}
			}
		} else {
			for x := col * scale; x < col*scale+scale; x++ {
				for y := row * scale; y < row*scale+scale; y++ {
					img.Set(x, y, bgColor)
				}
			}
		}
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

//========================================================================
// HSL
//========================================================================

const (
	// hueMax is the maximum allowed value for Hue in the HSL color model.
	hueMax = 360
	// saturationMax is the maximum allowed value for Saturation in the HSL
	// color model.
	saturationMax = 100
	// lightnessMax is the maximum allowed value for lightnessMax in the HSL
	// color model.
	lightnessMax = 100
	// rgbaMax is the maximum allowed value for any R, G, B, A property value.
	rgbaMax = 255
)

// HSL is a color model representation based on RGB. HSL facilitates the
// generation of colors that look similar between themselves by changing the
// value of Hue H while keeping Saturation S and Lightness L the same.
type HSL struct {
	// Hue [0, 360]
	H uint32
	// Saturation [0, 100]
	S uint32
	// Lightness [0, 100]
	L uint32
}

// RGBA conversion
func (hsl HSL) RGBA() (r, g, b, a uint32) {
	h := 1.0 / float64(hueMax) * float64(hsl.H)
	s := float64(hsl.S) / float64(saturationMax)
	l := float64(hsl.L) / float64(lightnessMax)
	r, g, b = hslToRgb(h, s, l)
	a = rgbaMax
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a |= a << 8
	return
}

func hslToRgb(h, s, l float64) (uint32, uint32, uint32) {
	var q, p float64
	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = (l + s) - (l * s)
		}
		p = (2 * l) - q
		r = hueToRgb(p, q, h+(1.0/3.0))
		g = hueToRgb(p, q, h)
		b = hueToRgb(p, q, h-(1.0/3.0))
	}

	return uint32(r * rgbaMax), uint32(g * rgbaMax), uint32(b * rgbaMax)
}

func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	} else if t > 1 {
		t--
	}
	switch {
	case 6*t < 1:
		return p + (q-p)*6*t
	case 2*t < 1:
		return q
	case 3*t < 2:
		return p + (q-p)*((2.0/3.0)-t)*6
	}
	return p
}
