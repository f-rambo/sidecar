package interfaces

import (
	"context"
	"encoding/json"
	"io"
	"os"

	v1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/f-rambo/ship/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SystemInterface struct {
	v1alpha1.UnimplementedSystemInterfaceServer
	systemUc *biz.SystemUsecase
	log      *log.Helper
}

func NewSystemInterface(systemUc *biz.SystemUsecase, logger log.Logger) *SystemInterface {
	return &SystemInterface{
		systemUc: systemUc,
		log:      log.NewHelper(logger),
	}
}

// conn, err := grpc.DialInsecure(
// 	ctx,
// 	grpc.WithEndpoint(fmt.Sprintf("%s:%d", "localhost", 9001)),
// )
// if err != nil {
// 	return nil, err
// }
// defer conn.Close()
// client := systemv1alpha1.NewSystemInterfaceClient(conn)
// stream, err := client.GetLogs(ctx)
// if err != nil {
// 	return nil, err
// }

// // 创建一个goroutine来接收来自服务端的消息
// go func() {
// 	for {
// 		msg, err := stream.Recv()
// 		if err == io.EOF {
// 			return
// 		}
// 		if err != nil {
// 			log.Fatalf("Error receiving message: %v", err)
// 		}
// 		fmt.Printf("Server: %s\n", msg.Log)
// 	}
// }()

// i := 0
// for {
// 	err = stream.Send(&systemv1alpha1.LogRequest{
// 		TailLines: 10,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	time.Sleep(1 * time.Second)
// 	if i >= 10 {
// 		stream.CloseSend()
// 		break
// 	}
// 	i++
// }

func (c *SystemInterface) GetLogs(stream v1alpha1.SystemInterface_GetLogsServer) error {
	var lastReadPos int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if req.TailLines == 0 {
			req.TailLines = 30
		}
		reqJson, err := json.Marshal(req)
		if err != nil {
			return err
		}
		clusterLogReq := &v1alpha1.LogRequest{}
		err = json.Unmarshal(reqJson, clusterLogReq)
		if err != nil {
			return err
		}
		clusterLogPath, err := utils.GetPackageStorePathByNames("log", "ship.log")
		if err != nil {
			return err
		}
		if ok := utils.IsFileExist(clusterLogPath); !ok {
			return errors.New("cluster log does not exist")
		}

		file, err := os.Open(clusterLogPath)
		if err != nil {
			return err
		}
		defer file.Close()

		var logs string
		if lastReadPos == 0 {
			// Read the last 30 lines
			logs, err = utils.ReadLastNLines(file, int(req.TailLines))
			if err != nil {
				return err
			}
		} else {
			// Read from the last read position
			_, err = file.Seek(lastReadPos, io.SeekStart)
			if err != nil {
				return err
			}
			newLogs, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			logs = string(newLogs)
		}

		// If logs are empty, send a "." character
		if logs == "" {
			logs = "."
		}

		err = stream.Send(&v1alpha1.LogResponse{
			Log: logs,
		})
		if err != nil {
			return err
		}

		lastReadPos, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}
	}
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
