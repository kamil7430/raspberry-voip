package display

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
)

const (
	i2cAddress   = 0x27
	i2cBus       = 1
	displayWidth = 16
)

type DisplayController struct {
	ShowVerificationCodeChan chan *ShowVerificationCodeDetails
	IncomingCallChan         chan *IncomingCallDetails
	InCallChan               chan *InCallDetails
	CallFinishedChan         chan *CallFinishedDetails
	lcd                      *hd44780.Lcd
}

type ShowVerificationCodeDetails struct {
	Time time.Time
	Code string
}

type IncomingCallDetails struct {
	DisplayName string
}

type InCallDetails struct {
	DisplayName string
	CallStart   time.Time
}

type CallFinishedDetails struct{}

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
		ShowVerificationCodeChan: make(chan *ShowVerificationCodeDetails, 1),
		IncomingCallChan:         make(chan *IncomingCallDetails, 1),
		InCallChan:               make(chan *InCallDetails, 1),
		CallFinishedChan:         make(chan *CallFinishedDetails, 1),
		lcd:                      lcd,
	}
}

func (c *DisplayController) EventLoop() {
	var svc *ShowVerificationCodeDetails
	var icc *IncomingCallDetails
	var ic *InCallDetails
	var cf *CallFinishedDetails

	for {
		// blocking receive from channels
		select {
		case svc = <-c.ShowVerificationCodeChan:
		case icc = <-c.IncomingCallChan:
			ic = nil // ic and icc are exclusive
		case ic = <-c.InCallChan:
			icc = nil // as above
		case cf = <-c.CallFinishedChan:
			icc = nil
			ic = nil
		}

		// TODO: redraw the screen every second during call
		// TODO: timeout the verification code showing
		// probably should use a go routine

		if svc != nil {
			c.drawSvc(svc, icc, ic)
		} else if icc != nil {
			c.drawIcc(icc)
		} else if ic != nil {
			c.drawIc(ic)
		} else if cf != nil {
			c.drawCf(cf)
		}
	}
}

func (c *DisplayController) showMsg(text string, line hd44780.ShowOptions) {
	if len(text) > displayWidth {
		log.Fatal("Provided text is too long")
	}
	if line != hd44780.SHOW_LINE_1 && line != hd44780.SHOW_LINE_2 {
		log.Fatal("Invalid line number")
	}

	err := c.lcd.ShowMessage(text, line)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *DisplayController) drawSvc(svc *ShowVerificationCodeDetails, icc *IncomingCallDetails, ic *InCallDetails) {
	if icc != nil || ic != nil {
		if icc != nil {
			c.showMsg(center(icc.DisplayName), hd44780.SHOW_LINE_1)
		} else {
			c.showMsg(center(ic.DisplayName), hd44780.SHOW_LINE_1)
		}

		c.showMsg(center(fmt.Sprintf(
			"Code: %s",
			svc.Code,
		)), hd44780.SHOW_LINE_2)
	} else {
		c.showMsg(center("Verify by Code:"), hd44780.SHOW_LINE_1) // 15 chars; "Verification Code" is 17 chars ;c
		c.showMsg(center(svc.Code), hd44780.SHOW_LINE_2)
	}
}

func (c *DisplayController) drawIcc(icc *IncomingCallDetails) {
	c.showMsg(center("Incoming Call:"), hd44780.SHOW_LINE_1) // 14 chars
	c.showMsg(center(icc.DisplayName), hd44780.SHOW_LINE_2)
}

func (c *DisplayController) drawIc(ic *InCallDetails) {
	duration := time.Since(ic.CallStart)
	timeString := fmt.Sprintf(
		"%03d:%02d",
		int(duration.Minutes()),
		int(duration.Seconds()),
	)

	c.showMsg(center(timeString), hd44780.SHOW_LINE_1) // 6 chars
	c.showMsg(center(ic.DisplayName), hd44780.SHOW_LINE_2)
}

func (c *DisplayController) drawCf(cf *CallFinishedDetails) {
	c.showMsg(center("Call"), hd44780.SHOW_LINE_1)
	c.showMsg(center("Finished"), hd44780.SHOW_LINE_2)
}

func center(text string) string {
	spacesCount := (displayWidth - len(text)) / 2
	spaces := strings.Repeat(" ", spacesCount)
	return spaces + text + spaces
}
