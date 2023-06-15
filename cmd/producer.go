/*****************************************************************************
*
*	File			: producer.go
*
* 	Created			: 27 Aug 2021
*
*	Description		: Creates Fake JSON Payload, originally posted onto Kafka topic, modified to post onto FeatureSpace event API endpoint
*
*	Modified		: 27 Aug 2021	- Start
*					: 24 Jan 2023   - Mod for applab sandbox
*					: 20 Feb 2023	- repackaged for TFM 2.0 load generation, we're creating fake FS Payment events messages onto
*					:				- FeatureSpace API endpoint
*
*					: 09 May 2023	- Enhancement: we can now read a directory of files (JSON structured), each as a paymentEvent and
*					:				- post the contents onto the API endpoint, allowing for pre designed schenarios to be used/tested
*					:				- NOTE, Original idea/usage was posting payloads onto Kafka topic, thus the fake data gen,
*					:				- With the new usage of reading scenario files allot of the bits below can be removed...
*					:				- Also removed the Prometheus instrumentation form this version as it will mostly be used to input/post a
*					: 				- coupld of files, not a big batch that needs to be timed/measure.
*
*					: 12 May 2023 	- Moved all environment variables from .exps environment export file to *_app.json file, this works better with a App
*									- destined for a desktop vs a app for a K8S server which prefers environment vars.
*					:				- https://onexlab-io.medium.com/golang-config-file-best-practise-d27d6a97a65a
*
*					: 24 May 2023	- Moved the seed data to a json structure seed.json thats ready in and then utilised instead of the seed package
*					:				- Modifying the payment structure to be aligned with Kiveshan's excell spread sheet, Makes for better down the line
*					:				- Fake data generation.
*					:				- also introduced the min and max transaction values.
*					:				- This required that I split the paymentNRT andpaymentRT into 2 dif functions as they are VERY different.
*
*
*	By				: George Leonard (georgelza@gmail.com)
*
*	jsonformatter 	: https://jsonformatter.curiousconcept.com/#
*
*
*
*****************************************************************************/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/TylerBrock/colorjson"
	"github.com/tkanos/gonfig"
	glog "google.golang.org/grpc/grpclog"

	// My Types/Structs/functions
	"cmd/types"

	"crypto/tls"
	"crypto/x509"
)

var (
	grpcLog  glog.LoggerV2
	validate = validator.New()
	varSeed  types.TPSeed
	vGeneral types.Tp_general
)

func init() {

	// Keeping it very simple
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

	grpcLog.Infoln("###############################################################")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Project   : TFM 2.0")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Comment   : FeatureSpace Scenario Publisher / Fake Data Generator")
	grpcLog.Infoln("#             : To be Event/Alert publisher")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   By        : George Leonard (georgelza@gmail.com)")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Date/Time :", time.Now().Format("2006-01-02 15:04:05"))
	grpcLog.Infoln("#")
	grpcLog.Infoln("###############################################################")
	grpcLog.Infoln("")
	grpcLog.Infoln("")

}

func loadConfig(params ...string) types.Tp_general {

	vGeneral := types.Tp_general{}
	env := "dev"
	if len(params) > 0 {
		env = params[0]
	}

	path, err := os.Getwd()
	if err != nil {
		grpcLog.Fatalln("Problem retrieving current path: %s", err)

	}

	//	fileName := fmt.Sprintf("%s/%s_app.json", path, env)
	fileName := fmt.Sprintf("%s/%s_app.json", path, env)
	err = gonfig.GetConf(fileName, &vGeneral)
	if err != nil {
		grpcLog.Fatalln("Error Reading Config File: ", err)

	} else {

		vHostname, err := os.Hostname()
		if err != nil {
			grpcLog.Fatalln("Can't retrieve hostname %s", err)

		}
		vGeneral.Hostname = vHostname

		vGeneral.Cert_file = path + "/" + vGeneral.Cert_dir + "/" + vGeneral.Cert_file
		vGeneral.Cert_key = path + "/" + vGeneral.Cert_dir + "/" + vGeneral.Cert_key

		if vGeneral.Json_to_file == 1 {
			vGeneral.Output_path = path + "/" + vGeneral.Output_path

		} else {
			vGeneral.Output_path = ""

		}

		if vGeneral.Json_from_file == 1 {
			vGeneral.Input_path = path + "/" + vGeneral.Input_path

		} else {
			vGeneral.Input_path = ""

		}
	}

	if vGeneral.EchoConfig == 1 {
		printConfig(vGeneral)
	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("*")
		grpcLog.Infoln("* Config:")
		grpcLog.Infoln("* Current path:", path)
		grpcLog.Infoln("* Config File :", fileName)
		grpcLog.Infoln("*")

	}

	return vGeneral
}

func loadSeed(file string) types.TPSeed {

	var vSeed types.TPSeed

	path, err := os.Getwd()
	if err != nil {
		grpcLog.Fatalln("Problem retrieving current path: %s", err)

	}

	fileName := fmt.Sprintf("%s/%s", path, file)
	err = gonfig.GetConf(fileName, &vSeed)
	if err != nil {
		grpcLog.Fatalln("Error Reading Seed File: ", err)

	}

	v, err := json.Marshal(vSeed)
	if err != nil {
		grpcLog.Fatalln("Marchalling error: ", err)
	}

	if vGeneral.EchoSeed == 1 {
		prettyJSON(string(v))

	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("*")
		grpcLog.Infoln("* Seed :")
		grpcLog.Infoln("* Current path:", path)
		grpcLog.Infoln("* Seed File :", fileName)
		grpcLog.Infoln("*")

	}

	return vSeed
}

func printConfig(vGeneral types.Tp_general) {

	grpcLog.Info("****** General Parameters *****")
	grpcLog.Info("*")
	grpcLog.Info("* Hostname is\t\t\t", vGeneral.Hostname)
	grpcLog.Info("* Debug Level is\t\t", vGeneral.Debuglevel)
	grpcLog.Info("* Echo JSON is\t\t", vGeneral.Echojson)
	grpcLog.Info("*")
	grpcLog.Info("* Sleep Duration is\t\t", vGeneral.Sleep)
	grpcLog.Info("* Test Batch Size is\t\t", vGeneral.Testsize)
	grpcLog.Info("* Call FS API is\t\t", vGeneral.Call_fs_api)
	grpcLog.Info("* HTTP JSON POST URL is\t", vGeneral.Httpposturl)
	grpcLog.Info("* Cert file is\t\t", vGeneral.Cert_file)
	grpcLog.Info("* Cert key is\t\t\t", vGeneral.Cert_key)
	grpcLog.Info("* Event Type is\t\t", vGeneral.Eventtype)
	grpcLog.Info("* Output JSON to file is\t", vGeneral.Json_to_file)
	grpcLog.Info("* Output path is\t\t", vGeneral.Output_path)
	grpcLog.Info("* Read JSON from file is\t", vGeneral.Json_from_file)
	grpcLog.Info("* Input path is\t\t", vGeneral.Input_path)
	grpcLog.Info("* MinTransactionValue is\tR ", vGeneral.MinTransactionValue)
	grpcLog.Info("* MaxTransactionValue is\tR ", vGeneral.MaxTransactionValue)
	grpcLog.Info("* SeedFile is\t\t\t", vGeneral.SeedFile)
	grpcLog.Info("* EchoSeed is\t\t\t", vGeneral.EchoSeed)
	grpcLog.Info("*")
	grpcLog.Info("*******************************")

	grpcLog.Info("")

}

// Pretty Print JSON string
func prettyJSON(ms string) {

	var obj map[string]interface{}

	json.Unmarshal([]byte(ms), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	result, _ := f.Marshal(obj)
	fmt.Println(string(result))

}

// paymentNRT payload build
func contructPaymentNRTFromFake() (t_Payment map[string]interface{}) {

	// We just using gofakeit to pad the json document size a bit.
	//
	// https://github.com/brianvoe/gofakeit
	// https://pkg.go.dev/github.com/brianvoe/gofakeit

	gofakeit.Seed(time.Now().UnixNano())
	gofakeit.Seed(0)

	//tenantCount := len(varSeed.Tenants) - 1
	tenantCount := 4 // tenants - From Banks

	directionCount := len(varSeed.Direction) - 1
	cDirection := varSeed.Direction[gofakeit.Number(0, directionCount)]

	cTenant := gofakeit.Number(0, tenantCount) // tenants - to Bank
	jTenant := varSeed.Tenants[cTenant]
	jTenantBranchId := gofakeit.Number(jTenant.BranchRangeStart, jTenant.BranchRangeEnd)

	cTo := gofakeit.Number(0, tenantCount)
	jTo := varSeed.Tenants[cTo]
	jToBranchId := gofakeit.Number(jTo.BranchRangeStart, jTo.BranchRangeEnd)

	paymentFrequencyCount := len(varSeed.PaymentFrequency) - 1
	nPaymentFrequency := gofakeit.Number(0, paymentFrequencyCount)
	paymentFrequency := varSeed.PaymentFrequency[nPaymentFrequency]

	nAmount := gofakeit.Price(vGeneral.MinTransactionValue, vGeneral.MaxTransactionValue)
	t_amount := &types.TAmount{
		BaseCurrency: "zar",
		BaseValue:    nAmount,
		Currency:     "zar",
		Value:        nAmount,
	}

	// Are we using good or bad data
	var jMerchant types.TEntity
	if vGeneral.IntroBadEntity == 0 {
		entityCount := len(varSeed.GoodEntities) - 1
		cMerchant := gofakeit.Number(0, entityCount)
		jMerchant = varSeed.GoodEntities[cMerchant]

	} else {
		entityCount := len(varSeed.BadEntities) - 1
		cMerchant := gofakeit.Number(0, entityCount)
		jMerchant = varSeed.BadEntities[cMerchant]

	}

	// Are we using good or bad data
	var jAgent types.TAgent
	if vGeneral.IntroBadAgent == 0 {
		agentCount := len(varSeed.GoodAgent) - 1
		cAgent := gofakeit.Number(0, agentCount)
		jAgent = varSeed.GoodAgent[cAgent]

	} else {
		agentCount := len(varSeed.BadAgent) - 1
		cAgent := gofakeit.Number(0, agentCount)
		jAgent = varSeed.BadAgent[cAgent]

	}

	// We ust showing 2 ways to construct a JSON document to be Marshalled, this is the first using a map/interface,
	// followed by using a set of struct objects added together.
	t_Payment = map[string]interface{}{
		"accountAgentId":                 jAgent.Id,
		"accountAgentName":               jAgent.Name,
		"accountEntityId":                strconv.Itoa(rand.Intn(6)),
		"accountId":                      jMerchant.EntityId,
		"amount":                         t_amount,
		"chargeBearer":                   "SLEV",
		"counterpartyAgentId":            "",
		"counterpartyEntityId":           strconv.Itoa(gofakeit.Number(0, 9999)),
		"counterpartyId":                 strconv.Itoa(gofakeit.Number(10000, 19999)),
		"customerEntityId":               "customerEntityId_1",
		"customerId":                     jMerchant.EntityId,
		"creationDate":                   time.Now().Format("2006-01-02T15:04:05"),
		"destinationCountry":             "ZAF",
		"direction":                      cDirection,
		"eventId":                        uuid.New().String(),
		"eventTime":                      time.Now().Format("2006-01-02T15:04:05"),
		"eventType":                      "paymentNRT",
		"fromFIBranchId":                 jTenantBranchId,
		"fromId":                         jTenant.TenantId,
		"localInstrument":                "42",
		"msgStatus":                      "Success",
		"msgType":                        "RCCT",
		"numberOfTransactions":           1,
		"paymentClearingSystemReference": uuid.New().String(),
		"paymentFrequency":               paymentFrequency,
		"paymentMethod":                  "TRF",
		"paymentReference":               "sdfsfd",
		"remittanceId":                   "sdfsdsd",
		"requestExecutionDate":           time.Now().Format("2006-01-02"),
		"schemaVersion":                  1,
		"settlementClearingSystemCode":   "RTC",
		"settlementDate":                 time.Now().Format("2006-01-02"),
		"settlementMethod":               "CLRG",
		"tenantId":                       jTenant.TenantId,
		"toFIBranchId":                   jToBranchId,
		"toId":                           jTo.TenantId,
		"totalAmount":                    t_amount,
		"transactionId":                  uuid.New().String(),
	}

	return t_Payment
}

// paymentRT payload build
func contructPaymentRTFromFake() (t_Payment map[string]interface{}) {

	// We just using gofakeit to pad the json document size a bit.
	//
	// https://github.com/brianvoe/gofakeit
	// https://pkg.go.dev/github.com/brianvoe/gofakeit

	gofakeit.Seed(time.Now().UnixNano())
	gofakeit.Seed(0)

	nAmount := gofakeit.Price(vGeneral.MinTransactionValue, vGeneral.MaxTransactionValue)
	t_amount := &types.TAmount{
		BaseCurrency: "zar",
		BaseValue:    nAmount,
		Currency:     "zar",
		Value:        nAmount,
	}

	directionCount := len(varSeed.Direction) - 1
	cDirection := gofakeit.Number(0, directionCount)
	direction := varSeed.Direction[cDirection]

	//tenantCount := len(varSeed.Tenants) - 1
	tenantCount := 3

	chargeBearersCount := len(varSeed.ChargeBearers) - 1
	goodEntityCount := len(varSeed.GoodEntities) - 1
	//	badEntityCount := len(varSeed.BadEntities) - 1
	//	goodPayerCount := len(varSeed.GoodPayers) - 1
	//	badPayerCount := len(varSeed.BadPayers) - 1

	paymentFrequencyCount := len(varSeed.PaymentFrequency) - 1
	remittanceLocationMethodCount := len(varSeed.RemittanceLocationMethod) - 1
	settlementMethodCount := len(varSeed.SettlementMethod) - 1
	transactionTypesCount := len(varSeed.TransactionTypesRt) - 1
	verificationResultCount := len(varSeed.VerificationResult) - 1

	cId := gofakeit.Number(0, tenantCount)
	jToID := varSeed.Tenants[cId] // tenants - to Bank
	jToBranchId := gofakeit.Number(jToID.BranchRangeStart, jToID.BranchRangeEnd)
	jTenant := varSeed.Tenants[gofakeit.Number(0, tenantCount)]
	jTenantBranchId := gofakeit.Number(jTenant.BranchRangeStart, jTenant.BranchRangeEnd)

	nChargeBearers := gofakeit.Number(0, chargeBearersCount)
	//cCounterPartyAgent := gofakeit.Number(0, agentsCount) // Agents
	nCounterParty := gofakeit.Number(0, goodEntityCount) // Agents
	nPaymentFrequency := gofakeit.Number(0, paymentFrequencyCount)
	nRemittanceLocationMethod := gofakeit.Number(0, remittanceLocationMethodCount)
	nSettlementMethod := gofakeit.Number(0, settlementMethodCount)
	nTransactionTypesRt := gofakeit.Number(3, transactionTypesCount) // 3 onwards is RT types
	nVerificationResult := gofakeit.Number(0, verificationResultCount)

	chargeBearers := varSeed.ChargeBearers[nChargeBearers]
	paymentFrequency := varSeed.PaymentFrequency[nPaymentFrequency]
	xtransactionTypes := varSeed.TransactionTypesRt[nTransactionTypesRt]
	xverificationResult := varSeed.VerificationResult[nVerificationResult]

	// Are we using good or bad data
	var jAgent types.TAgent
	var jCounterPartyAgent types.TAgent
	var jInstructedAgent types.TAgent
	var jInstructingAgent types.TAgent
	var jIntermediaryAgent1Id types.TAgent
	var jIntermediaryAgent2Id types.TAgent
	var jIntermediaryAgent3Id types.TAgent
	var jPayer types.TPayer

	agentCount := len(varSeed.GoodAgent) - 1

	cAgent := gofakeit.Number(0, agentCount)
	jAgent = varSeed.GoodAgent[cAgent]

	cCounterPartyAgent := gofakeit.Number(0, agentCount)
	jCounterPartyAgent = varSeed.GoodAgent[cCounterPartyAgent]

	cInstructedAgent := gofakeit.Number(0, agentCount)
	jInstructedAgent = varSeed.GoodAgent[cInstructedAgent]

	cInstructingAgent := gofakeit.Number(0, agentCount)
	jInstructingAgent = varSeed.GoodAgent[cInstructingAgent]

	cIntermediaryAgent1Id := gofakeit.Number(0, agentCount)
	jIntermediaryAgent1Id = varSeed.GoodAgent[cIntermediaryAgent1Id]

	cIntermediaryAgent2Id := gofakeit.Number(0, agentCount)
	jIntermediaryAgent2Id = varSeed.GoodAgent[cIntermediaryAgent2Id]

	cIntermediaryAgent3Id := gofakeit.Number(0, agentCount)
	jIntermediaryAgent3Id = varSeed.GoodAgent[cIntermediaryAgent3Id]

	payerCount := len(varSeed.GoodPayers) - 1
	cPayer := gofakeit.Number(0, payerCount)
	jPayer = varSeed.GoodPayers[cPayer]

	// We ust showing 2 ways to construct a JSON document to be Marshalled, this is the first using a map/interface,
	// followed by using a set of struct objects added together.
	t_Payment = map[string]interface{}{
		"accountAgentId":   jAgent.Id,
		"accountAgentName": jAgent.Name,
		"accountAgentAddress": map[string]interface{}{
			"addressLine1":       "accountAgentAddress_addressLine1_1",
			"addressLine2":       "accountAgentAddress_addressLine2_1",
			"addressLine3":       "accountAgentAddress_addressLine3_1",
			"addressType":        "accountAgentAddress_addressType_1",
			"country":            "acc",
			"countrySubDivision": "accountAgentAddress_countrySubDivision_1",
			"postalCode":         "accountAgentAddress_postalCode_1",
			"townName":           "accountAgentAddress_townName_1",
		},
		"accountId":           uuid.New().String(),
		"accountBalanceAfter": t_amount,
		//	"accountEntityId":     "",
		"accountAddress": map[string]interface{}{
			"addressLine1": "accountAddress_addressLine1_1",
			"addressLine2": "accountAddress_addressLine2_1",
			"postalCode":   "accountAddress_postalCode_1",
			"townName":     "accountAddress_townName_1",
			"country":      "ZAF",
		},
		"accountName": map[string]interface{}{
			"fullName":   jPayer.Name.FullName,
			"namePrefix": jPayer.Name.NamePrefix,
			"surname":    jPayer.Name.Surname,
		},
		"amount":       t_amount,
		"cardEntityId": strconv.Itoa(gofakeit.CreditCardNumber()) + "CarSupplier",
		"cardId":       strconv.Itoa(gofakeit.CreditCardNumber()),
		"channel":      "",
		"chargeBearer": chargeBearers,
		"counterpartyAddress": map[string]interface{}{
			"addressLine1": "counterpartyAddress_addressLine1_1",
			"addressLine2": "counterpartyAddress_addressLine2_1",
			"addressLine3": "counterpartyAddress_addressLine3_1",
			"country":      "ZAF",
			"postalCode":   "e1234",
		},
		"counterpartyAgentId":      jCounterPartyAgent.Id,
		"counterpartyAgentName":    jCounterPartyAgent.Name,
		"counterpartyAgentAddress": jCounterPartyAgent.Address,
		"counterpartyId":           varSeed.GoodEntities[nCounterParty].Id,
		"counterpartyName":         map[string]interface{}{},
		"counterpartyEntityId":     varSeed.GoodEntities[nCounterParty].EntityId,
		"creationDate":             time.Now().Format("2006-01-02T15:04:05"),
		"customerEntityId":         jPayer.AccountNumber,
		"customerId":               jPayer.Id,
		"decorationId":             varSeed.Decoration,
		"destinationCountry":       "ZAF",
		"device": map[string]interface{}{
			"anonymizerInUseFlag": "device_anonymizerInUseFlag_1",
			"deviceFingerprint":   "device_deviceFingerprint_1",
			"deviceIMEI":          "device_deviceIMEI_1",
			"sessionLatitude":     5391.835331574868,
			"sessionLongitude":    6200.391276149532,
		},
		"deviceEntityId":                      "",
		"deviceId":                            "IMEACode",
		"direction":                           direction,
		"eventId":                             uuid.New().String(),
		"eventTime":                           time.Now().Format("2006-01-02T15:04:05"),
		"eventType":                           "paymentRT",
		"finalPaymentDate":                    time.Now().Format("2006-01-02"),
		"firstPaymentDate":                    time.Now().Format("2006-01-02"),
		"fromFIBranchId":                      jTenantBranchId,
		"fromId":                              jTenant.TenantId,
		"instructedAgentId":                   jInstructedAgent.Id,
		"instructedAgentName":                 jInstructedAgent.Name,
		"instructedAgentAddress":              jInstructedAgent.Address,
		"instructingAgentId":                  jInstructingAgent.Id,
		"instructingAgentName":                jInstructingAgent.Name,
		"instructingAgentAddress":             jInstructingAgent.Address,
		"intermediaryAgent1Id":                jIntermediaryAgent1Id.Id,
		"intermediaryAgent1Name":              jIntermediaryAgent1Id.Name,
		"intermediaryAgent1AccountId":         jIntermediaryAgent1Id.AccountId,
		"intermediaryAgent1Address":           jIntermediaryAgent1Id.Address,
		"intermediaryAgent2Id":                jIntermediaryAgent2Id.Id,
		"intermediaryAgent2Name":              jIntermediaryAgent2Id.Name,
		"intermediaryAgent2AccountId":         jIntermediaryAgent2Id.AccountId,
		"intermediaryAgent2Address":           jIntermediaryAgent2Id.Address,
		"intermediaryAgent3Id":                jIntermediaryAgent3Id.Id,
		"intermediaryAgent3Name":              jIntermediaryAgent3Id.Name,
		"intermediaryAgent3AccountId":         jIntermediaryAgent3Id.AccountId,
		"intermediaryAgent3Address":           jIntermediaryAgent3Id.Address,
		"localInstrument":                     "",
		"msgStatus":                           "New",
		"msgStatusReason":                     "New Payee",
		"msgType":                             "2100",
		"numberOfTransactions":                1,
		"paymentClearingSystemReference":      "",
		"paymentFrequency":                    paymentFrequency,
		"paymentMethod":                       "",
		"paymentReference":                    "",
		"remittanceId":                        "",
		"remittanceLocationElectronicAddress": "remittanceLocationElectronicAddress_1",
		"remittanceLocationMethod":            varSeed.RemittanceLocationMethod[nRemittanceLocationMethod],
		"requestExecutionDate":                time.Now().Format("2006-01-02"),
		"schemaVersion":                       1,
		"serviceLevelCode":                    "",
		"settlementClearingSystemCode":        "",
		"settlementDate":                      time.Now().Format("2006-01-02"),
		"settlementMethod":                    varSeed.SettlementMethod[nSettlementMethod],
		"tenantId":                            jTenant.TenantId,
		"toFIBranchId":                        jToBranchId,
		"toId":                                jToID.TenantId,
		"totalAmount":                         t_amount,
		"transactionId":                       uuid.New().String(),
		"transactionType":                     xtransactionTypes,
		"ultimateAccountAddress": map[string]interface{}{
			"addressLine1": "ultimateCounterpartyAddress_addressLine1_1",
			"addressLine2": "ultimateCounterpartyAddress_addressLine2_1",
			"addressLine3": "ultimateCounterpartyAddress_addressLine3_1",
			"postalCode":   "2342342",
			"country":      "ZAF",
		},
		"ultimateAccountId": "",
		"ultimateAccountName": map[string]interface{}{
			"fullName":  "ultimateAccountName_fullName_1",
			"givenName": "ultimateAccountName_givenName_1",
			"surname":   "ultimateAccountName_surname_1",
		},
		"ultimateCounterpartyAddress": map[string]interface{}{
			"addressLine1":       "ultimateCounterpartyAddress_addressLine1_1",
			"addressLine2":       "ultimateCounterpartyAddress_addressLine2_1",
			"countrySubDivision": "ultimateCounterpartyAddress_countrySubDivision_1",
			"latitude":           8534.86390567603,
			"longitude":          6842.924811568308,
			"postalCode":         "ultimateCounterpartyAddress_postalCode_1",
			"townName":           "ultimateCounterpartyAddress_townName_1",
			"country":            "ZAF",
		},
		"ultimateCounterpartyId": "",
		"ultimateCounterpartyName": map[string]interface{}{
			"fullName":  "ultimateCounterpartyName_fullName_1",
			"givenName": "ultimateCounterpartyName_givenName_1",
			"surname":   "ultimateCounterpartyName_surname_1",
		},
		"unstructuredRemittanceInformation": "",
		"verificationResult":                xverificationResult,
		"verificationType": map[string]interface{}{
			"aa": "verificationType_aa_1",
		},
	}

	return t_Payment
}

func ReadJSONFile(varRec string) []byte {

	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(varRec)
	if err != nil {
		grpcLog.Fatalln("Error when opening file: ", err)

	}
	return content
}

func isJSON(content []byte) (isJson bool) {

	var t_Payment interface{}

	err := json.Unmarshal(content, &t_Payment)
	if err != nil {
		isJson = false

	} else {
		isJson = true

	}

	return isJson
}

func contructPaymentFromJSON(varRec string) (t_Payment map[string]interface{}) {

	content := ReadJSONFile(varRec)

	err := json.Unmarshal(content, &t_Payment)
	if err != nil {
		grpcLog.Errorln("Unmarshall error ", err)

	}

	return t_Payment
}

// Query database and get the record set to work with - For now we're mimicing a fake EFT query/fetch
func fetchRecords() {

	if vGeneral.Debuglevel > 1 {
		grpcLog.Info("**** Quering Backend database ****")

	}

	// Execute a large sql #1 execute
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10000) // if vGeneral.sleep = 10000, 10 second
	if vGeneral.Debuglevel > 1 {
		grpcLog.Info("EFT SQL Sleeping Millisecond - Simulating long database fetch...", n)

	}

	time.Sleep(time.Duration(n) * time.Millisecond)

	if vGeneral.Debuglevel > 1 {
		grpcLog.Info("**** Backend dataset retrieved ****")

	}
}

// Return list of files located in input_path to be repackaged as JSON payloads and posted onto the API endpoint
func fetchJSONRecords(input_path string) (records map[int]string, count int) {

	count = 0

	m := make(map[int]string)

	// https://yourbasic.org/golang/list-files-in-directory/
	// Use the ioutil.ReadDir function in package io/ioutil. It returns a sorted slice containing elements of type os.FileInfo.
	files, err := ioutil.ReadDir(input_path)
	if err != nil {
		grpcLog.Errorln("Problem retrieving list of input files: %s", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			m[count] = file.Name()
			count++

		}
	}

	records = m

	return records, count
}

func runLoader() {

	// Initialize the vGeneral struct variable - This holds our configuration settings.
	vGeneral = loadConfig("dev")

	// Lets get Seed Data from the specified seed file
	varSeed = loadSeed(vGeneral.SeedFile)

	// Create client with Cert once
	// https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang

	caCert, err := ioutil.ReadFile(vGeneral.Cert_file)
	if err != nil {
		grpcLog.Errorln("Problem reading :", vGeneral.Cert_file, " Error :", err)

	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, _ := tls.LoadX509KeyPair(vGeneral.Cert_file, vGeneral.Cert_key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true, // Self Signed cert
				Certificates:       []tls.Certificate{cert},
			},
		},
	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Info("**** LETS GO Processing ****")
		grpcLog.Infoln("")

	}

	////////////////////////////////////////////////////////////////////////
	// Lets fecth the records that need to be pushed to the fs api end point
	var todo_count = 0
	var returnedRecs map[int]string
	if vGeneral.Json_from_file == 0 { // Build Fake Record - atm we're generating the data, eventually we might fetch via SQL

		// False SQL fetch / sleep
		fetchRecords()

		// As we're still faking it:
		todo_count = vGeneral.Testsize // this will be recplaced by the value of todo_count from above.

	} else { // Build Record set from data fetched from JSON files in input_path

		// this will return an map of files names, each being a JSON document
		returnedRecs, todo_count = fetchJSONRecords(vGeneral.Input_path)

		if vGeneral.Debuglevel > 1 {
			grpcLog.Infoln("Checking input event files (Making sure it's valid JSON)...")
			grpcLog.Infoln("")

		}

		var weFailed bool = false
		for count := 0; count < todo_count; count++ {

			// Build the entire JSON Payload document, either a fake record or from a input/scenario JSON file

			filename := vGeneral.Input_path + "/" + returnedRecs[count]

			contents := ReadJSONFile(filename)
			if !isJSON(contents) {
				weFailed = true
				grpcLog.Infoln(filename, "=> FAIL")

			} else {
				grpcLog.Infoln(filename, "=> Pass")

			}

		}
		if weFailed {
			os.Exit(1)
		}
		grpcLog.Infoln("")

	}

	if vGeneral.Debuglevel > 1 {
		grpcLog.Infoln("Number of records to Process", todo_count) // just doing this to prefer a unused error

	}

	// now we loop through the results, building a json document based on FS requirements and then post it, for this code I'm posting to
	// Confluent Kafka topic, but it's easy to change to have it post to a API endpoint.

	// this is to keep record of the total batch run time
	vStart := time.Now()

	for count := 0; count < todo_count; count++ {

		if vGeneral.Debuglevel > 1 {
			grpcLog.Infoln("")
			grpcLog.Infoln("Record                :", count+1)

		}

		// We're going to time every record and push that to prometheus
		txnStart := time.Now()

		var t_Payload map[string]interface{}

		// Build the entire JSON Payload document, either a fake record or from a input/scenario JSON file
		if vGeneral.Json_from_file == 0 { // Build Fake Record

			if vGeneral.Eventtype == "paymentNRT" {
				// They are just to different to have kept in one function, so split them into 2 seperate specific use case functions.
				t_Payload = contructPaymentNRTFromFake()

			} else if vGeneral.Eventtype == "paymentRT" {
				t_Payload = contructPaymentRTFromFake()

			}
		} else {
			// returnedRecs is a map of file names, each filename is JSON document which contains a FS Payment event,
			// At this point we simply post the contents of the payment onto the end point, and record the response.
			filename := vGeneral.Input_path + "/" + returnedRecs[count]
			if vGeneral.Debuglevel > 2 {
				grpcLog.Infoln("Source Event          :", filename)

			}
			t_Payload = contructPaymentFromJSON(filename)

			// we update/refresh the eventID, to ensure we don't get duplicate id's at POST time
			t_Payload["eventId"] = uuid.New().String()
			if vGeneral.Debuglevel > 1 {
				grpcLog.Infoln("eventId assigned      :", t_Payload["eventId"])

			}
		}

		valueBytes, err := json.Marshal(t_Payload)
		if err != nil {
			grpcLog.Errorln("Marchalling error: ", err)

		}

		if vGeneral.Debuglevel > 1 && vGeneral.Echojson == 1 {
			grpcLog.Infoln("Output Payload   	:")
			prettyJSON(string(valueBytes))
		}

		var body []byte
		var tBody map[string]interface{}

		if vGeneral.Call_fs_api == 1 { // POST to API endpoint

			// Demo environment only available:
			// 07:00 to 19:00
			apiStart := time.Now()

			// https://golangtutorial.dev/tips/http-post-json-go/
			request, err := http.NewRequest("POST", vGeneral.Httpposturl, bytes.NewBuffer(valueBytes))
			if err != nil {
				grpcLog.Errorln("http.NewRequest error: ", err)

			}

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			response, err := client.Do(request)
			if err != nil {
				grpcLog.Errorln("client.Do error: ", err)

			}
			defer response.Body.Close()

			// Did we call the API, how long did it take, do this here before we write to a file that will impact this time
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("API Call Time         :", time.Since(apiStart).Seconds(), "Sec")

			}

			body, _ = ioutil.ReadAll(response.Body)
			if vGeneral.Debuglevel > 2 {
				grpcLog.Infoln("response Payload      :")
				grpcLog.Infoln("response Status       :", response.Status)
				grpcLog.Infoln("response Headers      :", response.Header)

				if response.Status == "200 OK" {
					// it's a paymentNT - SUCCESS
					// it's a paymentRT and we have a very big body

					json.Unmarshal(body, &tBody)
					if vGeneral.Echojson == 1 {
						grpcLog.Infoln("response Body        :")
						prettyJSON(string(body))

					} else {
						grpcLog.Infoln("response Body         : JSON Printing Disabled!")

					}

				} else if response.Status == "204 No Content" {
					// it's a paymentNRT - SUCCESS
					// lets build a body of the header and some additional information

					grpcLog.Infoln("response Body         : paymentNRT")
					tBody = map[string]interface{}{
						"eventId":         t_Payload["eventId"],
						"eventType":       t_Payload["eventType"],
						"responseStatus":  response.Status,
						"responseHeaders": response.Header,
						"processTime":     time.Now().UTC(),
					}

				} else {
					// oh sh$t, its not a success so now to try and build a body to fault fix later

					grpcLog.Infoln("response Body        :", string(body))

					grpcLog.Infoln("response Result         : FAILED POST")
					tBody = map[string]interface{}{
						"eventId":         t_Payload["eventId"],
						"eventType":       t_Payload["eventType"],
						"responseResult":  "FAILED POST",
						"responseBody":    string(body),
						"responseStatus":  response.Status,
						"responseHeaders": response.Header,
						"processTime":     time.Now().UTC(),
					}
				}

			}

		}

		// even if we post to FS API or not, we want isolated control if we output to the json file.
		if vGeneral.Json_to_file == 1 {

			fileStart := time.Now()

			//...................................
			// Writing struct type to a JSON file
			//...................................
			// Writing
			// https://www.golangprograms.com/golang-writing-struct-to-json-file.html
			// https://www.developer.com/languages/json-files-golang/
			// Reading
			// https://medium.com/kanoteknologi/better-way-to-read-and-write-json-file-in-golang-9d575b7254f2

			tagId := t_Payload["eventId"]

			loc_in := fmt.Sprintf("%s/%s.json", vGeneral.Output_path, tagId)
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("Output Event          :", loc_in)

			}

			fd, err := json.MarshalIndent(t_Payload, "", " ")
			if err != nil {
				grpcLog.Errorln("MarshalIndent error", err)

			}

			err = ioutil.WriteFile(loc_in, fd, 0644)
			if err != nil {
				grpcLog.Errorln("ioutil.WriteFile error", err)

			}

			// Did we call the API endpoint above... if yes then do these steps
			if vGeneral.Call_fs_api == 1 { // we need to call the API to get a output/response on paymentRT events

				loc_out := fmt.Sprintf("%s/%s-out.json", vGeneral.Output_path, tagId)
				if vGeneral.Debuglevel > 0 {
					grpcLog.Infoln("engineResponse        :", loc_out)

				}

				fj, err := json.MarshalIndent(tBody, "", " ")
				if err != nil {
					grpcLog.Errorln("MarshalIndent error", err)

				}

				err = ioutil.WriteFile(loc_out, fj, 0644)
				if err != nil {
					grpcLog.Errorln("ioutil.WriteFile error", err)

				}
			}

			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("JSON to File Time     :", time.Since(fileStart).Seconds(), "Sec")

			}
		}

		if vGeneral.Debuglevel > 0 {
			grpcLog.Infoln("Total Time            :", time.Since(txnStart).Seconds(), "Sec")

		}

		//////////////////////////////////////////////////
		// THIS IS SLEEP BETWEEN RECORD POSTS
		//
		// if 0 then sleep is disabled otherwise
		//
		// lets get a random value 0 -> vGeneral.sleep, then delay/sleep as up to that fraction of a second.
		// this mimics someone thinking, as if this is being done by a human at a keyboard, for batcvh file processing we don't have this.
		// ie if the user said 200 then it implies a randam value from 0 -> 200 milliseconds.
		//////////////////////////////////////////////////

		if vGeneral.Sleep != 0 {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(vGeneral.Sleep) // if vGeneral.sleep = 1000, then n will be random value of 0 -> 1000  aka 0 and 1 second
			if vGeneral.Debuglevel >= 2 {
				grpcLog.Infof("Going to sleep for    : %d Milliseconds\n", n)

			}
			time.Sleep(time.Duration(n) * time.Millisecond)
		}

	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("")
		grpcLog.Infoln("**** DONE Processing ****")
		grpcLog.Infoln("")

		if vGeneral.Debuglevel >= 1 {
			vEnd := time.Now()
			grpcLog.Infoln("Start      : ", vStart)
			grpcLog.Infoln("End        : ", vEnd)
			grpcLog.Infoln("Duration   : ", vEnd.Sub(vStart))
			grpcLog.Infoln("Records    : ", vGeneral.Testsize)
			grpcLog.Infoln("")
		}
	}

} // runEFTLoader()

func main() {

	grpcLog.Info("****** Starting           *****")

	runLoader()

	grpcLog.Info("****** Completed          *****")

}
