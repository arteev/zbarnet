package zbar

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/arteev/zbarnet/barcode"
	"github.com/arteev/zbarnet/config"
	"github.com/arteev/zbarnet/logger"
)

//Capture barcode from stdout
type capture struct {
	cfg         *config.Config
	ChanBarCode chan barcode.BarCode
}

//Write it io.Writer
func (c *capture) Write(p []byte) (n int, err error) {
	logger.Trace.Println("Capture write start")
	logger.Debug.Printf("Data: %q\n", p)
	defer logger.Trace.Println("Capture write done")
	bc, err := barcode.MustBarCode(p, c.cfg.ZBar.GetMode())
	logger.Debug.Println("New barcode:", bc)
	if err != nil {
		if err != barcode.ErrSkipRootTag {
			logger.Warn.Printf("Failed parse barcode:%v \n", err)
			fmt.Println(err)
		}
	} else {
		if c.cfg.Output == "raw" {
			logger.Trace.Println("Capture write raw output")
			fmt.Print(string(p))
		}
		c.ChanBarCode <- *bc
	}
	return len(p), nil
}

//Run zbarcam and waiting incoming barcode
func Run(chb chan barcode.BarCode, kill chan bool, cfg *config.Config) error {
	logger.Trace.Println("Run start")
	defer logger.Trace.Println("Run end")
	args := append(cfg.ZBar.Args, cfg.ZBar.Device)
	done := make(chan bool)
	defer close(done)
	stop := false
	for {
		logger.Debug.Printf("Run cmd %s args:%v", cfg.ZBar.Location, args)
		cmd := exec.Command(cfg.ZBar.Location, args[:]...)
		cmd.Stdout = &capture{ChanBarCode: chb, cfg: cfg}
		stderr := &bytes.Buffer{}
		cmd.Stderr = io.Writer(stderr)
		if err := cmd.Start(); err != nil {
			logger.Error.Println(err)
			return err
		}
		//waiting chan kill for correct close zbarcam
		go func() {
			for {
				select {
				case <-kill:
					logger.Trace.Println("Kill zbarcam")
					if e := cmd.Process.Kill(); e != nil {
						logger.Warn.Println(e)
					}
					stop = true
					break
				case <-done:
					logger.Trace.Println("Chan done")
					return
				}
			}
		}()
		if e := cmd.Wait(); e != nil {
			logger.Warn.Println(e)
		}
		done <- true
		if stop {
			break
		}
		// pause if zbarcam run or not exists etc.
		if !cmd.ProcessState.Success() {
			time.Sleep(time.Second * 3)
			logger.Error.Println(stderr.String())
		}
	}
	return nil
}
