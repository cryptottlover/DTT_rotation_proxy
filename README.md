# DTT_rotation_proxy
Golang Web Server Proxies with automatic rotation

# How install
Download https://go.dev/ </br>
Edit rotationproxy.go there you need enter a your proxy with rotation api link

//proxy, api link rotate, how many second need wait for rotate, time.Now()
myProxies = append(myProxies, Proxy{"socks5://login:password@1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()})
myProxies = append(myProxies, Proxy{"socks5://1.2.3.4:5052", "https://api.mail.com/rotate?accessId=i9g24gmd9mppyzdphwspsoho", 6, false, time.Now()})


# How run
Open cmd (Win + R > cmd.exe)
Write: "cd  [path to folder with rotationproxy.go]". Command example:
cd C:\Users\admin\Desktop

And run: go run rotationproxy.go
