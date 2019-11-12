module github.com/tietang/props

go 1.12

//被墙的原因，替换golang.org源为github.com源
replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.37.2
	golang.org/x/build => github.com/golang/build v0.0.0-20190327004547-5a2224f3eb52
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190321205749-f0864edee7f3
	golang.org/x/image => github.com/golang/image v0.0.0-20190321063152-3fc05d484e9f
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190319155245-9487ef54b94a
	golang.org/x/net => github.com/golang/net v0.0.0-20190327025741-74e053c68e29
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190319182350-c85d3e98c914
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190312170614-0655857e383f
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190327011446-79af862e6737
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.5.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190321212433-e79c0c59cdb5
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.1
)

require (
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-ini/ini v1.51.0
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20180628210949-0892b62f0d9f // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/hashicorp/consul/api v1.2.0
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/pascaldekloe/goe v0.1.0 // indirect
	github.com/prometheus/common v0.7.0
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/assertions v0.0.0-20180301161246-7678a5452ebe // indirect
	github.com/smartystreets/goconvey v0.0.0-20170602164621-9e8dc3f972df
	github.com/smartystreets/gunit v0.0.0-20180314194857-6f0d6275bdcd // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tietang/go-utils v0.1.3
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/valyala/fasttemplate v1.1.0
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	golang.org/x/sys v0.0.0-20191110163157-d32e6e3b99c4 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/genproto v0.0.0-20190605220351-eb0b1bdb6ae6 // indirect
	google.golang.org/grpc v1.21.1 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

//exclude (
//	google.golang.org/appengine v1.5.0
//	google.golang.org/api v0.0.0
//)
