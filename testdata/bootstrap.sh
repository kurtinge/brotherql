#!/usr/bin/env bash
# bootstrap.sh — regenerates testdata/*.png and testdata/*.golden.bin
#
# The golden files are byte-for-byte references produced by the Python
# brother_ql library against the same input image. They serve as a stable
# correctness oracle for the Go raster encoder and protocol builder.
#
# This script:
#   1. Creates a local Python venv at testdata/.venv (gitignored)
#   2. Installs Pillow + brother_ql + pyusb
#   3. Generates a deterministic 1-bit PNG test input
#   4. Runs brother_ql against it to dump the canonical print job bytes
#   5. Extracts the raster portion as a separate golden file
#
# Re-run this script if you add new label types or test inputs.

set -euo pipefail

cd "$(dirname "$0")"
VENV=".venv"
PNG="text-label.png"
RASTER_GOLDEN="text-label.golden.bin"
MODELS=("QL-700" "QL-710W")

if [ ! -d "$VENV" ]; then
    echo "Creating venv at testdata/$VENV..."
    python3 -m venv "$VENV"
fi

echo "Installing Python dependencies..."
"$VENV/bin/pip" install --quiet pillow brother_ql pyusb

echo "Generating $PNG (pure 1-bit, no antialiasing)..."
"$VENV/bin/python3" - <<'PY'
from PIL import Image, ImageDraw
img = Image.new('L', (696, 100), 255)
d = ImageDraw.Draw(img)
d.rectangle([(0, 0), (695, 99)], outline=0, width=2)
d.text((20, 30), 'TEST', fill=0)
bw = img.point(lambda v: 0 if v < 128 else 255, mode='L')
bw.convert('1', dither=Image.Dither.NONE).save('text-label.png')
PY

echo "Generating per-model print-job goldens and $RASTER_GOLDEN from brother_ql..."
"$VENV/bin/python3" - "${MODELS[@]}" <<'PY'
import sys
from brother_ql.raster import BrotherQLRaster
from brother_ql.conversion import convert

models = sys.argv[1:]
raster_written = False
for model in models:
    qlr = BrotherQLRaster(model)
    qlr.exception_on_warning = True
    convert(qlr=qlr, images=['text-label.png'], label='62', cut=True, compress=False)
    data = qlr.data
    out = f'print-job.{model}.golden.bin'
    with open(out, 'wb') as f:
        f.write(data)
    print(f'  {out}: {len(data)} bytes')
    if not raster_written:
        raster = bytearray()
        i = 0
        while i < len(data) - 92:
            if data[i] == 0x67 and data[i+1] == 0x00 and data[i+2] == 0x5A:
                raster.extend(data[i+3:i+93])
                i += 93
            else:
                i += 1
        with open('text-label.golden.bin', 'wb') as f:
            f.write(raster)
        print(f'  text-label.golden.bin: {len(raster)} bytes')
        raster_written = True
PY

echo "Done. Run 'go test ./...' to verify goldens still match."
