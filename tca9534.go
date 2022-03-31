package tca9534

// Tinygo driver for TCA9534 GPIO expander chip, used in e.g. SparkFun DEV-17047

import "machine"

const BASE_ADDR = 0x27
const CMD_INPUT_PORT = 0x00
const CMD_OUTPUT_PORT = 0x01
const CMD_INVERSION = 0x02
const CMD_CONFIGURATION = 0x03

type TCA9534 struct {
	Addr uint16
	I2c  *machine.I2C
}

func (t *TCA9534) xeq(cmd, tx byte) (rx byte, err error) {
	txbuf := []byte{cmd, tx}
	rxbuf := make([]byte, 1)
	err = t.I2c.Tx(t.Addr, txbuf, rxbuf)
	rx = rxbuf[0]
	return
}

func (t *TCA9534) Invert(bits byte) (rx byte, err error) {
	rx, err = t.xeq(CMD_INVERSION, bits)
	return
}

func (t *TCA9534) Put(bits byte) (rx byte, err error) {
	rx, err = t.xeq(CMD_OUTPUT_PORT, bits)
	return
}

func (t *TCA9534) Get() (rx byte, err error) {
	rx, err = t.xeq(CMD_INPUT_PORT, 0x00)
	return
}

// Config sets pin configurations to input or output. Bit 0 is pin 0,
// and so on.  A high bit sets the pin to output mode, and a low bit
// sets the pin to input mode.
func (t *TCA9534) Config(conf byte) (rx byte, err error) {
	rx, err = t.xeq(CMD_CONFIGURATION, conf)
	return
}

func (t *TCA9534) Read(pin int) (bit bool, err error) {
	bits, err := t.Get()
	if err != nil {
		return
	}
	// discard high bits
	pin &= 0x07
	// read bit value
	bit = (bits & (1 << pin)) > 0
	return
}

func (t *TCA9534) Write(pin int, bit bool) (err error) {
	bits, err := t.Get()
	if err != nil {
		return
	}
	// discard high bits
	pin &= 0x07
	if bit {
		// turn bit on
		bits |= 1 << pin
	} else {
		// turn bit off
		bits &^= 1 << pin
	}
	_, err = t.Put(bits)
	return
}
