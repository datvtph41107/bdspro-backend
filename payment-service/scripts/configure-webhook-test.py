#!/usr/bin/env python3
"""Select inline VietQR with a bank account created in SePay Webhooks Test mode."""

import argparse
import re
import secrets
from pathlib import Path


def configure(env_file: Path, bank_bin: str, account_number: str) -> None:
    if not re.fullmatch(r"[0-9]{6}", bank_bin):
        raise ValueError("Mã BIN ngân hàng phải gồm 6 chữ số.")
    if not re.fullmatch(r"[0-9]{5,19}", account_number):
        raise ValueError("Số tài khoản Test mode phải gồm 5–19 chữ số.")
    if not env_file.is_file():
        raise ValueError("Chưa có payment-service/.env; hãy chạy make setup trước.")

    updates = {
        "SEPAY_GATEWAY_ENVIRONMENT": "sandbox",
        "SEPAY_BANK_TRANSFER_FLOW": "vietqr",
        "SEPAY_WEBHOOK_ENVIRONMENT": "test",
        "SEPAY_BANK_TRANSFER_ENABLED": "true",
        "SEPAY_CARD_ENABLED": "false",
        "SEPAY_QR_BANK_BIN": bank_bin,
        "SEPAY_QR_ACCOUNT_NUMBER": account_number,
    }
    lines = env_file.read_text().splitlines()
    # Preserve the existing webhook secret; never print it into task output.
    has_key = any(
        re.match(r"^\s*(?:export\s+)?SEPAY_API_KEY\s*=", line)
        and line.split("=", 1)[1].strip().strip("\"'")
        for line in lines
    )
    if not has_key:
        updates["SEPAY_API_KEY"] = secrets.token_hex(32)
    rewritten = []
    written = set()
    for line in lines:
        match = re.match(r"^\s*(?:export\s+)?([A-Z0-9_]+)\s*=", line)
        key = match.group(1) if match else None
        if key in updates:
            if key not in written:
                rewritten.append(f"{key}={updates[key]}")
                written.add(key)
        else:
            rewritten.append(line)
    rewritten.extend(f"{key}={value}" for key, value in updates.items() if key not in written)
    env_file.write_text("\n".join(rewritten) + "\n")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--bank-bin", required=True)
    parser.add_argument("--account-number", required=True)
    args = parser.parse_args()
    try:
        configure(Path(__file__).resolve().parents[1] / ".env", args.bank_bin.strip(), args.account_number.strip())
    except ValueError as error:
        parser.error(str(error))
    print("Đã chọn VietQR + SePay Webhooks Test mode trong payment-service/.env.")
    print("Khởi động lại Payment và Gateway; tạo phiên đăng ký mới để nhận QR.")
    print("Tạo webhook Test mode tới <public-gateway>/v2/payment/sepay/webhook,")
    print("chọn JSON, tiền vào và API Key trùng SEPAY_API_KEY trong payment-service/.env.")
    print("Mô phỏng giao dịch trên SePay; không chuyển tiền thật bằng QR thử nghiệm.")
