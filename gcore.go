package main

import (
	"context"
	_ "embed"
	"strings"
	"time"

	"github.com/G-Core/gcore-dns-sdk-go"
	"github.com/kokizzu/gotro/L"
)

//go:embed .token
var apiToken string

func main() {
	apiToken = strings.TrimSpace(apiToken)

	sdk := dnssdk.NewClient(dnssdk.PermanentAPIKeyAuth(apiToken), func(client *dnssdk.Client) {
		client.Debug = true
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const zoneName = `benalu2.dev`

	_, err := sdk.CreateZone(ctx, zoneName)
	if err != nil && !strings.Contains(err.Error(), `already exists`) {
		L.PanicIf(err, `sdk.CreateZone`)
	}

	zoneResp, err := sdk.Zone(ctx, zoneName)
	L.PanicIf(err, `sdk.Zone`)
	L.Describe(zoneResp)

	err = sdk.AddZoneRRSet(ctx,
		zoneName,        // zone
		`www.`+zoneName, // name
		`A`,             // rrtype
		[]dnssdk.ResourceRecord{ // https://apidocs.gcore.com/dns#tag/rrsets/operation/CreateRRSet
			{
				Content: []any{
					`194.233.65.174`,
				},
			},
		},
		120, // TTL
	)
	L.PanicIf(err, `AddZoneRRSet`)

	rr, err := sdk.RRSet(ctx, zoneName, `www.`+zoneName, `A`)
	L.PanicIf(err, `sdk.RRSet`)
	L.Describe(rr)

	err = sdk.DeleteRRSet(ctx, zoneName, `www.`+zoneName, `A`)
	L.PanicIf(err, `sdk.DeleteRRSet`)

	err = sdk.DeleteZone(ctx, zoneName)
	L.PanicIf(err, `sdk.DeleteZone`)
}
