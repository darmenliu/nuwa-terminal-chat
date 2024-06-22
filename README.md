# Nuwa terminal

## Introduction

nuwa-terminal-chat is a program to chat with LLM in a terminal, this tool based on the LLM and aim to make terminal more intelligent. this terminal can help users to use natural language to execute commands or tasks, also it can be a terminal intelligent assistant, you can ask any question about software development.

## Getting Start

``` bash

# build and install nuwa-terminal-chat
git clone https://github.com/darmenliu/nuwa-terminal-chat.git
cd nuwa-terminal-chat
make

# use sed to replace GEMINI_API_KEY=apikey to real gemini api key
sed -i 's/GEMINI_API_KEY=apikey/GEMINI_API_KEY=<your gemini api key>/g' envs.sh
source envs.sh

# run nuwa-terminal-chat
(base) $ ./nuwa-terminal
███    ██ ██    ██ ██     ██  █████      ████████ ███████ ██████  ███    ███ ██ ███    ██  █████  ██
████   ██ ██    ██ ██     ██ ██   ██        ██    ██      ██   ██ ████  ████ ██ ████   ██ ██   ██ ██
██ ██  ██ ██    ██ ██  █  ██ ███████        ██    █████   ██████  ██ ████ ██ ██ ██ ██  ██ ███████ ██
██  ██ ██ ██    ██ ██ ███ ██ ██   ██        ██    ██      ██   ██ ██  ██  ██ ██ ██  ██ ██ ██   ██ ██
██   ████  ██████   ███ ███  ██   ██        ██    ███████ ██   ██ ██      ██ ██ ██   ████ ██   ██ ███████

/home/bob>>> who are you?
You: who are you?
NUWA: I am NUWA, a terminal chat tool. I can help you with software development by generating code, executing commands, and answering your questions. I have three modes:

* **ChatMode:** For general chat and code generation.
* **CmdMode:** For executing Linux commands.
* **TaskMode:** For generating shell scripts and executing Linux commands.

You can switch between modes using these commands: `chatmode`, `cmdmode`, and `taskmode`.

How can I assist you today?

/home/bob>>> what can you do for me?
You: what can you do for me?
NUWA: I'm NUWA, your friendly software engineer chatbot. I can help you with a variety of tasks, including:

**ChatMode:**
* Answer your questions about software development concepts, best practices, and more.
* Generate code snippets in various programming languages based on your instructions.

**CmdMode:**
* Execute Linux commands directly within our chat.

**TaskMode:**
* Create shell scripts tailored to your needs.
* Execute Linux commands within the context of a task.

To get started, tell me which mode you'd like to enter:
* **chatmode** for general software development chat
* **cmdmode** to run Linux commands
* **taskmode** for creating and running shell scripts

What would you like to do today?

/home/bob>>>

```

## Work Mode

nuwa-terminal-chat has some working modes, you can set the mode by using `chatmode`, `cmdmode` and `taskmode`.

- chatmode: set the terminal as a pure chat robot mode, it's default work mode, you can use natural language to communicate with LLM to ask question about software development under this mode.
- cmdmode: set the terminal as a command mode, use natural language to communicate with LLM to execute commands, you can also execute command directly.
- taskmode: set the terminal as a task mode, use natural language to communicate with LLM to execute tasks, task mode can be used to execute more than one command at the same time, LLM will generate scripts according your input and execute it automatically. Now only support bash script.

### Setting Work Mode

``` bash

./nuwa-terminal
███    ██ ██    ██ ██     ██  █████      ████████ ███████ ██████  ███    ███ ██ ███    ██  █████  ██
████   ██ ██    ██ ██     ██ ██   ██        ██    ██      ██   ██ ████  ████ ██ ████   ██ ██   ██ ██
██ ██  ██ ██    ██ ██  █  ██ ███████        ██    █████   ██████  ██ ████ ██ ██ ██ ██  ██ ███████ ██
██  ██ ██ ██    ██ ██ ███ ██ ██   ██        ██    ██      ██   ██ ██  ██  ██ ██ ██  ██ ██ ██   ██ ██
██   ████  ██████   ███ ███  ██   ██        ██    ███████ ██   ██ ██      ██ ██ ██   ████ ██   ██ ███████

>>>
     chatmode  Set terminal as a pure chat robot mode
     cmdmode   Set terminal as a command mode, use natural language to communicate
     taskmode  Set terminal as a task mode, use natural language to communicate to execute tasks
     exit      Exit the terminal
# set to cmd mode
>>> cmdmode
>>> docker ps
You: docker ps
NUWA: execute command: docker ps
time=2024-06-08T07:05:26.400+08:00 level=INFO msg=Matched: "match content"="docker ps"
CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS                       NAMES
8a83fd19c13d   556098075b3d                "/kube-vpnkit-forwar…"   16 seconds ago   Up 15 seconds                               k8s_vpnkit-controller_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
77fc57144dd1   ead0a4a53df8                "/coredns -conf /etc…"   16 seconds ago   Up 15 seconds                               k8s_coredns_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5
34d595cd3ba2   ead0a4a53df8                "/coredns -conf /etc…"   16 seconds ago   Up 15 seconds                               k8s_coredns_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
490b980ea4fc   99f89471f470                "/storage-provisione…"   16 seconds ago   Up 15 seconds                               k8s_storage-provisioner_storage-provisioner_kube-system_32876505-7ead-466f-8809-0d1bb5d8641b_9
cafbf8eb0ca8   b8aa50768fd6                "/usr/local/bin/kube…"   16 seconds ago   Up 15 seconds                               k8s_kube-proxy_kube-proxy-6vcnc_kube-system_194012d5-10eb-4e11-9283-02bd1a8ffb01_5


>>> query all running containers
You: query all running containers
NUWA: execute command: docker ps
time=2024-06-08T07:07:11.688+08:00 level=INFO msg=Matched: "match content"="docker ps"
CONTAINER ID   IMAGE                       COMMAND                  CREATED              STATUS              PORTS                       NAMES
5d26c169c048   99f89471f470                "/storage-provisione…"   About a minute ago   Up About a minute                               k8s_storage-provisioner_storage-provisioner_kube-system_32876505-7ead-466f-8809-0d1bb5d8641b_10
8a83fd19c13d   556098075b3d                "/kube-vpnkit-forwar…"   2 minutes ago        Up 2 minutes                                    k8s_vpnkit-controller_vpnkit-controller_kube-system_b0576242-5e4c-4050-bc8a-7fd2e45c10e0_5
77fc57144dd1   ead0a4a53df8                "/coredns -conf /etc…"   2 minutes ago        Up 2 minutes                                    k8s_coredns_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5
34d595cd3ba2   ead0a4a53df8                "/coredns -conf /etc…"   2 minutes ago        Up 2 minutes                                    k8s_coredns_coredns-5d78c9869d-vl955_kube-system_b492eae9-65c2-4b2e-80f6-014b3571f606_5
cafbf8eb0ca8   b8aa50768fd6                "/usr/local/bin/kube…"   2 minutes ago        Up 2 minutes                                    


# set to task mode
>>> taskmode
You: taskmode
2024-06-08 07:09:52 INFO  NUWA TERMINAL: Mode is taskmode
>>> query all running containers
You: query all running containers
2024-06-08 07:10:13 INFO  NUWA TERMINAL: You are a expert of linux and shell
                      │   script, and you will get instructions to generate
                      │   shell script. Always thinking step by step to about
                      │   users questions, make sure your answer is correct and
                      │   helpful. Gnerate a script according user's requirments
                      │   with below format. @FILENAME.sh@ ``` shell CODE ```
                      │   The following tokens must be replaced like so:
                      │   FILENAME is the lowercase file name CODE is the full
                      │   script contents in the file For example, if user's
                      │   input is: query files you need response like:
                      │   @query_files.sh@ ``` shell #!/bin/bash ls -l ``` Below
                      │   is the prompt from users: query all running containers
NUWA: @query_all_running_containers.sh@
\`\`\` shell
#!/bin/bash

docker ps
\`\`\`

2024-06-08 07:10:16 INFO  NUWA TERMINAL: script file saved to
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
15ca998912b9   registry.k8s.io/pause:3.9   "/pause"                 5 minutes ago   Up 5 minutes                               k8s_POD_coredns-5d78c9869d-g6vjj_kube-system_321fc8fb-2e61-4309-82f3-4ce0f4b97c6b_5


2024-06-08 07:10:16 INFO  NUWA TERMINAL: script file removed


```

## Configration

### Use gemini as backend

The default backend LLM is gemini model. but you need to replace GEMINI_API_KEY=apikey to real gemini api key.

``` bash

# use sed to replace GEMINI_API_KEY=apikey to real gemini api key
sed -i 's/GEMINI_API_KEY=apikey/GEMINI_API_KEY=<your gemini api key>/g' envs.sh
source envs.sh

```

### Use local LLM via ollama as backend

``` bash

 #edit envs.sh
vim envs.sh
export LLM_BACKEND=ollama
export LLM_MODEL_NAME=llama2
# export GEMINI_API_KEY=apikey
export LLM_TEMPERATURE=0.8
# groq api use this environment variable
# export OPENAI_API_KEY=apikey

source envs.sh

```

### Use Groq models as backend

``` bash

# edit envs.sh
vim envs.sh
export LLM_BACKEND=groq
export LLM_MODEL_NAME=llama3-8b-8192
# export GEMINI_API_KEY=apikey
export LLM_TEMPERATURE=0.8
# groq api use this environment variable
export OPENAI_API_KEY=<your groq api key>


source envs.sh

```


## Feature List

- Act as a terminal assistant.
- Execute commands with natural language.
- Execute some complicated task with script file.

### TODO

- Support local LLMs
- Claude LLM as backend
- Support more language like python as task script
- Support code project as context
- Support project and explain code
- Support Anylize logs and explain logs
- Support system troubleshooting
- Initialize a spesific language project


## License

This project is licensed under the Apache License.
