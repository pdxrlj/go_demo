package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Demo struct {
	a string
}

func main() {
	v := "object://test:123456@obs/region/bucket/path"
	bbb := filepath.Join(v, "z", "x", "y")
	fmt.Println(bbb)
	return

	cc := "zy_x"
	pre := ""
	var tmp []string
	glue := false
	strings.FieldsFunc(cc, func(r rune) bool {
		if glue {
			fmt.Println("glue")
			glue = false
			pre += "_" + string(r)
			tmp = append(tmp, pre)
			return true
		}

		if string(r) == "_" {
			glue = true
			tmp = tmp[:len(tmp)-1]
			return true
		}

		pre = string(r)
		tmp = append(tmp, string(r))

		return true
	})
	fmt.Printf("%#v\n", tmp)
	g := filepath.Join(tmp...)
	fmt.Println(g)
	z := 10
	x := 1
	y := 4
	bb := strings.NewReplacer("z", cast.ToString(z), "x", cast.ToString(x), "y", cast.ToString(y)).Replace(g)
	fmt.Println(bb)

	//for _, i := range cc {
	//	fmt.Println(string(i))
	//}

	return

	d := &Demo{}
	d.a = "a"
	c := d.a
	fmt.Println(c)
	d.a = "b"
	fmt.Println(c)
	return

	t := time.Now()
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Second*5, errors.New("timeout"))
	//defer cancelFunc()

ForEnd:
	select {
	case <-time.After(time.Second * 20):
		fmt.Println("time out")
		//cancelFunc()
	case <-ctx.Done():
		fmt.Println("ctx done")
		s := time.Now().Sub(t).Seconds()
		fmt.Println(s)
		break ForEnd
	}

	fmt.Println("end")
	return

	// GOEXPERIMENT=loopvar
	//wg := errgroup.Group{}
	//wg.SetLimit(10)
	//for i := 0; i < 10; i++ {
	//	wg.Go(func() error {
	//		time.Sleep(1)
	//		fmt.Println("Hello:" + strconv.Itoa(i))
	//		return nil
	//	})
	//}
	//if err := wg.Wait(); err != nil {
	//	fmt.Println(err)
	//}
	//a := min(10, 1)
	//fmt.Printf("max:%d\n", a)
	//
	//c := map[string]string{"a": "b", "c": "d", "e": "f"}
	//maps.DeleteFunc(c, func(key string, value string) bool {
	//	return key == "a"
	//})
	//fmt.Println(c)
	//slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	//	AddSource: true,
	//	Level:     slog.LevelInfo,
	//	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	//		if a.Key == "time" {
	//			a.Value = slog.StringValue(time.Now().Format(time.DateTime))
	//		}
	//		return a
	//	},
	//})))
	//slog.Info("Hello World!")
	//g := slog.Group("d")

}

var _ slog.Handler = (*CustomSlogHandler)(nil)

type CustomSlogHandler struct {
	slog.Handler
	buf *bytes.Buffer
	ch  chan []byte
}

func (c *CustomSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.Handler.Enabled(ctx, level)
}

func (c *CustomSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	c.buf.WriteString(record.Level.String())
	c.buf.WriteByte(' ')
	c.buf.WriteString(record.Message)
	//c.Handler.Handle()
	return nil
}

func (c *CustomSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func (c *CustomSlogHandler) WithGroup(name string) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func SlogHandler() {
	//slog.SetDefault(CustomSlogHandler{})
}
