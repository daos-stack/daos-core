//
// (C) Copyright 2018-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/daos-stack/daos/src/control/build"
	ctlpb "github.com/daos-stack/daos/src/control/common/proto/ctl"
	mgmtpb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"github.com/daos-stack/daos/src/control/events"
	"github.com/daos-stack/daos/src/control/lib/control"
	"github.com/daos-stack/daos/src/control/lib/netdetect"
	"github.com/daos-stack/daos/src/control/logging"
	"github.com/daos-stack/daos/src/control/pbin"
	"github.com/daos-stack/daos/src/control/security"
	"github.com/daos-stack/daos/src/control/server/config"
	"github.com/daos-stack/daos/src/control/server/engine"
	"github.com/daos-stack/daos/src/control/server/storage/bdev"
	"github.com/daos-stack/daos/src/control/server/storage/scm"
	"github.com/daos-stack/daos/src/control/system"
)

const (
	iommuPath        = "/sys/class/iommu"
	minHugePageCount = 128
)

func cfgHasBdev(cfg *config.Server) bool {
	for _, engineCfg := range cfg.Engines {
		if len(engineCfg.Storage.Bdev.DeviceList) > 0 {
			return true
		}
	}

	return false
}

func iommuDetected() bool {
	// Simple test for now -- if the path exists and contains
	// DMAR entries, we assume that's good enough.
	dmars, err := ioutil.ReadDir(iommuPath)
	if err != nil {
		return false
	}

	return len(dmars) > 0
}

func raftDir(cfg *config.Server) string {
	if len(cfg.Engines) == 0 {
		return "" // can't save to SCM
	}
	return filepath.Join(cfg.Engines[0].Storage.SCM.MountPoint, "control_raft")
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("Hostname() failed: %s", err.Error())
	}
	return hn
}

// Start is the entry point for a daos_server instance.
func Start(log *logging.LeveledLogger, cfg *config.Server) error {
	err := cfg.Validate(log)
	if err != nil {
		return errors.Wrapf(err, "%s: validation failed", cfg.Path)
	}

	// Temporary notification while the feature is still being polished.
	if len(cfg.AccessPoints) > 1 {
		log.Info("\n*******\nNOTICE: Support for multiple access points is an alpha feature and is not well-tested!\n*******\n\n")
	}

	// Backup active config.
	config.SaveActiveConfig(log, cfg)

	if cfg.HelperLogFile != "" {
		if err := os.Setenv(pbin.DaosAdminLogFileEnvVar, cfg.HelperLogFile); err != nil {
			return errors.Wrap(err, "unable to configure privileged helper logging")
		}
	}

	if cfg.FWHelperLogFile != "" {
		if err := os.Setenv(pbin.DaosFWLogFileEnvVar, cfg.FWHelperLogFile); err != nil {
			return errors.Wrap(err, "unable to configure privileged firmware helper logging")
		}
	}

	faultDomain, err := getFaultDomain(cfg)
	if err != nil {
		return err
	}
	log.Debugf("fault domain: %s", faultDomain.String())

	// Create the root context here. All contexts should
	// inherit from this one so that they can be shut down
	// from one place.
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	controlAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.ControlPort))
	if err != nil {
		return errors.Wrap(err, "unable to resolve daos_server control address")
	}

	bdevProvider := bdev.DefaultProvider(log)
	runningUser, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "unable to lookup current user")
	}

	iommuDisabled := !iommuDetected()
	// Perform an automatic prepare based on the values in the config file.
	prepReq := bdev.PrepareRequest{
		// Default to minimum necessary for scan to work correctly.
		HugePageCount: minHugePageCount,
		TargetUser:    runningUser.Username,
		PCIAllowlist:  strings.Join(cfg.BdevInclude, " "),
		PCIBlocklist:  strings.Join(cfg.BdevExclude, " "),
		DisableVFIO:   cfg.DisableVFIO,
		DisableVMD:    cfg.DisableVMD || cfg.DisableVFIO || iommuDisabled,
		// TODO: pass vmd include/white list
	}

	if cfgHasBdev(cfg) {
		// The config value is intended to be per-engine, so we need to adjust
		// based on the number of engines.
		prepReq.HugePageCount = cfg.NrHugepages * len(cfg.Engines)

		// Perform these checks to avoid even trying a prepare if the system
		// isn't configured properly.
		if runningUser.Uid != "0" {
			if cfg.DisableVFIO {
				return FaultVfioDisabled
			}

			if iommuDisabled {
				return FaultIommuDisabled
			}
		}
	}

	log.Debugf("automatic NVMe prepare req: %+v", prepReq)
	if _, err := bdevProvider.Prepare(prepReq); err != nil {
		log.Errorf("automatic NVMe prepare failed (check configuration?)\n%s", err)
	}

	hugePages, err := getHugePageInfo()
	if err != nil {
		return errors.Wrap(err, "unable to read system hugepage info")
	}

	if cfgHasBdev(cfg) {
		// Double-check that we got the requested number of huge pages after prepare.
		if hugePages.Free < prepReq.HugePageCount {
			return FaultInsufficientFreeHugePages(hugePages.Free, prepReq.HugePageCount)
		}
	}

	var dbReplicas []*net.TCPAddr
	for _, ap := range cfg.AccessPoints {
		apAddr, err := net.ResolveTCPAddr("tcp", ap)
		if err != nil {
			return config.FaultConfigBadAccessPoints
		}
		dbReplicas = append(dbReplicas, apAddr)
	}

	// If this daos_server instance ends up being the MS leader,
	// this will record the DAOS system membership.
	sysdb, err := system.NewDatabase(log, &system.DatabaseConfig{
		Replicas:   dbReplicas,
		RaftDir:    raftDir(cfg),
		SystemName: cfg.SystemName,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create system database")
	}
	membership := system.NewMembership(log, sysdb)
	scmProvider := scm.DefaultProvider(log)
	harness := NewEngineHarness(log).WithFaultDomain(faultDomain)

	// Create rpcClient for inter-server communication.
	cliCfg := control.DefaultConfig()
	cliCfg.TransportConfig = cfg.TransportConfig
	rpcClient := control.NewClient(
		control.WithConfig(cliCfg),
		control.WithClientLogger(log))

	// Create event distributor.
	eventPubSub := events.NewPubSub(ctx, log)
	defer eventPubSub.Close()

	// Init management RPC subsystem.
	mgmtSvc := newMgmtSvc(harness, membership, sysdb, rpcClient, eventPubSub)

	// Forward published actionable events (type RASTypeStateChange) to the
	// management service leader, behavior is updated on leadership change.
	eventForwarder := control.NewEventForwarder(rpcClient, cfg.AccessPoints)
	eventPubSub.Subscribe(events.RASTypeStateChange, eventForwarder)
	// Log events on the host that they were raised (and first published) on.
	eventLogger := control.NewEventLogger(log)
	eventPubSub.Subscribe(events.RASTypeAny, eventLogger)

	var netDevClass uint32

	netCtx, err := netdetect.Init(context.Background())
	if err != nil {
		return err
	}
	defer netdetect.CleanUp(netCtx)

	// On a NUMA-aware system, emit a message when the configuration
	// may be sub-optimal.
	numaCount := netdetect.NumNumaNodes(netCtx)
	if numaCount > 0 && len(cfg.Engines) > numaCount {
		log.Infof("NOTICE: Detected %d NUMA node(s); %d-server config may not perform as expected", numaCount, len(cfg.Engines))
	}

	// Create a closure to be used for joining engine instances.
	joinInstance := func(ctx context.Context, req *control.SystemJoinReq) (*control.SystemJoinResp, error) {
		req.SetHostList(cfg.AccessPoints)
		req.SetSystem(cfg.SystemName)
		req.ControlAddr = controlAddr
		return control.SystemJoin(ctx, rpcClient, req)
	}

	var allStarted sync.WaitGroup

	for idx, engineCfg := range cfg.Engines {
		// Provide special handling for the ofi+verbs provider.
		// Mercury uses the interface name such as ib0, while OFI uses the
		// device name such as hfi1_0 CaRT and Mercury will now support the
		// new OFI_DOMAIN environment variable so that we can specify the
		// correct device for each.
		if strings.HasPrefix(engineCfg.Fabric.Provider, "ofi+verbs") && !engineCfg.HasEnvVar("OFI_DOMAIN") {
			deviceAlias, err := netdetect.GetDeviceAlias(netCtx, engineCfg.Fabric.Interface)
			if err != nil {
				return errors.Wrapf(err, "failed to resolve alias for %s", engineCfg.Fabric.Interface)
			}
			envVar := "OFI_DOMAIN=" + deviceAlias
			engineCfg.WithEnvVars(envVar)
		}

		// If the configuration specifies that we should explicitly set
		// hugepage values per engine instance, do it. Otherwise, let
		// SPDK/DPDK figure it out.
		if cfg.SetHugepages {
			// If we have multiple engine instances with block devices, then
			// apportion the hugepage memory among the instances.
			engineCfg.Storage.Bdev.MemSize = hugePages.FreeMB() / len(cfg.Engines)
			// reserve a little for daos_admin
			engineCfg.Storage.Bdev.MemSize -= engineCfg.Storage.Bdev.MemSize / 16
		}

		// Indicate whether VMD devices have been detected and can be used.
		engineCfg.Storage.Bdev.VmdDisabled = bdevProvider.IsVMDDisabled()

		bp, err := bdev.NewClassProvider(log, engineCfg.Storage.SCM.MountPoint, &engineCfg.Storage.Bdev)
		if err != nil {
			return err
		}

		engine := NewEngineInstance(log, bp, scmProvider, joinInstance, engine.NewRunner(log, engineCfg)).
			WithHostFaultDomain(faultDomain)
		if err := harness.AddInstance(engine); err != nil {
			return err
		}

		// Register callback to publish I/O Engine process exit events.
		engine.OnInstanceExit(publishInstanceExitFn(eventPubSub.Publish, hostname(), engine.Index()))

		var onceReady sync.Once
		allStarted.Add(1)
		engine.OnReady(func(_ context.Context) error {
			// Indicate that engine has been started, only do this
			// the first time that the engine starts as shared
			// memory persists between engine restarts.
			onceReady.Do(func() {
				allStarted.Done()
			})
			return nil
		})

		if idx == 0 {
			netDevClass, err = cfg.GetDeviceClassFn(engineCfg.Fabric.Interface)
			if err != nil {
				return err
			}

			if !sysdb.IsReplica() {
				continue
			}

			// Start the system db after instance 0's SCM is
			// ready.
			var onceStorageReady sync.Once
			engine.OnStorageReady(func(ctx context.Context) (err error) {
				onceStorageReady.Do(func() {
					err = errors.Wrap(sysdb.Start(ctx),
						"failed to start system db",
					)
				})
				return
			})

			if !sysdb.IsBootstrap() {
				continue
			}

			// For historical reasons, we reserve rank 0 for the first
			// instance on the raft bootstrap server. This implies that
			// rank 0 will always be associated with a MS replica, but
			// it is not guaranteed to always be the leader.
			engine.joinSystem = func(ctx context.Context, req *control.SystemJoinReq) (*control.SystemJoinResp, error) {
				if sb := engine.getSuperblock(); !sb.ValidRank {
					engine.log.Debug("marking bootstrap instance as rank 0")
					req.Rank = 0
					sb.Rank = system.NewRankPtr(0)
				}
				return joinInstance(ctx, req)
			}
		}
	}

	go func() {
		allStarted.Wait()

		if cfg.TelemetryPort == 0 {
			return
		}

		log.Debug("starting Prometheus exporter")
		if err := startPrometheusExporter(ctx, log, cfg.TelemetryPort, harness.Instances()); err != nil {
			log.Errorf("failed to start prometheus exporter: %s", err)
		}
	}()

	// Create and setup control service.
	controlService := NewControlService(log, harness, bdevProvider, scmProvider, cfg, eventPubSub)
	if err := controlService.Setup(); err != nil {
		return errors.Wrap(err, "setup control service")
	}

	// Create and start listener on management network.
	lis, err := net.Listen("tcp4", controlAddr.String())
	if err != nil {
		return errors.Wrap(err, "unable to listen on management interface")
	}

	// Create new grpc server, register services and start serving.
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		unaryErrorInterceptor,
		unaryStatusInterceptor,
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		streamErrorInterceptor,
	}
	tcOpt, err := security.ServerOptionForTransportConfig(cfg.TransportConfig)
	if err != nil {
		return err
	}
	srvOpts := []grpc.ServerOption{tcOpt}

	uintOpt, err := unaryInterceptorForTransportConfig(cfg.TransportConfig)
	if err != nil {
		return err
	}
	if uintOpt != nil {
		unaryInterceptors = append(unaryInterceptors, uintOpt)
	}
	sintOpt, err := streamInterceptorForTransportConfig(cfg.TransportConfig)
	if err != nil {
		return err
	}
	if sintOpt != nil {
		streamInterceptors = append(streamInterceptors, sintOpt)
	}
	srvOpts = append(srvOpts, []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}...)

	grpcServer := grpc.NewServer(srvOpts...)
	ctlpb.RegisterCtlSvcServer(grpcServer, controlService)

	mgmtSvc.clientNetworkCfg = &config.ClientNetworkCfg{
		Provider:        cfg.Fabric.Provider,
		CrtCtxShareAddr: cfg.Fabric.CrtCtxShareAddr,
		CrtTimeout:      cfg.Fabric.CrtTimeout,
		NetDevClass:     netDevClass,
	}
	mgmtpb.RegisterMgmtSvcServer(grpcServer, mgmtSvc)

	tSec, err := security.DialOptionForTransportConfig(cfg.TransportConfig)
	if err != nil {
		return err
	}
	sysdb.ConfigureTransport(grpcServer, tSec)
	sysdb.OnLeadershipGained(func(ctx context.Context) error {
		log.Infof("MS leader running on %s", hostname())
		mgmtSvc.startJoinLoop(ctx)

		// Stop forwarding events to MS and instead start handling
		// received forwarded (and local) events.
		eventPubSub.Reset()
		eventPubSub.Subscribe(events.RASTypeAny, eventLogger)
		eventPubSub.Subscribe(events.RASTypeStateChange, membership)
		eventPubSub.Subscribe(events.RASTypeStateChange, sysdb)
		eventPubSub.Subscribe(events.RASTypeStateChange, events.HandlerFunc(func(ctx context.Context, evt *events.RASEvent) {
			switch evt.ID {
			case events.RASSwimRankDead:
				// Mark the rank as unavailable for membership in
				// new pools, etc. Do group update on success.
				if err := membership.MarkRankDead(system.Rank(evt.Rank)); err == nil {
					mgmtSvc.reqGroupUpdate(ctx)
				}
			}
		}))

		return nil
	})
	sysdb.OnLeadershipLost(func() error {
		log.Infof("MS leader no longer running on %s", hostname())

		// Stop handling received forwarded (in addition to local)
		// events and start forwarding events to the new MS leader.
		eventPubSub.Reset()
		eventPubSub.Subscribe(events.RASTypeAny, eventLogger)
		eventPubSub.Subscribe(events.RASTypeStateChange, eventForwarder)

		return nil
	})

	go func() {
		_ = grpcServer.Serve(lis)
	}()
	defer grpcServer.Stop()

	log.Infof("%s v%s (pid %d) listening on %s", build.ControlPlaneName, build.DaosVersion, os.Getpid(), controlAddr)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		// SIGKILL I/O Engine immediately on exit.
		// TODO: Re-enable attempted graceful shutdown of I/O Engines.
		sig := <-sigChan
		log.Debugf("Caught signal: %s", sig)

		shutdown()
	}()

	return errors.Wrapf(harness.Start(ctx, sysdb, eventPubSub, cfg), "%s exited with error", build.DataPlaneName)
}
