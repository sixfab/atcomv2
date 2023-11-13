package atcom

type SupportedModem struct {
	vid     string
	pid     string
	vendor  string
	product string
	ifs     string
}

var supportedModems = []SupportedModem{
	// Quectel
	{"2c7c", "0125", "Quectel", "EC25", "if02"},
	{"2c7c", "0121", "Quectel", "EC21", "if02"},
	{"2c7c", "0296", "Quectel", "BG96", "if02"},
	{"2c7c", "0700", "Quectel", "BG95", "if02"},
	{"2c7c", "0306", "Quectel", "EP06", "if02"},
	{"2c7c", "0800", "Quectel", "RM5XXQ", "if02"},
	// Telit
	{"1bc7", "1201", "Telit", "LE910Cx RMNET", "if04"},
	{"1bc7", "1203", "Telit", "LE910Cx RNDIS", "if05"},
	{"1bc7", "1204", "Telit", "LE910Cx MBIM", "if05"},
	{"1bc7", "1206", "Telit", "LE910Cx ECM", "if05"},
	{"1bc7", "1031", "Telit", "LE910Cx ThreadX RMNET", "if02"},
	{"1bc7", "1033", "Telit", "LE910Cx ThreadX ECM", "if02"},
	{"1bc7", "1034", "Telit", "LE910Cx ThreadX RMNET", "if00"},
	{"1bc7", "1035", "Telit", "LE910Cx ThreadX ECM", "if00"},
	{"1bc7", "1036", "Telit", "LE910Cx ThreadX OPTION ONLY", "if00"},
	{"1bc7", "1101", "Telit", "ME910C1", "if01"},
	{"1bc7", "1102", "Telit", "ME910C1", "if01"},
	{"1bc7", "1052", "Telit", "FN980 RNDIS", "if05"},
	{"1bc7", "1050", "Telit", "FN980 RMNET", "if04"},
	{"1bc7", "1051", "Telit", "FN980 MBIM", "if05"},
	{"1bc7", "1053", "Telit", "FN980 ECM", "if05"},
	// Thales
	{"1e2d", "0069", "Thales/Cinterion", "PLSx3", "if04"},
	{"1e2d", "006f", "Thales/Cinterion", "PLSx3", "if04"},
}
