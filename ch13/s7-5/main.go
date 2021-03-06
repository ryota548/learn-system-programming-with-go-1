package main

import (
    "fmt"
    "sync"
    //"runtime"
)

// sync.Pool オブジェクトのキャッシュ
func main() {
    // Poolを作成。Newで新規作成時のコードを実装
    var count int
    pool := sync.Pool {
        // このNewはだれ??
        New: func() interface {} {
            count++
            return fmt.Sprintf("created: %d", count)
        },
    }

    // 追加した要素から受け取れる
    // プールが空だと新規作成
    pool.Put("manualy added: 1")
    pool.Put("manualy added: 2")
    // GCを呼ぶと追加された要素が消える
    // runtime.GC()
    fmt.Println(pool.Get())
    fmt.Println(pool.Get())
    fmt.Println(pool.Get()) //ここは新規作成
}
