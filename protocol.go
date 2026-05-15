package brotherql

import "encoding/binary"

// buildPrintJob assembles the full Brother QL command sequence for printing
// the given raster bytes. Output is the exact byte stream to send over USB.
//
// Sequence layout (per Brother QL raster command reference, mirrored
// from brother_ql Python reference):
//  0. Switch to raster   — ESC i a 0x01            (only if model.NeedsModeSwitch)
//  1. Invalidate         — 200 zero bytes
//  2. Initialize         — ESC @
//  2b. Switch to raster  — ESC i a 0x01            (only if model.NeedsModeSwitch)
//  3. Status request     — ESC i S
//  4. Print info         — ESC i z {validity media_type width length lines:4LE start fixed}
//  5. Auto-cut mode      — ESC i M {0x40 if cut}    (only if AutoCut)
//  6. Cut every N pages  — ESC i A 0x01             (only if AutoCut)
//  7. Advanced mode      — ESC i K {0x08 if cut at end}
//  8. Margin             — ESC i d 0x23 0x00        (35 dots)
//  9. Compression off    — M 0x00
// 10. Raster rows        — repeated: 0x67 0x00 0x5A <90 bytes>
// 11. Print command      — 0x1A (with cut) or 0x0C (no cut)
func buildPrintJob(raster []byte, label LabelType, opts PrintOptions, model modelInfo) []byte {
	var buf []byte

	// 0. Switch to raster mode (newer models default to template/ESC-P mode).
	if model.NeedsModeSwitch {
		buf = append(buf, 0x1B, 0x69, 0x61, 0x01)
	}

	// 1. Invalidate
	buf = append(buf, make([]byte, 200)...)

	// 2. Initialize
	buf = append(buf, 0x1B, 0x40)

	// 2b. Switch to raster mode (some models require it after init too).
	if model.NeedsModeSwitch {
		buf = append(buf, 0x1B, 0x69, 0x61, 0x01)
	}

	// 3. Status request
	buf = append(buf, 0x1B, 0x69, 0x53)

	// 4. Print information
	// Validity flags: bit7 always, bit6 high quality, bit3 mlength set,
	// bit2 mwidth set, bit1 mtype set = 0xCE.
	rasterLines := uint32(len(raster) / bytesPerRow)
	const validityFlags byte = 0xCE
	var mediaType byte = 0x0A
	if label.HeightMM > 0 {
		mediaType = 0x0B
	}
	printInfo := []byte{
		0x1B, 0x69, 0x7A,
		validityFlags,
		mediaType,
		byte(label.WidthMM),
		byte(label.HeightMM),
	}
	rl := make([]byte, 4)
	binary.LittleEndian.PutUint32(rl, rasterLines)
	printInfo = append(printInfo, rl...)
	printInfo = append(printInfo, 0x00, 0x00)
	buf = append(buf, printInfo...)

	if opts.AutoCut {
		// 5. Auto-cut mode
		buf = append(buf, 0x1B, 0x69, 0x4D, 0x40)
		// 6. Cut every 1 page
		buf = append(buf, 0x1B, 0x69, 0x41, 0x01)
	}

	// 7. Advanced mode (cut at end)
	var advanced byte
	if opts.AutoCut {
		advanced |= 0x08
	}
	buf = append(buf, 0x1B, 0x69, 0x4B, advanced)

	// 8. Margin (35 dots)
	buf = append(buf, 0x1B, 0x69, 0x64, 0x23, 0x00)

	// (QL-700 does not support the M compression command; raster data follows
	// directly after margin.)

	// 9. Raster rows
	for i := 0; i < int(rasterLines); i++ {
		rowStart := i * bytesPerRow
		buf = append(buf, 0x67, 0x00, byte(bytesPerRow))
		buf = append(buf, raster[rowStart:rowStart+bytesPerRow]...)
	}

	// 11. Print command
	if opts.AutoCut {
		buf = append(buf, 0x1A)
	} else {
		buf = append(buf, 0x0C)
	}

	return buf
}
