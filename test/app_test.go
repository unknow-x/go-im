/**
  @author:kk
  @data:2021/11/14
  @note
**/
package test

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"im_app/config"
	"im_app/pkg/jwt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"testing"
	"time"
)

func init() {
	config.Initialize()
}

var addr = flag.String("addr", "127.0.0.1:9502", "http service address")

type TestHandler interface {
	Send(id int, msg []byte)
}

var ClientManager = ConnMap{
	client: make(map[int]*Client),
}

type ConnMap struct {
	client map[int]*Client
}
type Client struct {
	Index int
	Token string
	Conn  *websocket.Conn
	Queue chan []byte
	Mu    sync.Mutex
}

var numbers = 500

var appWg sync.WaitGroup

func TestApp(t *testing.T) {
	for i := 1; i < numbers; i++ {
		err := ClientManager.start(i, t)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("%d连接存储成功:", numbers)

	for j := 1; j < numbers; j++ {

		if conn, ok := ClientManager.client[j]; ok {
			//time.Sleep(time.Microsecond * 2)
			appWg.Add(1)
			conn.Send(j)

		}
	}

	appWg.Wait()
}

func (c *ConnMap) start(i int, t *testing.T) error {
	name := fmt.Sprintf("测试%d", i)
	token := jwt.GenerateToken(int64(i), name, "test", "2540463097@qq.com", 1)

	u := fmt.Sprintf("ws://%s/im/connect?token=%s", *addr, token)
	dialer := websocket.Dialer{
		NetDial:           nil,
		NetDialContext:    nil,
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   nil,
		HandshakeTimeout:  2 * time.Second,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		EnableCompression: false,
		Jar:               nil,
	}
	conn, _, err := dialer.Dial(u, nil)
	if err != nil {
		return err
	}

	mutexKey.Lock()
	c.client[i] = &Client{Index: i, Conn: conn, Token: token, Queue: make(chan []byte, 40)}
	mutexKey.Unlock()
	go c.client[i].write(t) //执行
	return nil
}

var mutexKey sync.Mutex

func (c *Client) write(t *testing.T) {
	// 关闭socket连接
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			t.Fatal(err)
		}
		appWg.Done()
	}(c.Conn)
	for {
		select {
		case message, ok := <-c.Queue:
			if !ok {
				// 关闭
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			//c.Mu.Lock()
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				t.Fatal(err)
			}
			//c.Mu.Unlock()
		}
	}
}

//消息推送
func (c *Client) Send(i int) {

	tId := random()
	data := fmt.Sprintf(`{"msg":"%s","from_id":%v,"to_id":%v,"status":0,"msg_type":%v,"channel_type":%v}`,
		"test", i, tId, 1, 1)
	//消息投递
	c.Queue <- []byte(data)
	c.Queue <- []byte(data)
	c.Queue <- []byte(data)
	close(c.Queue)
}

func random() int {
	max := big.NewInt(5000)
	i, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal("rand:", err)
	}
	return i.BitLen()
}
