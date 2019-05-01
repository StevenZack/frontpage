# FrontPage : A cross-platform front page for Go

# FUCK GUI ! JUST LAUNCH YOUR BROWSER

# Quick Start

Install frontpage
```shell
go get github.com/StevenZack/frontpage
```
Create a main.go
```go
package main

import (
	"fmt"
	"time"

	"github.com/StevenZack/frontpage"
)

func main() {
	fp := frontpage.New(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Document</title>
		<script src="/var.js"></script>
		<script>
			function showJSAlert(s){
				alert(s)
			}
		</script>
	</head>
	<body>
		<input type="button" value="goprint" onclick="goprint('call Go from js')">
	</body>
	</html>`)
	fp.BindFunc(goprint)
	go func() {
		time.Sleep(time.Second * 2)
		fp.Invoke("showJSAlert", "call js function from Go")
	}()
	fp.Run()
}

func goprint(s string) {
	fmt.Println(s)
}

```
Run
```shell
go run main.go
```

![screenshot](/res/get_started.png)

# How it works

1.Create a local HTTP server with random port

2.Bind Go Functions

3.Launch your browser, If you installed Chrome, we'll use app mode (No address bar)

4.Connect Javascript and Go by WebSocket

