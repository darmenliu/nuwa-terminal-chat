# Nuwa 终端

## 简介

nuwa-terminal-chat 是一个在终端中与大语言模型(LLM)聊天的程序,该工具基于LLM,旨在使终端更加智能。这个终端可以帮助用户使用自然语言执行命令或任务,同时它还可以作为终端智能助手,你可以询问任何关于软件开发的问题。

## 开始使用

``` bash

# 构建并安装 nuwa-terminal-chat
git clone https://github.com/darmenliu/nuwa-terminal-chat.git
cd nuwa-terminal-chat
make

# 使用 sed 将 LLM_API_KEY=apikey 替换为真实的 api key
sed -i 's/LLM_API_KEY=apikey/LLM_API_KEY=<你的 api key>/g' envs.sh
source envs.sh

# 运行 nuwa-terminal-chat
(base) $ ./nuwa-terminal
███    ██ ██    ██ ██     ██  █████      ████████ ███████ ██████  ███    ███ ██ ███    ██  █████  ██
████   ██ ██    ██ ██     ██ ██   ██        ██    ██      ██   ██ ████  ████ ██ ████   ██ ██   ██ ██
██ ██  ██ ██    ██ ██  █  ██ ███████        ██    █████   ██████  ██ ████ ██ ██ ██ ██  ██ ███████ ██
██  ██ ██ ██    ██ ██ ███ ██ ██   ██        ██    ██      ██   ██ ██  ██  ██ ██ ██  ██ ██ ██   ██ ██
██   ████  ██████   ███ ███  ██   ██        ██    ███████ ██   ██ ██      ██ ██ ██   ████ ██   ██ ███████

/home/bob>>> 你是谁?
你: 你是谁?
NUWA: 我是NUWA,一个终端聊天工具。我可以通过生成代码、执行命令和回答问题来帮助你进行软件开发。我有三种模式:

* **聊天模式:** 用于一般聊天和代码生成。
* **命令模式:** 用于执行Linux命令。
* **任务模式:** 用于生成shell脚本并执行Linux命令。

你可以使用以下命令在模式之间切换: `chatmode`, `cmdmode`, 和 `taskmode`。

今天我能为你做些什么?

/home/bob>>> 你能为我做什么?
你: 你能为我做什么?
NUWA: 我是NUWA,你友好的软件工程师聊天机器人。我可以帮助你完成各种任务,包括:

**聊天模式:**
* 回答你关于软件开发概念、最佳实践等问题。
* 根据你的指示生成各种编程语言的代码片段。

**命令模式:**
* 直接在我们的聊天中执行Linux命令。

**任务模式:**
* 创建根据你需求定制的shell脚本。
* 在任务上下文中执行Linux命令。

要开始,告诉我你想进入哪种模式:
* **chatmode** 用于一般软件开发聊天
* **cmdmode** 用于运行Linux命令
* **taskmode** 用于创建和运行shell脚本

今天你想做什么?

/home/bob>>>

```

## 工作模式

nuwa-terminal-chat 有一些工作模式,你可以使用 `chatmode`, `cmdmode`, `taskmode`, 和 `agentmode` 来设置模式。

- chatmode: 将终端设置为纯聊天机器人模式,这是默认工作模式,你可以在这种模式下使用自然语言与LLM交流,询问有关软件开发的问题。
- cmdmode: 将终端设置为命令模式,使用自然语言与LLM交流以执行命令,你也可以直接执行命令。
- taskmode: 将终端设置为任务模式,使用自然语言与LLM交流以执行任务,任务模式可以用于同时执行多个命令,LLM将根据你的输入生成脚本并自动执行。目前只支持bash脚本。
- agentmode: 将终端设置为代理模式,用于复杂任务和故障排除。在这种模式下,LLM可以通过执行脚本使用各种工具来完成任务。

### 设置工作模式

``` bash

./nuwa-terminal
███    ██ ██    ██ ██     ██  █████      ████████ ███████ ███���██  ███    ███ ██ ███    ██  █████  ██
████   ██ ██    ██ ██     ██ ██   ██        ██    ██      ██   ██ ████  ████ ██ ████   ██ ██   ██ ██
██ ██  ██ ██    ██ ██  █  ██ ███████        ██    █████   ██████  ██ ████ ██ ██ ██ ██  ██ ███████ ██
██  ██ ██ ██    ██ ██ ███ ██ ██   ██        ██    ██      ██   ██ ██  ██  ██ ██ ██  ██ ██ ██   ██ ██
██   ████  ██████   ███ ███  ██   ██        ██    ███████ ██   ██ ██      ██ ██ ██   ████ ██   ██ ███████

>>>
     chatmode  将终端设置为纯聊天机器人模式
     cmdmode   将终端设置为命令模式,使用自然语言交流
     taskmode  将终端设置为任务模式,使用自然语言交流执行任务
     agentmode 将终端设置为代理模式,用于复杂任务和故障排除
     exit      退出终端
# 设置为命令模式
>>> cmdmode
>>> docker ps
你: docker ps
NUWA: 执行命令: docker ps
time=2024-06-08T07:05:26.400+08:00 level=INFO msg=Matched: "match content"="docker ps"
CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS                       NAMES
8a83fd19c13d   556098075b3d                "/kube-vpnkit-forwar…"   16 seconds ago   Up 15 seconds                               k8s_vpnkit-controller_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
77fc57144dd1   ead0a4a53df8                "/coredns -conf /etc…"   16 seconds ago   Up 15 seconds                               k8s_coredns_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5
34d595cd3ba2   ead0a4a53df8                "/coredns -conf /etc…"   16 seconds ago   Up 15 seconds                               k8s_coredns_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
490b980ea4fc   99f89471f470                "/storage-provisione…"   16 seconds ago   Up 15 seconds                               k8s_storage-provisioner_storage-provisioner_kube-system_32876505-7ead-466f-8809-0d1bb5d8641b_9
cafbf8eb0ca8   b8aa50768fd6                "/usr/local/bin/kube…"   16 seconds ago   Up 15 seconds                               k8s_kube-proxy_kube-proxy-6vcnc_kube-system_194012d5-10eb-4e11-9283-02bd1a8ffb01_5


>>> 查询所有运行中的容器
你: 查询所有运行中的容器
NUWA: 执行命令: docker ps
time=2024-06-08T07:07:11.688+08:00 level=INFO msg=Matched: "match content"="docker ps"
CONTAINER ID   IMAGE                       COMMAND                  CREATED              STATUS              PORTS                       NAMES
5d26c169c048   99f89471f470                "/storage-provisione…"   About a minute ago   Up About a minute                               k8s_storage-provisioner_storage-provisioner_kube-system_32876505-7ead-466f-8809-0d1bb5d8641b_10
8a83fd19c13d   556098075b3d                "/kube-vpnkit-forwar…"   2 minutes ago        Up 2 minutes                                    k8s_vpnkit-controller_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
77fc57144dd1   ead0a4a53df8                "/coredns -conf /etc…"   2 minutes ago        Up 2 minutes                                    k8s_coredns_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5
34d595cd3ba2   ead0a4a53df8                "/coredns -conf /etc…"   2 minutes ago        Up 2 minutes                                    k8s_coredns_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
cafbf8eb0ca8   b8aa50768fd6                "/usr/local/bin/kube…"   2 minutes ago        Up 2 minutes                                    


# 设置为任务模式
>>> taskmode
你: taskmode
2024-06-08 07:09:52 INFO  NUWA TERMINAL: 模式是taskmode
>>> 查询所有运行中的容器
你: 查询所有运行中的容器
2024-06-08 07:10:13 INFO  NUWA TERMINAL: 你是Linux和shell脚本的专家,你将获得指令来生成shell脚本。
                      │   始终逐步思考用户的问题,确保你的回答是正确和有帮助的。根据用户的要求
                      │   使用以下格式生成脚本。@FILENAME.sh@ ``` shell CODE ```
                      │   以下标记必须按如下方式替换:
                      │   FILENAME 是小写文件名 CODE 是文件中的完整脚本内容
                      │   例如,如果用户的输入是: 查询文件 你需要回复如下:
                      │   @query_files.sh@ ``` shell #!/bin/bash ls -l ```
                      │   以下是用户的提示: 查询所有运行中的容器
NUWA: @query_all_running_containers.sh@
\`\`\` shell
#!/bin/bash

docker ps
\`\`\`

2024-06-08 07:10:16 INFO  NUWA TERMINAL: 脚本文件已保存到
                      │   /home/ubuntu/.nuwa-terminal/scripts/query_all_running_containers.sh
+ docker ps
CONTAINER ID   IMAGE                       COMMAND                  CREATED         STATUS         PORTS                       NAMES
5d26c169c048   99f89471f470                "/storage-provisione…"   4 minutes ago   Up 4 minutes                               k8s_storage-provisioner_storage-provisioner_kube-system_32876505-7ead-466f-8809-0d1bb5d8641b_10
8a83fd19c13d   556098075b3d                "/kube-vpnkit-forwar…"   5 minutes ago   Up 5 minutes                               k8s_vpnkit-controller_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
77fc57144dd1   ead0a4a53df8                "/coredns -conf /etc…"   5 minutes ago   Up 5 minutes                               k8s_coredns_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5
34d595cd3ba2   ead0a4a53df8                "/coredns -conf /etc…"   5 minutes ago   Up 5 minutes                               k8s_coredns_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
cafbf8eb0ca8   b8aa50768fd6                "/usr/local/bin/kube…"   5 minutes ago   Up 5 minutes                               k8s_kube-proxy_kube-proxy-6vcnc_kube-system_194012d5-10eb-4e11-9283-02bd1a8ffb01_5
61b0e516f8f5   registry.k8s.io/pause:3.9   "/pause"                 5 minutes ago   Up 5 minutes                               k8s_POD_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
505a6ac2ddb5   registry.k8s.io/pause:3.9   "/pause"                 5 minutes ago   Up 5 minutes                               k8s_POD_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
15ca998912b9   registry.k8s.io/pause:3.9   "/pause"                 5 minutes ago   Up 5 minutes                               k8s_POD_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce

# 设置为代理模式
>>> agentmode
你: agentmode
2024-06-08 08:30:00 INFO  NUWA TERMINAL: 模式是agentmode
>>> 分析系统日志以识别可能表明潜在问题的最近错误或警告。

思考: 为了识别系统日志中的最近错误或警告,我需要搜索包含关键词"error"或"warning"的日志文件条目。我将创建一个shell脚本,搜索最常见的日志文件(如`/var/log/syslog`, `/var/log/messages`, 和 `/var/log/kern.log`)中的这些关键词,并显示相关条目。

行动: ScriptExecutor
行动输入:

<NUWA生成的脚本>

观察结果

<脚本输出>

思考: 为了识别系统日志中的最近错误或警告,我需要搜索包含关键词"error"或"warning"的日志文件条目。我将创建一个shell脚本,搜索最常见的日志文件(如`/var/log/syslog`, `/var/log/messages`, 和 `/var/log/kern.log`)中的这些关键词,并显示相关条目。

行动: ScriptExecutor
行动输入:

<NUWA生成的脚本>

观察结果

<脚本输出>

NUWA: 系统日志中包含几个可能表明潜在问题的错误和警告。值得注意的错误包括:
1. **networkd-dispatcher 错误**: 网络接口的未知状态。
2. **systemd-xdg-autostart-generator 错误**: 解析 Exec= 行时出错。
3. **subiquity 错误**: 自动安装过程中的错误。
4. **networkd-dispatcher 警告**: systemd-networkd 未运行。

这些错误和警告应该进一步调查以确定根本原因并解决问题。
```

## 配置

### 使用本地 LLM 通过 ollama 作为后端

``` bash

 #编辑 envs.sh
vim envs.sh
export LLM_BACKEND=ollama
export LLM_MODEL_NAME=llama2
export LLM_API_KEY=apikey
export LLM_TEMPERATURE=0.8
export OLLAMA_SERVER_URL=http://localhost:8000

source envs.sh

```

### 使用 Groq 模型作为后端

``` bash

# 编辑 envs.sh
vim envs.sh
export LLM_BACKEND=groq
export LLM_MODEL_NAME=llama3-8b-8192
export LLM_API_KEY=<你的 groq api key>
export LLM_TEMPERATURE=0.8

source envs.sh

```

## 功能列表

- 充当终端助手。
- 使用自然语言执行命令。
- 使用脚本文件执行一些复杂任务。
- 使用代理模式执行复杂任务和故障排除。

### 待办事项

- 支持本地 LLMs
- Claude LLM 作为后端
- 支持更多语言如 Python 作为任务脚本
- 支持代码项目作为上下文
- 支持项目和解释代码
- 支持分析日志和解释日志
- 支持系统故障排除
- 初始化特定语言项目

## 许可证

本项目采用 Apache 许可证。
