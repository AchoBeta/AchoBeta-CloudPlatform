package cloud

type Container struct {
	Id        string `bson:"_id"`
	Name      string `bson:"name"`
	Status    int8   `bson:"status"` // 0 运行, 1 停止
	Username  string `bson:"username"`
	Password  string `bson:"password"`
	Image     string `bson:"image"`
	Ports     int    `bson:"ports"` // 每个容器分配10个连续端口，此字段为起始端口
	Ip        string `bson:"ip"`
	Param     Param  `bson:"param"`
	Command   string `bson:"command"`
	StartTime int64  `bson:"startTime"`
}

type Param struct {
	Env      []string `bson:"env"`      // 容器的环境变量
	Ports    []int    `bson:"ports"`    // 开放端口
	HostName string   `bson:"hostname"` // 主机名
	//Network  int8     `bson:"network"`  // 网络模式: 0 主机，1 网桥
	//Volumn   string   `bson:"volumn"`   // 数据卷
	//Restart  string   `bson:"restart"`  // 容器重启策略
	//Memory   string   `bson:"memory"`   // 分配给容器的最大内存
	//Cpus     string   `bson:"cpus"`     // 分配给容器的最大 cpu
}
