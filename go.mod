module gocat

require (
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/linuxkit/virtsock v0.0.0-20220523201153-1a23e78aa7a2
	github.com/mdlayher/vsock v1.1.1
	github.com/palantir/stacktrace v0.0.0-20161112013806-78658fd2d177
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mdlayher/socket v0.2.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220204135822-1c1b9b1eba6a // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f
	golang.org/x/text => golang.org/x/text v0.3.3
)

go 1.17
