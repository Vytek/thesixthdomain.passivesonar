package tsd_passivesonar

import (
	"errors"
	"math"
)

// Version exposes the current package version.
const Version = "0.0.2"

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

var (
	freshWaterMinTemp = 4.0
	freshWaterMaxTemp = 86.0
	freshWaterSoundSpeed = map[float64]float64{
		4.0: 1421.62, 17.5: 1474.38, 39.2: 1421.62, 63.5: 1474.38,
		4.5: 1423.90, 18.0: 1476.01, 40.1: 1423.90, 64.4: 1476.01,
		5.0: 1426.15, 18.5: 1477.62, 41.0: 1426.15, 65.3: 1477.62,
		5.5: 1428.38, 19.0: 1479.21, 41.9: 1428.38, 66.2: 1479.21,
		6.0: 1430.58, 19.5: 1480.77, 42.8: 1430.58, 67.1: 1480.77,
		6.5: 1432.75, 20.0: 1482.32, 43.7: 1432.75, 68.0: 1482.32,
		7.0: 1434.90, 20.5: 1483.84, 44.6: 1434.90, 68.9: 1483.84,
		7.5: 1437.02, 21.0: 1485.35, 45.5: 1437.02, 69.8: 1485.35,
		8.0: 1439.12, 21.5: 1486.83, 46.4: 1439.12, 70.7: 1486.83,
		8.5: 1441.19, 22.0: 1488.29, 47.3: 1441.19, 71.6: 1488.29,
		9.0: 1443.23, 22.5: 1489.74, 48.2: 1443.23, 72.5: 1489.74,
		9.5: 1445.25, 23.0: 1491.16, 49.1: 1445.25, 73.4: 1491.16,
		10.0: 1447.25, 23.5: 1492.56, 50.0: 1447.25, 74.3: 1492.56,
		10.5: 1449.22, 24.0: 1493.95, 50.9: 1449.22, 75.2: 1493.95,
		11.0: 1451.17, 24.5: 1495.32, 51.8: 1451.17, 76.1: 1495.32,
		11.5: 1453.09, 25.0: 1496.66, 52.7: 1453.09, 77.0: 1496.66,
		12.0: 1454.99, 25.5: 1497.99, 53.6: 1454.99, 77.9: 1497.99,
		12.5: 1456.87, 26.0: 1499.30, 54.5: 1456.87, 78.8: 1499.30,
		13.0: 1458.72, 26.5: 1500.59, 55.4: 1458.72, 79.7: 1500.59,
		13.5: 1460.55, 27.0: 1501.86, 56.3: 1460.55, 80.6: 1501.86,
		14.0: 1462.36, 27.5: 1503.11, 57.2: 1462.36, 81.5: 1503.11,
		14.5: 1464.14, 28.0: 1504.35, 58.1: 1464.14, 82.4: 1504.35,
		15.0: 1465.91, 28.5: 1505.56, 59.0: 1465.91, 83.3: 1505.56,
		15.5: 1467.65, 29.0: 1506.76, 59.9: 1467.65, 84.2: 1506.76,
		16.0: 1469.36, 29.5: 1507.94, 60.8: 1469.36, 85.1: 1507.94,
		16.5: 1471.06, 30.0: 1509.10, 61.7: 1471.06, 86.0: 1509.10,
		17.0: 1472.73, 62.6: 1472.70,
	}
)

// Sound speed in fresh water. 
// Values according to https://bathylogger.com/wp-content/uploads/2015/10/Speed-of-Sound-in-Freshwater.pdf
func GetSoundSpeedInFreshWater(temp float64) (float64, error) {
	if temp > freshWaterMaxTemp || temp < freshWaterMinTemp {
		return 0, errors.New("temp is out of range")
	}

	ln := math.MaxFloat64
	rn := -math.MaxFloat64

	for k := range freshWaterSoundSpeed {
		if math.Abs(temp-k) < math.Abs(temp-ln) && temp > k {
			ln = k
		}
		if math.Abs(temp-k) < math.Abs(temp-rn) && temp < k {
			rn = k
		}
	}

	ssln := freshWaterSoundSpeed[ln]
	ssrn := freshWaterSoundSpeed[rn]

	return ssln + (ssrn-ssln)*(temp-ln)/(rn-ln), nil
}

// CalcSoundSpeed calculates speed of sound in water according to Wilson formula
// See: https://rbr-global.com/speed-of-sound-in-water/
func CalcSoundSpeed(t, p, s float64) (float64, error) {
	if t < -4 || t > 30 {
		return 0, errors.New("t is out of range (-4 to 30)")
	}

	if p < 0.1 || p > 100 {
		return 0, errors.New("p is out of range (0.1 to 100 MPa)")
	}

	if s < 0 || s > 40 {
		return 0, errors.New("s is out of range (0 to 40)")
	}

	c0 := 1449.14
	Dct := 4.5721*t - 4.4532E-2*math.Pow(t, 2) - 2.6045E-4*math.Pow(t, 3) + 7.9851E-6*math.Pow(t, 4)
	Dcs := 1.39799*(s-35) - 1.69202E-3*math.Pow(s-35, 2)
	Dcp := 1.63432*p - 1.06768E-3*math.Pow(p, 2) + 3.73403E-6*math.Pow(p, 3) - 3.6332E-8*math.Pow(p, 4)
	Dcstp := (s-35)*(-1.1244E-2*t+7.7711E-7*math.Pow(t, 2)+7.85344E-4*p-
		1.3458E-5*math.Pow(p, 2)+3.2203E-7*p*t+1.6101E-8*math.Pow(t, 2)*p)+
		p*(-1.8974E-3*t+7.6287E-5*math.Pow(t, 2)+4.6176E-7*math.Pow(t, 3))+
		math.Pow(p, 2)*(-2.6301E-5*t+1.9302E-7*math.Pow(t, 2))+
		math.Pow(p, 3)*(-2.0831E-7*t)

	result := c0 + Dct + Dcs + Dcp + Dcstp

	return result, nil
}