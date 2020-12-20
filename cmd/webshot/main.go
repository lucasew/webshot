package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

var pwd string

type Settings struct {
    InputFile string
    Width int64
    Height int64
    OutputFile string
    Timeout int
    DevtoolsAddr string
}

func (s Settings) LaunchChrome(ctx context.Context) error {
    alternatives := []string{"chrome", "google-chrome", "google-chrome-stable"}
    var chosen string
    for _, alternative := range alternatives {
        _, err := exec.LookPath(alternative)
        if err == nil {
            chosen = alternative
            break
        }
    }
    if chosen == "" {
        return fmt.Errorf("chrome not found")
    }
    go func() {
        cmd := exec.Command(chosen, "--remote-debugging-port=9222")
        cmd.Run()
        <-ctx.Done()
        cmd.Process.Kill()
    }()
    return nil
}

func (s Settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Cache-Control", "no-cache")
    if (r.URL.Path == "/") {
        log.Printf("Loading '%s'", s.InputFile)
        http.ServeFile(w, r, s.InputFile)
    } else {
        url := filepath.Join(pwd, r.URL.Path[1:])
        log.Printf("Loading '%s'", url)
        http.ServeFile(w, r, url)
    }
}

var Options Settings

func init() {
    flag.StringVar(&Options.InputFile, "i", "index.html", "file to render")
    flag.StringVar(&Options.OutputFile, "f", "out.png", "where to save the png screenshot")
    flag.IntVar(&Options.Timeout, "t", 30, "timeout to do the actions")
    flag.StringVar(&Options.DevtoolsAddr, "d", "http://localhost:9222", "chrome devtools address for remote control")
    flag.Int64Var(&Options.Width, "w", 1280, "viewport width")
    flag.Int64Var(&Options.Height, "h", 720, "viewport height")
    flag.Parse()

    var err error
    Options.InputFile, err = filepath.Rel(pwd, Options.InputFile)
    handleErr(err)
    pwd, err = os.Getwd()
    handleErr(err)
    log.Printf("Current directory: %s", pwd)
    go http.ListenAndServe(":49490", Options)
}

func handleErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    log.Printf("Setting up...")
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Options.Timeout)*time.Second)
    cdpCtx, _ := chromedp.NewContext(ctx)
    defer cancel()
    var imgBuf []byte
    handleErr(chromedp.Run(cdpCtx, chromedp.Tasks{
        emulation.SetDeviceMetricsOverride(Options.Width, Options.Height, 1, false),
        chromedp.Navigate("http://localhost:49490"),
        chromedp.Reload(),
        chromedp.WaitVisible("html", chromedp.ByQuery),
        chromedp.Screenshot("html", &imgBuf, chromedp.NodeVisible),
    }))
    ioutil.WriteFile(Options.OutputFile, imgBuf, 0644)
}
