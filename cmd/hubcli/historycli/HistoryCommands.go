package historycli

import (
	"context"
	"fmt"
	"sort"

	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/hiveot/hub/lib/hubclient"
	"github.com/hiveot/hub/lib/svcconfig"
	"github.com/hiveot/hub/pkg/history"
	"github.com/hiveot/hub/pkg/history/capnpclient"
)

func HistoryInfoCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name: "histinfo",
		//Aliases:   []string{"hin"},
		Usage:     "Show history store info",
		Category:  "history",
		ArgsUsage: "(no args)",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 0 {
				return fmt.Errorf("no arguments expected")
			}
			err := HandleHistoryInfo(ctx, f)
			return err
		},
	}
}

func HistoryListCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "histevents <pubID> <thingID>",
		Aliases:   []string{"hev", "lev"},
		Usage:     "List historical events",
		UsageText: "List the history of events from a Thing by its publisher and Thing ID",
		Category:  "history",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 2 {
				return fmt.Errorf("publisherID and thingID expected")
			}
			err := HandleListEvents(ctx, f, cCtx.Args().First(), cCtx.Args().Get(1), 30)
			return err
		},
	}
}

func HistoryLatestCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "histlatest <pubID> <thingID>",
		Usage:     "List latest values of a thing",
		UsageText: "List the latest value of each property/event of a thing by its publisher/thing ID",
		Aliases:   []string{"hla"},
		Category:  "history",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 2 {
				return fmt.Errorf("publisherID and thingID expected")
			}
			err := HandleListLatestEvents(ctx, f, cCtx.Args().First(), cCtx.Args().Get(1))
			return err
		},
	}
}
func HistoryRetainCommand(ctx context.Context, f svcconfig.AppFolders) *cli.Command {
	return &cli.Command{
		Name:      "histretained",
		Aliases:   []string{"hrt"},
		Usage:     "List retained events",
		UsageText: "List the events that are retained in the history store",
		Category:  "history",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 0 {
				return fmt.Errorf("no arguments expected")
			}
			err := HandleListRetainedEvents(ctx, f)
			return err
		},
	}
}

func HandleHistoryInfo(ctx context.Context, f svcconfig.AppFolders) error {
	var hist history.IHistoryService
	var rd history.IReadHistory

	conn, err := hubclient.ConnectToService(history.ServiceName, f.Run)
	if err == nil {
		hist = capnpclient.NewHistoryCapnpClient(ctx, conn)
		rd, err = hist.CapReadHistory(ctx, "hubcli", "", "")
	}
	if err != nil {
		return err
	}
	info := rd.Info(ctx)

	fmt.Println(fmt.Sprintf("ID:          %s", info.Id))
	fmt.Println(fmt.Sprintf("Size:        %d", info.DataSize))
	fmt.Println(fmt.Sprintf("Nr Records   %d", info.NrRecords))
	fmt.Println(fmt.Sprintf("Engine       %s", info.Engine))

	rd.Release()
	return conn.Close()
}

// HandleListEvents lists the history content
func HandleListEvents(ctx context.Context, f svcconfig.AppFolders, publisherID, thingID string, limit int) error {
	var hist history.IHistoryService
	var rd history.IReadHistory

	conn, err := hubclient.ConnectToService(history.ServiceName, f.Run)
	if err == nil {
		hist = capnpclient.NewHistoryCapnpClient(ctx, conn)
		rd, err = hist.CapReadHistory(ctx, "hubcli", publisherID, thingID)
	}
	if err != nil {
		return err
	}
	eventName := ""
	cursor := rd.GetEventHistory(ctx, eventName)
	fmt.Println("PublisherID    ThingID            Timestamp                    Event           Value (truncated)")
	fmt.Println("-----------    -------            ---------                    -----           ---------------- ")
	count := 0
	for tv, valid := cursor.Last(); valid && count < limit; tv, valid = cursor.Prev() {
		count++
		utime, err := dateparse.ParseAny(tv.Created)

		if err != nil {
			logrus.Infof("Parsing time failed '%s': %s", tv.Created, err)
		}

		fmt.Printf("%-14s %-18s %-28s %-15s %-30s\n",
			tv.PublisherID,
			tv.ThingID,
			utime.Format("02 Jan 2006 15:04:05 MST"),
			tv.ID,
			tv.ValueJSON,
		)
	}
	rd.Release()
	err = conn.Close()
	return err
}

// HandleListRetainedEvents lists the events that are retained
func HandleListRetainedEvents(ctx context.Context, f svcconfig.AppFolders) error {

	var hist history.IHistoryService
	var mngRet history.IManageRetention

	conn, err := hubclient.ConnectToService(history.ServiceName, f.Run)
	if err == nil {
		hist = capnpclient.NewHistoryCapnpClient(ctx, conn)
		mngRet, err = hist.CapManageRetention(ctx, "hubcli")
	}
	if err != nil {
		return err
	}
	evList, _ := mngRet.GetEvents(ctx)
	sort.Slice(evList, func(i, j int) bool {
		return evList[i].Name < evList[j].Name
	})

	fmt.Printf("Events (%2d)      days     publishers                     Things                         Excluded\n", len(evList))
	fmt.Println("----------       ----     ----------                     ------                         -------- ")
	for _, evRet := range evList {

		fmt.Printf("%-16.16s %-8d %-30.30s %-30.30s %-30.30s\n",
			evRet.Name,
			evRet.RetentionDays,
			fmt.Sprintf("%s", evRet.Publishers),
			fmt.Sprintf("%s", evRet.Things),
			fmt.Sprintf("%s", evRet.Exclude),
		)
	}
	mngRet.Release()
	err = conn.Close()
	return err
}

func HandleListLatestEvents(
	ctx context.Context, f svcconfig.AppFolders, publisherID, thingID string) error {
	var hist history.IHistoryService
	var readHist history.IReadHistory

	conn, err := hubclient.ConnectToService(history.ServiceName, f.Run)
	if err == nil {
		hist = capnpclient.NewHistoryCapnpClient(ctx, conn)
		readHist, err = hist.CapReadHistory(ctx, "hubcli", publisherID, thingID)
	}
	if err != nil {
		return err
	}
	props := readHist.GetProperties(ctx, nil)

	fmt.Println("Event ID         Publisher       Thing                Created                     Value")
	fmt.Println("----------         ---------       -----                -------                     -----")
	for _, prop := range props {
		utime, _ := dateparse.ParseAny(prop.Created)

		fmt.Printf("%-18.18s %-15.15s %-20s %-27s %s\n",
			prop.ID,
			prop.PublisherID,
			prop.ThingID,
			//utime.Format("02 Jan 2006 15:04:05 -0700"),
			utime.Format("02 Jan 2006 15:04:05 MST"),
			prop.ValueJSON,
		)
	}
	readHist.Release()
	conn.Close()
	return nil
}
