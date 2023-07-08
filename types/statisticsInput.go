package types

type StatisticsRequestInput struct {
	From  *int64 `in:"query=from"`
	Until *int64 `in:"query=until"`
}
