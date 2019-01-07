package ping

import (
	"fmt"
	"os/exec"
)

func Ping(destination string, option Option) (Result, error) {
	option.Validate()

	params := make([]string, 0)
	params = append(params, "/C")
	params = append(params, fmt.Sprintf("chcp 437 && ping %s -n %d -w %d", destination, option.GetCount(), option.GetDeadline()*1000))

	pCmd := exec.Command("cmd.exe", params...)
	output, err := pCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseWindowsPing(output)
}
