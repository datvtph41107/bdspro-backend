#!/usr/bin/env python3
"""Run a temporary HTTPS tunnel forwarding only the SePay bank webhook."""

import http.client
import subprocess
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path

WEBHOOK_PATH = "/v2/payment/sepay/webhook"


class WebhookProxy(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != WEBHOOK_PATH:
            self.send_error(404)
            return
        try:
            length = int(self.headers.get("Content-Length", "0"))
        except ValueError:
            self.send_error(400)
            return
        if not 0 < length <= 65536 or self.headers.get("Transfer-Encoding"):
            self.send_error(400)
            return
        self.connection.settimeout(15)
        connection = http.client.HTTPConnection("127.0.0.1", 8000, timeout=20)
        try:
            body = self.rfile.read(length)
            if len(body) != length:
                self.send_error(400)
                return
            connection.request("POST", WEBHOOK_PATH, body, {
                "Content-Type": self.headers.get("Content-Type", "application/json"),
                "Authorization": self.headers.get("Authorization", ""),
            })
            response = connection.getresponse()
            payload = response.read()
            self.send_response(response.status)
            self.send_header("Content-Type", response.getheader("Content-Type", "application/json"))
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)
        except (OSError, http.client.HTTPException):
            self.send_error(502, "Local Gateway unavailable")
        finally:
            connection.close()

    def log_message(self, *_):
        pass


if __name__ == "__main__":
    root = Path(__file__).resolve().parents[2]
    cloudflared = root / ".tmp/development/tools/cloudflared"
    if not cloudflared.is_file():
        raise SystemExit("Install cloudflared at .tmp/development/tools/cloudflared first.")
    server = ThreadingHTTPServer(("127.0.0.1", 18080), WebhookProxy)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    print("SePay Test webhook tunnel: append " + WEBHOOK_PATH + " to the HTTPS URL below.", flush=True)
    try:
        subprocess.run([
            str(cloudflared), "tunnel", "--no-autoupdate", "--protocol", "http2",
            "--url", "http://127.0.0.1:18080",
        ], check=True)
    finally:
        server.shutdown()
        server.server_close()
