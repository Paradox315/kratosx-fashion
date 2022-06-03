package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	pb "kratosx-fashion/api/system/v1"
	"runtime"
)

type MonitorService struct {
	pb.UnimplementedMonitorServer

	log *log.Helper
}

func NewMonitorService(logger log.Logger) *MonitorService {
	return &MonitorService{
		log: log.NewHelper(log.With(logger, "service", "monitor")),
	}
}

func (s *MonitorService) GetRuntimeInfo(ctx context.Context, req *pb.EmptyRequest) (*pb.RuntimeReply, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &pb.RuntimeReply{
		Host:     info.Hostname,
		Os:       info.OS,
		Platform: info.Platform,
		Arch:     info.KernelArch,
		Version:  runtime.Version(),
		Compiler: runtime.Compiler,
		Cpus:     uint32(runtime.NumCPU()),
	}, nil
}

func (s *MonitorService) GetDiskInfo(ctx context.Context, req *pb.EmptyRequest) (*pb.DiskReply, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	return &pb.DiskReply{
		Path:        usage.Path,
		Total:       usage.Total,
		Free:        usage.Free,
		Used:        usage.Used,
		UsedPercent: float32(usage.UsedPercent),
	}, nil
}
