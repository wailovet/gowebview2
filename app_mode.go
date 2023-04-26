package gowebview2

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "embed"
)

type StorageInfterface interface {
	GetItem(key string) string
	SetItem(key string, value string)
	RemoveItem(key string)
	Clear()
}

type FileStorage struct {
	filename string
	lock     sync.RWMutex
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func NewFileStorage(filename string) *FileStorage {

	return &FileStorage{
		filename: filename,
	}
}

func (f *FileStorage) GetItem(key string) string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if !fileExists(f.filename) {
		ioutil.WriteFile(f.filename, []byte("{}"), 0644)
	}
	data, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return "null"
	}

	var m map[string]string

	err = json.Unmarshal(data, &m)
	if err != nil {
		return "null"
	}

	return m[key]
}

func (f *FileStorage) SetItem(key string, value string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if !fileExists(f.filename) {
		ioutil.WriteFile(f.filename, []byte("{}"), 0644)
	}
	data, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return
	}

	var m map[string]string

	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}

	m[key] = value

	data, err = json.Marshal(m)

	err = ioutil.WriteFile(f.filename, data, 0644)

	if err != nil {
		return
	}
}

func (f *FileStorage) RemoveItem(key string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	data, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return
	}

	var m map[string]string

	err = json.Unmarshal(data, &m)

	if err != nil {

		return
	}

	delete(m, key)

	data, err = json.Marshal(m)

	err = ioutil.WriteFile(f.filename, data, 0644)

	if err != nil {

		return
	}
}

func (f *FileStorage) Clear() {
	f.lock.Lock()
	defer f.lock.Unlock()

	err := ioutil.WriteFile(f.filename, []byte("{}"), 0644)

	if err != nil {

		return
	}
}

type AppMode struct {
	loadAppFilename string
	srcFilesystem   fs.FS
	srcDirSub       string
	storage         StorageInfterface
}

func NewAppMode(loadAppFilename string) (*AppMode, error) {
	return &AppMode{
		loadAppFilename: loadAppFilename,
	}, nil
}

func NewAppModeWithMemory(fs fs.FS, sub string) (*AppMode, error) {

	return &AppMode{
		srcDirSub:     sub,
		srcFilesystem: fs,
	}, nil

}

func getFreePort() int {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

//go:embed static/vue.js
var vueJs string

//go:embed static/jquery.js
var jqueryJs string

//go:embed static/right-menu.js
var rightMenuJs string

//go:embed static/ajaxhook.min.js
var ajaxhookJs string

func writeContentType(w http.ResponseWriter, path string) {

	if strings.HasSuffix(path, ".js") {
		w.Header().Add("Content-Type", "text/javascript")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(path, ".html") {
		w.Header().Add("Content-Type", "text/html")
	} else if strings.HasSuffix(path, ".png") {
		w.Header().Add("Content-Type", "image/png")
	} else if strings.HasSuffix(path, ".jpg") {
		w.Header().Add("Content-Type", "image/jpg")
	} else if strings.HasSuffix(path, ".gif") {
		w.Header().Add("Content-Type", "image/gif")
	} else if strings.HasSuffix(path, ".svg") {
		w.Header().Add("Content-Type", "image/svg+xml")
	} else if strings.HasSuffix(path, ".ico") {
		w.Header().Add("Content-Type", "image/x-icon")
	} else if strings.HasSuffix(path, ".json") {
		w.Header().Add("Content-Type", "application/json")
	} else if strings.HasSuffix(path, ".woff") {
		w.Header().Add("Content-Type", "application/font-woff")
	} else if strings.HasSuffix(path, ".woff2") {
		w.Header().Add("Content-Type", "application/font-woff2")
	} else if strings.HasSuffix(path, ".ttf") {
		w.Header().Add("Content-Type", "application/font-ttf")
	} else if strings.HasSuffix(path, ".eot") {
		w.Header().Add("Content-Type", "application/vnd.ms-fontobject")
	} else if strings.HasSuffix(path, ".otf") {
		w.Header().Add("Content-Type", "application/font-otf")
	} else if strings.HasSuffix(path, ".xml") {
		w.Header().Add("Content-Type", "application/xml")
	} else if strings.HasSuffix(path, ".pdf") {
		w.Header().Add("Content-Type", "application/pdf")
	} else if strings.HasSuffix(path, ".mp4") {
		w.Header().Add("Content-Type", "video/mp4")
	} else {
		w.Header().Add("Content-Type", "application/octet-stream")
	}

}

func (f *AppMode) startUpHTTP() (*http.Server, int) {

	freePort := getFreePort()
	if f.srcFilesystem != nil {
		fss, err := fs.Sub(f.srcFilesystem, f.srcDirSub)

		if err != nil {
			panic(err)
		}

		fs.WalkDir(fss, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if d.IsDir() {
				return nil
			}
			// log.Println(path)

			http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
				writeContentType(w, path)
				fp, err := fss.Open(path)
				defer fp.Close()
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}

				data, err := ioutil.ReadAll(fp)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}

				w.Write(data)

			})
			return nil
		})

	} else {
		fs := http.FileServer(http.Dir(f.loadAppFilename))
		http.Handle("/", fs)
	}

	server := &http.Server{
		Addr: fmt.Sprintf("localhost:%d", freePort),
	}
	go func() {
		server.ListenAndServe()
	}()
	return server, freePort
}

func (f *AppMode) SetStorage(storage StorageInfterface) {
	f.storage = storage
}

func (f *AppMode) InitEnvironment(w WebView, args map[string]interface{}) {
	jsSrc := ajaxhookJs + "\n\n"
	for k, v := range args {
		val, _ := json.Marshal(v)
		jsSrc = jsSrc + fmt.Sprintf("window.%s = %s;", k, string(val))
	}
	w.Init(jsSrc)
}

type AppModeConfig struct {
	Width          uint
	Height         uint
	Title          string
	Debug          bool
	InitJsSrc      string
	InitJsFiles    []string
	Hint           string
	GlobalVariable map[string]interface{}
}

func (f *AppMode) Run(args AppModeConfig) {
	server, freePort := f.startUpHTTP()

	width := args.Width
	if width == 0 {
		width = 1080
	}

	height := args.Height
	if height == 0 {
		height = 860
	}

	debug := args.Debug
	initJs := args.InitJsSrc

	title := args.Title

	w := NewWithOptions(WebViewOptions{
		Debug:     debug,
		AutoFocus: true,
		WindowOptions: WindowOptions{
			IconId: 2, // icon resource id
			Center: true,
			Width:  width,
			Title:  title,
			Height: height,
		},
	})

	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer func() {
		w.Destroy()
		server.Close()
	}()
	w.SetTitle(title)

	switch args.Hint {
	case "fixed":
		w.SetSize(int(width), int(height), HintFixed)
	case "min":
		w.SetSize(int(width), int(height), HintMin)
	case "max":
		w.SetSize(int(width), int(height), HintMax)
	case "none":
		w.SetSize(int(width), int(height), HintNone)
	}

	f.InitEnvironment(w, args.GlobalVariable)
	w.Init(initJs)
	w.Navigate(fmt.Sprintf("http://127.0.0.1:%d/index.html", freePort))
	w.Run()
}
