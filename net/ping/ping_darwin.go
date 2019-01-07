package ping

import (
	"fmt"
	"os/exec"
)

func Ping(destination string, option Option) (Result, error) {
	option.Validate()

	params := make([]string, 0)
	params = append(params, destination)
	params = append(params, fmt.Sprintf("-c %d", option.GetCount()))
	params = append(params, fmt.Sprintf("-i %d", option.GetInterval()))

	pCmd := exec.Command("ping", params...)
	output, err := pCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseDarwinPing(output)
}
