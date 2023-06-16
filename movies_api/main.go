package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 图片访问方式：https://image.tmdb.org/t/p/w780/8FhKnPpql374qyyHAkZDld93IUw.jpg
// https://image.tmdb.org/t/p/{size}/{path}
// {size}：表示图像的尺寸，您可以选择以下选项之一：w92, w154, w185, w342, w500, w780, original。选择适合您需求的尺寸。
// {path}：表示TMDB返回的图像的相对路径，例如/example_path.jpg。

func main() {

	//u := "https://api.themoviedb.org/3/discover/movie?include_adult=true&include_video=true&language=zh&page=1&sort_by=popularity.desc"
	//u := "https://api.themoviedb.org/3/movie/840326?include_adult=true&include_video=true&language=zh&page=1&sort_by=popularity.desc"
	//u := "https://api.themoviedb.org/3/search/movie?query=极寒之城&api_key=62d94d1cbdb6e5b05fe83c0e449f79e3&language=zh"
	//u := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=zh&page=1&sort_by=popularity.desc"
	u := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=zh&page=1&sort_by=popularity.desc"

	proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI2MmQ5NGQxY2JkYjZlNWIwNWZlODNjMGU0NDlmNzllMyIsInN1YiI6IjY0Nzc0NTE0MjU1ZGJhMDBhOWEyMTUzZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.82bP1htsojps9KHeaGlW4ftfOy-wskQbBN8A_WHlLVc")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(response.Body)

	fmt.Println(string(body))

}
