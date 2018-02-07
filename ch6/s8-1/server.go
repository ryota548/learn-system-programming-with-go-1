package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "net/http/httputil"
    "strings"
    "time"
    "io"
    "bytes"
    "compress/gzip"
)

// クライアントはgzipを受け入れ可能か
func isGzipAcceptable(request *http.Request) bool {
    return strings.Index(strings.Join(request.Header["Accept-Encoding"], ","),"gzip")  != -1
}

// 1セッションの胥吏をする
func processSession(conn net.Conn) {
    fmt.Printf("Accept %v\n", conn.RemoteAddr())
    defer conn.Close()
    // Accept後のソケットで何度も応答を返すためループ
    for {
        // タイムアウトを設定
        conn.SetReadDeadline(time.Now().Add(5 * time.Second))
        // リクエストを読み込む
        request, err := http.ReadRequest(bufio.NewReader(conn))
        if err != nil {
            // タイムアウトもしくはソケットクローズ時は終了
            // それ以外はエラーになる
            neterr, ok := err.(net.Error)
            if ok && neterr.Timeout() {
                fmt.Println("Timeout")
                break
            } else if err == io.EOF {
                break
            }
            panic(err)
        }

        // リクエストを表示
        dump, err := httputil.DumpRequest(request, true)
        if err != nil{
            panic(err)
        }
        fmt.Println(string(dump))

        // レスポンスを書き込む
        response := http.Response{
            StatusCode: 200,
            ProtoMajor: 1,
            ProtoMinor: 1,
            Header: make(http.Header),
        }
        if isGzipAcceptable(request) {
            content := "Hello World(gzipped)\n"
            // コンテンツをgzip化して転送
            var buffer bytes.Buffer
            writer := gzip.NewWriter(&buffer)
            io.WriteString(writer, content)
            writer.Close()
            response.Body = ioutil.NopCloser(&buffer)
            response.ContentLength = int64(len(content))
        }
        response.Write(conn)
    }
}

func main () {
    listener, err := net.Listen("tcp", "localhost:8888")
    if err != nil {
        panic(err)
    }
    fmt.Println("Server is running at localhost:8888")
    for {
        conn, err := listener.Accept()
        if err != nil {
            panic(err)
        }
        go processSession(conn)
    }
}
