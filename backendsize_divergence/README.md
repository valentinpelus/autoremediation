# auto-remediation

This repository host the source code of our auto-remediation solution : Remediate

It is composed of multiple packages making it able to correct recurring alerts on our onprem and cloud/kube clusters


## Table of Contents
- [What's included](https://github.m6web.fr/sysadmin/auto-remediation#whats-included)
- [Versioning](https://github.m6web.fr/sysadmin/auto-remediation#whats-included)
- [How to use it](https://github.m6web.fr/sysadmin/auto-remediation#how-to-use-it)
- [Contributing code](https://github.m6web.fr/sysadmin/auto-remediation#contributing-code)
- [How to develop features](https://github.m6web.fr/sysadmin/auto-remediation#how-to-develop-features/fix)

### What's inclued



### Versioning

- For each Remediate release, the major version (first digit) will be incremented by 1 
- New functionnality will result in the minor version (second digit) changing. PRs that are cherry-picked will result in an update to the master branch with a corresponding new tag changing the patch version.
- Bugfixes will result in the patch version (third digit) changing.

#### Branches and tags

We will create a new branch and tag for each increment in the minor version number. We will create only a new tag for each increment in the patch version number.

### How to use it

This application needs to be build with golang into a binary and deploy on staging environment followed by production environment if validated.

### Contributing code

Please send pull requests against this repository prefixing it with :
- feat/ if it's a new feature to the code base
- fix/ if it's a fix of an actual feature or bug

### How to develop features/fix

To have a Go environment
  ```sh
  sudo dnf -y update
  sudo dnf -y install go
  ```

To check if we golang has been sucessfully installed :
  ```sh
  go version
  # Should send a similar output :
  go version go1.18.3 linux/amd64
  ```

You can check also with the go env command :
  ```sh
  go env
  # Should deliver a similar output :
  GO111MODULE=""
  GOARCH="amd64"
  GOBIN=""
  GOCACHE="/home/fedora/.cache/go-build"
  GOENV="/home/fedora/.config/go/env"
  GOEXE=""
  GOEXPERIMENT=""
  GOFLAGS=""
  GOHOSTARCH="amd64"
  GOHOSTOS="linux"
  GOINSECURE=""
  GOMODCACHE="/home/fedora/go/pkg/mod"
  GONOPROXY=""
  GONOSUMDB=""
  GOOS="linux"
  GOPATH="/home/fedora/go"
  GOPRIVATE=""
  GOPROXY="direct"
  GOROOT="/usr/lib/golang"
  GOSUMDB="off"
  GOTMPDIR=""
  GOTOOLDIR="/usr/lib/golang/pkg/tool/linux_amd64"
  GOVCS=""
  GOVERSION="go1.18.3"
  GCCGO="gccgo"
  GOAMD64="v1"
  AR="ar"
  CC="gcc"
  CXX="g++"
  CGO_ENABLED="1"
  GOMOD="/dev/null"
  GOWORK=""
  CGO_CFLAGS="-g -O2"
  CGO_CPPFLAGS=""
  CGO_CXXFLAGS="-g -O2"
  CGO_FFLAGS="-g -O2"
  CGO_LDFLAGS="-g -O2"
  PKG_CONFIG="pkg-config"
  GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build3590682015=/tmp/go-build -gno-record-gcc-switches"
  ```

  We need to make some ajustment before begining by setting m6web go proxy to query package, for that we need to edit go env file :
  ```sh
  mkdir -p $HOME/.config/go
  vim $HOME/.config/go/env
  # Paste those two values 
  GONOSUMDB=github.m6web.fr
  GOPROXY=https://goproxy.services.bedrock.tech,https://proxy.golang.org,direct
  ```


  ``` sh
  git clone git@github.m6web.fr:sysadmin/auto-remediation.git
  cd auto-remediation
  ```

Keep in mind that we need to standardize at maximum our code base, in most use case we want to avoid at maximum specific func to a problem.

For exemple, we need to delete a pod in backend size divergence state in the alerting, we will not make a function named 
``` 
func deleteBackendSizeDivergencePod {}
```
But instead focusing more on the basic usage of this func 
```
func deletePod {}
```

This will allow us to delete the pod we want by passing the right arguments abd reuse this function for other problem solving at ease without needing any additional development
