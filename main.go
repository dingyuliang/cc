package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cc/controller"
	"github.com/ksarch-saas/cc/frontend"
	"github.com/ksarch-saas/cc/inspector"
	"github.com/ksarch-saas/cc/log"
	"github.com/ksarch-saas/cc/meta"
	"github.com/ksarch-saas/cc/streams"
	"github.com/ksarch-saas/cc/topo"
	"github.com/ksarch-saas/cc/utils"
)

var (
	appName     string
	localRegion string
	seeds       string
	zkHosts     string
	zkService   string
	httpPort    int
	wsPort      int
)

func init() {
	flag.StringVar(&appName, "appname", "", "app name")
	flag.StringVar(&localRegion, "local-region", "", "local region")
	flag.StringVar(&seeds, "seeds", "", "redis cluster seeds, seperate by comma")
	flag.StringVar(&zkHosts, "zkhosts", "", "zk hosts, seperate by comma")
	flag.StringVar(&zkService, "zkservice", "", "naming service locate zookeeper")
	flag.IntVar(&httpPort, "http-port", 0, "http port")
	flag.IntVar(&wsPort, "ws-port", 0, "ws port")
}

func main() {
	flag.Parse()

	seedNodes := []*topo.Node{}
	for _, addr := range strings.Split(seeds, ",") {
		glog.Info(addr)
		n := topo.NewNodeFromString(addr)
		if n == nil {
			glog.Fatal("invalid seeds %s", addr)
		}
		seedNodes = append(seedNodes, n)
	}
	if httpPort == 0 {
		glog.Fatal("invalid http port")
		flag.PrintDefaults()
	}
	if wsPort == 0 {
		glog.Fatal("invalid websocket port")
		flag.PrintDefaults()
	}
	if zkHosts == "" {
		glog.Info("zkHosts not set, trying to locate by naming service")
		var err error
		zkHosts, err = utils.ResolvZkFromNaming(zkService)
		if err != nil {
			glog.Fatal("Resolve zkhosts failed")
		}
		glog.Infof("Resolve zkhosts from naming service, result is %s", zkHosts)
	}

	initCh := make(chan error)
	go meta.Run(appName, localRegion, httpPort, wsPort, zkHosts, seedNodes, initCh)
	err := <-initCh
	if err != nil {
		glog.Fatal(err)
	}

	streams.StartAllStreams()
	streams.LogStream.Sub(log.WriteFileHandler)
	streams.LogStream.Sub(log.WriteRingBufferHandler)

	sp := inspector.NewInspector()
	go sp.Run()

	c := controller.NewController()
	fe := frontend.NewFrontEnd(c, httpPort, wsPort)
	fe.Run()
}
