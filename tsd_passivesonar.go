package tsd_passivesonar

import "math"

// Version exposes the current package version.
const Version = "0.0.1"

// https://www.betasom.it/forum/index.php?/topic/49098-sonar-computazioni-del-livello-del-rumore-del-mare-in-visual-basic/
// Converted to Golang
func NL(x float64, ss string) (NL float64) {
	
	// calcolo del livello spettrale del rumore del mare
	k := 5 / (20 * math.Log(2) / math.Log(10))
	Y := 20 * math.Log(math.Pow(x, k)) / math.Log(10)

	var livdB float64
	switch ss {
	case "SS=0":
		livdB = 55 - Y - 10.8 // per ss=0
	case "SS=1/2":
		livdB = 55 - Y - 4.7 // per ss=1/2
	case "SS=1":
		livdB = 55 - Y // per ss=1
	case "SS=2":
		livdB = 55 - Y + 6.8 // per ss=2
	case "SS=4":
		livdB = 55 - Y + 11.6 // per ss=4
	case "SS=6":
		livdB = 55 - Y + 15 // per ss=6
	}

	//Label1.Caption = fmt.Sprintf("%.1f", livdB)
	return livdB
}