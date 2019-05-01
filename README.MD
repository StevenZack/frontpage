# FrontPage : Simple GUI Framework for Go

# Quick Start

Install frontpage
```shell
go get github.com/StevenZack/frontpage
```
Create a main.go
```go
package main

import (
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
	</head>
	<body>
		<input type="button" name="show" id="show" value="show">
	</body>
	</html>`)
	fp.Run()
}
```
Run
