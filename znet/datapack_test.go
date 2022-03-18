package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	/**
	模拟的服务端
	*/
	// 1 创建 socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// 创建一个 goroutine 负责从客户端处理业务
	go func() {
		// 2 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// -----> 拆包的过程 <-----
				// 定义一个拆包对象dp
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg有数据，需要进行第二次读取
						// 第二次从conn读，根据head的dataLen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						// 根据dataLen再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						// 完整的消息已经读取完毕
						fmt.Println("---> Recv MsgId: ", msg.Id, ", dataLen = ", msg.DataLen, "data = ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/**
	模拟的客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}
	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一起发送
	// 封装msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	fmt.Println(sendData1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	// 封装msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 6,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}