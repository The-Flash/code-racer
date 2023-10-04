package exec_utils

import (
	"context"

	"github.com/The-Flash/code-racer/internal/cappedbuffer"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type ExecCmdConfig struct {
	*types.ExecConfig
	StdOutSizeLimit int
	StdErrSizeLimit int
}

// ExecCmd execute a command in an executor
func ExecCmd(cmd []string, config ExecCmdConfig, cli *client.Client, containerID string) (stdout cappedbuffer.CappedBuffer, stderr cappedbuffer.CappedBuffer, err error) {
	execCreateResp, err := cli.ContainerExecCreate(context.Background(), containerID, types.ExecConfig{
		Tty:          config.Tty,
		Cmd:          cmd,
		WorkingDir:   config.WorkingDir,
		AttachStderr: config.AttachStderr,
		AttachStdout: config.AttachStdout,
		Detach:       config.Detach,
	})
	if err != nil {
		return
	}
	execAttachResp, err := cli.ContainerExecAttach(context.Background(), execCreateResp.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		return
	}
	defer execAttachResp.Close()
	stdout = *cappedbuffer.New(make([]byte, 0, config.StdOutSizeLimit), config.StdOutSizeLimit)
	stderr = *cappedbuffer.New(make([]byte, 0, config.StdErrSizeLimit), config.StdErrSizeLimit)

	if err != nil {
		return
	}
	_, err = stdcopy.StdCopy(&stdout, &stderr, execAttachResp.Reader)
	return
}
