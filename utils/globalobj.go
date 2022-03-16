package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切有关Zinx框架的全局参数，供其他模块调用，一些参数是可以通过zinx.json由用户进行配置

type GlobalObj struct {
	//Server
	TcpSever ziface.IServer //当前Zinx全局的Server对象
	Host     string         //当前服务器主机监听的IP
	TcpPort  int            //当前服务器主机监听的端口
	Name     string         //当前服务器的名称

	//Zinx
	Version        string //当前Zinx版本号
	MaxConn        int    //当前服务器允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值
}

// GlobalObject 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

//Reload 从zinx.json中去加载用户自定义的配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	//如果文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从conf/zinx.json文件中去加载用户自定义的配置
	GlobalObject.Reload()
}
