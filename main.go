package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

type colorName struct {
	name string
}

type config struct {
	// Three main colors
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`

	// Message colors
	Good           string `json:"good"`
	Caution        string `json:"caution"`
	Alarm          string `json:"alarm"`
	DefaultGood    []int  `json:"rgbgood"`    // RGB value
	DefaultCaution []int  `json:"rgbcaution"` // RGB value
	DefaultAlarm   []int  `json:"rgbalarm"`   // RGB value

	// Lucents are color with opacity
	Lucent string `json:"lucent"`

	// Shades are shades of gray.
	Shades string `json:"shades"`
}

type rgb struct {
	red   float64
	green float64
	blue  float64
}

type hsl struct {
	hue        float64
	saturation float64
	luminance  float64
}

func calcrgb(f1, f2, t float64) float64 {

	if 6*t < 1.0 {
		return f2 + (f1-f2)*6*t
	}

	if 2*t < 1.0 {
		return f1
	}

	if 3*t < 2.0 {
		return f2 + (f1-f2)*(0.666-t)*6
	}
	return f2
}

func hslToRgb(c hsl) rgb {
	var newrgb rgb
	//
	// W A R N I N G
	// THIS FUNCTION "hslToRgb" IS NOT WORKING RIGHT!!!!
	// Given a HUE higher than say 100 or 180 not sure, but all goes to hell.
	// will remove comment when fixed.

	// First check saturation,
	// Zero saturation means that hue is irrelevant, it's gray.
	// (AKA: achromatic)
	//For example H = 0, S = 0 and L = 40%, so L/100 we get 0.4,
	// 0.4 * 255 = 102, so R = 102, G = 102 and B = 102
	if c.saturation == 0 {
		newrgb.red = ((c.luminance / 100) * 255)
		newrgb.green = ((c.luminance / 100) * 255)
		newrgb.blue = ((c.luminance / 100) * 255)
		return newrgb
	}

	var formula1 float64 = 0
	c.luminance = c.luminance / 100
	c.saturation = c.saturation / 100
	if c.luminance < 0.5 {
		formula1 = c.luminance * (1 + c.saturation)
	} else {
		formula1 = (c.luminance + c.saturation) - (c.luminance * c.saturation)
	}

	var formula2 float64 = 2*c.luminance - formula1

	c.hue = c.hue / 360

	newrgb.red = calcrgb(formula1, formula2, c.hue+0.333) * 255
	newrgb.green = calcrgb(formula1, formula2, c.hue) * 255
	newrgb.blue = calcrgb(formula1, formula2, c.hue-0.333) * 255

	return newrgb
}

func rgbToHsl(from rgb) hsl {
	from.red = from.red / 255
	from.green = from.green / 255
	from.blue = from.blue / 255
	min := math.Min(from.red, math.Min(from.green, from.blue))
	max := math.Max(from.red, math.Max(from.green, from.blue))

	var hue float64 = 0.0
	var saturation = 0.0

	// set luminance
	var luminance = (max + min) / 2

	var delta = (max - min)

	if max == min { // it's gray, no hue or saturation
		hue = 0
		saturation = 0
	} else {
		if luminance > 0.5 {
			saturation = (delta / (2 - max - min))
		} else {
			saturation = delta / (max + min)
		}

		switch max {
		case from.red:
			hue = (from.green - from.blue) / delta

		case from.green:
			hue = (from.blue-from.red)/delta + 2

		case from.blue:
			hue = (from.red-from.green)/delta + 4

		}

		hue = math.Round(hue * 60)
		saturation = (saturation * 100)
		luminance = (luminance * 100)

	}

	if hue < 0 {
		hue = hue + 360
	}
	if hue > 360 {
		hue = hue - 360
	}

	if saturation > 100 {
		saturation = 100
	}

	if luminance > 100 {
		luminance = 100
	}

	return hsl{hue, saturation, luminance}
}

func outputHSL(name string, i int, c hsl) {
	fmt.Printf("	--%s%02d: hsl(%v, %v%%, %v%%); \n", name, i,
		int(math.Round(c.hue)), int(math.Round(c.saturation)), int(math.Round(c.luminance)))
}
func outputHSLA(name string, i int, c hsl, op float64) {
	fmt.Printf("	--%s%02d: hsl(%v, %v%%, %v%%, %1.1f); \n", name, i,
		int(math.Round(c.hue)), int(math.Round(c.saturation)), int(math.Round(c.luminance)), op)
}
func outputRGB(name string, c rgb) {
	fmt.Printf("  --%s: #%02x%02x%02x;\n", name,
		int(math.Round(c.red)), int(math.Round(c.green)), int(math.Round(c.blue)))
}

func randRGB() rgb {
	rand.Seed(time.Now().Unix())
	x := float64(rand.Intn(10))
	rand.Seed(time.Now().Unix())
	green := float64(rand.Intn(255))
	blue := float64(rand.Intn(255))
	rand.Seed(time.Now().Unix())
	red2 := float64(rand.Intn(245))
	x = x + red2
	return rgb{x, green, blue}
}

func newHsl(c hsl, delta int) hsl {
	newColor := hsl{c.hue, c.saturation, c.luminance}
	newColor.hue = float64(delta) + newColor.hue
	if newColor.hue < 0 {
		newColor.hue = newColor.hue + 360
	}
	if newColor.hue > 360 {
		newColor.hue = newColor.hue - 360
	}
	return newColor
}
func outputShades(color hsl, name, prefix, lucentname string) {
	// 12 Shades of a color
	x := 20
	for i := 0; i < 12; i++ {
		color.luminance = float64(x)
		outputHSL(name, i+1, color)
		x = x + 6
	}

	// 3 Lumens of a color
	color.luminance = float64(50)
	for i := 0; i < 3; i++ {
		x := float64((float64(i) + float64(1.0)) / float64(10))
		outputHSLA(prefix+lucentname, i+1, color, x)
	}
}

func outputGrays(grayname, lucentname string) {
	// 12 Shades of Gray
	for i := 0; i < 9; i++ {
		x := (i + 1) * 12
		if x > 100 {
			x = 100
		}
		outputHSL(grayname, i+1, hsl{float64(0), float64(0), float64(x)})
	}

	// 8 Lumens of Gray
	for i := 0; i < 8; i++ {
		x := float64((float64(i) + float64(1.0)) / float64(8))
		outputHSLA(lucentname, i+1, hsl{float64(0), float64(0), float64(50)}, x)
	}
}

func getRgb() rgb {
	rand.Seed(time.Now().Unix())
	// variables declaration
	var red int
	var green int
	var blue int

	// flags declaration using flag package
	flag.IntVar(&red, "red", -1, "Specify RGB red value, default is 255")
	flag.IntVar(&green, "green", -1, "Specify RGB green value, default is 255")
	flag.IntVar(&blue, "blue", -1, "Specify RGB blue value, default is 255")

	flag.Parse() // after declaring flags we need to call it

	randomRGB := randRGB()
	//fmt.Printf("%v %v %v\n", randomRGB.red, randomRGB.green, randomRGB.blue)

	// check if cli params match
	if red > 255 || red < 0 {
		fmt.Printf("/* red generated randomly. */\n")
		red = int(randomRGB.red)
	}
	if green > 255 || green < 0 {
		fmt.Printf("/* green generated randomly. */\n")
		green = int(randomRGB.green)
	}
	if blue > 255 || blue < 0 {
		fmt.Printf("/* blue generated randomly.*/ \n")
		blue = int(randomRGB.blue)
	}
	return newRgb(red, green, blue)
}

func newRgb(r int, g int, b int) rgb {
	return rgb{float64(r), float64(g), float64(b)}

}

func readConf() (config, error) {
	// json data
	var obj config

	// read file
	data, err := ioutil.ReadFile("./rootcss.json")
	if err != nil {
		fmt.Print(err)
		return obj, err
	}

	// Unmarshal json.
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		return obj, err
	}

	return obj, nil
}

func main() {
	cfg, err := readConf()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}

	fmt.Println(":root {")
	fmt.Println("	--round-radius: 170960px;")

	// Every design must have some gray-scale
	outputGrays(cfg.Shades, cfg.Lucent)

	// This is the primary color, randomly generated or from the cmd line.
	rgb := getRgb()
	outputRGB("basecolor", rgb)

	// Convert to HSL, easier to manipulate.
	primary1 := rgbToHsl(rgb)
	outputShades(primary1, cfg.Primary, cfg.Primary[0:3], cfg.Lucent)

	// Create a secondary color by altering the primary hue 30 degrees.
	secondary := newHsl(primary1, 30)
	outputShades(secondary, cfg.Secondary, cfg.Secondary[0:3], cfg.Lucent)

	// Create a tertiary color by altering the primary hue 195 degrees.
	// Why 195? Halfway between primary & secondary is +15 from primary.
	// 180 degrees from there is 195 from primary.
	tertiary := newHsl(primary1, 195)
	outputShades(tertiary, cfg.Tertiary, cfg.Tertiary[0:3], cfg.Lucent)

	// The following colors are for messaging:
	//	> Good messages
	//  > Caution messages
	//  > Alarm messages.
	// I decided on using traffic light colors, you may choose other colors.

	// To conform better,
	// message colors will have the same saturation as the primary color.
	good := rgbToHsl(newRgb(cfg.DefaultGood[0],
		cfg.DefaultGood[1],
		cfg.DefaultGood[2]))
	good.saturation = primary1.saturation
	outputShades(good, cfg.Good, cfg.Good[0:3], cfg.Lucent)

	caution := rgbToHsl(newRgb(cfg.DefaultCaution[0],
		cfg.DefaultCaution[1],
		cfg.DefaultCaution[2]))
	caution.saturation = primary1.saturation
	outputShades(caution, cfg.Caution, cfg.Caution[0:3], cfg.Lucent)

	alarm := rgbToHsl(newRgb(cfg.DefaultAlarm[0],
		cfg.DefaultAlarm[1],
		cfg.DefaultAlarm[2]))
	alarm.saturation = primary1.saturation
	outputShades(alarm, cfg.Alarm, cfg.Alarm[0:3], cfg.Lucent)

	fmt.Println("}")

}
