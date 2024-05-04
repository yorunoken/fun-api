package utils

type ModsEnum map[string]int

func GetModsEnum(mods []string) int {
	ModsEnum := ModsEnum{
		"NF":  1,
		"EZ":  2,
		"TD":  4,
		"HD":  8,
		"HR":  16,
		"SD":  32,
		"DT":  64,
		"RX":  128,
		"HT":  256,
		"NC":  512,
		"FL":  1024,
		"AT":  2048,
		"SO":  4096,
		"AP":  8192,
		"PF":  16384,
		"4K":  32768,
		"5K":  65536,
		"6K":  131072,
		"7K":  262144,
		"8K":  524288,
		"FI":  1048576,
		"RD":  2097152,
		"CN":  4194304,
		"TP":  8388608,
		"K9":  16777216,
		"KC":  33554432,
		"1K":  67108864,
		"3K":  134217728,
		"2K":  268435456,
		"SV2": 536870912,
		"MR":  1073741824,
	}

	var count int

	for _, mod := range mods {
		if val, ok := ModsEnum[mod]; ok {
			count += val
			if mod == "NC" {
				count += ModsEnum["NC"] + ModsEnum["DT"]
			}
			if mod == "PF" {
				count += ModsEnum["PF"] + ModsEnum["SD"]
			}
		}
	}

	return count
}
