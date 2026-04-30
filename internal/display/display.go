package display

import (
	"log"
	"time"

	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
)

const (
	i2cAddress = 0x27
	i2cBus     = 1
)

type DisplayController struct {
	ShowVerificationCodeChan chan time.Time
	IncomingCallChan         chan string
	InCallChan               chan string
	lcd                      *hd44780.Lcd
}

func NewDisplayController() DisplayController {
	// TODO: i2c connection leak -- no Close invocation
	i2cConn, err := i2c.NewI2C(i2cAddress, i2cBus)
	if err != nil {
		log.Fatal(err)
	}

	lcd, err := hd44780.NewLcd(i2cConn, hd44780.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}
	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	return DisplayController{
		ShowVerificationCodeChan: make(chan time.Time, 1),
		IncomingCallChan:         make(chan string, 1),
		InCallChan:               make(chan string, 1),
		lcd:                      lcd,
	}
}

func (c *DisplayController) EventLoop() {
	for {

	}
}
