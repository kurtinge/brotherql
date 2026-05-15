package brotherql

// modelInfo identifies a supported Brother QL printer by USB IDs and the
// protocol flags it requires. Entries are matched during USB enumeration.
type modelInfo struct {
	Name            string
	VendorID        uint16
	ProductID       uint16
	NeedsModeSwitch bool // emit ESC i a 0x01 "switch to raster mode" commands
}

// supportedModels lists every QL printer the library recognizes.
// Add new entries here when extending model support.
var supportedModels = []modelInfo{
	{Name: "QL-700", VendorID: 0x04F9, ProductID: 0x2042, NeedsModeSwitch: false},
	{Name: "QL-710W", VendorID: 0x04F9, ProductID: 0x2043, NeedsModeSwitch: true},
}

// findModel returns the registered model matching the given USB IDs.
func findModel(vid, pid uint16) (modelInfo, bool) {
	for _, m := range supportedModels {
		if m.VendorID == vid && m.ProductID == pid {
			return m, true
		}
	}
	return modelInfo{}, false
}
