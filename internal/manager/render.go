package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jibbolo/cwdash/internal/widget"
)

var additionalFields = map[string]interface{}{
	"width":    600,
	"height":   300,
	"start":    "-PT336H",
	"end":      "P0D",
	"timezone": "+0200",
}

func (d *DashboardManager) RenderWidget(widget widget.Widget) ([]byte, error) {
	props := widget.Properties
	for k, v := range additionalFields {
		props[k] = v
	}
	propsBytes, err := json.Marshal(props)
	if err != nil {
		return []byte{}, fmt.Errorf("can't marshal props: %v ", err)
	}
	resp, err := d.cw.GetMetricWidgetImage(context.Background(), &cloudwatch.GetMetricWidgetImageInput{
		MetricWidget: aws.String(string(propsBytes)),
	})
	if err != nil {
		return []byte{}, fmt.Errorf("cw failed to render widget: %v ", err)
	}
	return resp.MetricWidgetImage, nil
}
