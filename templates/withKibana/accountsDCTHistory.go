package withKibana

// AccountsDCTHistory will hold the configuration for the accountsdcthistory index
var AccountsDCTHistory = Object{
	"index_patterns": Array{
		"accountsdcthistory-*",
	},
	"settings": Object{
		"number_of_shards":   5,
		"number_of_replicas": 0,
		"opendistro.index_state_management.rollover_alias": "accountsdcthistory",
	},
	"mappings": Object{
		"properties": Object{
			"address": Object{
				"type": "keyword",
			},
			"balance": Object{
				"type": "keyword",
			},
			"identifier": Object{
				"type": "text",
			},
			"isSender": Object{
				"type": "boolean",
			},
			"isSmartContract": Object{
				"type": "boolean",
			},
			"shardID": Object{
				"type": "long",
			},
			"timestamp": Object{
				"type":   "date",
				"format": "epoch_second",
			},
			"token": Object{
				"type": "text",
			},
			"tokenNonce": Object{
				"type": "double",
			},
		},
	},
}
