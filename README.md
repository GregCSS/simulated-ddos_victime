# 🛜🕸️ DDoS Simulation
This project simulates a DDoS (Distributed Denial of Service) attack on an intentionally vulnerable victim, in this case a basic e-commerce web application setup using *Node.js* and *SQLite*.

## 🪲 Vulnerability being demo-ed
This simulation primarily addresses vulnerabilities leading to **Resource Exhaustion** such as:
* [CWE-770: Allocation of Resources Without Limits or Throttling](https://cwe.mitre.org/data/definitions/770.html)
* [OWASP API4:2023 Unrestricted Resource Consumption](https://api-security.owasp.org/editions/2023/en/0xa4-unrestricted-resource-consumption/)

Computing resources is utilized for any sort of work by software and hardware, and thus there are ways to exploit this maliciously. For example, unregulated computationally heavy requests to a service spikes up memory and CPU usage and eventually enough of these requests lead to exhausting resource limits, leading to crashes and failures.

Documented cases include:
* [CVE-2021-22029: VMware Workspace ONE UEM REST API](https://app.opencve.io/cve/CVE-2021-22029)
* [CVE-2025-48375: Schule Missing Rate Limiting on OTP Email Requests](https://app.opencve.io/cve/CVE-2025-48375)

## ⚙️ How this simulation works
This section details how the simulation has been designed.

### 1. Virtual Docker Container (Bot)
This project uses **Docker**, acting as the test environment, to spin up multiple containers with scripts for:
* Establishing a Websocket connection to the Bot Master (C2 Layer)
* Bash scripts to attack, using relevant tools such as `hping3`, a given target

### 2. HTTP Server managing Websocket Connections (Master)
When a container (bot) starts, it attempts to establish a Websocket connection with the C2 layer and keeps trying until one is established. The C2 Layer has an HTTP server that manages Websocket connections, and sends/broadcasts messages to
bots.

The bot then receives target IPs from the C2 layer and executes scripts accordingly to attack the destination.

### 3. Vulnerable E-Commerce Web Application
The web application is implemented using *Node.js* and *SQLite*, running on a container with low specs so as to crash faster for the sake of demonstration. The web application provides a **search feature with filtering options, which the botnet will exploit by sending volumous amounts of complex queries** which requires the database to work significantly more until resources are exhausted.

## 🔧 Possible Countermeasures 

## 📦 Prerequisites
Make sure you have the following dependencies installed.

* [Golang 1.27.0](https://go.dev/doc/install)
* [Docker 29.7.2](https://www.docker.com/get-started/)

## 💻 Setup & Usage
Follow these steps to get your development environment setup.

1. **Clone the repository**
    ```bash
    git clone https://github.com/dotping-me/simulated-ddos.git && cd simulated-ddos/
    ```

2. **Attacker Setup**
    
    Make sure to run the following scripts from within `/attacker`:
    ```bash
    cd attacker/ # From project root
    ```
    * **2.1. Build binairies**
        ```bash
        go buid -o bin/c2 ./cmd/c2/main.go
        ```

3. **Bot and Botnet Setup**
    
    Make sure to run the following scripts from within `/attacker`:
    ```bash
    cd attacker/ # From project root
    ```

    * **3.1. Start Docker daemon**
        ```bash
        sudo systemctl start docker
        ```

    * **3.2. Build Docker image**
        ```bash
        docker build --no-cache -f ./bot/Dockerfile -t bot_image .
        ```

    * **3.3. Ensure firewall does not get in the way**
        ```bash
        ../scripts/add_ufw_rule.sh # scripts/ found from Project root
        ```

        ***Note:*** 
        * *Manually check (using `ip -4 addr`) and modify variables in `add_ufw_rule.sh` and `delete_ufw_rule.sh`*
        * *Idea behind spinning up a bot is:*
            ```bash
            sudo docker run -d \
                --add-host=host.docker.internal:host-gateway \
                bot_image \
                --master ws://host.docker.internal:8080/connect
            ```
        * To monitor container resource usage, use:
            ```bash
            sudo docker stats
            ```

3. **Victim Setup**

    Make sure to run the following scripts from within `/victim`:
    ```bash
    cd victim/ # From project root
    ```

4. **Start C2 HTTP server**
    ```bash
    ./attacker/bin/c2 # From project root
    ```

5. **Spin up botnet**
    ```bash
    ./scripts/create_botnet.sh 10 # From project root and takes N bots as arg
    ```