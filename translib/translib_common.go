package translib

type PMEntityType int
const (
	CHASSIS             PMEntityType = iota
	LINECARD
	FAN
	PSU
	OCH
	TRANSCEIVER_C
	TRANSCEIVER_L
	LOGICAL_CHANNEL_OTU
	LOGICAL_CHANNEL_ODU
	LOGICAL_CHANNEL_ETH
	OPTICAL_DEVICE
	UNKNOWN
)

type PMTimePeriod int
const (
	MIN_15         PMTimePeriod = iota // 0
	HOUR_24                            // 1
	NOT_APPLICABLE                     // 2
)

const (
	PMCurrent          = "current"
	PMHistory15min     = "15_pm_history"
	PMHistory24h       = "24_pm_history"
	PMCurrent15min     = "15_pm_current"
	PMCurrent24h       = "24_pm_current"
)

const (
	CUSTOM_TIME_FORMAT = "2006-01-02T15:04:05Z-07:00"
)

type QueryMode int
const (
	All QueryMode = iota
	Restconf
	Telemetry
)
