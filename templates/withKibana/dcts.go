package withKibana

// DCTs will hold the configuration for the dcts index
var DCTs = Object{
	"index_patterns": Array{
		"dcts-*",
	},
	"settings": Object{
		"number_of_shards":   3,
		"number_of_replicas": 0,
	},
	"mappings": Object{
		"properties": Object{
			"name": Object{
				"type": "keyword",
			},
			"ticker": Object{
				"type": "keyword",
			},
			"token": Object{
				"type": "text",
			},
			"issuer": Object{
				"type": "keyword",
			},
			"currentOwner": Object{
				"type": "keyword",
			},
			"numDecimals": Object{
				"type": "long",
			},
			"type": Object{
				"type": "keyword",
			},
			"timestamp": Object{
				"type":   "date",
				"format": "epoch_second",
			},
			"ownersHistory": Object{
				"type": "nested",
				"properties": Object{
					"timestamp": Object{
						"index":  "false",
						"type":   "date",
						"format": "epoch_second",
					},
					"address": Object{
						"type": "keyword",
					},
				},
			},
			"paused": Object{
				"type": "boolean",
			},
			"properties": Object{
				"properties": Object{
					"canMint": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canBurn": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canUpgrade": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canTransferNFTCreateRole": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canAddSpecialRoles": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canPause": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canFreeze": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canWipe": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canChangeOwner": Object{
						"index": "false",
						"type":  "boolean",
					},
					"canCreateMultiShard": Object{
						"index": "false",
						"type":  "boolean",
					},
				},
			},
			"roles": Object{
				"type": "nested",
				"properties": Object{
					"DCTRoleLocalBurn": Object{
						"type": "keyword",
					},
					"DCTRoleLocalMint": Object{
						"type": "keyword",
					},
					"DCTRoleNFTAddQuantity": Object{
						"type": "keyword",
					},
					"DCTRoleNFTAddURI": Object{
						"type": "keyword",
					},
					"DCTRoleNFTBurn": Object{
						"type": "keyword",
					},
					"DCTRoleNFTCreate": Object{
						"type": "keyword",
					},
					"DCTRoleNFTUpdateAttributes": Object{
						"type": "keyword",
					},
					"DCTTransferRole": Object{
						"type": "keyword",
					},
				},
			},
		},
	},
}
