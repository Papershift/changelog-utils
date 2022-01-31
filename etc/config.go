package etc

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const CONFIG_FILENAME = ".ch.conf.hcl"

type UserConfig struct {
	GithubUsername string `hcl:"github_username"`
	Fullname       string `hcl:"fullname"`
}

type Config struct {
	User UserConfig `hcl:"user,block"`
}

func ConfigExists() bool {
	_, err := os.Stat(fmt.Sprintf("changelogs/%s", CONFIG_FILENAME))
	return err == nil
}

func MakeConfigFile(conf Config) error {
	path := fmt.Sprintf("changelogs/%s", CONFIG_FILENAME)
	file := hclwrite.NewEmptyFile()

	userBlock := file.Body().AppendNewBlock("user", nil)
	userBody := userBlock.Body()

	userBody.SetAttributeValue("github_username", cty.StringVal(conf.User.GithubUsername))
	userBody.SetAttributeValue("fullname", cty.StringVal(conf.User.Fullname))

	err := os.WriteFile(path, file.Bytes(), 0644)

	if err == nil {
		fmt.Printf(">> %s written\n", path)
	}

	return err
}

func ReadConfig() (Config, error) {
	var conf Config
	path := fmt.Sprintf("changelogs/%s", CONFIG_FILENAME)

	err := hclsimple.DecodeFile(path, nil, &conf)

	return conf, err
}
