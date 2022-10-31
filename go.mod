module github.com/eleven-sh/hetzner-cloud-provider

go 1.18

replace github.com/eleven-sh/eleven v0.0.0 => ../eleven

replace github.com/eleven-sh/agent v0.0.0 => ../agent

require (
	github.com/eleven-sh/agent v0.0.0
	github.com/eleven-sh/eleven v0.0.0
	github.com/golang/mock v1.6.0
	github.com/hetznercloud/hcloud-go v1.35.2
	github.com/mikesmitty/edkey v0.0.0-20170222072505-3356ea4e686a
	github.com/pelletier/go-toml v1.9.5
	golang.org/x/crypto v0.0.0-20220313003712-b769efc7c000
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gosimple/slug v1.12.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.8-0.20211004125949-5bd84dd9b33b // indirect
	golang.org/x/tools v0.1.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
