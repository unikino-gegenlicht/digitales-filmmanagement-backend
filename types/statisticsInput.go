package types

import "time"

type StatisticsRequestInput struct {
	From  *time.Time `in:"query=from"`
	Until *time.Time `in:"query=until"`
}
