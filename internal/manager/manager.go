package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jibbolo/cwdash/internal/aws"
	"github.com/jibbolo/cwdash/internal/manager/widget"
)

// Dashboard is the main application struct
type DashboardManager struct {
	mux        sync.RWMutex
	dashboards map[string]*Dashboard
}

func New() *DashboardManager {
	return &DashboardManager{
		dashboards: make(map[string]*Dashboard),
	}
}

func (d *DashboardManager) GetDashboard(name string) *Dashboard {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return d.dashboards[name]
}

func (d *DashboardManager) DashboardList() (list []string) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	for name, _ := range d.dashboards {
		list = append(list, name)
	}
	return
}

func (d *DashboardManager) LoadDashboard(ctx context.Context, dashboardName string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	resp, err := aws.CloudWatch().GetDashboard(ctx, &cloudwatch.GetDashboardInput{
		DashboardName: &dashboardName,
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
	resp, err := aws.CloudWatch().ListDashboards(ctx, &cloudwatch.ListDashboardsInput{})
	if err != nil {
		return fmt.Errorf("failed to describe dashboard: %v ", err)
	}
	for _, e := range resp.DashboardEntries {
		if err := d.LoadDashboard(ctx, *e.DashboardName); err != nil {
			log.Printf("can't load dashboard %s: %v", *e.DashboardName, err)
		}
	}
	return nil
}

func (d *DashboardManager) RenderGraph(ctx context.Context, w io.Writer, dashboardName string, widgetIndex int) (int, error) {
	dashboard := d.GetDashboard(dashboardName)
	if widgetIndex >= len(dashboard.Widgets) {
		return 0, fmt.Errorf("widget not found")
	}
	widget := dashboard.Widgets[widgetIndex]
	body, err := widget.Render(ctx)
	if err != nil {
		return 0, err
	}
	return w.Write(body)
}

type Dashboard struct {
	Name    string
	Widgets []widget.Widget `json:"widgets,omitempty"`
}
