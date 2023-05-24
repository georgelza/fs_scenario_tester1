package seed

func GetTenants() map[int]string {

	ar := make(map[int]string)

	// Registered Banks
	ar[1] = "absa"     // ABSA Bank  Case Sensitive
	ar[2] = "fnb"      // FNB
	ar[3] = "nedbank"  // Nedbank
	ar[4] = "standard" // Standard Bank
	ar[5] = "Capitect Bank Limited"
	ar[6] = "Bank Zero"
	ar[7] = "Thyme Bank"
	ar[8] = "Discovery Bank Limited"
	ar[9] = "African Bank Limited"
	ar[10] = "Bidvest Bank Limited"
	ar[11] = "Grindrod Bank Limited"
	ar[12] = "Investec Bank Limited"
	ar[13] = "Ithala"
	ar[14] = "Sasfin Bank Limited"
	ar[15] = "Ubank Limited"
	ar[16] = "FirstRand"
	ar[17] = "Imperial Bank South Africa"
	ar[18] = "Mercantile Bank Limited"
	ar[19] = "GBS Mutual Bank"
	ar[20] = "VBS Mutual Bank"
	ar[21] = "Finbond Mutual Bank"
	ar[22] = "Development Bank of Southern Africa"
	ar[23] = "Land and Agricultural Development Bank of South Africa"
	ar[24] = "Postbank"
	ar[25] = "Grobank"
	ar[26] = "CitiBank"
	ar[27] = "Teba Bank Limited"
	// Registered FSI

	// Registered Non Banks

	return ar
}

func GetTransType() map[int]string {

	ar := make(map[int]string)

	ar[1] = "transactionType_1"
	// ar[2] = "EFT"
	// ar[2] = "AC Collection"
	// ar[3] = "RTC"
	// ar[4] = "PayByProxy"
	// ar[5] = "PayByAccount"
	// ar[6] = "RequestToPay"

	return ar
}

func GetDirection() map[int]string {

	ar := make(map[int]string)

	ar[1] = "inbound"
	ar[2] = "outbound"

	return ar
}

func GetGoodEntityId() map[int]string {

	ar := make(map[int]string)

	ar[1] = "Pick n Pay"
	ar[2] = "Spar"
	ar[3] = "Checkers"
	ar[4] = "Starbucks"
	ar[5] = "Shoprite Holdings"

	ar[6] = "Billabong"
	ar[7] = "Quicksilver"

	ar[8] = "Nike"
	ar[9] = "iStore"
	ar[10] = "Incredible Connection"
	ar[11] = "Samsung"
	ar[12] = "Sportsmans Warehouse"

	//ar[12] = "@Home"
	//ar[12] = "MrPriceHome"

	ar[13] = "Builders"
	ar[14] = "BUCO"
	ar[15] = "ACDC"
	ar[16] = "BuildIt"

	ar[17] = "Discem"
	ar[18] = "Clicks"

	ar[19] = "Edgars"
	ar[20] = "Truworths"
	ar[21] = "Markhams"
	ar[22] = "MrPrice"
	ar[23] = "Ackermans"

	ar[24] = "CycleLab"
	ar[25] = "World of Golf"

	ar[26] = "PnA"
	ar[27] = "CNA"
	ar[28] = "Exclusive Books"

	// Petroleum
	ar[29] = "SASOL"
	ar[30] = "Engen"
	ar[31] = "Total"
	ar[32] = "BP"

	ar[33] = "MTN"
	ar[34] = "Telkom"
	ar[35] = "Vodacom"
	ar[36] = "Cell C"
	ar[37] = "Afrihost"
	ar[38] = "WebAfrica"
	ar[39] = "Websquad"
	ar[40] = "Cool Ideas"
	ar[41] = "Frog foot"
	ar[42] = "MWeb"
	ar[43] = "Neotel"
	ar[44] = "Liquid Telecom"
	ar[45] = "Netcare"

	ar[46] = "Cape Union Mart"

	ar[47] = "M-Net"
	ar[48] = "SABC"
	ar[49] = "Netflix"

	// AirLines
	ar[50] = "SAA"
	ar[51] = "Airlink"
	ar[52] = "CemAir"
	ar[53] = "FlySafAir"
	ar[54] = "Lift Air"

	// Food
	ar[55] = "Woolworths"
	ar[56] = "Food Lover's Market"
	ar[57] = "Burger King"
	ar[58] = "Chicken Licken"
	ar[59] = "Dominos pizza"
	ar[60] = "Debonairs Pizza"
	ar[61] = "John Dory's"
	ar[62] = "Wimpy"
	ar[63] = "Krispy Kreme"
	ar[64] = "Kauai"
	ar[65] = "Mugg & Bean"
	ar[66] = "Mozambik"
	ar[67] = "Mochachos"
	ar[68] = "Nando's"
	ar[69] = "News Cafe"
	ar[70] = "Panarottis"
	ar[71] = "Pizza Hut"
	ar[72] = "Popeyes"
	ar[73] = "Papachinos"
	ar[74] = "Roman's Pizza"
	ar[75] = "Roco Mamas"
	ar[76] = "Steers"
	ar[77] = "Subway"
	ar[78] = "Vida e Caffè"

	ar[79] = "Spur"
	ar[80] = "McDonals"
	ar[81] = "Kentucky Fried Chicken"
	ar[82] = "Something Fishy"
	ar[83] = "DunkingDonuts"
	ar[84] = "Milky Lane"
	ar[85] = "Meat company"

	return ar
}

func GetBadEntityId() map[int]string {

	ar := make(map[int]string)

	ar[1] = "Alibaba"
	ar[2] = "Solar Guru"
	ar[3] = "Geni Solutions"
	ar[4] = "Shoes for All"
	ar[5] = "Off Grid Solutions"
	ar[6] = "Cell Phone Central"
	ar[7] = "IT Answers"

	return ar
}

type TPayer struct {
	Name          string `json:"name, omitempty"`
	TenantId      string `json:"tenantid, omitempty"`
	AccountNumber string `json:"accountnumber, omitempty"`
}

func GetGoodPayerId() map[int]TPayer {

	ar := make(map[int]TPayer)

	ar[1] = TPayer{Name: "George Simpson", TenantId: "fnb", AccountNumber: "23434545345"}
	ar[2] = TPayer{Name: "Mark Botha", TenantId: "absa", AccountNumber: "23434545345"}
	ar[3] = TPayer{Name: "Gert Van Der Merwe", TenantId: "standard bank", AccountNumber: "23434545345"}
	ar[4] = TPayer{Name: "Ciljee Klopper", TenantId: "fnb", AccountNumber: "23434545345"}
	ar[5] = TPayer{Name: "Mari Engelbrecght", TenantId: "capitec", AccountNumber: "23434545345"}
	ar[6] = TPayer{Name: "Manyny Groenewaldt", TenantId: "discovery", AccountNumber: "23434545345"}
	ar[7] = TPayer{Name: "Mark Geldenhuys", TenantId: "absa", AccountNumber: "23434545345"}
	ar[8] = TPayer{Name: "Neil Smith", TenantId: "fnb", AccountNumber: "23434545345"}

	return ar
}

func GetBadPayerId() map[int]TPayer {

	ar := make(map[int]TPayer)

	ar[1] = TPayer{Name: "Mark Lourens", TenantId: "absa", AccountNumber: "23434545345"}
	ar[2] = TPayer{Name: "Neil Botha", TenantId: "fnb", AccountNumber: "23434545345"}

	return ar
}
