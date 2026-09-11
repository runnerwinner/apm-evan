只要你改了 Dockerfile（或它 `COPY` 到的目录结构），就要重新 build 对应镜像并重建容器。

你这个项目最实用的命令是按影响范围来：

1. 只改了某一个服务相关目录  
   例如只影响 ordersvc：

```bash
cd /home/evan/backend/apm-evan
docker compose -f conf/docker-compose.yml up -d --build ordersvc
```

2. 改了通用 Dockerfile（你现在是三服务共用）  
   建议重建三个业务服务：

```bash
cd /home/evan/backend/apm-evan
docker compose -f conf/docker-compose.yml up -d --build ordersvc skusvc usrsvc
```

3. 想确保一定用最新镜像（不吃旧容器）

```bash
cd /home/evan/backend/apm-evan
docker compose -f conf/docker-compose.yml up -d --build --force-recreate ordersvc skusvc usrsvc
```

4. 如果你还担心缓存导致旧内容没更新

```bash
cd /home/evan/backend/apm-evan
docker compose -f conf/docker-compose.yml build --no-cache ordersvc skusvc usrsvc
docker compose -f conf/docker-compose.yml up -d ordersvc skusvc usrsvc
```

验证是否部署成功：

1. 看状态

```bash
docker compose -f conf/docker-compose.yml ps ordersvc skusvc usrsvc
```

2. 看日志

```bash
docker compose -f conf/docker-compose.yml logs -f --tail=100 ordersvc skusvc usrsvc
```

一句话：改 Dockerfile 后，执行 `up -d --build`（按受影响服务范围）就是标准重新部署方式。

-- -////////////////////////////////////////////////////////////////////////////////

如果没改 Dockerfile，只是改了代码文件，也需要重新构建镜像再重启对应服务。

直接用这条就够了：

```bash
cd /home/evan/backend/apm-evan
docker compose -f conf/docker-compose.yml up -d --build ordersvc skusvc usrsvc
```

按范围更精细地用：

1. 只改了 ordersvc 代码

```bash
docker compose -f conf/docker-compose.yml up -d --build ordersvc
```

2. 改了 dogapm（三服务都依赖）

```bash
docker compose -f conf/docker-compose.yml up -d --build ordersvc skusvc usrsvc
```

3. 改了 protos（通常也影响多个服务）

```bash
docker compose -f conf/docker-compose.yml up -d --build ordersvc skusvc usrsvc
```

什么时候可以不用 `--build`：

1. 你只是改了容器运行参数（例如 compose 里环境变量），没改代码镜像内容
2. 这时可用：

```bash
docker compose -f conf/docker-compose.yml up -d ordersvc
```

你当前这套（代码有变动）建议默认带 `--build`，最稳。
