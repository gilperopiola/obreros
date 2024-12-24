package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/obreros/core"
	"github.com/gilperopiola/obreros/toolbox"

	"go.uber.org/zap"
)

type App struct {
	*core.Config
	*toolbox.Toolbox
	*time.Ticker

	Obreros []Obrero // Workers
}

var _ core.Toolbox = &toolbox.Toolbox{}

// This will be called by main.go on init.
func NewApp() (runAppFunc, cleanUpFunc) {

	app := &App{
		Config:  &core.Config{},     // 🆕🗺️
		Toolbox: &toolbox.Toolbox{}, // 🆕🛠️
		Obreros: []Obrero{
			NewObreroDeArgenpills(), // 🆕🚧💊
		},
	}

	firstLog("Welcome to Obreros! 🚧🔧🔨", 500).thenLog("Where dreams. come. true.", 300).thenLog("Ahre xd weno cargo la config bancame", 0)

	func() {
		app.Config = core.LoadConfig()
		core.SetupLogger(&app.LoggerCfg)
	}()
	firstLog("FUuUuAA!! Ya toy re configurado corte así re Ok wachín!", 250).thenLog("Donde están las wachas? 🤔", 700)
	firstLog("No hay wachas? Bueno, entonces vamos a cargar el Toolbox 🛠️", 500)
	func() {
		app.Toolbox = toolbox.Setup(app.Config)
	}()

	firstLog("Bueno Toolbox 10/10 working ahora vienen unas boludeces y dsp ya corremos los obreros", 0)
	func() {
		app.Toolbox.AddCleanupFunc(app.CloseDB)
		app.Toolbox.AddCleanupFuncWithErr(zap.L().Sync)
	}()

	return app.Run, app.Toolbox.Cleanup
}

func (a *App) Run() {
	a.Ticker = time.NewTicker(time.Duration(20) * time.Second)
	defer a.Ticker.Stop()

	firstLog("Running Obreros in 3...", 1000).thenLog("2...", 1000).thenLog("1...", 1000)

	a.Tick()

	for range a.Ticker.C {
		a.Tick()
	}
}

func (a *App) Tick() {
	for _, obrero := range a.Obreros {
		zap.L().Info("🚧 Obrero trabajando 🚧")
		go obrero.Laburar(god.NewCtx(), a.DBTool)
	}
}

// This down here is madness. It's a chainable logger. It's a joke. It's a joke that got out of hand.
// It's a joke that got out of hand and now it's a chainable logger.
// Aaaand a joke.
// Get it? It's a chainable hand of loggers, that's a GOT OUT Of-----
// It's a chainable logger!
// That's a logger.
// Madness.

type chainLogger interface {
	thenLog(msg string, waitMs int) chainLogger
}

type chainLogFn func(msg string, waitMs int) chainLogger

func (chain chainLogFn) thenLog(msg string, waitMs int) chainLogger {
	return firstLog(msg, waitMs)
}

// Basically this returns a chainLogFn that can be called with .thenLog() to chain more logs.
// Madness.
func firstLog(msg string, waitMs int) chainLogger {
	log.Println("🌟 " + msg)
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
	return new(chainLogFn)
}

type Obrero interface {
	Laburar(ctx god.Ctx, dbTool core.DBTool)
}

var _ Obrero = &ObreroDeArgenpills{}

type ObreroDeArgenpills struct {
	URLs       []string
	HTTPClient *http.Client
}

func NewObreroDeArgenpills() *ObreroDeArgenpills {
	return &ObreroDeArgenpills{
		URLs:       []string{"https://argenpills.org/forumdisplay.php?fid=4"},
		HTTPClient: &http.Client{},
	}
}

func (o ObreroDeArgenpills) Laburar(ctx god.Ctx, dbTool core.DBTool) {
	for _, url := range o.URLs {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			zap.L().Error("Argenpills - error creating request", zap.Error(err))
		}

		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("Accept-Language", "en-US,en;q=0.9")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("Cookie", "agmybb[lastvisit]=1720928337; agsid=2d4320c6a3d168677f344710eb8130d8; agmybb[lastactive]=1720928387; agloginattempts=1; agmybbuser=17017_pdNR3InY7jMnMRNTll23hbCR8YPfqEHB6ilkuC3TwuSMeLT8aD; agmybb[announcements]=0")
		req.Header.Add("Referer", "https://argenpills.org/")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-Site", "same-origin")
		req.Header.Add("Sec-Fetch-User", "?1")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
		req.Header.Add("sec-ch-ua", `"Google Chrome";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`)
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sec-ch-ua-platform", `"Windows"`)

		resp, err := o.HTTPClient.Do(req)
		if err != nil {
			zap.L().Error("Argenpills - error sending request", zap.Error(err))
		}
		defer resp.Body.Close()

		status := resp.Status
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("Argenpills - error reading response", zap.Error(err))
		}

		title, _ := god.GetSubstringBetween("<title>", "</title>", string(body))

		zap.L().Info("Argenpills",
			zap.String("url", url),
			zap.String("title", title),
			zap.String("status", status),
		)

		dbTool.InsertWebpage(ctx, url, title, string(body))
	}
	zap.L().Info("🚧 Argenpills terminó su jornada 🚧")
}

// NewApp returns a runAppFunc and a cleanUpFunc - so the caller can first run
// the Workers and then release gracefully all used resources when it's done.
type runAppFunc func()
type cleanUpFunc func()
