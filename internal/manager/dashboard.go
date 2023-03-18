package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jibbolo/cwdash/internal/widget"
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

func (d *DashboardManager) DashboardList() (list []string) {
	return d.dashboardList
}

func (d *DashboardManager) GetDashboard(name string) *Dashboard {
	return d.dashboards[name]
}

func (d *DashboardManager) InitAWS() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w ", err)
	}
	d.cw = cloudwatch.NewFromConfig(cfg)
	return nil
}

func (d *DashboardManager) RefreshBody(dashboardName string) error {
	println("refreshing body for " + dashboardName)

	resp, err := d.cw.GetDashboard(context.Background(), &cloudwatch.GetDashboardInput{
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

func (d *DashboardManager) RefreshDashboardList() error {
	resp, err := d.cw.ListDashboards(context.Background(), &cloudwatch.ListDashboardsInput{})
	if err != nil {
		return fmt.Errorf("failed to describe dashboard: %v ", err)
	}
	d.dashboardList = make([]string, 0, len(resp.DashboardEntries))
	for _, e := range resp.DashboardEntries {
		d.dashboardList = append(d.dashboardList, *e.DashboardName)
	}
	return nil
}

func (d *DashboardManager) RenderGraph(w io.Writer, dashboardName string, widgetIndex int) (int, error) {
	if d.dashboards[dashboardName] == nil {
		if err := d.RefreshBody(dashboardName); err != nil {
			return 0, fmt.Errorf("renderWidget can't refresh body: %v ", err)
		}
	}
	widget := d.dashboards[dashboardName].Widgets[widgetIndex]
	body, err := d.RenderWidget(widget)
	if err != nil {
		return 0, err
	}
	return w.Write(body)
}

type Dashboard struct {
	Widgets []widget.Widget `json:"widgets,omitempty"`
}
