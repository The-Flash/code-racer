package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/The-Flash/code-racer/internal/api"
	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/internal/execution"
	"github.com/The-Flash/code-racer/internal/file_system"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/The-Flash/code-racer/internal/scheduler"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"github.com/sarulabs/di/v2"
)

var (
	manifestPtr    = flag.String("f", "", "Path to manifest file")
	mountPointPtr  = flag.String("m", "", "Path to mount point")
	runnersPathPtr = flag.String("r", "", "Path to runners directory")
)

func main() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	diBuilder, err := di.NewBuilder()
	if err != nil {
		log.Fatal("Could not load IoC container")
	}

	// di for docker client connection
	diBuilder.Add(di.Def{
		Name: names.DiDockerClientProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			cli, err := client.NewClientWithOpts(client.FromEnv)
			if err != nil {
				log.Fatal("could not connect to docker ", err)
			}

			return cli, err
		},
		Close: func(obj interface{}) error {
			cli := obj.(*client.Client)
			cli.Close()
			return nil
		},
	})

	// di for rest api
	diBuilder.Add(di.Def{
		Name: names.DiAPIProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return api.NewAPI(ctn)
		},
	})

	// di for runtime manager
	diBuilder.Add(di.Def{
		Name: names.DiRuntimeManagerProvider,
		Build: func(ctn di.Container) (v interface{}, err error) {
			v = new(runtime_manager.RuntimeManager)
			v.(*runtime_manager.RuntimeManager).NewRuntimeManager(ctn)
			return
		},
	})

	// di for manifest
	diBuilder.Add(di.Def{
		Name: names.DiManifestProvider,
		Build: func(ctn di.Container) (m interface{}, err error) {
			m = new(manifest.Manifest)
			obj := m.(*manifest.Manifest)
			err = obj.Load(*manifestPtr)
			return
		},
	})

	// di for config
	diBuilder.Add(di.Def{
		Name: names.DiConfigProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			c := config.NewConfig(
				*manifestPtr,
			)
			c.ManifestPath = *manifestPtr
			c.FsMount.MountSourcePath = *mountPointPtr
			c.FsMount.MountTargetPath = "/code-racer"

			c.RunnersMount.MountSourcePath = *runnersPathPtr
			c.RunnersMount.MountTargetPath = "/runners"
			return c, nil
		},
	})

	// di for files
	diBuilder.Add(di.Def{
		Name: names.DiFileProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			config := ctn.Get(names.DiConfigProvider).(*config.Config)
			fp := file_system.NewFileProvider(config)
			return fp, nil
		},
	})

	// di for executor
	diBuilder.Add(di.Def{
		Name: names.DiExecutorProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			e := new(execution.Executor)
			e.Setup(ctn)
			return e, nil
		},
	})

	diBuilder.Add(di.Def{
		Name: names.DiSchedulerProvider,
		Build: func(ctn di.Container) (s interface{}, err error) {
			m := ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
			rm := ctn.Get(names.DiRuntimeManagerProvider).(*runtime_manager.RuntimeManager)
			s = scheduler.Load(m.Runtimes, rm)
			return
		},
	})

	ctn := diBuilder.Build()

	api := ctn.Get(names.DiAPIProvider).(*api.API)

	rtm := ctn.Get(names.DiRuntimeManagerProvider).(*runtime_manager.RuntimeManager)

	go rtm.Run()

	go api.ListenAndServeBlocking()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigs

	ctn.DeleteWithSubContainers()
}
