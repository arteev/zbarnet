package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arteev/zbarnet/actions/httpclient"
	"github.com/arteev/zbarnet/barcode"
	"github.com/arteev/zbarnet/config"
	"github.com/arteev/zbarnet/logger"
	"github.com/arteev/zbarnet/svc"
	"github.com/arteev/zbarnet/version"
	"github.com/arteev/zbarnet/zbar"
	"github.com/codegangsta/cli"
	"github.com/kardianos/service"
)

//A Application cli
type Application struct {
	terminated   bool
	done         chan bool
	chanKillZBar chan bool
	svc          service.Service
	cli          *cli.App
	cfg          *config.Config
}

//New create new application
func New() *Application {
	app := &Application{
		cfg: config.MustConfig(""),
	}
	app.svc, _ = svc.New(&service.Config{
		Name:        "zbarnet",
		DisplayName: "ZBarNet",
		Description: "ZBarNet http client for zbarcam",
		Arguments:   []string{"-start"}},
		app.dowork)
	app.newCli()
	return app
}

func (a *Application) defineFlags() *Application {
	a.cli.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "verbose",
			Usage: "set specific debug output level",
		},
		cli.BoolFlag{
			Name:  "start,s",
			Usage: "start service(daemon) mode",
		},
	}
	return a
}

func (a *Application) beforeAction() *Application {
	a.cli.Before = func(c *cli.Context) error {
		var curLevel int
		if c.IsSet("verbose") {
			curLevel = c.Int("verbose")
		}

		ln := logger.LevelByOrder(curLevel)
		logger.Init(ln.Level, os.Stdout, os.Stdout, os.Stderr, os.Stdout, os.Stdout)
		logger.Info.Printf("Verbose level: %s", ln.Name)
		return nil
	}
	return a
}

func (a *Application) logerror(err error) {
	if err != nil {
		logger.Error.Println(err)
		//if a.Verbose {
		fmt.Println(err)
		//}
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (a *Application) defineActions() *Application {
	a.cli.Action = func(c *cli.Context) {
		if c.Bool("start") || c.Bool("s") {
			logger.Error = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
			if !service.Interactive() &&
				a.cfg.HTTP.Enabled &&
				a.cfg.ZBar.Enabled &&
				!stringInSlice("--nodisplay", a.cfg.ZBar.Args) {
				a.cfg.ZBar.Args = append(a.cfg.ZBar.Args, "--nodisplay")
			}
			a.logerror(a.svc.Run())
		} else {
			fmt.Printf("Hit Ctrl+C to [EXIT]\n")
			a.dowork()
			a.Wait()
		}
	}

	a.cli.Commands = []cli.Command{
		{
			Name:        "service",
			Aliases:     []string{"svc", "daemon"},
			Usage:       "sevice(daemon) managment",
			Description: "sevice(daemon) managment",
			Subcommands: []cli.Command{
				{
					Name:   "install",
					Usage:  "install service(daemon)",
					Action: a.handlerSvcExec("Service installed", a.svc.Install),
				},
				{
					Name:   "uninstall",
					Usage:  "uninstall service(daemon)",
					Action: a.handlerSvcExec("Service uninstalled", a.svc.Uninstall),
				},
				{
					Name:   "start",
					Usage:  "start service(daemon)",
					Action: a.handlerSvcExec("Service started", a.svc.Start),
				},
				{
					Name:   "stop",
					Usage:  "stop service(daemon)",
					Action: a.handlerSvcExec("Service stoped", a.svc.Stop),
				},
				{
					Name:   "restart",
					Usage:  "restart service(daemon)",
					Action: a.handlerSvcExec("Service restarted", a.svc.Restart),
				},
			},
		},
	}

	return a
}

func (a *Application) newCli() {
	a.cli = cli.NewApp()
	a.cli.Version = version.Version
	a.cli.Copyright = version.Copyright
	a.cli.Name = version.Name
	a.cli.Author = version.Author
	a.defineFlags().
		defineActions().
		beforeAction()
}

func (a *Application) sendOverhttp(bc *barcode.BarCode) {
	logger.Trace.Println("sendOverhttp start")
	req := httpclient.New(a.cfg.HTTP.Method,
		a.cfg.HTTP.URL,
		a.cfg.HTTP.APIKey,
		a.cfg.HTTP.APIKeyHeader)

	if a.cfg.Once {
		if err := req.Send(bc); err != nil {
			logger.Error.Println(err)
		}
	} else {
		go func() {
			if err := req.Send(bc); err != nil {
				logger.Error.Println(err)
			}
		}()
	}
	logger.Trace.Println("sendOverhttp done")
}

func (a *Application) processBarCode(bc *barcode.BarCode) {
	if a.cfg.Output == "json" {
		logger.Trace.Println("dowork write json into Stdout")
		jn, e := bc.ToJSON()
		if e == nil {
			fmt.Fprintln(os.Stdout, jn)
		} else {
			logger.Error.Println(e)
		}
	}
	if a.cfg.HTTP.Enabled {
		logger.Trace.Println("dowork send http")
		a.sendOverhttp(bc)
	}
}

func (a *Application) dowork() {
	a.terminated = false
	chb := make(chan barcode.BarCode)
	a.done = make(chan bool)
	a.chanKillZBar = make(chan bool)
	go func() {
		for {
			select {
			case bc := <-chb:
				logger.Trace.Println("dowork income chan barcode")
				logger.Debug.Println("barcode. Type:", bc.Type, "Data:", string(bc.Data))
				a.processBarCode(&bc)
				if a.cfg.Once {
					logger.Trace.Println("dowork once execution, doing done()")
					a.Done()
					return
				}
				break
			case <-a.done:
				logger.Trace.Println("dowork chan done")
				//a.Done()
				return
			}
		}
	}()

	if service.Interactive() {
		go func() {
			err := zbar.Run(chb, a.chanKillZBar, a.cfg)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		err := zbar.Run(chb, a.chanKillZBar, a.cfg)
		if err != nil {
			panic(err)
		}

	}
}

func (a *Application) handlerSvcExec(out string, fsvc func() error) func(c *cli.Context) {
	return func(c *cli.Context) {
		logger.Trace.Println("handlerSvcExec", out)
		if e := fsvc(); e != nil {
			panic(e)
		} else {
			//a.PrintVerbose(out)
			fmt.Println(out)
		}
	}
}

//Run application
func (a *Application) Run() {
	err := a.cli.Run(os.Args)
	if err != nil {
		panic(err)
	}
	logger.Trace.Printf("Run done")
}

//Done correct close application
func (a *Application) Done() {
	logger.Trace.Println("Done start")
	if a.terminated {
		return
	}
	a.chanKillZBar <- true
	logger.Trace.Println("Done sleep waiting kill zbar instance")
	time.Sleep(time.Millisecond * 200)
	close(a.chanKillZBar)
	close(a.done)
	a.terminated = true
	logger.Trace.Println("Done end")
}

//Wait while finished application
func (a *Application) Wait() {
	logger.Trace.Println("Wait start")
	<-a.done
	logger.Trace.Println("Wait end")
}
