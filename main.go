package main

import (
	"math/rand"
	"time"

	_ "github.com/rancher/k3s/pkg/cloudprovider"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider/app"
	"k8s.io/cloud-provider/app/config"
	"k8s.io/cloud-provider/options"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/klog/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ccmOptions, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	fss := cliflag.NamedFlagSets{}
	command := app.NewCloudControllerManagerCommand(ccmOptions, cloudInitializer, app.DefaultInitFuncConstructors, fss, wait.NeverStop)

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		klog.Fatalf("Fatal error running RKE2 Cloud Provider: %v", err)
	}
}

func cloudInitializer(config *config.CompletedConfig) cloudprovider.Interface {
	cloudConfig := config.ComponentConfig.KubeCloudShared.CloudProvider

	// initialize cloud provider with the cloud provider name and config file provided
	cloud, err := cloudprovider.InitCloudProvider(cloudConfig.Name, cloudConfig.CloudConfigFile)
	if err != nil {
		klog.Fatalf("RKE2 Cloud Provider could not be initialized: %v", err)
	}
	if cloud == nil {
		klog.Fatalf("Cloud provider was nil")
	}

	return cloud
}
