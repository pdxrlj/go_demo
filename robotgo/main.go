package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/keycode"
)

func main() {
	//x, y := robotgo.GetScreenSize()
	//fmt.Println("color----", x, y)
	//active := robotgo.GetActive()
	//fmt.Printf("active: %v\n", active)
	//
	////for {
	////	cx, cy := robotgo.GetMousePos()
	////	fmt.Printf("pos: %v %v\n", cx, cy)
	////	time.Sleep(time.Second)
	////}
	//
	////sx, sy := robotgo.GetScaleSize()
	////fmt.Println("scale size----", sx, sy)
	////robotgo.ShowAlert("demo", "robotgo", "好的")
	//
	//robotgo.MoveSmooth(522, 467)
	//robotgo.MoveClick(522, 467)
	//robotgo.TypeStr("// hello world")
	//
	//all, err := robotgo.ReadAll()
	//
	//fmt.Println("readall-----", all, err)
	//
	//for u, _ := range keycode.Keycode {
	//	hook.AddEvent(u)
	//}
	//events := hook.Start()
	//defer hook.End()
	//select {
	//case event := <-events:
	//	fmt.Printf("event: %v\n", event)
	//}

	ss := make([]string, len(keycode.Keycode))
	for s, _ := range keycode.Keycode {
		ss = append(ss, s)
	}

	go func() {
		if event := robotgo.AddEvents("", ss...); event {
			println("you press...", "a")
		}
	}()

	select {}

}

// hello world
