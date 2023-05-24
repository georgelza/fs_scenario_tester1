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
	CreateNewAccount    string  // Do we want to use the same AccountNumber for each transaction of a client, or spin up a new AccountNumber every run - Not thought through yet....
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

type TPaymentRT struct {
	AccountAddress      TPAddress `json:"accountAddress,omitempty"`
	AccountAgentAddress TPAddress `json:"accountAgentAddress,omitempty"`
	AccountAgentId      string    `json:"accountAgentId,omitempty"`
	AccountAgentName    string    `json:"accountAgent,omitempty"`
}

type Tenant struct {
	Name     string `json:"name,omitempty"`
	TenantId string `json:"tenantid,omitempty"`
}

type GoodEntity struct {
	Name          string `json:"name,omitempty"`
	EntityId      string `json:"entityid,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type BadEntity struct {
	Name          string `json:"name,omitempty"`
	EntityId      string `json:"entityid,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type GoodPayer struct {
	Name          string `json:"name,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type BadPayer struct {
	Name          string `json:"name,omitempty"`
	TenantId      string `json:"tenantid,omitempty"`
	AccountNumber string `json:"accountnumber,omitempty"`
}

type Tp_seed struct {
	Tenants          []Tenant     `json:"tenants,omitempty"`
	GoodEntities     []GoodEntity `json:"goodentities,omitempty"`
	BadEntities      []BadEntity  `json:"badentities,omitempty"`
	GoodPayers       []GoodPayer  `json:"goodpayers,omitempty"`
	BadPayers        []BadPayer   `json:"badpayers,omitempty"`
	Direction        []string     `json:"direction,omitempty"`
	TransactionTypes []string     `json:"transactiontypes,omitempty"`
}
