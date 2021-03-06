/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package config

import (
	"fmt"
	"github.com/aliyun/aliyun-cli/cli"
	"github.com/aliyun/aliyun-cli/i18n"
	"os"
	"text/tabwriter"
)

func NewConfigureListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list",
		Short: i18n.T("list all config profile", "列出所有配置集"),
		Run: func(c *cli.Context, args []string) error {
			doConfigureList()
			return nil
		},
	}
}

func doConfigureList() {
	conf, err := LoadConfiguration()
	if err != nil {
		cli.Errorf("ERROR: load configure failed: %v\n", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', 0)
	fmt.Fprint(w, "Profile\t| Credential \t| Valid\t| Region\t| Language\n")
	fmt.Fprint(w, "---------\t| ------------------\t| -------\t| ----------------\t| --------\n")
	for _, pf := range conf.Profiles {
		name := pf.Name
		if name == conf.CurrentProfile {
			name = name + " *"
		}
		err := pf.Validate()
		valid := "Valid"
		if err != nil {
			valid = "Invalid"
		}

		cred := ""
		switch pf.Mode {
		case AK:
			cred = "AK:" + "***" + GetLastChars(pf.AccessKeyId, 3)
		case StsToken:
			cred = "StsToken:" + "***" + GetLastChars(pf.AccessKeyId, 3)
		case RamRoleArn:
			cred = "RamRoleArn:" + "***" + GetLastChars(pf.AccessKeyId, 3)
		case EcsRamRole:
			cred = "EcsRamRole:" + pf.RamRoleName
		case RsaKeyPair:
			cred = "RsaKeyPair:" + pf.KeyPairName
		}
		fmt.Fprintf(w, "%s\t| %s\t| %s\t| %s\t| %s\n", name, cred, valid, pf.RegionId, pf.Language)
	}
	w.Flush()
}
