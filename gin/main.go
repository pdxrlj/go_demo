package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password any    `json:"-"`
}

func (p Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	return json.Marshal(&struct {
		Alias
		Password any `json:"password"`
	}{
		Alias:    (Alias)(p),
		Password: p.Password,
	})
}

func TestMap(m map[string]string) {
	(m)["a"] = "gggggggggg"
}

func main() {
	b := make(map[string]string)
	b["a"] = "b"
	b["c"] = "d"
	TestMap(b)
	fmt.Println(b)
	return

	p := Person{
		Name:     "Alice",
		Age:      25,
		Password: []byte("123456"),
	}

	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("JSON序列化失败:", err)
		return
	}

	fmt.Println(string(data))

	return

	url := "https://www.baidu.com"

	before, _, found := strings.Cut(url, "?")
	fmt.Println("before:", before, "found:", found)
	return

	engine := gin.Default()
	engine.GET("/:id/:name", func(context *gin.Context) {
		id := context.Param("id")
		name := context.Param("name")
		fmt.Println("id:", id, "name:", name)
		sex := context.Query("sex")
		context.String(200, "id:"+id+" name:"+name+" sex:"+sex)
	})

	engine.GET("/index.html", func(context *gin.Context) {
		file, err := os.Open("index.html")
		if err != nil {
			context.String(404, "%s", err.Error())
			return
		}
		all, err := io.ReadAll(file)
		if err != nil {
			context.String(404, "%s", err.Error())
			return
		}
		context.String(200, "%s", all)
	})

	log.Fatal(engine.Run(":18085"))
}
