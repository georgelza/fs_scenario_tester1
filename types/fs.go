package types

// Structs - the values we bring in from *app.json configuration file
type Tp_general struct {
	EchoConfig          int
	Hostname            string
	Loglevel            string
	Debuglevel          int
	Echojson            int
	Testsize            int    // Used to limit number of records posted, over rided when reading test cases from input_source,
	Sleep               int    // sleep time between API post
	Httpposturl         string // FeatureSpace API URL
	Cert_dir            string
	Cert_file           string
	Cert_key            string
	Call_fs_api         int     // Do we call the API endpoint
	Eventtype           string  // paymentRT or paymentNRT
	Json_to_file        int     // Do we output JSON to file in output_path
	Output_path         string  // output location
	Json_from_file      int     // Do we read JSON from input_path directory and post to FS API endpoint
	Input_path          string  // Where are my scenario JSON files located
	MinTransactionValue float64 // Min value if the fake transaction
	MaxTransactionValue float64 // Max value of the fake transaction
	SeedFile            string  // Which seed file to read in
	EchoSeed            int     // 0/1 Echo the seed data to terminal
}

// FS engineResponse components
type TPamount struct {
	BaseCurrency string  `json:"baseCurrency,omitempty"`
	BaseValue    float64 `json:"baseValue,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	Value        float64 `json:"value,omitempty"`
}

type TPaymentNRT = struct {
	AccountAgentId                 string   `json:"accountAgentId,omitempty"`
	AccountEntityId                string   `json:"accountEntityId,omitempty"`
	AccountAgentName               string   `json:"accountAgentName,omitempty"`
	AccountId                      string   `json:"accountId,omitempty"`
	Amount                         TPamount `json:"amount,omitempty"`
	ChargeBearer                   string   `json:"chargeBearer,omitempty"` // hardcode SLEV
	CounterpartyAgentId            string   `json:"counterpartyAgentId,omitempty"`
	CounterpartyAgentName          string   `json:"counterpartyAgentName,omitempty"`
	CounterpartyEntityId           string   `json:"counterpartyEntityId,omitempty"`
	CounterpartyId                 string   `json:"counterpartyId,omitempty"`
	CreationDate                   string   `json:"creationDate,omitempty"`
	DestinationCountry             string   `json:"destinationCountry,omitempty"` // hardcode ZAF
	Direction                      string   `json:"direction,omitempty"`
	EventId                        string   `json:"eventId,omitempty"`
	EventTime                      string   `json:"eventTime,omitempty"`
	EventType                      string   `json:"eventType,omitempty"`
	FromFIBranchId                 string   `json:"fromFIBranchId,omitempty"`
	FromId                         string   `json:"fromId,omitempty"`
	LocalInstrument                string   `json:"localInstrument,omitempty"`
	MsgStatus                      string   `json:"msgStatus,omitempty"` // hardcode Success
	MsgType                        string   `json:"msgType,omitempty"`   // hardcode RCCT
	NumberOfTransactions           int      `json:"numberOfTransactions,omitempty"`
	PaymentClearingSystemReference string   `json:"paymentClearingSystemReference,omitempty"`
	PaymentMethod                  string   `json:"paymentMethod,omitempty"` // hardcode TRF
	PaymentReference               string   `json:"paymentReference,omitempty"`
	RemittanceId                   string   `json:"remittanceId,omitempty"`
	RequestExecutionDate           string   `json:"requestExecutionDate,omitempty"`
	SchemaVersion                  int      `json:"schemaVersion,omitempty"`
	SettlementClearingSystemCode   string   `json:"settlementClearingSystemCode,omitempty"` // hardcode RTC
	SettlementDate                 string   `json:"settlementDate,omitempty"`
	SettlementMethod               string   `json:"settlementMethod,omitempty"` // hardcode CLRG
	TenantId                       string   `json:"tenantId,omitempty"`
	ToFIBranchId                   string   `json:"toFIBranchId,omitempty"`
	ToId                           string   `json:"toId,omitempty"`
	TotalAmount                    TPamount `json:"totalAmount,omitempty"`
	TransactionId                  string   `json:"transactionId,omitempty"`
	TransactionType                string   `json:"transactionType,omitempty"`
	// FS Modifications required for these fields
	CounterpartyIDaccounttype  string `json:"counterpartyIDaccounttype,omitempty"`
	AccountIDaccounttype       string `json:"accountIDaccounttype,omitempty"`
	Usercode                   string `json:"usercode,omitempty"` // Hardcode RTC0000
	CounterpartyIDSuspenseflag string `json:"counterpartyIDSuspenseflag,omitempty"`
	AccountIDSuspenseflag      string `json:"accountIDSuspenseflag,omitempty"`
	EntryClass                 string `json:"entry,omitempty"` // hardcode RTC42
	UnpaidReasonCode           string `json:"unpaidReasonCode,omitempty"`
}

type TPAddress struct {
}
type TVerificationType struct {
}
type TPaymentRT struct {
	AccountAddress                      TPAddress         `json:"accountAddress,omitempty"`
	AccountAgentAddress                 TPAddress         `json:"accountAgentAddress,omitempty"`
	AccountAgentId                      string            `json:"accountAgentId,omitempty"`
	AccountAgentName                    string            `json:"accountAgent,omitempty"`
	AccountBalanceAfter                 string            `json:"accountBalanceAfter,omitempty"`
	AccountEntityId                     string            `json:"accountEntity"`
	AccountId                           string            `json:"accountId,omitempty"`
	AccountName                         string            `json:"accountName,omitempty"`
	Amount                              TPamount          `json:"amount,omitempty"`
	CardEntityId                        string            `json:"cardEntityId,omitempty"`
	CardId                              string            `json:"cardId,omitempty"`
	Channel                             string            `json:"channel,omitempty"`
	ChargeBearer                        string            `json:"chargeBearer,omitempty"`
	CounterpartyAddress                 TPAddress         `json:"CounterpartyAddress,omitempty"`
	CounterpartyAgentAddress            TPAddress         `json:"counterpartyAgentAddress,omitempty"`
	CounterpartyAgentId                 string            `json:"counterpartyAgentId,omitempty"`
	CounterpartyAgentName               string            `json:"counterPartyAgentName"`
	CounterpartyEntityId                string            `json:"counterPartyEntityId,omitempty"`
	CounterpartyId                      string            `json:"counterPartyId,omitempty"`
	CounterpartyName                    string            `json:"counterPartyName,omitempty"`
	CreationDate                        string            `json:"creationDate"`
	CustomerEntityId                    string            `json:"customerEntityId,omitempty"`
	CustomerId                          string            `json:"customerId,omitempty"`
	DecorationId                        string            `json:"decorationId,omitempty"`
	DestinationCountry                  string            `json:"destinationCountry,omitempty"`
	Device                              string            `json:"device,omitempty"`
	DeviceEntityId                      string            `json:"deviceEntityId,omitempty"`
	DeviceId                            string            `json:"deviceId,omitempty"`
	Direction                           string            `json:"direction,omitempty"`
	EventId                             string            `json:"eventId,omitempty"`
	EventTime                           string            `json:"eventTime,omitempty"`
	EventType                           string            `json:"eventType,omitempty"`
	FinalPaymentDate                    string            `json:"finalPaymentDate,omitempty"`
	FromFIBranchId                      string            `json:"fromFIBranchId,omitempty"`
	FromId                              string            `json:"fromId,omitempty"`
	InstructedAgentAddress              TPAddress         `json:"instructedAgentAddress,omitempty"`
	InstructedAgentId                   string            `json:"instructedAgentId,omitempty"`
	InstructedAgentName                 string            `json:"instructedAgentName,omitempty"`
	InstructingAgentAddress             TPAddress         `json:"instructingAgentAddress,omitempty"`
	InstructingAgentId                  string            `json:"instructingAgentId,omitempty"`
	InstructingAgentName                string            `json:"instructingAgentName,omitempty"`
	IntermediaryAgent1AccountId         string            `json:"intermediaryAgent1AccountId,omitempty"`
	IntermediaryAgent1Address           TPAddress         `json:"intermeduartAgent1Address,omitempty"`
	IntermediaryAgent1Id                string            `json:"intermediaryAgent1Id,omitempty"`
	IntermediaryAgent1Name              string            `json:"intermediaryAgent1Name,omitempty"`
	IntermediaryAgent2AccountId         string            `json:"intermediaryAgent2AccountId,omitempty"`
	IntermediaryAgent2Address           TPAddress         `json:"intermeduartAgent2Address,omitempty"`
	IntermediaryAgent2Id                string            `json:"intermediaryAgent2Id,omitempty"`
	IntermediaryAgent2Name              string            `json:"intermediaryAgent2Name,omitempty"`
	IntermediaryAgent3AccountId         string            `json:"intermediaryAgent3AccountId,omitempty"`
	IntermediaryAgent3Address           TPAddress         `json:"intermeduartAgent3Address,omitempty"`
	IntermediaryAgent3Id                string            `json:"intermediaryAgent3Id,omitempty"`
	IntermediaryAgent3Name              string            `json:"intermediaryAgent3Name,omitempty"`
	LocalInstrument                     string            `json:"localInstrument,omitempty"`
	MsgStatus                           string            `json:"msgStatus,omitempty"`
	MsgStatusReason                     string            `json:"msgStatusReason,omitempty"`
	MsgType                             string            `json:"msgType,omitempty"`
	NumberOfTransactions                int               `json:"numberOfTransactions,omitempty"`
	PaymentClearingSystemReference      string            `json:"paymentClearingSystemReference,omitempty"`
	PaymentFrequency                    string            `json:"paymentFrequency,omitempty"`
	PaymentMethod                       string            `json:"paymentMethod,omitempty"`
	PaymentReference                    string            `json:"paymentReference,omitempty"`
	RemittanceId                        string            `json:"payRemittanceIdmentMethod,omitempty"`
	RemittanceLocationElectronicAddress string            `json:"remittanceLocationElectronicAddress,omitempty"`
	RemittanceLocationMethod            string            `json:"remittanceLocationMethod,omitempty"`
	RequestExecutionDate                string            `json:"paymentMRequestExecutionDateethod,omitempty"`
	SchemaVersion                       int               `json:"schemaVersion,omitempty"`
	ServiceLevelCode                    string            `json:"serviceLevelCode,omitempty"`
	SettlementClearingSystemCode        string            `json:"settlementClearingSystemCode,omitempty"`
	SettlementDate                      string            `json:"settlementDate,omitempty"`
	SettlementMethod                    string            `json:"settlementMethod,omitempty"`
	TenantId                            string            `json:"tenantId,omitempty"`
	ToFIBranchId                        string            `json:"toFIBranchId,omitempty"`
	ToId                                string            `json:"toId,omitempty"`
	TotalAmount                         TPamount          `json:"totalAmount,omitempty"`
	TransactionId                       string            `json:"transactionId,omitempty"`
	TransactionType                     string            `json:"transactionType,omitempty"`
	UltimateAccountAddress              TPAddress         `json:"ultimateAccountAddress,omitempty"`
	UltimateAccountId                   string            `json:"ultimateAccountId,omitempty"`
	UltimateAccountName                 string            `json:"ultimateAccountName,omitempty"`
	UltimateCounterpartyAddress         TPAddress         `json:"ultimateCounterpartyAddress,omitempty"`
	UltimateCounterpartyId              string            `json:"ultimateCounterpartyId,omitempty"`
	UltimateCounterpartyName            string            `json:"ultimateCounterpartyName,omitempty"`
	UnstructuredRemittanceInformation   string            `json:"unstructuredRemittanceInformation,omitempty"`
	VerificationResult                  string            `json:"verificationResult,omitempty"`
	VerificationType                    TVerificationType `json:"verificationType,omitempty"`
}

type Tp_Tenant struct {
	Name     string `json:"name,omitempty"`
	TenantId string `json:"tenantid,omitempty"`
}

type Tp_Entity struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	EntityId      string `json:"entityid,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type Tp_Payer struct {
	Name          string `json:"name,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type Tp_Address struct {
	Street   string `json:"street,omitempty"`
	City     string `json:"city,omitempty"`
	Province string `json:"Province,omitempty"`
	Code     string `json:"code,omitempty"`
}

type Tp_Agent struct {
	Name      string     `json:"name,omitempty"`
	Id        string     `json:"id,omitempty"`
	Address   Tp_Address `json:"address,omitempty"`
	AccountId string     `json:"accountId,omitempty"`
}

type Tp_seed struct {
	Direction                []string    `json:"direction,omitempty"`
	TransactionTypes         []string    `json:"transactiontypes,omitempty"`
	ChargeBearers            []string    `json:"chargeBearers,omitempty"`
	RemittanceLocationMethod []string    `json:"remittanceLocationmethod,omitempty"`
	SettlementMethod         []string    `json:"settlementMethod,omitempty"`
	VerificationResult       []string    `json:"verificationResult,omitempty"`
	PaymentFrequency         []string    `json:"paymentFrequency,omitempty"`
	Agent                    []Tp_Agent  `json:"agent,omitempty"`
	Tenants                  []Tp_Tenant `json:"tenants,omitempty"`
	GoodEntities             []Tp_Entity `json:"goodentities,omitempty"`
	BadEntities              []Tp_Entity `json:"badentities,omitempty"`
	GoodPayers               []Tp_Payer  `json:"goodpayers,omitempty"`
	BadPayers                []Tp_Payer  `json:"badpayers,omitempty"`
}
