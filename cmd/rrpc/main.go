package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spiral/grpc"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"

	// grpc commands
	_ "github.com/spiral/grpc/cmd/rrpc/grpc"
)

func main() {
	rr.Container.Register(grpc.ID, &grpc.Service{})

	rr.Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	rr.Execute()
}
