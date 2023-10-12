module zero-provers/server

go 1.21

require (
	github.com/0xPolygon/polygon-edge v1.1.0
	github.com/rs/zerolog v1.29.1
	github.com/spf13/cobra v1.7.0
	google.golang.org/grpc v1.58.1
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/umbracle/ethgo v0.1.4-0.20230810113823-c9c19bcd8a1e // indirect
	github.com/valyala/fastjson v1.6.3 // indirect
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/umbracle/fastrlp v0.1.1-0.20230504065717-58a1b8a9929d // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
)

// Use polygon-edge@feat/zero last commit.
// https://github.com/0xPolygon/polygon-edge/tree/feat/zero
replace github.com/0xPolygon/polygon-edge => github.com/0xPolygon/polygon-edge v1.1.1-0.20230929152933-907104765c64
