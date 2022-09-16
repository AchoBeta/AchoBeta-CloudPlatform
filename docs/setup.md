## 环境搭建
运行目录 /project/script
### 1. Docker
若服务器上已有 Docker，则可以跳过此步骤。
- ./sop.py --docker install     安装 docker
- ./sop.py --docker update      更新 docker
### 2. Database
搭建数据库环境。
- `./sop.py --database start`   搭建并启动 database
- `./sop.py --database stop`    停止 database
- `./sop.py --database start`   重启 database
- `./sop.py --database update`  更新 database
### 3. Webssh
- `./sop.py --webssh start`     搭建并启动 webssh
- `./sop.py --webssh stop`      停止 webssh
- `./sop.py --webssh start`     重启 webssh
- `./sop.py --webssh update`    更新 webssh
### 3. Abcp
待完善
### 4. 快速执行所有服务
- `./sop.py --all start`    启动所有服务
- `./sop.py --all stop`     停止所有服务
- `./sop.py --all restart`  重启所有服务