bluebell_notes

#### 1.1 `./scaffold_demo.exe` 因为MySQL和Redis运行failed

**发现错误init mysql failed**
终端`go build`后成功出现一个scaffold_demo.exe文件， 在终端输入`./scaffold_demo.exe` 运行项目，
但发现错误init mysql failed, error: dial tcp 127.0.0.1:13306: connectex: No connection could be made because the target machine actively refused it.
解决：
在powershell中输入以下:（记得先启动docker，并注意密码名称和config里面写的一致，不一致就删掉容器重写）

```bash
            docker run -d `
            --name mysql-scaffold `
            -p 13306:3306 `
            -e MYSQL_ROOT_PASSWORD=root `
            -e MYSQL_DATABASE=scaffold_db `
            mysql:8.0
```

解析：-d (Detach): 代表“后台运行”。如果不加这个参数，你的终端窗口会被 MySQL 的运行日志占满，一旦关闭终端，数据库也就停止了。加上 -d 后，容器会在后台静默运行。
-p 13306:3306 (最关键的一步)，原理: 端口映射（Host Port : Container Port）。3306: 是容器内部 MySQL 服务默认监听的端口。13306: 是你电脑（宿主机）暴露出来的端口。

**发现错误init redis failes**
在终端输入`./scaffold_demo.exe` 运行项目但又发现错误init redis failes, error: dial tcp 127.0.0.1:16379: connectex: No connection could be made because the target machine actively refused it.
在powershell中输入以下（记得先启动docker，并注意密码名称和config里面写的一致，不一致就删掉容器重写）

```bash
        docker run -d `
        --name redis-scaffold `
        -p 16379:6379 `
        redis:7.0
```

**睡觉Docker 里的数据库要关吗？**
你的 MySQL 和 Redis 容器是在后台（-d 模式）运行的，即便你关掉 VS Code 甚至关掉电脑，只要你不手动停止，它们下次开机通常会随 Docker Desktop 自动启动。
如果你想彻底释放电脑内存，可以在 PowerShell 输入：`docker stop mysql-scaffold redis-scaffold`
下次启动Docker后，在终端输入 `docker ps` 如果列表里没有 mysql-scaffold，说明它们停了。输入 `docker start mysql-scaffold redis-scaffold` 唤醒它们。

### 1.2 **直接关闭Vscode和使用ctrl+C的区别**

Ctrl + C (Server exiting)：
        原理：这发送了一个 SIGINT（中断）信号给你的 Go 程序。
        过程：你的程序接收到信号后，主动停止接收新请求，并开始执行你在 main.go 中写的收尾工作（比如 db.Close() 或 redis.Close()），最后打印出 "Server exiting" 并退出。
        优点：数据最安全，不会因为数据库连接突然中断导致事务异常。
直接关闭 VS Code (Shutdown Server ...)
        原理：VS Code 会向它启动的所有子进程发送一个关闭指令。
        过程：你的程序识别到了环境变化（通常是 context 取消或 SIGTERM 信号），触发了另一段清理逻辑。
        优点：即使你偷懒直接关窗口，程序也会尽力帮你把“烂摊子”收拾好。
        
端口冲突：如果看到 bind: address already in use，说明你今晚的 Shutdown Server 并没有彻底杀死进程（偶尔会发生）。这时候去任务管理器杀掉 scaffold_demo.exe 即可。

### 1.3指定配置文件的路径

**方式一：SetCongigFile直接指定路径**
`viper.SetConfigFile("./config.yaml)` //相对路径
或
`viper.SetConfigFile("/Users/Ziying/.../config.yaml)` //绝对路径

**方式二：指定文件名和位置，viper自行查找文件**
`viper.SetConfigName("config")` //指定文件名(不带后缀)

+ `viper.AddConfigPath(".") `//指定路径（相对）, 配置文件位置可以是多个
  `viper.AddConfigPath("./config") `
  此方式需要注意的是，最好不要有同名的配置以防内容报错，如同时有config.json和config.yaml

**另外**
`viper.SetConfigType("yaml")`   
配合远程配置中心使用，解析配置前需要知道配置类型，告诉Viper是使用json还是yaml等什么格式去解析

### 1.4 Dockerfile + docker-compose.yml

**Dockerfile** ——「我这个程序怎么被打包」
用什么系统, 怎么编译, 运行时需要哪些文件, 最终怎么启动
👉 这是“镜像级”的事

**docker-compose.yml** ——「这些服务怎么一起跑」
MySQL, Redis,Go 服务,网络, 端口映射, 启动顺序
👉 这是“环境级”的事

**config.yaml** ——「程序运行时的业务配置」
服务端口, JWT 过期时间, 日志级别,数据库连接参数, Redis DB、连接池
👉 这是“应用级”的事, 它是给 Go 程序读的，不是给 Docker 读的。

4. 启动环境
   **step1.确保准备工作**
   你在项目根目录，里面有：
   Dockerfile
   docker-compose.yml
   wait-for.sh  (可执行权限 chmod +x wait-for.sh)
   conf/config.yaml
   init.sql
   **step2.停掉旧容器（可选，但推荐）**
   `docker-compose down`
   避免端口冲突或者残留容器影响
   **step3.构建并启动所有服务**
   `docker-compose up --build -d`
   **step4. 有东西版本老旧下线了，所以报错了**
   重建镜像
   `docker-compose build --no-cache`
   `docker-compose up -d`
   **step5.查看容器是否启动成功**
   `docker ps`
   **step6.查看 Go 服务日志（确认启动成功）**
   `docker-compose logs -f bluebell_app`
   **step7测试服务**
   · Go HTTP 接口：
   假设你的 Go 服务监听 8084：
   浏览器访问：http://localhost:8888
   或`curl http://localhost:8888/health`（如果项目有健康检查接口 /health 或 /ping）
   · MySQL 连接测试
   `docker exec -it bluebell-mysql8019-1 mysql -uroot -proot1234` bash
   `SHOW DATABASES;` sql
   · Redis 连接测试
   `docker exec -it redis507_1 redis-cli` bash
   测试`PING` 返回：PONG
   **Step8.停止所有服务**
   `docker-compose down`

5. docker-compose down 做了什么

停止并删除容器、网络、默认的卷（如果加了 -v 参数会删卷），

不会删除镜像（image），所以镜像还在本地。

2️⃣ docker-compose up -d 的行为

默认情况下，Docker Compose 会检查镜像是否存在：

如果镜像存在 → 直接用现有镜像启动容器，不重新 build

如果镜像不存在 → 先 build，再启动容器

所以如果你之前 build 过镜像，docker-compose up -d 是 不会重新 build 的。

3️⃣ 什么时候需要强制 build

docker-compose build --no-cache 或 docker-compose up --build 会强制重建镜像

通常在修改了 Dockerfile 或依赖（比如 Go 模块）时才需要。

**docker-compose up --build -d**

作用：先 build 镜像（如果有变化的话）然后 启动容器
-d 表示后台运行（detached mode）
使用场景：Dockerfile 或依赖有改动， 想确保启动的容器用的是最新的镜像
例：修改了 Go 代码或安装了新依赖
特点：只会 build 有变化的服务镜像，比 --no-cache 快，因为会使用缓存
**docker-compose build --no-cache**
作用：只 build 镜像，不启动容器， --no-cache 表示 完全不使用缓存，每一步都重新构建
使用场景：你怀疑镜像缓存导致构建不干净，或者依赖源已经更新，想完全重新构建
特点：构建速度慢，因为不使用缓存；不会启动容器，需要单独用 docker-compose up 启动

6. 开发阶段
   先只在docker里运行mysql和redis,go还是在本地运行go run:
   **启动数据库和 Redis：** `docker-compose -f docker-compose.yml up -d` 
   或 简化版（只有一个 compose 文件）`docker-compose up -d`
   **本地运行 Go：** `go run . ./conf/dev.yaml`
   **浏览器访问或者POSTMAN** http://localhost:8084

7. 关闭程序:

7.1 停止本地 Go 程序

Ctrl + C

7.2 停止 Docker 容器，但保留数据

```
docker-compose -f docker-compose.yml stop
```

7.3  下次开发：

`docker-composse -f docker-compose.yml start`   # 启动 MySQL/Redis
`go run . ./conf/dev.yaml`                     # 启动 Go

8. 清除旧容器

8.1 停止并删除旧容器

`docker-compose -f docker-compose.yml down`
· `down`会停止并删除 当前 Compose 文件启动的所有容器
· 默认不会删除 Volume，除非你加末尾加 `-v`
`docker-compose -f docker-compose.yml down -v`删除容器+卷

8.2 删除无用的旧容器（可选）

· 查看所有容器：
`docker ps -a`
· 删除无用容器：
`docker rm <container_id_or_name>`
· 删除无用卷：
`docker volume ls`
`docker volume rm <volume_name>`

8.3使用新的docker-compose.yml 启动

· 启动数据库和 Redis： `docker-compose -f docker-compose.yml up -d` 
  或 简化版（只有一个 compose 文件）`docker-compose up -d`
· 本地运行 Go： `go run . ./conf/dev.yaml`
· 浏览器访问或者POSTMAN http://localhost:8084

9. 有关-f/-p/-v/--v等问题
   -f 选文件，-p 起名字，up 和 down 参数必须一模一样,否则删不干净
   **-p 同一份 yml，起多套环境**
   `docker-compose -p bluebell_dev up -d` 和`docker-compose -p bluebell_dev down`
   `docker-compose -p bluebell_test up -d` 和`docker-compose -p bluebell_test down`
   **-f 指定哪个yml, 不写-f 就默认只会用 docker-compose.yml**
   通常有docker-compose.yml,docker-compose.dev.yml,docker-compose.prod.yml
   `docker-compose -f docker-compose.dev.yml up -d`
   `docker-compose -f docker-compose.dev.yml down`
   **-f和-p都写**
   `docker-compose -f docker-compose.dev.yml -p bluebell_dev up -d`
   `docker-compose -f docker-compose.dev.yml -p bluebell_dev down`
   **-v的作用**

   `docker-compose -f docker-compose.dev.yml -p bluebell_dev down -v`

   down：停容器 + 删容器 + 删 network
   -v：只删除这个 compose 用到的 volume

   services:
     mysql:
    volumes:

      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:

### 1.5 在根目录添加了.gitattributes和.gitignore执行

`git add --renormalize .` 别忘了末尾的.
`git commit -m "Add .gitattributes and .gitignore; normalize line endings"`

在根目录添加了.gitattributes和.gitignore执行 git add --renormalize . 别忘了末尾的. git commit -m "Add .gitattributes and .gitignore; normalize line endings" git顺序,在github创建一个名为blubell的repository： git init git checkout -b main git add .gitignore .gitattributes git add . git commit -m "chore: initial commit with .gitignore and .gitattributes" git remote add origin https://github.com/Zoe6486/bluebell.git git push -u origin main
新建一个功能分支： git checkout -b feature/signup 本地开发完成后 push 到远程：git push -u origin feature/login Pull Request (PR)： 功能完成后发 PR 到 main， 团队成员 Code Review，CI/CD 自动跑：语法检查 (lint)， 自动单元/集成测试，构建打包，部署到测试环境（可选） Merge / Release：PR 审核通过后 merge 到 main，main 可以打 release tag：v1.0.0，CI/CD 自动部署到生产环境 功能分支合并后可以删除： git branch -d feature/login # 本地删除 git push origin --delete feature/login # 远程删除

### 1.6 创建新的分支

Step 1: 确保 main 最新
git checkout main
git pull origin main

Step 2: 从 main 新建 feature 分支
git checkout -b feature/signup

Step 3: 在 feature/signup 开发

###### 修改文件 后 

git add . 
git commit -m "feat: signup"

Step 4: push 分支到远程
git push -u origin feature/signup

Step 5: 发 PR 到 main → CI + Code Review

Step 6: PR 审核通过 → merge 到 main
注意：merge 前不要切回 main 改动功能代码！！！

###### merge 完后：

1. 切回本地 main

git checkout main

2. 拉取远程最新

git pull origin main

3. 确认状态

git log --oneline -5   # 看最近提交，应该看到 merge commit

4. 删除远程 feature 分支

git push origin --delete feature/zzy

5. 删除本地 feature 分支

git branch -d feature/zyy (-d会检查merge是否成功)

git branch -D feature/zzy（强制删除）

6. 确认分支已经删除：

git branch          # 查看本地分支
git branch -r       # 查看远程分支

7.清理残留引用（可选，但推荐）：

git fetch -p

8. 确认 main 分支最新：

git checkout main
git pull origin main
git log --oneline -5

9. 再创建新的分支接着开发：

git checkout -b feature/abc





### 1.7 安装 golang-migrate（CLI）

https://github.com/golang-migrate/migrate/releases找到migrate.windows-amd64.zip下载

解压后得到migrate.exe,把它放到C:\Windows\System32\下 或放到C:\tools\migrate\migrate.exe然后把 `C:\tools\migrate` 加到 **PATH**

验证是否成功：打开 **新的** PowerShell / CMD输入 `migrate -version`看是否输出4.19.1类似版本

在项目根目录下创建一个 `migrations` 文件夹，然后运行命令生成第一对迁移脚本：

```bash
migrate create -ext sql -dir ./migrations -seq create_users_table
```

`migrations` 目录下生成了两个文件：

`000001_create_users_table.up.sql` (升级：写建表语句)

`000001_create_users_table.down.sql` (回滚：写删表语句)

写入sql语句

在项目根目录下创建cmd/migrate/main.go

```go
package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://migrations",
		"mysql://root:root1234@tcp(127.0.0.1:23306)/db_bluebell?multiStatements=true",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("migrate up success")
}

```



执行命令：`go run cmd/migrate/main.go`出现success就好了

postman测试：一个用户signup成功后，如果想查看数据库，直接打开DBeaver

### 1️⃣ 新建连接

1. 打开 DBeaver → 点击 **Database → New Database Connection**。
2. 选择 **MySQL**（如果没看到可以在搜索框里输入 “MySQL”）。
3. 点击 **Next**。

------

### 2️⃣ 填写连接信息

根据你的 Docker Compose 配置：

| 配置项   | 值          |
| -------- | ----------- |
| Host     | 127.0.0.1   |
| Port     | 23306       |
| Database | db_bluebell |
| Username | root        |
| Password | root1234    |

------

### 3️⃣ 测试连接

- 点击 **Test Connection**
- 如果成功，就可以点击 **Finish** 保存连接,
- 可能会提示缺失mysql-connector-j-8.0.33.jar + protobuf-java-3.21.9.jar，按照提示download就好了