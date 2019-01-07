package ping

import (
	"fmt"
	"log"
	"os/exec"
)

func Ping(destination string, option Option) (Result, error) {
	option.Validate()

	params := make([]string, 0)
	params = append(params, destination)
	params = append(params, fmt.Sprintf("-c %d", option.GetCount()))
	params = append(params, fmt.Sprintf("-i %d", option.GetInterval()))

	pCmd := exec.Command("ping", params...)
	log.Print(pCmd.Args)
	output, err := pCmd.CombinedOutput()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return parseLinuxPing(output)
}
