# DTT_rotation_proxy
Golang Web Server Proxies with automatic rotation

# How install
Download https://go.dev/ </br> </br>
Edit rotationproxy.go there you need enter a your proxy with rotation api link </br> </br>
//Example: </br>
//proxy, api link rotate, how many second need wait for rotate, time.Now() </br>
myProxies = append(myProxies, Proxy{"socks5://login:password@1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()}) </br>
myProxies = append(myProxies, Proxy{"socks5://1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()}) </br>

# How run
Open cmd (Win + R > cmd.exe) </br>
Write: "cd  [path to folder with rotationproxy.go]". Command example: </br>
cd C:\Users\admin\Desktop </br>
 </br>
And run: go run rotationproxy.go
