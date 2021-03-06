package gcpool

import (
	"io"
	"log"
	"net/http"
	"testing"

	"bufio"

	"golang.org/x/net/websocket"
)

// 节点连接池
var GO_CONN_POOL *Pool

func Test(t *testing.T) {

	GO_CONN_POOL = NewPool()         // 创建连接池
	GO_CONN_POOL.Register("default") // 注册连接组
	GO_CONN_POOL.Start()             // 启动服务

	// 创建 HTTP + WebSocket 服务
	http.Handle("/hello", websocket.Handler(HelloHandler))
	// 启动服务 ...
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}

// WS Handler
func HelloHandler(ws *websocket.Conn) {

	// WS 获取请求参数 ...
	if err := ws.Request().ParseForm(); err != nil {
		return
	}

	// 作为唯一标识符
	id := ws.Request().FormValue("id")

	// 保存连接
	GO_CONN_POOL.GetConn("default").Add(id, ws)
	// 断开移除
	defer GO_CONN_POOL.GetConn("default").Del(id)

	// 读取 ... 阻塞
	r := bufio.NewReader(ws)
	for {
		// 按行读取 ... (JSON)
		_, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				// 异常 ...
			}
			// 异常时跳出循环 ... 断开连接 ...
			break
		}
	}

}
