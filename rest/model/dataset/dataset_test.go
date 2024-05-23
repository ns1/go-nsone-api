package dataset

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = time.Unix(time.Now().Unix(), 0)

var mockedDataset = &Dataset{
	ID:   "dt-id",
	Name: "dt-name",
	Datatype: &Datatype{
		Type:  DatatypeTypeNumQueries,
		Scope: DatatypeScopeAccount,
		Data:  map[string]string{"foo": "bar"},
	},
	Repeat: &Repeat{
		Start:        UnixTimestamp(now),
		RepeatsEvery: RepeatsEveryMonth,
		EndAfterN:    1,
	},
	Timeframe: &Timeframe{
		Aggregation: TimeframeAggregationBillingPeriod,
		Cycles:      getInt32Ptr(3),
		From:        getTimestampPtr(UnixTimestamp(now)),
		To:          getTimestampPtr(UnixTimestamp(now)),
	},
	ExportType: ExportTypeCSV,
	Reports: []*Report{
		{
			ID:        "dt-report-id",
			Status:    ReportStatusAvailable,
			Start:     UnixTimestamp(now),
			End:       UnixTimestamp(now),
			CreatedAt: UnixTimestamp(now),
		},
	},
	RecipientEmails: []string{"datasets@ns1.com"},
	CreatedAt:       UnixTimestamp(now),
	UpdatedAt:       UnixTimestamp(now),
}

func getTimestampPtr(v UnixTimestamp) *UnixTimestamp {
	return &v
}

func getInt32Ptr(v int32) *int32 {
	return &v
}

func TestNewDataset(t *testing.T) {
	dt := NewDataset(
		mockedDataset.ID,
		mockedDataset.Name,
		NewDatatype(
			mockedDataset.Datatype.Type,
			mockedDataset.Datatype.Scope,
			mockedDataset.Datatype.Data,
		),
		NewRepeat(
			mockedDataset.Repeat.Start,
			mockedDataset.Repeat.RepeatsEvery,
			mockedDataset.Repeat.EndAfterN,
		),
		NewTimeframe(
			mockedDataset.Timeframe.Aggregation,
			mockedDataset.Timeframe.Cycles,
			mockedDataset.Timeframe.From,
			mockedDataset.Timeframe.To,
		),
		mockedDataset.ExportType,
		[]*Report{
			NewReport(
				mockedDataset.Reports[0].ID,
				mockedDataset.Reports[0].Status,
				mockedDataset.Reports[0].Start,
				mockedDataset.Reports[0].End,
				mockedDataset.Reports[0].CreatedAt,
			),
		},
		mockedDataset.RecipientEmails,
		mockedDataset.CreatedAt,
		mockedDataset.UpdatedAt,
	)

	assert.Equal(t, mockedDataset, dt)
}

func TestDatasetUnmarshall(t *testing.T) {
	dt := &Dataset{}

	dtJSON, err := json.Marshal(mockedDataset)
	require.NoError(t, err)

	err = json.Unmarshal(dtJSON, dt)
	require.NoError(t, err)

	assert.Equal(t, mockedDataset, dt)
}
