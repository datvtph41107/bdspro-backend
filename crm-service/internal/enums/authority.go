package enums

type EAuth int16

const (
	AuthCustomerCreate      EAuth = 100
	AuthCustomerUpdate      EAuth = 101
	AuthCustomerDelete      EAuth = 102
	AuthCustomerView        EAuth = 103
	AuthCustomerNote        EAuth = 104
	AuthCustomerAssign      EAuth = 105
	AuthCustomerSwitchStage EAuth = 106
	AuthContactCreate       EAuth = 200
	AuthContactUpdate       EAuth = 201
	AuthContactDelete       EAuth = 202
	AuthContactView         EAuth = 203
	AuthPipelineCreate      EAuth = 300
	AuthPipelineUpdate      EAuth = 301
	AuthPipelineDelete      EAuth = 302
	AuthPipelineView        EAuth = 303
	AuthStageCreate         EAuth = 400
	AuthStageUpdate         EAuth = 401
	AuthStageDelete         EAuth = 402
	AuthStageView           EAuth = 403
	AuthRuleCreate          EAuth = 500
	AuthRuleUpdate          EAuth = 501
	AuthRuleDelete          EAuth = 502
	AuthRuleView            EAuth = 503
)