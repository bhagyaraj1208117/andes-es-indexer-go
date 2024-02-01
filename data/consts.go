package data

const (
	AtSeparator                = "@"
	GasRefundForRelayerMessage = "gas refund for relayer"

	// MaxDCTValueLength defines the maximum length for an DCT value that can be parsed
	MaxDCTValueLength = 100
	// MaxFieldLength defines the maximum length for a keyword field, approximating the maximum length of the keyword type.
	MaxFieldLength = 30000

	// MaxKeywordFieldLengthBeforeBase64Encoding defines the maximum length for a keyword field that will be base64 encoded
	MaxKeywordFieldLengthBeforeBase64Encoding = 22500
)
