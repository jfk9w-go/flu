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
	if err != nil {
		t.Fatal(errors.Wrap(err, "create graphite mock"))
	}
	defer mock.Close()

	client := metrics.NewGraphiteClient(mock.Address, 0)
	defer client.Close()

	client.Counter("counter", nil).Add(1)
	if err := client.Flush(time.Unix(500, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "counter 1.000000000 500\n", <-mock.In)

	if err := client.Flush(time.Unix(550, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	client.Counter("counter", nil).Add(5)
	client.Counter("counter", nil).Add(1)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "counter 6.000000000 600\n", <-mock.In)

	client.Counter("counter", metrics.Labels{"label_a", "A", "label_b", "B"}).Add(5)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "A.B.counter 5.000000000 600\n", <-mock.In)

	client.WithPrefix("counter").
		Counter("hits", metrics.Labels{"label_a", "A", "label_b", "B"}).Add(5)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "counter.A.B.hits 5.000000000 600\n", <-mock.In)
}

func TestGraphiteClient_Gauge(t *testing.T) {
	mock, err := testutil.RunMockServer("tcp")
	if err != nil {
		t.Fatal(errors.Wrap(err, "create graphite mock"))
	}
	defer mock.Close()

	client := metrics.NewGraphiteClient(mock.Address, 0)
	defer client.Close()

	client.Gauge("gauge", nil).Set(1)
	if err := client.Flush(time.Unix(500, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "gauge 1.000000000 500\n", <-mock.In)

	if err := client.Flush(time.Unix(550, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "gauge 1.000000000 550\n", <-mock.In)

	client.Gauge("gauge", nil).Set(5)
	client.Gauge("gauge", nil).Set(1)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

	assert.Equal(t, "gauge 1.000000000 600\n", <-mock.In)

	client.Gauge("gauge", metrics.Labels{"label_a", "A", "label_b", "B"}).Set(5)
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

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
	if err := client.Flush(time.Unix(600, 0)); err != nil {
		t.Fatal(errors.Wrap(err, "flush metrics"))
	}

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
