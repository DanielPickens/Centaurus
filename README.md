# Centaurus

[Install](https://github.com/centaurus/centaurus?tab=readme-ov-file#battery-install)
| [Guide](https://github.com/centaurus/centaurus?tab=readme-ov-file#books-guide)
| [Releases](https://github.com/danielpickens/centaurusreleases)
| [Source Code](https://github.com/centaurus/centaurus)

A single binary to manage your multiple kubernetes clusters.

**centaurus** provides a simple and rich real time interface to manage and investigate your clusters.


**Key features of centaurus include:**

* **Single binary deployment:** centaurus can be easily deployed as a single binary, eliminating the need for complex configurations.
* **Browser-based access:** centaurus can be accessed directly from your favorite web browser, providing a seamless user experience.
* **Real-time cluster monitoring:** centaurus offers a rich, real-time interface that displays the current state of your Kubernetes clusters, allowing you to quickly identify and address issues.
* **Cluster management:** centaurus enables you to manage multiple Kubernetes clusters from a single pane of glass, reducing the overhead of switching between different tools and interfaces.
* **Detailed cluster insights:** centaurus provides comprehensive insights into your Kubernetes clusters, manifest info of your pods, services, config and others.


> [!Important]
> Please keep in mind that centaurus is still under active development.

# :battery: Install

#### Docker

```shell
docker run -p 7080:7080 -v centaurus:/.centaurus ghcr.io/danielpickens/centaurus:latest
```

> To access local kind cluster you can use "--network host" docker flag.

#### Helm

```bash
helm install centaurus oci://ghcr.io/centaurus/charts/centaurus -n centaurus-system --create-namespace
```

> With helm centaurus runs on port `8443` with self-signed certificates. [see charts](https://github.com/danielpickens/centaurustree/main/charts/centaurus)

#### Homebrew

```shell
brew install centaurus/tap/centaurus
```

#### Snap

```shell
sudo snap install centaurus
```

#### Arch Linux

```shell
yay -S centaurus-bin
```

#### Winget 

```shell
winget install --id=centaurus.centaurus -e
```

#### Scoop

```shell
scoop bucket add centaurus https://github.com/danielpickens/centaurus/scoop-bucket.git
scoop install centaurus
```

#### Binary

**MacOS**
[Binary](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Darwin_all.tar.gz) ( Multi-Architecture )

**Linux (Binaries)**
[amd64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Linux_x86_64.tar.gz) | [arm64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Linux_arm64.tar.gz) | [i386](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Linux_i386.tar.gz)

**Windows (Exe)**
[amd64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Windows_x86_64.zip) | [arm64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Windows_arm64.zip) | [i386](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Windows_i386.zip)

**FreeBSD (Binaries)**
[amd64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Freebsd_x86_64.tar.gz) | [arm64](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Freebsd_arm64.tar.gz) | [i386](https://github.com/danielpickens/centaurusreleases/latest/download/centaurus_Freebsd_i386.tar.gz)

Manually
Download the pre-compiled binaries from the [Release!](https://github.com/danielpickens/centaurusreleases) page and copy them to the desired location or system path.

> [!TIP] 
> After installation, you can access **centaurus** at `http://localhost:7080`
>
>  If you're running it in a Kubernetes cluster or on an on-premises server, we recommend using **HTTPS**.
>  When not used over HTTP/2 SSE suffers from a limitation to the maximum number of open connections. [Mozzila](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)⤴
>
>  You can start **centaurus** with **HTTPS** using the following command:
>
> ```
> $ centaurus --certFile=/path/to/cert.pem --keyFile=/path/to/key.pem
> ```

# :books: Guide

### Flags

Since centaurus runs as binary there are few of flag you can use.

```shell
> centaurus --help

Usage:
  centaurus [flags]
  centaurus [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version of centaurus

Flags:
      --certFile string        absolute path to certificate file
  -h, --help                   help for centaurus
      --k8s-client-burst int   Maximum burst for throttle (default 50)
      --k8s-client-qps int     maximum QPS to the master from client (default 50)
      --keyFile string         absolute path to key file
      --no-open-browser        Do not open the default browser
  -p, --port string            port to listen on (default ":7080")

Use "centaurus [command] --help" for more information about a command.
```

### Setting up HTTPS locally

You can use your own certificates or create new local trusted certificates using [mkcert](https://github.com/FiloSottile/mkcert)⤴.

> [!Important]
> You'll need to install [mkcert](https://github.com/FiloSottile/mkcert)⤴ separately.

1. Install mkcert on your computer.
2. Run the following command in your terminal or command prompt:

`mkcert centaurus.test localhost 127.0.0.1 ::1`

3. This command will generate two files: a certificate file and a key file (the key file will have `-key.pem` at the end of its name).
4. To use these files with **centaurus**, use `--certFile=` and `--keyFile=` flags.

```shell
centaurus --certFile=centaurus.test+3.pem --keyFile=centaurus.test+3-key.pem
```

**When using Docker**

When using Docker, you can attach volumes and provide certificates by using specific flags. 

In the following example, we mount the current directory from your host to the `/.certs` directory inside the Docker container:

```shell
docker run -p 7080:7080 \
    -v centaurus:/.centaurus \
    -v $(pwd):/.certs \
    ghcr.io/centaurus/centaurus:latest \
    --certFile=/.certs/centaurus.test+3.pem \
    --keyFile=/.certs/centaurus.test+3-key.pem
```

### Custom Port

You can run **centaurus** on any port you like. by using `-p` or `--port` flag

```shell
centaurus -p 7080
```

# Contribution

This project welcomes your PR and issues. For example, refactoring, adding features, etc.

If you need any help, you can contact me from the above Developers sections.



# License

Centaurus is licensed under [Apache License, Version 2.0](./LICENSE)
