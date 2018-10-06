package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/cathalgarvey/fmtless"
	"github.com/fatih/color"
)

type Containers []*Container

type Container struct {
	ID      string    `json:"Id"`
	Created time.Time `json:"Created"`
	Path    string    `json:"Path"`
	Args    []string  `json:"Args"`
	State   struct {
		Status     string    `json:"Status"`
		Running    bool      `json:"Running"`
		Paused     bool      `json:"Paused"`
		Restarting bool      `json:"Restarting"`
		OOMKilled  bool      `json:"OOMKilled"`
		Dead       bool      `json:"Dead"`
		Pid        int       `json:"Pid"`
		ExitCode   int       `json:"ExitCode"`
		Error      string    `json:"Error"`
		StartedAt  time.Time `json:"StartedAt"`
		FinishedAt time.Time `json:"FinishedAt"`
		Health     struct {
			Status        string `json:"Status"`
			FailingStreak int    `json:"FailingStreak"`
			Log           []struct {
				Start    time.Time `json:"Start"`
				End      time.Time `json:"End"`
				ExitCode int       `json:"ExitCode"`
				Output   string    `json:"Output"`
			} `json:"Log"`
		} `json:"Health"`
	} `json:"State"`
	Image           string      `json:"Image"`
	ResolvConfPath  string      `json:"ResolvConfPath"`
	HostnamePath    string      `json:"HostnamePath"`
	HostsPath       string      `json:"HostsPath"`
	LogPath         string      `json:"LogPath"`
	Name            string      `json:"Name"`
	RestartCount    int         `json:"RestartCount"`
	Driver          string      `json:"Driver"`
	Platform        string      `json:"Platform"`
	MountLabel      string      `json:"MountLabel"`
	ProcessLabel    string      `json:"ProcessLabel"`
	AppArmorProfile string      `json:"AppArmorProfile"`
	ExecIDs         interface{} `json:"ExecIDs"`
	HostConfig      struct {
		Binds           []string `json:"Binds"`
		ContainerIDFile string   `json:"ContainerIDFile"`
		LogConfig       struct {
			Type   string `json:"Type"`
			Config struct {
			} `json:"Config"`
		} `json:"LogConfig"`
		NetworkMode  string `json:"NetworkMode"`
		PortBindings struct {
			Nine200TCP []struct {
				HostIP   string `json:"HostIp"`
				HostPort string `json:"HostPort"`
			} `json:"9200/tcp"`
		} `json:"PortBindings"`
		RestartPolicy struct {
			Name              string `json:"Name"`
			MaximumRetryCount int    `json:"MaximumRetryCount"`
		} `json:"RestartPolicy"`
		AutoRemove           bool          `json:"AutoRemove"`
		VolumeDriver         string        `json:"VolumeDriver"`
		VolumesFrom          []interface{} `json:"VolumesFrom"`
		CapAdd               interface{}   `json:"CapAdd"`
		CapDrop              interface{}   `json:"CapDrop"`
		DNS                  interface{}   `json:"Dns"`
		DNSOptions           interface{}   `json:"DnsOptions"`
		DNSSearch            interface{}   `json:"DnsSearch"`
		ExtraHosts           interface{}   `json:"ExtraHosts"`
		GroupAdd             interface{}   `json:"GroupAdd"`
		IpcMode              string        `json:"IpcMode"`
		Cgroup               string        `json:"Cgroup"`
		Links                interface{}   `json:"Links"`
		OomScoreAdj          int           `json:"OomScoreAdj"`
		PidMode              string        `json:"PidMode"`
		Privileged           bool          `json:"Privileged"`
		PublishAllPorts      bool          `json:"PublishAllPorts"`
		ReadonlyRootfs       bool          `json:"ReadonlyRootfs"`
		SecurityOpt          interface{}   `json:"SecurityOpt"`
		UTSMode              string        `json:"UTSMode"`
		UsernsMode           string        `json:"UsernsMode"`
		ShmSize              int           `json:"ShmSize"`
		Runtime              string        `json:"Runtime"`
		ConsoleSize          []int         `json:"ConsoleSize"`
		Isolation            string        `json:"Isolation"`
		CPUShares            int           `json:"CpuShares"`
		Memory               int           `json:"Memory"`
		NanoCpus             int           `json:"NanoCpus"`
		CgroupParent         string        `json:"CgroupParent"`
		BlkioWeight          int           `json:"BlkioWeight"`
		BlkioWeightDevice    interface{}   `json:"BlkioWeightDevice"`
		BlkioDeviceReadBps   interface{}   `json:"BlkioDeviceReadBps"`
		BlkioDeviceWriteBps  interface{}   `json:"BlkioDeviceWriteBps"`
		BlkioDeviceReadIOps  interface{}   `json:"BlkioDeviceReadIOps"`
		BlkioDeviceWriteIOps interface{}   `json:"BlkioDeviceWriteIOps"`
		CPUPeriod            int           `json:"CpuPeriod"`
		CPUQuota             int           `json:"CpuQuota"`
		CPURealtimePeriod    int           `json:"CpuRealtimePeriod"`
		CPURealtimeRuntime   int           `json:"CpuRealtimeRuntime"`
		CpusetCpus           string        `json:"CpusetCpus"`
		CpusetMems           string        `json:"CpusetMems"`
		Devices              interface{}   `json:"Devices"`
		DeviceCgroupRules    interface{}   `json:"DeviceCgroupRules"`
		DiskQuota            int           `json:"DiskQuota"`
		KernelMemory         int           `json:"KernelMemory"`
		MemoryReservation    int           `json:"MemoryReservation"`
		MemorySwap           int           `json:"MemorySwap"`
		MemorySwappiness     interface{}   `json:"MemorySwappiness"`
		OomKillDisable       bool          `json:"OomKillDisable"`
		PidsLimit            int           `json:"PidsLimit"`
		Ulimits              interface{}   `json:"Ulimits"`
		CPUCount             int           `json:"CpuCount"`
		CPUPercent           int           `json:"CpuPercent"`
		IOMaximumIOps        int           `json:"IOMaximumIOps"`
		IOMaximumBandwidth   int           `json:"IOMaximumBandwidth"`
		MaskedPaths          []string      `json:"MaskedPaths"`
		ReadonlyPaths        []string      `json:"ReadonlyPaths"`
	} `json:"HostConfig"`
	GraphDriver struct {
		Data struct {
			LowerDir  string `json:"LowerDir"`
			MergedDir string `json:"MergedDir"`
			UpperDir  string `json:"UpperDir"`
			WorkDir   string `json:"WorkDir"`
		} `json:"Data"`
		Name string `json:"Name"`
	} `json:"GraphDriver"`
	Mounts []struct {
		Type        string `json:"Type"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Mode        string `json:"Mode"`
		RW          bool   `json:"RW"`
		Propagation string `json:"Propagation"`
	} `json:"Mounts"`
	Config struct {
		Hostname     string `json:"Hostname"`
		Domainname   string `json:"Domainname"`
		User         string `json:"User"`
		AttachStdin  bool   `json:"AttachStdin"`
		AttachStdout bool   `json:"AttachStdout"`
		AttachStderr bool   `json:"AttachStderr"`
		ExposedPorts struct {
			Nine200TCP struct {
			} `json:"9200/tcp"`
			Nine300TCP struct {
			} `json:"9300/tcp"`
		} `json:"ExposedPorts"`
		Tty         bool     `json:"Tty"`
		OpenStdin   bool     `json:"OpenStdin"`
		StdinOnce   bool     `json:"StdinOnce"`
		Env         []string `json:"Env"`
		Cmd         []string `json:"Cmd"`
		Healthcheck struct {
			Test     []string `json:"Test"`
			Interval int64    `json:"Interval"`
			Timeout  int64    `json:"Timeout"`
			Retries  int      `json:"Retries"`
		} `json:"Healthcheck"`
		ArgsEscaped bool   `json:"ArgsEscaped"`
		Image       string `json:"Image"`
		Volumes     struct {
			UsrShareElasticsearchConfigElasticsearchYml struct {
			} `json:"/usr/share/elasticsearch/config/elasticsearch.yml"`
			UsrShareElasticsearchConfigLog4J2Properties struct {
			} `json:"/usr/share/elasticsearch/config/log4j2.properties"`
		} `json:"Volumes"`
		WorkingDir string      `json:"WorkingDir"`
		Entrypoint []string    `json:"Entrypoint"`
		OnBuild    interface{} `json:"OnBuild"`
		Labels     struct {
			BuildDate                       string `json:"build-date"`
			ComDockerComposeConfigHash      string `json:"com.docker.compose.config-hash"`
			ComDockerComposeContainerNumber string `json:"com.docker.compose.container-number"`
			ComDockerComposeOneoff          string `json:"com.docker.compose.oneoff"`
			ComDockerComposeProject         string `json:"com.docker.compose.project"`
			ComDockerComposeService         string `json:"com.docker.compose.service"`
			ComDockerComposeVersion         string `json:"com.docker.compose.version"`
			License                         string `json:"license"`
			Maintainer                      string `json:"maintainer"`
			Name                            string `json:"name"`
			Vendor                          string `json:"vendor"`
		} `json:"Labels"`
	} `json:"Config"`
	NetworkSettings struct {
		Bridge                 string `json:"Bridge"`
		SandboxID              string `json:"SandboxID"`
		HairpinMode            bool   `json:"HairpinMode"`
		LinkLocalIPv6Address   string `json:"LinkLocalIPv6Address"`
		LinkLocalIPv6PrefixLen int    `json:"LinkLocalIPv6PrefixLen"`
		Ports                  struct {
			Nine200TCP []struct {
				HostIP   string `json:"HostIp"`
				HostPort string `json:"HostPort"`
			} `json:"9200/tcp"`
			Nine300TCP interface{} `json:"9300/tcp"`
		} `json:"Ports"`
		SandboxKey             string      `json:"SandboxKey"`
		SecondaryIPAddresses   interface{} `json:"SecondaryIPAddresses"`
		SecondaryIPv6Addresses interface{} `json:"SecondaryIPv6Addresses"`
		EndpointID             string      `json:"EndpointID"`
		Gateway                string      `json:"Gateway"`
		GlobalIPv6Address      string      `json:"GlobalIPv6Address"`
		GlobalIPv6PrefixLen    int         `json:"GlobalIPv6PrefixLen"`
		IPAddress              string      `json:"IPAddress"`
		IPPrefixLen            int         `json:"IPPrefixLen"`
		IPv6Gateway            string      `json:"IPv6Gateway"`
		MacAddress             string      `json:"MacAddress"`
		Networks               struct {
			BeamDefault struct {
				IPAMConfig          interface{} `json:"IPAMConfig"`
				Links               interface{} `json:"Links"`
				Aliases             []string    `json:"Aliases"`
				NetworkID           string      `json:"NetworkID"`
				EndpointID          string      `json:"EndpointID"`
				Gateway             string      `json:"Gateway"`
				IPAddress           string      `json:"IPAddress"`
				IPPrefixLen         int         `json:"IPPrefixLen"`
				IPv6Gateway         string      `json:"IPv6Gateway"`
				GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
				MacAddress          string      `json:"MacAddress"`
				DriverOpts          interface{} `json:"DriverOpts"`
			} `json:"beam_default"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

func ExecuteHealthCheck(object string) error {
	ps := exec.Command("docker")
	ps.Args = append(ps.Args, "ps", "-a", "-q", "--filter", fmt.Sprintf("name=%s", object))
	b, err := ps.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error getting container id: %s", err)
	}

	containerId := strings.TrimSpace(string(b))

	inspect := exec.Command("docker")
	inspect.Args = append(inspect.Args, "inspect", containerId)
	b, err = inspect.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error inspecting container: %s (%s)", err, string(b))
	}

	c := Containers{}
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	fmt.Printf("Waiting for %s to pass healthcheck ... ", strings.TrimPrefix(c[0].Name, "/"))
	healthy := c[0].State.Health.Status == "healthy"

	for !healthy {
		c = nil

		inspect := exec.Command("docker")
		inspect.Args = append(inspect.Args, "inspect", containerId)
		b, err := inspect.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error inspecting container: %s (%s)", err, string(b))
		}

		c = Containers{}
		if err := json.Unmarshal(b, &c); err != nil {
			return err
		}

		if c[0].State.Health.Status == "healthy" {
			healthy = true
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	color.Green("done")
	return nil
}

func ExecuteDockerCommand(home string, envVars []*string, composeFiles []string, command []string, objects []string) error {
	cmd := exec.Command("docker-compose")
	cmd.Dir = home

	for _, envVar := range envVars {
		cmd.Env = append(cmd.Env, *envVar)
	}

	for _, composeFile := range composeFiles {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, composeFile)
	}

	var cleanedObjects []string
	for _, object := range objects {
		if object == "all" {
			objects = []string{}
			break
		}

		cleanedObjects = append(cleanedObjects, strings.Split(object, " ")...)
	}

	cmd.Args = append(cmd.Args, command...)
	cmd.Args = append(cmd.Args, cleanedObjects...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for {
			select {
			case sig := <-signalChan:
				if err := cmd.Process.Signal(sig); err != nil {
					panic(err)
				}
				return
			case <-doneChan:
				close(doneChan)
				close(signalChan)
				signal.Reset(os.Interrupt)
				return
			}
		}
	}()

	cmd.Wait()
	doneChan <- struct{}{}

	return nil
}
