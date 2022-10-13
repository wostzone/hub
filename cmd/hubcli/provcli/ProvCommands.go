package provcli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/hiveot/hub/internal/listener"
	"github.com/hiveot/hub/internal/svcconfig"
	"github.com/hiveot/hub/pkg/provisioning"
	"github.com/hiveot/hub/pkg/provisioning/capnpclient"
)

// ProvisioningCommands returns the provisioning handling commands
// This requires the provisioning service to run.
func ProvisioningCommands(ctx context.Context, f svcconfig.AppFolders) *cli.Command {

	cmd := &cli.Command{
		//hub prov add|list  <deviceID> <secret>

		Name:  "prov",
		Usage: "IoT device provisioning",
		Subcommands: cli.Commands{
			ProvisionAddOOBSecretsCommand(ctx, f),
			ProvisionApproveRequestCommand(ctx, f),
			ProvisionGetPendingRequestsCommand(ctx, f),
			ProvisionGetApprovedRequestsCommand(ctx, f),
		},
	}

	return cmd
}

// ProvisionAddOOBSecretsCommand
// prov add  <deviceID> <oobsecret>
func ProvisionAddOOBSecretsCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add an out-of-band device provisioning secret for automatic provisioning",
		ArgsUsage: "<deviceID> <oobSecret>",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 2 {
				return fmt.Errorf("expected 2 arguments. Got %d instead", cCtx.NArg())
			}
			err := HandleAddOobSecret(ctx, f,
				cCtx.Args().Get(0),
				cCtx.Args().Get(1))
			fmt.Println("Adding secret for device: ", cCtx.Args().First())
			return err
		},
	}
}

// ProvisionApproveRequestCommand
// prov approve <deviceID>
func ProvisionApproveRequestCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "approve",
		Usage:     "Approve a pending provisioning request",
		ArgsUsage: "<deviceID> ",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				return fmt.Errorf("expected 1 arguments. Got %d instead", cCtx.NArg())
			}
			deviceID := cCtx.Args().First()
			err := HandleApproveRequest(ctx, f, deviceID)
			return err
		},
	}
}

// ProvisionGetApprovedRequestsCommand
// prov approved
func ProvisionGetApprovedRequestsCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "approved",
		Usage:     "Get a list of approved provisioning requests",
		ArgsUsage: "(no arguments)",
		Action: func(cCtx *cli.Context) error {
			err := HandleGetApprovedRequests(ctx, f)
			return err
		},
	}
}

// ProvisionGetPendingRequestsCommand
// prov approved
func ProvisionGetPendingRequestsCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "pending",
		Usage:     "Get a list of pending provisioning requests",
		ArgsUsage: "(no arguments)",
		Action: func(cCtx *cli.Context) error {
			err := HandleGetPendingRequests(ctx, f)
			return err
		},
	}
}

// HandleAddOobSecret invokes the out-of-band provisioning service to add a provisioning secret
//  deviceID is the ID of the device whose secret to set
//  secret to set
func HandleAddOobSecret(ctx context.Context, f svcconfig.AppFolders, deviceID string, secret string) error {
	var pc provisioning.IProvisioning
	var secrets []provisioning.OOBSecret

	conn, err := listener.CreateClientConnection(f.Run, provisioning.ServiceName)
	if err == nil {
		pc, err = capnpclient.NewProvisioningCapnpClient(ctx, conn)
	}
	if err != nil {
		return err
	}
	manage := pc.CapManageProvisioning()

	secrets = []provisioning.OOBSecret{
		{
			DeviceID:  deviceID,
			OobSecret: secret,
		},
	}
	err = manage.AddOOBSecrets(ctx, secrets)

	return err
}

// HandleApproveRequest
//  deviceID is the ID of the device to approve
func HandleApproveRequest(ctx context.Context, f svcconfig.AppFolders, deviceID string) error {
	var pc provisioning.IProvisioning

	conn, err := listener.CreateClientConnection(f.Run, provisioning.ServiceName)
	if err == nil {
		pc, err = capnpclient.NewProvisioningCapnpClient(ctx, conn)
	}
	pc.CapManageProvisioning()
	if err != nil {
		return err
	}
	manage := pc.CapManageProvisioning()
	err = manage.ApproveRequest(ctx, deviceID)

	return err
}

// HandleGetApprovedRequests
func HandleGetApprovedRequests(ctx context.Context, f svcconfig.AppFolders) error {
	var pc provisioning.IProvisioning

	conn, err := listener.CreateClientConnection(f.Run, provisioning.ServiceName)
	if err == nil {
		pc, err = capnpclient.NewProvisioningCapnpClient(ctx, conn)
	}
	if err != nil {
		return err
	}
	manage := pc.CapManageProvisioning()
	provStatus, err := manage.GetApprovedRequests(ctx)
	fmt.Printf("Client ID              Request Time      Assigned\n")
	fmt.Printf("--------------------   ------------      --------\n")
	for _, provStatus := range provStatus {
		// a certificate is assigned when generated
		assigned := provStatus.ClientCertPEM != ""
		fmt.Printf("%20s  %s, %s\n",
			provStatus.DeviceID, provStatus.RequestTime, assigned)
	}

	return err
}

// HandleGetPendingRequests
func HandleGetPendingRequests(ctx context.Context, f svcconfig.AppFolders) error {
	var pc provisioning.IProvisioning

	conn, err := listener.CreateClientConnection(f.Run, provisioning.ServiceName)
	if err == nil {
		pc, err = capnpclient.NewProvisioningCapnpClient(ctx, conn)
	}
	if err != nil {
		return err
	}
	manage := pc.CapManageProvisioning()
	provStatus, err := manage.GetPendingRequests(ctx)
	fmt.Printf("Client ID              Request Time\n")
	fmt.Printf("--------------------   ------------\n")
	for _, provStatus := range provStatus {
		// a certificate is assigned when generated
		fmt.Printf("%20s  %s\n",
			provStatus.DeviceID, provStatus.RequestTime)
	}

	return err
}
