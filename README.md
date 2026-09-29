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
This project uses **Docker** (Docker Network) to spin up multiple containers with scripts for:
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
