# #!/bin/bash

# # install_zbar_unicode.sh
# # Script tự động build ZBar trên macOS + cài pyzbar
# # Hỗ trợ UTF-8 / tiếng Việt, xử lý lỗi autopoint và Apple Silicon

# set -e

# echo "=== 1. Cài dependencies ==="
# brew install autoconf automake libtool pkg-config libjpeg libpng gettext libiconv
# brew link --force gettext || true

# # Cài Python libs cần thiết
# pip3.10 install --upgrade pip
# pip3.10 install qrcode pillow

# echo "=== 2. Clone ZBar source ==="
# if [ -d "zbar" ]; then
#     echo "Thư mục zbar đã tồn tại, sẽ xóa và clone lại"
#     rm -rf zbar
# fi
# git clone https://github.com/mchehab/zbar.git
# cd zbar

# echo "=== 3. Reset repo để tránh autopoint lỗi ==="
# git reset --hard

# echo "=== 4. Export flags cho Apple Silicon (fix _iconv_open) ==="
# export LDFLAGS="-L/opt/homebrew/opt/libiconv/lib"
# export CPPFLAGS="-I/opt/homebrew/opt/libiconv/include"
# export PKG_CONFIG_PATH="/opt/homebrew/opt/libiconv/lib/pkgconfig"

# echo "=== 5. Generate configure + Makefile (fix autopoint) ==="
# autoreconf -fi

# echo "=== 6. Configure build ==="
# ./configure --enable-video=no --enable-python=yes

# echo "=== 7. Build ZBar ==="
# make

# echo "=== 8. Install ZBar ==="
# sudo make install

# Export library path cho Python/macOS
export DYLD_LIBRARY_PATH=/usr/local/lib:$DYLD_LIBRARY_PATH
echo "DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH"

echo "=== 9. Cài pyzbar ==="
# pip3.10 install --force-reinstall pyzbar pillow
# pip3.10 install --force-reinstall pyzbar==0.3.13 Pillow==10.2.0
pip3.10 install --force-reinstall pyzbar==0.1.9 Pillow==9.5.0 qrcode

echo "=== 10. Test decode QR tiếng Việt ==="
# Tạo QR test file
TEST_QR="../qrcode_test.png"
python3.10 - <<EOF
import qrcode
from PIL import Image
data = "Hiếu"
qr = qrcode.QRCode(version=1, box_size=10, border=4)
qr.add_data(data)
qr.make(fit=True)
img = qr.make_image(fill_color="black", back_color="white")
img.save("$TEST_QR")
print("Test QR code saved to $TEST_QR")

from pyzbar.pyzbar import decode
img2 = Image.open("$TEST_QR")
result = decode(img2)
if result:
    print("Decoded text from QR:", result[0].data.decode("utf-8"))
else:
    print("Không đọc được QR")
EOF

echo "=== Hoàn tất ==="
