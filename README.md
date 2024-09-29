# Nuwa terminal

## Introduction

nuwa-terminal-chat is a program to chat with LLM in a terminal, this tool based on the LLM and aim to make terminal more intelligent. this terminal can help users to use natural language to execute commands or tasks, also it can be a terminal intelligent assistant, you can ask any question about software development.

## Getting Start

``` bash

# build and install nuwa-terminal-chat
git clone https://github.com/darmenliu/nuwa-terminal-chat.git
cd nuwa-terminal-chat
make

# use sed to replace LLM_API_KEY=apikey to real api key
sed -i 's/LLM_API_KEY=apikey/LLM_API_KEY=<your api key>/g' envs.sh
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

nuwa-terminal-chat has some working modes, you can set the mode by using `chatmode`, `cmdmode`, `taskmode`, and `agentmode`.

- chatmode: set the terminal as a pure chat robot mode, it's default work mode, you can use natural language to communicate with LLM to ask question about software development under this mode.
- cmdmode: set the terminal as a command mode, use natural language to communicate with LLM to execute commands, you can also execute command directly.
- taskmode: set the terminal as a task mode, use natural language to communicate with LLM to execute tasks, task mode can be used to execute more than one command at the same time, LLM will generate scripts according your input and execute it automatically. Now only support bash script.
- agentmode: set the terminal as an agent mode for complex tasks and troubleshooting. In this mode, LLM can use various tools to complete tasks via executing scripts.

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
     agentmode Set terminal as an agent mode for complex tasks and troubleshooting
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

# set to agent mode
>>> agentmode
You: agentmode
2024-06-08 08:30:00 INFO  NUWA TERMINAL: Mode is agentmode
>>> Analyze the system logs to identify any recent errors or warnings that might indicate a potential issue.

Thought: To identify recent errors or warnings in the system logs, I need to search through the log files for entries that contain the keywords "error" or "warning". I will create a shell script that searches through the most common log files (such as `/var/log/syslog`, `/var/log/messages`, and `/var/log/kern.log`) for these keywords and displays the relevant entries.

Action: ScriptExecutor
Action_input:

@analyze_logs.sh@
``` shell
#!/bin/bash

# Define the log files to search
LOG_FILES=(/var/log/syslog /var/log/messages /var/log/kern.log)

# Define the keywords to search for
KEYWORDS=("error" "warning")

# Function to search logs for keywords
search_logs() {
    for log_file in "${LOG_FILES[@]}"; do
        if [ -f "$log_file" ]; then
            echo "Searching $log_file for errors and warnings..."
            for keyword in "${KEYWORDS[@]}"; do
                echo "Entries containing '$keyword':"
                grep -i "$keyword" "$log_file" | tail -n 10
                echo "----------------------------------------"
            done
        else
            echo "Log file $log_file does not exist."
        fi
    done
}

# Execute the search
search_logs
```

Observation
2024-09-29 09:08:43 INFO  Matched:
                      └ match content for tool name:: Action: ScriptExecutor
Action_input:
2024-09-29 09:08:43 INFO  Start to parse the script from input:
                      └ input: Question: Analyze the system logs to identify any recent errors or warnings that might indicate a potential issue.

Thought: To identify recent errors or warnings in the system logs, I need to search through the log files for entries that contain the keywords "error" or "warning". I will create a shell script that searches through the most common log files (such as `/var/log/syslog`, `/var/log/messages`, and `/var/log/kern.log`) for these keywords and displays the relevant entries.

Action: ScriptExecutor
Action_input:

@analyze_logs.sh@
``` shell
#!/bin/bash

# Define the log files to search
LOG_FILES=(/var/log/syslog /var/log/messages /var/log/kern.log)

# Define the keywords to search for
KEYWORDS=("error" "warning")

# Function to search logs for keywords
search_logs() {
    for log_file in "${LOG_FILES[@]}"; do
        if [ -f "$log_file" ]; then
            echo "Searching $log_file for errors and warnings..."
            for keyword in "${KEYWORDS[@]}"; do
                echo "Entries containing '$keyword':"
                grep -i "$keyword" "$log_file" | tail -n 10
                echo "----------------------------------------"
            done
        else
            echo "Log file $log_file does not exist."
        fi
    done
}

# Execute the search
search_logs
```

Observation
time=2024-09-29T09:08:43.649+08:00 level=INFO msg="Adding file to source file dict" content="@analyze_logs.sh@\n``` shell\n#!/bin/bash\n\n# Define the log files to search\nLOG_FILES=(/var/log/syslog /var/log/messages /var/log/kern.log)\n\n# Define the keywords to search for\nKEYWORDS=(\"error\" \"warning\")\n\n# Function to search logs for keywords\nsearch_logs() {\n    for log_file in \"${LOG_FILES[@]}\"; do\n        if [ -f \"$log_file\" ]; then\n            echo \"Searching $log_file for errors and warnings...\"\n            for keyword in \"${KEYWORDS[@]}\"; do\n                echo \"Entries containing '$keyword':\"\n                grep -i \"$keyword\" \"$log_file\" | tail -n 10\n                echo \"----------------------------------------\"\n            done\n        else\n            echo \"Log file $log_file does not exist.\"\n        fi\n    done\n}\n\n# Execute the search\nsearch_logs\n```"
time=2024-09-29T09:08:43.649+08:00 level=INFO msg=Matched: "match_content "=analyze_logs.sh
time=2024-09-29T09:08:43.649+08:00 level=INFO msg=Matched: "match_content "="#!/bin/bash\n\n# Define the log files to search\nLOG_FILES=(/var/log/syslog /var/log/messages /var/log/kern.log)\n\n# Define the keywords to search for\nKEYWORDS=(\"error\" \"warning\")\n\n# Function to search logs for keywords\nsearch_logs() {\n    for log_file in \"${LOG_FILES[@]}\"; do\n        if [ -f \"$log_file\" ]; then\n            echo \"Searching $log_file for errors and warnings...\"\n            for keyword in \"${KEYWORDS[@]}\"; do\n                echo \"Entries containing '$keyword':\"\n                grep -i \"$keyword\" \"$log_file\" | tail -n 10\n                echo \"----------------------------------------\"\n            done\n        else\n            echo \"Log file $log_file does not exist.\"\n        fi\n    done\n}\n\n# Execute the search\nsearch_logs"
2024-09-29 09:08:43 INFO  Start to execute the script:
                      └ scriptfile: /home/yliu/.nuwa-terminal/scripts/analyze_logs.sh
2024-09-29 09:08:43 INFO  Script executed successfully, output:
                      └ output: + LOG_FILES=(/var/log/syslog /var/log/messages /var/log/kern.log)
+ KEYWORDS=("error" "warning")
+ search_logs
+ for log_file in "${LOG_FILES[@]}"
+ '[' -f /var/log/syslog ']'
+ echo 'Searching /var/log/syslog for errors and warnings...'
Searching /var/log/syslog for errors and warnings...
+ for keyword in "${KEYWORDS[@]}"
+ echo 'Entries containing '\''error'\'':'
Entries containing 'error':
+ grep -i error /var/log/syslog
+ tail -n 10
Sep 29 08:13:43 PC-20220101GUTO networkd-dispatcher[208]: ERROR:Unknown state for interface NetworkctlListState(idx=1, name='lo', type='loopback', operational='n/a', administrative='unmanaged'): n/a
Sep 29 08:13:43 PC-20220101GUTO networkd-dispatcher[208]: ERROR:Unknown state for interface NetworkctlListState(idx=2, name='eth0', type='ether', operational='n/a', administrative='unmanaged'): n/a
Sep 29 08:13:43 PC-20220101GUTO systemd-xdg-autostart-generator[422]: Not generating service for XDG autostart app-at\x2dspi\x2ddbus\x2dbus@autostart.service, error parsing Exec= line: No such file or directory
Sep 29 08:23:46 PC-20220101GUTO subiquity_event.300[300]: subiquity/Error/apply_autoinstall_config
Sep 29 08:23:46 PC-20220101GUTO subiquity_event.300[300]: subiquity/Error/apply_autoinstall_config
+ echo ----------------------------------------
----------------------------------------
+ for keyword in "${KEYWORDS[@]}"
+ echo 'Entries containing '\''warning'\'':'
Entries containing 'warning':
+ grep -i warning /var/log/syslog
+ tail -n 10
Sep 29 08:13:43 PC-20220101GUTO networkd-dispatcher[296]: WARNING: systemd-networkd is not running, output will be incomplete.
+ echo ----------------------------------------
----------------------------------------
+ for log_file in "${LOG_FILES[@]}"
+ '[' -f /var/log/messages ']'
+ echo 'Log file /var/log/messages does not exist.'
Log file /var/log/messages does not exist.
+ for log_file in "${LOG_FILES[@]}"
+ '[' -f /var/log/kern.log ']'
+ echo 'Searching /var/log/kern.log for errors and warnings...'
Searching /var/log/kern.log for errors and warnings...
+ for keyword in "${KEYWORDS[@]}"
+ echo 'Entries containing '\''error'\'':'
Entries containing 'error':
+ grep -i error /var/log/kern.log
+ tail -n 10
+ echo ----------------------------------------
----------------------------------------
+ for keyword in "${KEYWORDS[@]}"
+ echo 'Entries containing '\''warning'\'':'
Entries containing 'warning':
+ grep -i warning /var/log/kern.log
+ tail -n 10
+ echo ----------------------------------------
----------------------------------------

NUWA:  The system logs contain several errors and warnings that might indicate potential issues. Notable errors include:
1. **networkd-dispatcher errors**: Unknown state for network interfaces.
2. **systemd-xdg-autostart-generator error**: Error parsing Exec= line.
3. **subiquity errors**: Errors during the autoinstall process.
4. **networkd-dispatcher warning**: systemd-networkd is not running.

These errors and warnings should be investigated further to determine the root cause and resolve the issues.
```

## Configration

### Use local LLM via ollama as backend

``` bash

 #edit envs.sh
vim envs.sh
export LLM_BACKEND=ollama
export LLM_MODEL_NAME=llama2
export LLM_API_KEY=apikey
export LLM_TEMPERATURE=0.8
export OLLAMA_SERVER_URL=http://localhost:8000

source envs.sh

```

### Use Groq models as backend

``` bash

# edit envs.sh
vim envs.sh
export LLM_BACKEND=groq
export LLM_MODEL_NAME=llama3-8b-8192
export LLM_API_KEY=<your groq api key>
export LLM_TEMPERATURE=0.8

source envs.sh

```


## Feature List

- Act as a terminal assistant.
- Execute commands with natural language.
- Execute some complicated task with script file.
- Use agent mode to execute complex tasks and troubleshooting.

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