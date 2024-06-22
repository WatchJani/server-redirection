package consistent_hashing

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"sort"
)

const ConfigFile = "../config.json"

type ConsistentHashing struct {
	*Configuration
	VirtualServers []VirtualServer
	Namer          int
}

type Server struct {
	Addr string `json:"ip"`
	Name string `json:"name"`
}

type VirtualServer struct {
	Name      string
	HashValue uint32
	*Server
}

type Configuration struct {
	VirtualNodeServer int      `json:"nodeNumber"`
	Servers           []Server `json:"servers"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Servers: make([]Server, 0),
	}
}

func (c *ConsistentHashing) NewVirtualServer(server *Server) VirtualServer {
	c.Namer = (c.Namer + 1) % c.VirtualNodeServer
	name := fmt.Sprintf("%s_%d", server.Name, c.Namer)

	return VirtualServer{
		Name:      name,
		HashValue: crc32.ChecksumIEEE([]byte(name)),
		Server:    server,
	}
}

func NewConsistentHashing() *ConsistentHashing {
	return &ConsistentHashing{
		Configuration:  NewConfiguration(),
		VirtualServers: make([]VirtualServer, 0),
	}
}

// load configuration file and decode
func (c *ConsistentHashing) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c.Configuration); err != nil {
		return err
	}

	return nil
}

// create all virtual server in ring
func (c *ConsistentHashing) Create() {
	for i := 0; i < len(c.Servers); i++ {
		for j := 0; j < c.VirtualNodeServer; j++ {
			c.VirtualServers = append(c.VirtualServers, c.NewVirtualServer(&c.Servers[i]))
		}
	}
}

func (c *ConsistentHashing) Sort() {
	sort.Slice(c.VirtualServers, func(i, j int) bool {
		return c.VirtualServers[i].HashValue < c.VirtualServers[j].HashValue
	})
}

func (c *ConsistentHashing) Load(path string) error {
	if err := c.LoadConfig(path); err != nil {
		return err
	}

	c.Create() //create all virtual servers
	c.Sort()   //sort all virtual servers

	return nil
}

func Search(arr []VirtualServer, target uint32) int {
	left, right := 0, len(arr)-1
	result := len(arr)

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid].HashValue >= target {
			result = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return result
}

func (c *ConsistentHashing) FindServer(key string) *Server {
	newHash := crc32.ChecksumIEEE([]byte(key))

	virtualServerIndex := Search(c.VirtualServers, newHash)
	return c.VirtualServers[virtualServerIndex].Server
}
