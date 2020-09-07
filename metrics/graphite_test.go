package metrics_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jfk9w-go/flu/metrics"
	"github.com/jfk9w-go/flu/testutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGraphiteClient_Counter(t *testing.T) {
	mock, err := testutil.RunMockServer("tcp")
	assert.Nil(t, err)
	defer mock.Close()

	client := metrics.NewGraphiteClient(mock.Address, 0)
	defer client.Close()

	client.Counter("counter", nil).Add(1)
	err = client.Flush(time.Unix(500, 0))
	assert.Nil(t, err)

	assert.Equal(t, "counter 1.000000000 500\n", <-mock.In)
	err = client.Flush(time.Unix(550, 0))
	assert.Nil(t, err)

	client.Counter("counter", nil).Add(5)
	client.Counter("counter", nil).Add(1)
	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	assert.Equal(t, "counter 6.000000000 600\n", <-mock.In)

	client.Counter("counter", metrics.Labels{"label_a", "A", "label_b", "B"}).Add(5)
	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	assert.Equal(t, "A.B.counter 5.000000000 600\n", <-mock.In)

	client.WithPrefix("counter").
		Counter("hits", metrics.Labels{"label_a", "A", "label_b", "B"}).Add(5)
	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	assert.Equal(t, "counter.A.B.hits 5.000000000 600\n", <-mock.In)
}

func TestGraphiteClient_Gauge(t *testing.T) {
	mock, err := testutil.RunMockServer("tcp")
	assert.Nil(t, err)
	defer mock.Close()

	client := metrics.NewGraphiteClient(mock.Address, 0)
	defer client.Close()

	client.Gauge("gauge", nil).Set(1)
	err = client.Flush(time.Unix(500, 0))
	assert.Nil(t, err)

	assert.Equal(t, "gauge 1.000000000 500\n", <-mock.In)

	err = client.Flush(time.Unix(550, 0))
	assert.Nil(t, err)

	assert.Equal(t, "gauge 1.000000000 550\n", <-mock.In)

	client.Gauge("gauge", nil).Set(5)
	client.Gauge("gauge", nil).Set(1)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "gauge 1.000000000 600\n", <-mock.In)

	client.Gauge("gauge", metrics.Labels{"label_a", "A", "label_b", "B"}).Set(5)
	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	expected := []string{
		"",
		"A.B.gauge 5.000000000 600",
		"gauge 1.000000000 600",
	}

	actual := strings.Split(<-mock.In, "\n")
	sort.Strings(actual)
	assert.Equal(t, expected, actual)

	client.WithPrefix("gauge").
		Gauge("hits", metrics.Labels{"label_a", "A", "label_b", "B"}).Set(5)
	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	expected = []string{
		"",
		"A.B.gauge 5.000000000 600",
		"gauge 1.000000000 600",
		"gauge.A.B.hits 5.000000000 600",
	}

	actual = strings.Split(<-mock.In, "\n")
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}

func TestGraphiteClient_Histogram(t *testing.T) {
	mock, err := testutil.RunMockServer("tcp")
	assert.Nil(t, err)
	defer mock.Close()

	client := metrics.NewGraphiteClient(mock.Address, 0)
	defer client.Close()
	client.HistogramBucketFormat = "%.2f"

	histogram := client.Histogram("histogram", nil, []float64{0, 0.25, 0.5, 0.75, 1})
	for _, value := range []float64{-0.1, 0, 0.1, 0.3, 0.4, 0.45, 0.5, 0.55, 0.75, 1, 1.1} {
		histogram.Observe(value)
	}

	err = client.Flush(time.Unix(600, 0))
	assert.Nil(t, err)

	actual := strings.Split(<-mock.In, "\n")
	expected := []string{
		"histogram.0_00 1.000000000 600",
		"histogram.0_25 2.000000000 600",
		"histogram.0_50 3.000000000 600",
		"histogram.0_75 2.000000000 600",
		"histogram.1_00 1.000000000 600",
		"histogram.inf 2.000000000 600",
		"",
	}

	assert.Equal(t, actual, expected)
}
