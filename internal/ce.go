package internal

import (
	"errors"
	"fmt"
)

type ceFactory struct {
	args map[string]interface{}
}

func NewCeFactory(args map[string]interface{}) *ceFactory {
	return &ceFactory{
		args: args,
	}
}

type containerEngine interface {
	// Append run arg - environment variables
	AppendEnvVar(key string, value string)
	// Append build arg - container build arg
	AppendBuildArg(name string, value string)
	// Append run arg - host/container volume map
	AppendVolMap(hostVol string, containerVol string, mapAttrs string)
	// Append run arg - host/container port/address
	AppendPortMap(hostPort string, containerPort string, hostAddr string)
	// Constructs and returns a build image command
	GetBuildCmd() []string
	// Constructs and returns a run container command
	GetRunCmd(containerName string, entryPoint string, image string, entryPointArgs ...string) []string
	// Returns ce executable name (e.g. podman)
	GetExecName() string
	// Todo: Add func here as necessary
}

func (cf *ceFactory) Create() (containerEngine, error) {
	if cf.args["ceName"] == "podman" {
		return NewPodman(), nil
	} else {
		return nil, errors.New("container engine name is not supported")
	}
}

type podman struct {
	envVars      [][]string
	volMaps      [][]string
	volMapsAttr  map[string]string
	portMaps     [][]string
	portMapAddrs map[string]string
	buildArgs    [][]string
}

func NewPodman() *podman {
	return &podman{}
}

func (p *podman) AppendEnvVar(key string, value string) {
	env := []string{key, value}
	p.envVars = append(p.envVars, env)
}

func (p *podman) ToEnvVarArgs() []string {
	args := []string{}
	for _, val := range p.envVars {
		k := val[0]
		v := val[1]
		envVarDef := fmt.Sprintf("%s=%s", k, v)
		args = append(args, "-e", envVarDef)
	}
	return args
}

func (p *podman) AppendVolMap(hostVol string, containerVol string, mapAttrs string) {
	vol := []string{hostVol, containerVol}
	p.volMaps = append(p.volMaps, vol)

	if p.volMapsAttr == nil {
		p.volMapsAttr = map[string]string{hostVol: mapAttrs}
	} else {
		p.volMapsAttr[hostVol] = mapAttrs
	}

}

func (p *podman) ToVolMapArgs() []string {
	args := []string{}
	for _, val := range p.volMaps {
		hostVol := val[0]
		contVol := val[1]
		mapAttrs := p.volMapsAttr[hostVol]
		volMap := fmt.Sprintf("%s:%s:%s", hostVol, contVol, mapAttrs)
		args = append(args, "-v", volMap)
	}
	return args
}

func (p *podman) AppendPortMap(hostPort string, containerPort string, hostAddr string) {
	port := []string{hostPort, containerPort}
	p.portMaps = append(p.portMaps, port)

	if p.portMapAddrs == nil {
		p.portMapAddrs = map[string]string{hostPort: hostAddr}
	} else {
		p.portMapAddrs[hostPort] = hostAddr
	}
}

func (p *podman) ToPortMapArgs() []string {
	args := []string{}
	for _, val := range p.portMaps {
		hostPort := val[0]
		containerPort := val[1]
		hostAddr := p.portMapAddrs[hostPort]
		portMap := fmt.Sprintf("%s:%s:%s", hostAddr, hostPort, containerPort)
		args = append(args, "-p", portMap)
	}
	return args
}

func (p *podman) AppendBuildArg(name string, value string) {
	buildArg := []string{name, value}
	p.buildArgs = append(p.buildArgs, buildArg)
}

func (p *podman) ToBuildArgs() []string {
	args := []string{}
	for _, val := range p.buildArgs {
		name := val[0]
		value := val[1]
		buildArg := fmt.Sprintf("%s=%s", name, value)
		args = append(args, "--build-arg", buildArg)
	}
	return args
}

func (p *podman) GetRunCmd(containerName string, entryPoint string, image string, entryPointArgs ...string) []string {

	runCmd := []string{
		"run",
		"--name",
		containerName,
		"-it",
		"--privileged",
	}
	runCmd = append(runCmd, p.ToEnvVarArgs()...)
	runCmd = append(runCmd, p.ToPortMapArgs()...)
	runCmd = append(runCmd, p.ToVolMapArgs()...)
	ep := []string{
		"--entrypoint",
		entryPoint,
		image,
	}
	ep = append(ep, entryPointArgs...)
	runCmd = append(runCmd, ep...)
	return runCmd
}

func (p *podman) GetExecName() string {
	return "podman"
}

func (p *podman) GetEnvVars() [][]string {
	return p.envVars
}

func (p *podman) GetVolMaps() [][]string {
	return p.volMaps
}

func (p *podman) GetPortMaps() [][]string {
	return p.portMaps
}

func (p *podman) GetBuildCmd() []string {
	buildCmd := []string{
		"build",
		"-t",
		"hc:latest",
	}
	buildCmd = append(buildCmd, p.ToBuildArgs()...)
	buildCmd = append(buildCmd, ".")
	return buildCmd
}
