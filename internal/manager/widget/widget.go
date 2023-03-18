package widget

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jibbolo/cwdash/internal/aws"
	"github.com/russross/blackfriday/v2"
)

type Widget struct {
	Properties map[string]any `json:"properties,omitempty"`
}

func (w Widget) HasMarkdown() bool {
	_, ok := w.Properties["markdown"]
	return ok
}

func (w Widget) HasQuery() bool {
	_, ok := w.Properties["query"]
	return ok
}

func (w *Widget) RenderMarkdown() template.HTML {
	if markdown, ok := w.Properties["markdown"]; ok {
		return template.HTML(
			blackfriday.Run([]byte(markdown.(string))),
		)
	}
	return "Error: markdown unavailable"
}

var additionalFields = map[string]interface{}{
	"width":    600,
	"height":   300,
	"start":    "-PT336H",
	"end":      "P0D",
	"timezone": "+0200",
}

func (w Widget) Render(ctx context.Context) ([]byte, error) {
	props := w.Properties
	for k, v := range additionalFields {
		props[k] = v
	}
	propsBytes, err := json.Marshal(props)
	if err != nil {
		return []byte{}, fmt.Errorf("can't marshal props: %v ", err)
	}
	propsString := string(propsBytes)
	resp, err := aws.CloudWatch().GetMetricWidgetImage(ctx, &cloudwatch.GetMetricWidgetImageInput{
		MetricWidget: &propsString,
	})
	if err != nil {
		return []byte{}, fmt.Errorf("cw failed to render widget: %v ", err)
	}
	return resp.MetricWidgetImage, nil
}
