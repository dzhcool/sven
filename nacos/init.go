package nacos

// @date 2021-08-18 11:21

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dzhcool/sven/utils"
)

var (
	endpoint  string
	namespace string
	accessKey string
	secretKey string
)

func InitConfig(dataId, group string) (string, error) {
	endpoint = utils.Getenv("GP_NACOS_ENDPOINT", utils.Getenv("NACOS_ENDPOINT", ""))
	namespace = utils.Getenv("GP_NACOS_NAMESPACE", utils.Getenv("NACOS_NAMESPACE", ""))
	accessKey = utils.Getenv("GP_NACOS_ACCESSKEY", utils.Getenv("NACOS_ACCESSKEY", ""))
	secretKey = utils.Getenv("GP_NACOS_SECRETKEY", utils.Getenv("NACOS_SECRETKEY", ""))

	workPath := utils.Getenv("GP_WORK_PATH", utils.Getenv("WORK_PATH", ""))
	if workPath != "" {
		workPath = strings.TrimRight(workPath, "/") + "/"
	}
	confFile := workPath + "conf/" + dataId

	log.Printf("[nacos env] endpoint:%s namespace:%s accessKey:%s, secretKey:%s confFile:%s \n",
		endpoint, namespace, accessKey, secretKey, confFile)

	endpoint = strings.ToLower(strings.TrimSuffix(endpoint, "/"))
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}
	api := fmt.Sprintf("%s/nacos/v1/cs/configs?tenant=%s&dataId=%s&group=%s", endpoint, namespace, dataId, group)

	var err error
	for retry := 0; retry < 3; retry++ {
		if _, err = utils.HTTPDownload(api, confFile); err != nil {
			if err != nil {
				continue
			}
			break
		}
	}
	if err != nil {
		log.Fatalf("get nacos file failed: %s", err)
		os.Exit(1)
	}
	return confFile, err
}
