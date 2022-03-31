package main

import (
	"fmt"
	"machine"
	"time"

	"github.com/stevegt/tca9534"
)

func main() {

	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		fmt.Println("could not configure I2C:", err)
		return
	}

	gpio := &tca9534.TCA9534{Addr: tca9534.BASE_ADDR, I2c: i2c}

	// set pins 0-3 as input, 4-7 as output
	rx, err := gpio.Config(0x0f)
	fmt.Println("configure got", rx, err)

	for {
		// toggle pins 4-7
		for pin := 4; pin <= 7; pin++ {
			defer gpio.Write(pin, false)
			err = gpio.Write(pin, true)
			fmt.Println("pin", pin, "write true got", err)
			time.Sleep(time.Second)
			bit, err := gpio.Read(pin)
			fmt.Println("pin", pin, "read got", bit, err)
			time.Sleep(time.Second)
			err = gpio.Write(pin, false)
			fmt.Println("pin", pin, "write false got", err)
			time.Sleep(time.Second)
			bit, err = gpio.Read(pin)
			fmt.Println("pin", pin, "read got", bit, err)
			time.Sleep(time.Second)
		}

		// show bit array for pins 0-7 for 20 seconds
		//
		// There don't appear to be any pullups on the SparkFun DEV-17047
		// board, so reading a floating input pin may show arbitrary
		// or last value.
		for i := 0; i < 200; i++ {
			bits, err := gpio.Get()
			fmt.Printf("pins 7-0 %08b err %v\n", bits, err)
			time.Sleep(100 * time.Millisecond)
		}
	}

}
