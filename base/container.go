package base

type Container struct {
	Id      string `json:"_id"`
	Name    string `json:"name"`
	Status  int8   `json:"status"`
	ImageId string `json:"imageId"`
	Ports   int32  `json:"ports"` // 每个容器分配10个连续端口，此字段为起始端口
	Param   param  `json:"param"`
	Command string `json:"command"`
	HasSSH  bool   `json:"hasSSH"` // 是否已安装 ssh
	HasFTP  bool   `json:"hasFTP"` // 是否已安装 ftp
}

type param struct {
	Env      string   `json:"env"`      // 容器的环境变量
	Ports    []string `json:"ports"`    // 开放端口
	HostName string   `json:"hostname"` // 主机名
	Ip       string   `json:"ip"`       // ip
	Network  int8     `json:"network"`  // 网络模式: 0 主机，1 网桥
	Volumn   string   `json:"volumn"`   // 数据卷
	Restart  string   `json:"restart"`  // 容器重启策略
	Memory   string   `json:"memory"`   // 分配给容器的最大内存
	Cpus     string   `json:"cpus"`     // 分配给容器的最大 cpu
}
