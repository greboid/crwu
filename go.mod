module github.com/greboid/rwtccus

go 1.20

require (
	github.com/csmith/envflag v1.0.0
	github.com/docker/compose/v2 v2.16.0
	github.com/docker/docker v23.0.1+incompatible
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/render v1.0.2
	github.com/rs/zerolog v1.29.0
	github.com/samber/lo v1.37.0
	golang.org/x/sync v0.1.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/compose-spec/compose-go v1.11.0 // indirect
	github.com/distribution/distribution/v3 v3.0.0-20230217222803-a4d64adabcb9 // indirect
	github.com/docker/cli v23.0.1+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230213192124-5e25df0256eb // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)

replace (
	github.com/cucumber/godog => github.com/laurazard/godog v0.0.0-20220922095256-4c4b17abdae7
	github.com/docker/cli => github.com/docker/cli v20.10.3-0.20221013132413-1d6c6e2367e2+incompatible // 22.06 master branch
	github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221021173910-5aac513617f0+incompatible // 22.06 branch
	github.com/moby/buildkit => github.com/moby/buildkit v0.10.1-0.20220816171719-55ba9d14360a // same as buildx
	k8s.io/api => k8s.io/api v0.22.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.4
	k8s.io/client-go => k8s.io/client-go v0.22.4
)
