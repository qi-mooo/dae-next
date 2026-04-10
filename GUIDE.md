# dae-next Guide

`dae-next` 是一个把 `dae` runtime 和 next-ui WebUI 打包在一起的包装层项目。
当前目标是：

- 尽量少改 `dae`
- 把 HTTP controller / WebUI 外壳放在父仓库
- 默认直接产出带内置 UI 的 `dae-next` 二进制

## 仓库结构

- `dae/`
  `dae` 子模块。这里保留 runtime、provider、配置读写和最小的命令行注入点。
- `cmd/dae-next/`
  `dae-next` 可执行入口。
- `internal/controllerapi/`
  父仓库里的 HTTP controller、WebUI 文件服务和日志流。
- `internal/daecontroller/`
  把 `dae` runtime 接到父仓库 controller shell 的桥接层。
- `ui/`
  WebUI 静态资源目录。
- `.github/workflows/build.yml`
  当前仓库的自动构建 workflow。

## 二进制行为

`dae-next` 默认会把 `ui/` embed 进二进制。
运行时会按下面的优先级决定 WebUI 从哪里提供：

1. `DAE_WEBUI_DIR`
2. 可执行文件附近或当前工作目录附近的 `ui/`
3. 系统共享目录里的 `ui/`
   常见位置：
   `/usr/share/dae/ui`
   `/usr/local/share/dae/ui`
4. 二进制内置的 embed UI

这意味着：

- 默认部署时，只放一个 `dae-next` 二进制也能提供 `/ui/`
- 如果磁盘上存在新的 `ui/` 目录，会优先覆盖内置 UI

## 运行条件

WebUI 要能从 `/ui/` 提供，仍然需要 `dae` 的 external controller 打开：

```dae
global {
  external_controller: "0.0.0.0:9090"
  external_controller_secret: "password"
}
```

访问地址通常是：

```text
http://<external_controller>/ui/#token=<secret>
```

例如：

```text
http://10.0.0.253:9090/ui/#token=password
```

## 本地构建

### 为什么要先生成 eBPF 绑定

`dae/control` 依赖 `bpf2go` 生成的 Go 文件。
这些生成物在 `dae` 子模块里是忽略文件，不直接提交进仓库。

所以干净 checkout 下，如果你直接执行：

```bash
go build ./cmd/dae-next
```

会看到类似下面的错误：

```text
undefined: bpfObjects
undefined: bpfTuplesKey
```

这不是代码本身坏了，而是还没先生成 `dae` 的 eBPF 绑定。

### Linux 构建

推荐顺序：

```bash
GOWORK=off make -C dae ebpf CLANG=clang-15 STRIP=llvm-strip-15
make dae-next
```

产物：

```text
build/dae-next
```

### macOS 本地开发

如果你本机有 Homebrew LLVM，可以这样：

```bash
GOWORK=off make -C dae ebpf CLANG=/usr/local/opt/llvm/bin/clang STRIP=llvm-strip TARGET=bpfel
make dae-next
```

如果你的 LLVM 安装在 Apple Silicon 默认路径，也可以改成：

```bash
GOWORK=off make -C dae ebpf CLANG=/opt/homebrew/opt/llvm/bin/clang STRIP=llvm-strip TARGET=bpfel
make dae-next
```

### 为什么要显式加 `GOWORK=off`

如果你的上层目录存在 `go.work`，`dae` 子模块里的 `go generate` 可能会吃到外层 workspace，导致 `make -C dae ebpf` 失败。

所以在生成 eBPF 绑定时，统一显式使用：

```bash
GOWORK=off
```

## GitHub Actions

当前 workflow 会：

1. checkout 主仓库和 submodule
2. 安装 `clang-15` 和 `llvm-15`
3. 执行：

```bash
GOWORK=off make -C dae ebpf CLANG=clang-15 STRIP=llvm-strip-15
```

4. 再编译 `dae-next`

当前 CI 只发布带 embed UI 的产物：

- `dae-next-linux-amd64`
- `dae-next-linux-arm64`

不会再发布 `noembed` 变体。

## OpenWrt / 路由器部署

如果你使用默认 embed 构建，最简单的部署方式是只替换二进制：

```bash
scp dae-next root@10.0.0.253:/tmp/dae.new
ssh root@10.0.0.253 'cp /tmp/dae.new /usr/bin/dae && chmod 755 /usr/bin/dae && /etc/init.d/dae restart'
```

如果你想覆盖内置 UI，再额外同步 `ui/` 目录：

```bash
tar -C ui -cf - . | ssh root@10.0.0.253 'rm -rf /usr/share/dae/ui && mkdir -p /usr/share/dae/ui && tar -C /usr/share/dae/ui -xf -'
```

部署完成后检查：

```text
http://10.0.0.253:9090/ui/
```

和：

```text
http://10.0.0.253:9090/ui/#token=password
```

## 磁盘上的 UI 什么时候有用

虽然默认已经 embed 了 UI，但磁盘上的 `ui/` 目录仍然有两个常见用途：

- 本地开发时实时覆盖内置 UI
- 已部署环境里快速替换前端，不重新发二进制

也就是说，embed 是默认兜底，不是禁止磁盘覆盖。

## 常见问题

### 1. `undefined: bpfObjects`

原因通常是还没先生成 `dae` 的 eBPF 绑定。

处理方式：

```bash
GOWORK=off make -C dae ebpf CLANG=clang-15 STRIP=llvm-strip-15
make dae-next
```

### 2. `make -C dae ebpf` 提示 workspace/module 错误

通常是外层 `go.work` 干扰。

处理方式：

```bash
GOWORK=off make -C dae ebpf ...
```

### 3. `/ui/` 是 404

优先检查：

- `external_controller` 是否已配置
- `external_controller_secret` 是否非空
- 你运行的是否真的是当前新二进制
- 如果使用 `noembed` 本地自编版本，是否真的提供了磁盘上的 `ui/`

### 4. 页面能打开，但接口都是未登录

优先检查：

- URL 里是否带了 `#token=...`
- token 是否和 `external_controller_secret` 一致
- `external_controller` 监听地址和你访问的地址是否一致

### 5. 修改了 `ui/` 但页面没变化

如果你的部署目录里仍有旧 UI 资源，磁盘覆盖优先级高于 embed。
所以应先确认当前实际命中的 UI 来源，再决定是删除磁盘目录还是继续覆盖它。

## 当前推荐方式

如果没有特别理由，建议直接使用默认 embed 构建：

- CI 产物更简单
- 部署只需要一个二进制
- 仍然保留磁盘 `ui/` 覆盖能力

这也是当前仓库默认维护和测试的路径。
