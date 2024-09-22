package interfaces

import (
	"context"
	"io"
	"os"

	v1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SystemInterface struct {
	v1alpha1.UnimplementedSystemInterfaceServer
	systemUc *biz.SystemUsecase
	log      *log.Helper
	c        *conf.Server
}

func NewSystemInterface(systemUc *biz.SystemUsecase, logger log.Logger, c *conf.Server) *SystemInterface {
	return &SystemInterface{
		systemUc: systemUc,
		log:      log.NewHelper(logger),
		c:        c,
	}
}

func (c *SystemInterface) GetLogs(stream v1alpha1.SystemInterface_GetLogsServer) error {
	i := 0
	for {
		ctx, cancel := context.WithCancel(stream.Context())
		defer cancel()
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if i > 0 {
			c.log.Info("repeat message, don't need to process")
			continue
		}
		i++

		if req.TailLines == 0 {
			req.TailLines = 30
		}

		logpath, err := utils.GetLogFilePath(c.c.Name)
		if err != nil {
			return err
		}
		if ok := utils.IsFileExist(logpath); !ok {
			return errors.New("log file does not exist")
		}

		file, err := os.Open(logpath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read initial lines if TailLines is specified
		if req.TailLines > 0 {
			initialLogs, err := utils.ReadLastNLines(file, int(req.TailLines))
			if err != nil {
				return err
			}
			err = stream.Send(&v1alpha1.LogResponse{Log: initialLogs})
			if err != nil {
				return err
			}
		}

		// Move to the end of the file
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}

		// Start watching for new logs
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		err = watcher.Add(logpath)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						newLogs, err := readNewLines(file)
						if err != nil {
							return
						}
						if newLogs != "" {
							err = stream.Send(&v1alpha1.LogResponse{Log: newLogs})
							if err != nil {
								return
							}
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					c.log.Errorf("error watching log file: %v", err)
				case <-ctx.Done():
					c.log.Info("GetLogs stream closed by client")
					return
				}
			}
		}()
	}
}

func readNewLines(file *os.File) (string, error) {
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	newContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(newContent) > 0 {
		_, err = file.Seek(currentPos+int64(len(newContent)), io.SeekStart)
		if err != nil {
			return "", err
		}
		return string(newContent), nil
	}

	return "", nil
}

func (c *SystemInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*v1alpha1.Msg, error) {
	return &v1alpha1.Msg{Message: "pong"}, nil
}

func (c *SystemInterface) GetSystem(ctx context.Context, _ *emptypb.Empty) (*v1alpha1.System, error) {
	system, err := c.systemUc.GetSystem(ctx)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.System{
		Id:         system.ID,
		Os:         system.OS,
		Arch:       system.ARCH,
		Cpu:        system.CPU,
		Memory:     system.Memory,
		Gpu:        system.GPU,
		GpuSpec:    system.GpuSpec,
		DataDisk:   system.DataDisk,
		Kernel:     system.Kernel,
		Container:  system.Container,
		Kubelet:    system.Kubelet,
		KubeProxy:  system.KubeProxy,
		InternalIp: system.InternalIP,
	}, nil
}
