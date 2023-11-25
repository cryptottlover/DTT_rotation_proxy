/*
[запускается зенка]
создается колво потоков которое вы указали слева внизу
каждый поток кидает запрос на 127.0.0.1:1337/getproxy вашего вебсервера чтобы получить прокси
эта прокси уже в вашем вебсервере будет не доступна т.к. она уже в работе
когда зенка завершит поток она кинет запрост на 127.0.0.1:1337/endproxy?proxy=прокси и кидает запрос чтобы сделало ротейт
и прокси уже можно будет получить через N секунд (которое вы указали)

F.A.Q
1) Чтобы узнать время сколько секунд ждать пока поменяется IP прокси - просто перейдите по ссылке
2) Если все прокси заняты и потоки зенки не завершили работу, то будетж дать N секунд скок указал отдельно в коде
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Proxy struct {
	proxyurl      string
	rotationurl   string
	rotationdelay int
	busy          bool
	canbeused     time.Time
}

var myProxies []Proxy
var mutex sync.Mutex
var minWait int
var maxWait int

func main() {

	port := 1337 //вебсервер API с проксями

	/* сколько поток будет делать паузу, если все прокси заняты (либо уже не заняты, но время ротейта еще не прошо, то что вы указали 3 параметров после ссылки на апи ротейта каждой прокси) */

	minWait = 1000 //в миллисекундах
	maxWait = 3000 //в миллисекундах

	//прокси юрл, линк на ротейт (роейт  происходит после получения прокси для зенки), кол-во секунд которое ждать ротейта (в апи пишется обычно), ебать не должно че последнее

	myProxies = append(myProxies, Proxy{"socks5://login:password@1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()})
	myProxies = append(myProxies, Proxy{"socks5://1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()})

	http.HandleFunc("/getproxy", getProxy)
	http.HandleFunc("/endproxy", endProxy)

	fmt.Printf("Server is listening on port %d...\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func getProxy(w http.ResponseWriter, r *http.Request) {

	var selectedProxy Proxy
	mutex.Lock()
	defer mutex.Unlock()

	for i := range myProxies {
		if !myProxies[i].busy && time.Now().After(myProxies[i].canbeused) {
			selectedProxy = myProxies[i]

			myProxies[i].busy = true

			response := fmt.Sprintf(selectedProxy.proxyurl)
			log.Println("Thread got proxy", response)
			w.Write([]byte(response))
			return
		}
	}

	selectedProxy = findFastestProxy()

	if selectedProxy.rotationurl == "" {

		source := rand.NewSource(time.Now().UnixNano())
		randomGenerator := rand.New(source)
		grisha := randomGenerator.Intn(maxWait-minWait+1) + minWait
		log.Printf("Thread got all proxies is busy. Random Sleep: %d millisecond\n", grisha)
		response := fmt.Sprintf("%d", grisha)
		w.Write([]byte(response))

		fmt.Println("")
	} else {
		timeRemaining := int(time.Until(selectedProxy.canbeused).Milliseconds())
		log.Printf("Thread got not all proxies is busy. Soon end: %d millisecond\n", timeRemaining)
		response := fmt.Sprintf("%d", timeRemaining)
		w.Write([]byte(response))
	}

}

func endProxy(w http.ResponseWriter, r *http.Request) {

	proxyParam := r.URL.Query().Get("proxy")
	statusParam := r.URL.Query().Get("status")

	if proxyParam == "" || statusParam == "" {
		w.Write([]byte(`Where proxy and status http GET parameters`))
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i := range myProxies {
		if myProxies[i].proxyurl == proxyParam {

			myProxies[i].busy = false

			if statusParam == "true" {
				myProxies[i].canbeused = time.Now().Add(time.Second * time.Duration(myProxies[i].rotationdelay))
				log.Println("Thread got successfully rotated the proxy", proxyParam)
				go sendRequest(myProxies[i].rotationurl, i)
			} else {
				log.Println("Thread got successfully stopped the proxy (maybe got error in platform)", proxyParam)
			}

			break
		}
	}

	w.Write([]byte(`Good`))
}

func findFastestProxy() Proxy {

	var fastestProxy Proxy

	for i := range myProxies {
		if !myProxies[i].busy && myProxies[i].canbeused.After(time.Now()) {
			fastestProxy = myProxies[i]
			break
		}
	}

	return fastestProxy
}

func sendRequest(url string, proxyINDEX int) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("[%s] %s :: %d\n", myProxies[proxyINDEX].proxyurl, url, response.StatusCode)
}
