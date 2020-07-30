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
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	rrpc "github.com/spiral/php-grpc"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/cmd/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func init() {
	cobra.OnInitialize(func() {
		if rr.Debug {
			svc, _ := rr.Container.Get(rrpc.ID)
			if svc, ok := svc.(*rrpc.Service); ok {
				debug := &debugger{logger: rr.Logger}
				svc.AddListener(debug.listener)
			}
		}
	})
}

// listener provide debug callback for system events. With colors!
type debugger struct{ logger *logrus.Logger }

// listener listens to http events and generates nice looking output.
func (d *debugger) listener(event int, ctx interface{}) {
	if event == rrpc.EventUnaryCall {
		uc := ctx.(*rrpc.UnaryCallEvent)

		if uc.Error == nil {
			d.logger.Info(util.Sprintf(
				"<cyan+h>%s</reset> <green+h>Ok</reset> %s %s",
				getPeer(uc.Context),
				elapsed(uc.Elapsed()),
				uc.Info.FullMethod,
			))
			return
		}
		if st, ok := status.FromError(uc.Error); ok {
			d.logger.Error(util.Sprintf(
				"<cyan+h>%s</reset> %s %s %s <red>%s</reset>",
				getPeer(uc.Context),
				wrapStatus(st),
				elapsed(uc.Elapsed()),
				uc.Info.FullMethod,
				st.Message(),
			))
		} else {
			d.logger.Error(util.Sprintf(
				"<cyan+h>%s</reset> %s %s <red>%s</reset>",
				getPeer(uc.Context),
				elapsed(uc.Elapsed()),
				uc.Info.FullMethod,
				uc.Error.Error(),
			))
		}
	}

	if util.LogEvent(d.logger, event, ctx) {
		// handler by default debug package
		return
	}
}

func wrapStatus(st *status.Status) string {
	switch st.Code() {
	case codes.NotFound, codes.Canceled, codes.Unavailable:
		return util.Sprintf("<yellow+h>%s</reset>", st.Code().String())
	}

	return util.Sprintf("<red+h>%s</reset>", st.Code().String())
}

func getPeer(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if ok {
		return pr.Addr.String()
	}

	return "unknown"
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
