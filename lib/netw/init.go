package netw

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

const useragent = "github.com/laptopcat/netw v0.1 - "

type Probe struct {
	IPA    string
	ASN    int
	ASS    string
	IPP    string
	Worked bool
}

var httpc = &fasthttp.Client{
	DialDualStack: true,
}

var cfg struct {
	Email string
}

func acquire() *fasthttp.Request {
	req := fasthttp.AcquireRequest()

	req.Header.SetUserAgent(useragent + cfg.Email)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	return req
}

// sometimes ipv4 and ipv6 probes get messed up when I use graftcp with the tool
var mut sync.Mutex

func ProbeV4() (p Probe, err error) {
	req := acquire()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://v4.bgp.tools/whoami")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	mut.Lock()
	err = httpc.Do(req, resp)
	if err != nil {
		mut.Unlock()
		return
	}

	data, err := resp.BodyUncompressed()
	mut.Unlock()
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &p)
	return
}

func ProbeV6() (p Probe, err error) {
	req := acquire()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://v6.bgp.tools/whoami")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	mut.Lock()

	err = httpc.Do(req, resp)
	if err != nil {
		mut.Unlock()
		return
	}

	data, err := resp.BodyUncompressed()
	mut.Unlock()
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &p)
	return
}

func init() {
	d, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("can't get homedir:", err)
	}

	if d[len(d)-1] != '/' {
		d += "/"
	}
	d += ".netw"

	os.Mkdir(d, 0766)
	f, err := os.OpenFile(d+"/config.json", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		log.Fatalln("open config:", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln("read config:", err)
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			log.Fatalln("parse config:", err)
		}
	}

	if cfg.Email == "" {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input your E-mail").
					Description("bgp.tools requires you to identify yourself if you use their APIs. Read more here: https://bgp.tools/kb/api").
					Value(&cfg.Email),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatalln("failed to run form:", err)
		}

		enc := json.NewEncoder(f)
		enc.Encode(cfg)
		f.Close()
	}
}
