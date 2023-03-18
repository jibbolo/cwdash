package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jibbolo/cwdash/internal/manager/widget"
)

// Dashboard is the main application struct
type DashboardManager struct {
	cw            *cloudwatch.Client
	dashboards    map[string]*Dashboard
	dashboardList []string
}

func New() *DashboardManager {
	return &DashboardManager{
		dashboards:    make(map[string]*Dashboard),
		dashboardList: make([]string, 0),
	}
}

func (d *DashboardManager) DashboardList(ctx context.Context) (list []string) {
	return d.dashboardList
}

func (d *DashboardManager) GetDashboard(name string) *Dashboard {
	return d.dashboards[name]
}

func (d *DashboardManager) InitAWS(ctx context.Context) error {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w ", err)
	}
	d.cw = cloudwatch.NewFromConfig(cfg)
	return nil
}

func (d *DashboardManager) RefreshBody(ctx context.Context, dashboardName string) error {
	resp, err := d.cw.GetDashboard(ctx, &cloudwatch.GetDashboardInput{
		DashboardName: aws.String(dashboardName),
	})
	if err != nil {
		return fmt.Errorf("failed to describe dashboard: %v ", err)
	}
	var body Dashboard
	err = json.Unmarshal([]byte(*resp.DashboardBody), &body)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}
	d.dashboards[dashboardName] = &body
	return nil
}

func (d *DashboardManager) RefreshDashboardList(ctx context.Context) error {
	resp, err := d.cw.ListDashboards(ctx, &cloudwatch.ListDashboardsInput{})
	if err != nil {
		return fmt.Errorf("failed to describe dashboard: %v ", err)
	}
	d.dashboardList = make([]string, 0, len(resp.DashboardEntries))
	for _, e := range resp.DashboardEntries {
		d.dashboardList = append(d.dashboardList, *e.DashboardName)
	}
	return nil
}

func (d *DashboardManager) RenderGraph(ctx context.Context, w io.Writer, dashboardName string, widgetIndex int) (int, error) {
	if d.dashboards[dashboardName] == nil {
		if err := d.RefreshBody(ctx, dashboardName); err != nil {
			return 0, fmt.Errorf("renderWidget can't refresh body: %v ", err)
		}
	}
	widget := d.dashboards[dashboardName].Widgets[widgetIndex]
	body, err := d.RenderWidget(ctx, widget)
	if err != nil {
		return 0, err
	}
	return w.Write(body)
}

type Dashboard struct {
	Widgets []widget.Widget `json:"widgets,omitempty"`
}
