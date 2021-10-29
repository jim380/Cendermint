module github.com/jim380/Cendermint

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.41.3
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99 // indirect
	github.com/joho/godotenv v1.4.0
	github.com/prometheus/client_golang v1.8.0
	github.com/tendermint/tendermint v0.34.13 // indirect
	go.uber.org/zap v1.15.0
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
