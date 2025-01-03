package interfaces

import (
	"context"
	"errors"
	"io"
	"os"

	clusterApi "github.com/f-rambo/cloud-copilot/sidecar/api/cluster"
	"github.com/f-rambo/cloud-copilot/sidecar/api/common"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/biz"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/conf"
	"github.com/f-rambo/cloud-copilot/sidecar/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterInterface struct {
	clusterApi.UnimplementedClusterInterfaceServer
	uc  *biz.ClusterUsecase
	log *log.Helper
	c   *conf.Server
}

func NewClusterInterface(uc *biz.ClusterUsecase, c *conf.Server, logger log.Logger) *ClusterInterface {
	return &ClusterInterface{
		uc:  uc,
		log: log.NewHelper(logger),
		c:   c,
	}
}

func (c *ClusterInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*common.Msg, error) {
	return common.Response(), nil
}

func (c *ClusterInterface) Info(ctx context.Context, _ *emptypb.Empty) (*clusterApi.Cluster, error) {
	return &clusterApi.Cluster{
		Name: c.c.Name,
	}, nil
}

func (c *ClusterInterface) GetLogs(stream clusterApi.ClusterInterface_GetLogsServer) error {
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

		logpath := utils.GetLogFilePath()
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
			err = stream.Send(&clusterApi.LogResponse{Log: initialLogs})
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
							err = stream.Send(&clusterApi.LogResponse{Log: newLogs})
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
