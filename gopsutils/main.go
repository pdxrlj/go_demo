package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

func main() {
	stats, err := cpu.Info()
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", stats)
	fmt.Println(strings.Repeat("=", 150))
	misc, err := load.Misc()
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", misc)
	fmt.Println(strings.Repeat("=", 150))

	temperatures, err := host.SensorsTemperatures()
	if err != nil {
		return
	}
	for _, temperature := range temperatures {
		if temperature.Temperature > 0 {
			fmt.Printf("name:%v val:%+v\n", temperature.SensorKey, temperature)
		}
	}
	fmt.Println(strings.Repeat("=", 150))
}
