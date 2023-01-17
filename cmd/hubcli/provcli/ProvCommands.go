package provcli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/hiveot/hub/lib/hubclient"
	"github.com/hiveot/hub/lib/svcconfig"
	"github.com/hiveot/hub/pkg/provisioning"
	"github.com/hiveot/hub/pkg/provisioning/capnpclient"
)

// ProvisioningCommands returns the provisioning handling commands
// This requires the provisioning service to run.
func ProvisioningCommands(ctx context.Context, f svcconfig.AppFolders) *cli.Command {

	cmd := &cli.Command{
		//hub prov add|list  <deviceID> <secret>
		Name:    "provision",
		Aliases: []string{"pr"},
		Usage:   "IoT device provisioning",
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
		Name:      "addoob <deviceID> <secret>",
		Aliases:   []string{"ados"},
		Usage:     "Add a provisioning secret",
		UsageText: "Add an out-of-band device provisioning secret for automatic provisioning",
		Category:  "provisioning",
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
		Name:      "approveprov <deviceID>",
		Aliases:   []string{"appr"},
		Usage:     "Approve provisioning request",
		UsageText: "Approvide a pending provisioning request to issue a device authentication certificate",
		Category:  "provisioning",
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
		Name:      "listapproved",
		Aliases:   []string{"lap"},
		Usage:     "List approved provisioning requests",
		UsageText: "View a list of recent approved provisioning requests. ",
		Category:  "provisioning",
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
		Name:      "listpending",
		Aliases:   []string{"lip"},
		Usage:     "List pending provisioning requests",
		UsageText: "View a list of recent pending provisioning requests",
		Category:  "provisioning",
		Action: func(cCtx *cli.Context) error {
			err := HandleGetPendingRequests(ctx, f)
			return err
		},
	}
}

// HandleAddOobSecret invokes the out-of-band provisioning service to add a provisioning secret
//
//	deviceID is the ID of the device whose secret to set
//	secret to set
func HandleAddOobSecret(ctx context.Context, f svcconfig.AppFolders, deviceID string, secret string) error {
	var pc provisioning.IProvisioning
	var secrets []provisioning.OOBSecret

	conn, err := hubclient.CreateLocalClientConnection(provisioning.ServiceName, f.Run)
	if err == nil {
		pc = capnpclient.NewProvisioningCapnpClient(ctx, conn)
	}
	if err != nil {
		return err
	}
	manage, _ := pc.CapManageProvisioning(ctx, "hubcli")

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
//
//	deviceID is the ID of the device to approve
func HandleApproveRequest(ctx context.Context, f svcconfig.AppFolders, deviceID string) error {
	var pc provisioning.IProvisioning

	conn, err := hubclient.CreateLocalClientConnection(provisioning.ServiceName, f.Run)
	if err == nil {
		pc = capnpclient.NewProvisioningCapnpClient(ctx, conn)
		manage, _ := pc.CapManageProvisioning(ctx, "hubcli")
		err = manage.ApproveRequest(ctx, deviceID)
	}

	return err
}

func HandleGetApprovedRequests(ctx context.Context, f svcconfig.AppFolders) error {
	var pc provisioning.IProvisioning
	var provStatus []provisioning.ProvisionStatus

	conn, err := hubclient.CreateLocalClientConnection(provisioning.ServiceName, f.Run)
	if err == nil {
		pc = capnpclient.NewProvisioningCapnpClient(ctx, conn)
		manage, _ := pc.CapManageProvisioning(ctx, "hubcli")
		provStatus, err = manage.GetApprovedRequests(ctx)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Client ID              Request Time      Assigned\n")
	fmt.Printf("--------------------   ------------      --------\n")
	for _, provStatus := range provStatus {
		// a certificate is assigned when generated
		assigned := provStatus.ClientCertPEM != ""
		fmt.Printf("%20s  %s, %v\n",
			provStatus.DeviceID, provStatus.RequestTime, assigned)
	}

	return err
}

func HandleGetPendingRequests(ctx context.Context, f svcconfig.AppFolders) error {
	var pc provisioning.IProvisioning
	var provStatus []provisioning.ProvisionStatus

	conn, err := hubclient.CreateLocalClientConnection(provisioning.ServiceName, f.Run)
	if err == nil {
		pc = capnpclient.NewProvisioningCapnpClient(ctx, conn)
		manage, _ := pc.CapManageProvisioning(ctx, "hubcli")
		provStatus, err = manage.GetPendingRequests(ctx)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Client ID              Request Time\n")
	fmt.Printf("--------------------   ------------\n")
	for _, provStatus := range provStatus {
		// a certificate is assigned when generated
		fmt.Printf("%20s  %s\n",
			provStatus.DeviceID, provStatus.RequestTime)
	}

	return err
}
