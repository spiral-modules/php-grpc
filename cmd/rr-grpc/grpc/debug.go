// Copyright (c) 2018 SpiralScout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package grpc

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	rrpc "github.com/spiral/php-grpc"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/cmd/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
)

func init() {
	cobra.OnInitialize(func() {
		if rr.Debug {
			svc, _ := rr.Container.Get(rrpc.ID)
			if svc, ok := svc.(*rrpc.Service); ok {
				debug := &debugger{logger: rr.Logger}
				svc.AddListener(debug.listener)
				svc.AddOption(grpc.UnaryInterceptor(debug.interceptor))
			}
		}
	})
}

// listener provide debug callback for system events. With colors!
type debugger struct{ logger *logrus.Logger }

// listener listens to http events and generates nice looking output.
func (d *debugger) listener(event int, ctx interface{}) {
	if util.LogEvent(d.logger, event, ctx) {
		// handler by default debug package
		return
	}
}

// call info
func (d *debugger) interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)

	if err == nil {
		d.logger.Info(util.Sprintf(
			"<cyan+h>%s</reset> <green+h>Ok</reset> %s %s",
			d.getPeer(ctx),
			elapsed(time.Since(start)),
			info.FullMethod,
		))
	} else {
		if st, ok := status.FromError(err); ok {
			d.logger.Error(util.Sprintf(
				"<cyan+h>%s</reset> %s %s %s <red>%s</reset>",
				d.getPeer(ctx),
				d.wrapStatus(st),
				elapsed(time.Since(start)),
				info.FullMethod,
				st.Message(),
			))
		} else {
			d.logger.Error(util.Sprintf(
				"<cyan+h>%s</reset> %s %s <red>%s</reset>",
				d.getPeer(ctx),
				elapsed(time.Since(start)),
				info.FullMethod,
				err.Error(),
			))
		}
	}

	return resp, err
}

func (d *debugger) getPeer(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if ok {
		return pr.Addr.String()
	}

	return "unknown"
}

func (d *debugger) wrapStatus(st *status.Status) string {
	switch st.Code() {
	case codes.NotFound, codes.Canceled, codes.Unavailable:
		return util.Sprintf("<yellow+h>%s</reset>", st.Code().String())
	}

	return util.Sprintf("<red+h>%s</reset>", st.Code().String())
}

// fits duration into 5 characters
func elapsed(d time.Duration) string {
	var v string
	switch {
	case d > 100*time.Second:
		v = fmt.Sprintf("%.1fs", d.Seconds())
	case d > 10*time.Second:
		v = fmt.Sprintf("%.2fs", d.Seconds())
	case d > time.Second:
		v = fmt.Sprintf("%.3fs", d.Seconds())
	case d > 100*time.Millisecond:
		v = fmt.Sprintf("%.0fms", d.Seconds()*1000)
	case d > 10*time.Millisecond:
		v = fmt.Sprintf("%.1fms", d.Seconds()*1000)
	default:
		v = fmt.Sprintf("%.2fms", d.Seconds()*1000)
	}

	if d > time.Second {
		return util.Sprintf("<red>{%v}</reset>", v)
	}

	if d > time.Millisecond*50 {
		return util.Sprintf("<yellow>{%v}</reset>", v)
	}

	return util.Sprintf("<gray+hb>{%v}</reset>", v)
}
