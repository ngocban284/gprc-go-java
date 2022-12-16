package sample

import (
	// import package pb from generateProto folder
	// "pcbook/generateProto"
	pb "pcbook/generateProto"

	// import ptypes package
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  RandomKeyboard(),
		Backlit: RandomBool(),
	}

	return keyboard
}

func NewCPU() *pb.CPU {
	brand := RandomCPUBrand()
	name := RandomCPUName(brand)
	numberCores := RandomInt(2, 8)
	numberThreads := RandomInt(numberCores, 12)
	minGhz := RandomFloat(2.0, 3.5)
	maxGhz := RandomFloat(minGhz, 5.0)

	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}

	return cpu
}

func NewGPU() *pb.GPU {
	brand := RandomGPUBrand()
	name := RandomGPUName(brand)
	minGhz := RandomFloat(1.0, 1.5)
	maxGhz := RandomFloat(minGhz, 2.0)

	memory := &pb.Memory{
		Value: uint64(RandomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}

	return gpu
}

func NewRam() *pb.Memory {
	ram := &pb.Memory{
		Value: uint64(RandomInt(4, 32)),
		Unit:  pb.Memory_GIGABYTE,
	}

	return ram
}

func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(RandomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	return ssd
}

func NewHDD() *pb.Storage {
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(RandomInt(500, 2000)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	return hdd
}

func NewScreen() *pb.Screen {
	screen := &pb.Screen{
		Resolution: &pb.Screen_Resolution{
			Width:  uint32(RandomInt(1920, 3840)),
			Height: uint32(RandomInt(1080, 2160)),
		},
		SizeInch:   RandomFloat32(13.3, 17.3),
		Panel:      RandomScreenPanel(),
		Multitouch: RandomBool(),
	}

	return screen
}

func NewLaptop() *pb.Laptop {
	brand := RandomLaptopBrand()
	name := RandomLaptopName(brand)

	laptop := &pb.Laptop{
		Id:       RandomID(),
		Brand:    brand,
		Name:     name,
		PriceUsd: RandomFloat(1500, 3000),
		Cpu:      NewCPU(),
		Ram:      NewRam(),
		Gpus:     []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: RandomFloat(1.2, 2.3),
		},
		ReleaseYear: uint32(RandomInt(2015, 2020)),
		UpdatedAt:   timestamppb.Now(),
	}

	return laptop
}
