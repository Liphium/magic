package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/Liphium/magic/v3/util"
)

var _ mconfig.ServiceDriver = &RedisDriver{}

var redisLog *log.Logger = log.New(os.Stdout, "redis ", log.Default().Flags())

type RedisDriver struct {
	Image string `json:"image"`
}

func init() {
	mconfig.RegisterDriver(&RedisDriver{})
}

func NewDriver(image string) *RedisDriver {
	imageVersion := strings.Split(image, ":")[1]

	var supportedRedisVersions = []int{7, 8}

	supported := false
	imageMajor := mservices.GetImageMajorVersion(image)
	for _, version := range supportedRedisVersions {
		if imageMajor == version {
			supported = true
		}
	}
	if !supported {
		redisLog.Fatalln("ERROR: Version", imageVersion, "is currently not supported.")
	}

	return &RedisDriver{
		Image: image,
	}
}

func (rd *RedisDriver) Load(data string) (mconfig.ServiceDriver, error) {
	var driver RedisDriver
	if err := json.Unmarshal([]byte(data), &driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (rd *RedisDriver) Save() (string, error) {
	bytes, err := json.Marshal(rd)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (rd *RedisDriver) GetUniqueId() string {
	return "redis"
}

func (rd *RedisDriver) GetRequiredPortAmount() int {
	return 1
}

func (rd *RedisDriver) GetImage() string {
	return rd.Image
}

func (rd *RedisDriver) Host(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueStatic("127.0.0.1")
}

func (rd *RedisDriver) Port(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueFunction(func() string {
		for id, container := range ctx.Plan().Containers {
			if id == rd.GetUniqueId() {
				return fmt.Sprintf("%d", ctx.Plan().AllocatedPorts[container.Ports[0]])
			}
		}

		util.Log.Fatalln("ERROR: Couldn't find port for Redis container in plan!")
		return "not found"
	})
}
