package metrics_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jfk9w-go/flu"
	fluhttp "github.com/jfk9w-go/flu/http"
	"github.com/jfk9w-go/flu/metrics"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusClient_Counter(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(errors.Wrap(err, "find free port"))
	}

	address := fmt.Sprintf("http://localhost:%d/metrics", port)
	listener := metrics.NewPrometheusListener(address)
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_ = listener.Close(ctx)
	}()

	listener.Counter("counter", nil).Add(1)
	client := fluhttp.NewClient(nil)
	text := new(flu.PlainText)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP counter \n"+
		"# TYPE counter counter\n"+
		"counter 1\n", text.Value)

	listener.Counter("counter", nil).Add(2)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP counter \n"+
		"# TYPE counter counter\n"+
		"counter 3\n", text.Value)

	listener.WithPrefix("labels").Counter("counter", metrics.Labels{"A", "1", "B", "2"}).Add(1)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP counter \n"+
		"# TYPE counter counter\n"+
		"counter 3\n"+
		"# HELP labels_counter \n"+
		"# TYPE labels_counter counter\n"+
		"labels_counter{A=\"1\",B=\"2\"} 1\n", text.Value)
}

func TestPrometheusClient_Gauge(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(errors.Wrap(err, "find free port"))
	}

	address := fmt.Sprintf("http://localhost:%d/metrics", port)
	listener := metrics.NewPrometheusListener(address)
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_ = listener.Close(ctx)
	}()

	listener.Gauge("gauge", nil).Set(1)
	client := fluhttp.NewClient(nil)
	text := new(flu.PlainText)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP gauge \n"+
		"# TYPE gauge gauge\n"+
		"gauge 1\n", text.Value)

	listener.Gauge("gauge", nil).Set(2)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP gauge \n"+
		"# TYPE gauge gauge\n"+
		"gauge 2\n", text.Value)

	listener.WithPrefix("labels").Gauge("gauge", metrics.Labels{"A", "1", "B", "2"}).Add(1)
	err = client.GET(address).
		Execute().
		CheckStatus(http.StatusOK).
		DecodeBody(text).
		Error

	assert.Nil(t, err)
	assert.Equal(t, ""+
		"# HELP gauge \n"+
		"# TYPE gauge gauge\n"+
		"gauge 2\n"+
		"# HELP labels_gauge \n"+
		"# TYPE labels_gauge gauge\n"+
		"labels_gauge{A=\"1\",B=\"2\"} 1\n", text.Value)
}
