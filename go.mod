module github.com/albatross-org/sergeant

go 1.15

require (
	github.com/albatross-org/go-albatross v0.1.0
	github.com/dghubble/trie v0.0.0-20201011220304-ed6d6b8add55
	github.com/fatih/color v1.10.0
	github.com/gin-gonic/contrib v0.0.0-20201101042839-6a891bf89f19
	github.com/gin-gonic/gin v1.7.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mroth/weightedrand v0.4.1
	github.com/robotn/gohook v0.30.5
	github.com/segmentio/fasthash v1.0.3
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/albatross-org/go-albatross => ../go-albatross
