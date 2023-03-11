package core

func boolToInt(b bool) uint8 {
	if b {
		return 0x01
	}
	return 0x00
}
