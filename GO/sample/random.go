package sample

import (
	"math/rand"
	"time"

	// import pb package from generateProto folder
	pb "pcbook/generateProto"
	// import google uuid package
	uuid "github.com/google/uuid"
)

func init() {
	// use one seed for all random number
	rand.Seed(time.Now().UnixNano())
}

func RandomKeyboard() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	case 3:
		return pb.Keyboard_AZERTY
	default:
		return pb.Keyboard_UNKNOWN
	}
} //

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomCPUBrand() string {
	switch rand.Intn(3) {
	case 1:
		return "Intel"
	case 2:
		return "AMD"
	default:
		return "Other"
	}
}

func RandomCPUName(brand string) string {
	switch brand {
	case "Intel":
		switch rand.Intn(3) {
		case 1:
			return "i3"
		case 2:
			return "i5"
		default:
			return "i7"
		}
	case "AMD":
		switch rand.Intn(3) {
		case 1:
			return "Ryzen 3"
		case 2:
			return "Ryzen 5"
		default:
			return "Ryzen 7"
		}
	default:
		return "Other"
	}
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomGPUBrand() string {
	switch rand.Intn(3) {
	case 1:
		return "Nvidia"
	case 2:
		return "AMD"
	default:
		return "Other"
	}
}

func RandomGPUName(brand string) string {
	switch brand {
	case "Nvidia":
		switch rand.Intn(3) {
		case 1:
			return "GTX 1050"
		case 2:
			return "GTX 1060"
		default:
			return "GTX 1070"
		}
	case "AMD":
		switch rand.Intn(3) {
		case 1:
			return "Radeon 530"
		case 2:
			return "Radeon 550"
		default:
			return "Radeon 570"
		}
	default:
		return "Other"
	}
}

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomScreenPanel() pb.Screen_Panel {
	switch rand.Intn(3) {
	case 1:
		return pb.Screen_IPS
	case 2:
		return pb.Screen_OLED
	default:
		return pb.Screen_UNKNOWN
	}
}

func RandomID() string {
	return uuid.New().String()
}

func RandomLaptopBrand() string {
	switch rand.Intn(3) {
	case 1:
		return "Apple"
	case 2:
		return "Microsoft"
	default:
		return "Google"
	}
}

func RandomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		switch rand.Intn(3) {
		case 1:
			return "Macbook Pro"
		case 2:
			return "Macbook Air"
		default:
			return "Macbook"
		}
	case "Microsoft":
		switch rand.Intn(3) {
		case 1:
			return "Surface Pro"
		case 2:
			return "Surface Laptop"
		default:
			return "Surface Go"
		}
	case "Google":
		switch rand.Intn(3) {
		case 1:
			return "Pixelbook"
		case 2:
			return "Chromebook"
		default:
			return "Chromebook Pixel"
		}
	default:
		return "Other"
	}
}
