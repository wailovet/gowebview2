package gowebview2

import (
	"fmt"

	"github.com/wailovet/gofunc"
)

type HttpEngine interface {
	Run()
	SetPort(port string)
	Close()
}

func NuwaAppModelRun(httpEngine HttpEngine, w, h int) error {
	port := getFreePort()
	url := fmt.Sprint("http://127.0.0.1:", port)

	webview := NewWithOptions(WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: WindowOptions{
			IconId: 2, // icon resource id
			Center: true,
		},
	})
	if webview == nil {
		return fmt.Errorf("new webview failed")
	}
	defer func() {
		webview.Destroy()
		httpEngine.Close()
	}()

	httpEngine.SetPort(fmt.Sprint(port))
	gofunc.New(func() {
		httpEngine.Run()
	})

	webview.Navigate(url)
	webview.Run()
	return nil
}
