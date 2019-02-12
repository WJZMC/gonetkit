package server

import (
	"gonetkit/interfacer"
	"gonetkit/router"
	"testing"
)

const (
	QA_submit = iota
	QA_getreport
)

type Qasubmit struct {
	router.BaseRouter
}

func (r *Qasubmit) Handle(request interfacer.Requester) {

}

type Qareport struct {
	router.BaseRouter
}

func (r *Qareport) Handle(request interfacer.Requester) {

}

func TestServe_Serve(t *testing.T) {

	server := NewServe()

	var submit Qasubmit = Qasubmit{}
	var report Qareport = Qareport{}

	server.AddRouter(QA_submit, &submit)
	server.AddRouter(QA_getreport, &report)

	server.Serve()

}
