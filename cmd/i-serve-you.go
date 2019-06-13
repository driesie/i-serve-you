// Copyright (c) 2019 Oliver Wyman Digital
// Use of this source code is governed by a MIT Licence
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/lshift/i-serve-you/pkg/server"
)

func main() {
	config := server.NewConfig()
	portFlag := flag.Int("port", config.Port, "port to serve on")
	adminpPortFlag := flag.Int("adminPort", config.AdminPort, "port to serve admin api on")

	flag.Parse()

	config.Port = *portFlag
	config.AdminPort = *adminpPortFlag

	server.Start(config)
}
