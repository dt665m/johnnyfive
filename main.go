package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const ()

var (
	listenPort   = "443" //defaults to SSL currently
	disableHTTP2 bool
)

func init() {
	flag.BoolVar(&disableHTTP2, "h1", false, "h1 flag disables http2")
	port := os.Getenv("PORT")
	if len(port) > 0 {
		listenPort = port
	}
}

func main() {
	withTLSandH2()
}

func withTLSandH2() {
	fs := http.FileServer(http.Dir("/srv/files/"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handleWebRoot)
	mux.HandleFunc("/view", handleView)
	mux.HandleFunc("/view/events", subManager.SubscriptionHandler)
	mux.HandleFunc("/wechatcallback", WechatCallback)
	if lineServer != nil {
		mux.HandleFunc("/linecallback", LineCallback)
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			//tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	tlsproto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
	if disableHTTP2 {
		fmt.Println("Disabling HTTP2")
		tlsproto = nil
	}
	server := &http.Server{
		Addr:         ":" + listenPort,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: tlsproto, //UNCOMMENT FOR HTTP1.1
	}

	fmt.Println("Starting HTTP Server on port", server.Addr)
	server.ErrorLog = log.New(ioutil.Discard, "", log.Lshortfile)
	err := server.ListenAndServeTLS("/srv/fullchain.pem", "/srv/privkey.pem")
	if err != nil {
		fmt.Println("HTTP Server Failed: ", err)
	}

	subManager.Shutdown()
}

func handleWebRoot(w http.ResponseWriter, req *http.Request) {
	if disableHTTP2 {
		fmt.Printf("SERVING HTTP1.1 to %s, user-agent: %s\n", req.RemoteAddr, req.UserAgent())
	} else {
		fmt.Printf("ATTEMPTING TO SERVE HTTP2 to %s, user-agent: %s\n", req.RemoteAddr, req.UserAgent())
	}

	if len(req.Cookies()) == 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		random := r.Intn(1000)
		randomStr := strconv.Itoa(random)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := &http.Cookie{
			Name:    "randomUser" + randomStr,
			Value:   randomStr,
			Expires: expiration,
			Domain:  "armordt.me",
		}
		fmt.Printf("No Cookie Found, Setting Cookie: %+v\n", cookie)
		http.SetCookie(w, cookie)
	} else {
		for _, cookie := range req.Cookies() {
			fmt.Println("Cookie Found!")
			fmt.Printf("%+v\n", cookie)
		}
	}

	fmt.Fprintf(w, fakehttpdata)
}

func parseContent(content string) (map[string]int, error) {
	//r := strings.NewReader(content)
	return nil, errors.New("Not Implemented")
}

var fakehttpdata = `
<HEAD>
<TITLE>Armor Denny Tiffy!</TITLE>
</HEAD>
<BODY BGCOLOR="WHITE">
<CENTER>
<H1>Armor Denny Tiffy!!!</H1>

 

  <IMG SRC="/static/thisisarmordt.jpg">

 

 

  <H4>ARMORDT</H4>
`

/*
// DEPRICATED/TESTING ONLY
func VerifyWechat(token, timestamp, nonce, signature string) bool {
	s := []string{token, timestamp, nonce}
	sort.Strings(s)
	fmt.Println("Sorted: ", s)
	j := strings.Join(s[:], "")
	fmt.Println("Joined: ", j)
	h := sha1.New()
	h.Write([]byte(j))
	sumHex := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Sha1 Sum: ", sumHex)
	fmt.Println("Sha1 Signature: ", signature)
	return sumHex == signature
}


func handleCallback(w http.ResponseWriter, req *http.Request) {
	fmt.Println("######")
	fmt.Println(req.Host, req.RemoteAddr, req.UserAgent())
	fmt.Println(req.URL.RequestURI())
	if err := req.ParseForm(); err != nil {
		// handle error
		fmt.Println("error:", err)
	}

	for key, values := range req.PostForm {
		fmt.Println(key, values)
	}
	fmt.Println("######")

	if len(req.FormValue("echostr")) > 0 {
		handleCallbackVerify(w, req)
		return
	}

}


func handleCallbackVerify(w http.ResponseWriter, req *http.Request) {
	echostr := req.FormValue("echostr")
	fmt.Println(echostr)

	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	signature := req.FormValue("signature")
	if VerifyWechat(WECHAT_TOKEN, timestamp, nonce, signature) {
		w.Header().Set("echostr", echostr)
		fmt.Fprintf(w, echostr)
	} else {
		fmt.Println("Wechat Verification Failed!")
	}
}*/
