package fs

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type TurboConfigJSON struct {
	Base               string   `json:"baseBranch,omitempty"`
	GlobalDependencies []string `json:"globalDependencies,omitempty"`
	TurboCacheOptions  string   `json:"cacheOptions,omitempty"`
	Outputs            []string `json:"outputs,omitempty"`
	RemoteCacheUrl     string   `json:"remoteCacheUrl,omitempty"`
	Pipeline           map[string]Pipeline
}

type PPipeline struct {
	Outputs   *[]string `json:"outputs"`
	Cache     *bool     `json:"cache,omitempty"`
	DependsOn []string  `json:"dependsOn,omitempty"`
}

type Pipeline struct {
	Outputs   []string `json:"-"`
	Cache     *bool    `json:"cache,omitempty"`
	DependsOn []string `json:"dependsOn,omitempty"`
	PPipeline
}

func (c *Pipeline) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &c.PPipeline); err != nil {
		return err
	}
	// We actually need a nil value to be able to unmarshal the json
	// because we interpret the omission of outputs to be different
	// from an empty array. We can't use omitempty because it will
	// always unmarshal into an empty array which is not what we want.
	if c.PPipeline.Outputs != nil {
		c.Outputs = *c.PPipeline.Outputs
	}
	c.Cache = c.PPipeline.Cache
	c.DependsOn = c.PPipeline.DependsOn
	return nil
}

// PackageJSON represents NodeJS package.json
type PackageJSON struct {
	Name                   string            `json:"name,omitempty"`
	Version                string            `json:"version,omitempty"`
	Scripts                map[string]string `json:"scripts,omitempty"`
	Dependencies           map[string]string `json:"dependencies,omitempty"`
	DevDependencies        map[string]string `json:"devDependencies,omitempty"`
	OptionalDependencies   map[string]string `json:"optionalDependencies,omitempty"`
	PeerDependencies       map[string]string `json:"peerDependencies,omitempty"`
	PackageManager         string            `json:"packageManager,omitempty"`
	Os                     []string          `json:"os,omitempty"`
	Workspaces             Workspaces        `json:"workspaces,omitempty"`
	Private                bool              `json:"private,omitempty"`
	PackageJSONPath        string
	Hash                   string
	Dir                    string
	InternalDeps           []string
	UnresolvedExternalDeps map[string]string
	ExternalDeps           []string
	SubLockfile            YarnLockfile
	Turbo                  TurboConfigJSON `json:"turbo"`
	Mu                     sync.Mutex
	FilesHash              string
	ExternalDepsHash       string
}

type Workspaces []string

type WorkspacesAlt struct {
	Packages []string `json:"packages,omitempty"`
}

func (r *Workspaces) UnmarshalJSON(data []byte) error {
	var tmp = &WorkspacesAlt{}
	if err := json.Unmarshal(data, tmp); err == nil {
		*r = Workspaces(tmp.Packages)
		return nil
	}
	var tempstr = []string{}
	if err := json.Unmarshal(data, &tempstr); err != nil {
		return err
	}
	*r = tempstr
	return nil
}

// Parse parses package.json payload and returns structure.
func Parse(payload []byte) (*PackageJSON, error) {
	var packagejson *PackageJSON
	err := json.Unmarshal(payload, &packagejson)
	return packagejson, err
}

// ReadPackageJSON returns a struct of package.json
func ReadPackageJSON(path string) (*PackageJSON, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}
