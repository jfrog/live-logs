module github.com/jfrog/live-logs

go 1.15

require (
	github.com/c-bata/go-prompt v0.2.5 // indirect
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/jfrog/jfrog-cli-core v1.3.1
	github.com/jfrog/jfrog-client-go v0.19.1
	github.com/manifoldco/promptui v0.8.0
	github.com/mholt/archiver v2.1.0+incompatible // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/stretchr/testify v1.6.1
)

replace github.com/jfrog/jfrog-cli-core => github.com/jfrog/jfrog-cli-core v1.3.2-0.20210307144918-774813db58f7

replace github.com/jfrog/jfrog-client-go => github.com/jfrog/jfrog-client-go v0.19.2-0.20210307144103-d39e869a25e6
