# NetMon

NetMon is a network monitoring tool(hence the name), built using [Pro-bing](https://github.com/prometheus-community/pro-bing)
that can be used to detect network fluctuations and alert the relevant parties.
When launched, it constantly runs an `ICMP ECHO_REQUEST` command(ping) and
analyses the results of the command.
Alerts the responsible people when the latencies exceed a pre-defined amount.
It also has a web interface to visualize the results, this is shared to the 
people configured in config.yml

> **NOTE**
> 
> To use it you need to **clone** the repo not install it using `go install`

## Usage
### Requirements
<details>
<summary>
Environment variables 

(See [example env](.env.example))
</summary

- email n password (Get this from google)
- ngrok token
</details>

<details>
<summary>Config options

(See [example config](config.example.yml))
</summary>

- **max_latency** -> expected latencies This can be the maximum latencies the ISP states, or the average in your network
- **Timeout** -> How long to wait for the `ping` results before it cancels the `ping`
- emails of the people to alert (1 or more)
</details>

<details>
<summary>Optional</summary>

- **Logo** -> place this in `./web/static` directory
</details>



### Running it

Navigate to the root of the project

```bash
# Linux 
make start
## or
go build -o ./bin .
./bin/netmon

# Windows
make.bat 
## or
go build . 
netmon.exe
```

#### Tests

```bash
make test 
# or
go test -v ./...
```

### Known Issues

1. Graph dips to 0 a lot 
Try increasing the **max_latency** value
This happened more on windows
