#!/usr/bin/env python3
# Port Scanner escrito em Python3 - versão otimizada
import argparse
import ipaddress
import socket
import sys
from concurrent.futures import ThreadPoolExecutor, as_completed

ver = 1.1
print(f"DevSec 360 Port Scanner v{ver:.1f}")

def parse_ip(value):
    try:
        return str(ipaddress.ip_address(value))
    except ValueError:
        raise argparse.ArgumentTypeError(f"IP inválido: {value}")

def parse_ports(value):
    # aceita "80" ou "20-1024"
    if "-" in value:
        start, end = value.split("-", 1)
        try:
            s = int(start); e = int(end)
        except ValueError:
            raise argparse.ArgumentTypeError("Porta inválida. Use inteiro ou intervalo start-end.")
        if not (1 <= s <= 65535 and 1 <= e <= 65535 and s <= e):
            raise argparse.ArgumentTypeError("Intervalo de portas inválido (1-65535).")
        return list(range(s, e+1))
    else:
        try:
            p = int(value)
        except ValueError:
            raise argparse.ArgumentTypeError("Porta inválida. Use inteiro ou intervalo start-end.")
        if not (1 <= p <= 65535):
            raise argparse.ArgumentTypeError("Porta fora do intervalo (1-65535).")
        return [p]

def scan_port(host, port, timeout):
    try:
        # determina família conforme IP (v4/v6)
        family = socket.AF_INET6 if ":" in host else socket.AF_INET
        with socket.socket(family, socket.SOCK_STREAM) as s:
            s.settimeout(timeout)
            res = s.connect_ex((host, port))
            return port, (res == 0)
    except Exception:
        return port, False

def main():
    parser = argparse.ArgumentParser(
        description="Port scanner simples. Exemplo: python3 PortScanner.py 192.168.1.10 22-25",
        epilog="Exemplo: python3 PortScanner.py 192.168.1.10 80"
    )
    parser.add_argument("host", type=parse_ip, help="Endereço IP (IPv4/IPv6)")
    parser.add_argument("ports", type=parse_ports, help="Porta única (e.g. 80) ou intervalo (e.g. 20-1024)")
    parser.add_argument("--timeout", type=float, default=1.0, help="Timeout por conexão em segundos (default: 1.0)")
    parser.add_argument("--workers", type=int, default=100, help="Número máximo de threads (default: 100)")
    args = parser.parse_args()

    host = args.host
    ports = args.ports
    timeout = args.timeout
    workers = max(1, min(args.workers, len(ports)))

    print(f"Scanning Host: {host} ports: {ports[0]}{'-'+str(ports[-1]) if len(ports)>1 else ''}")

    open_ports = []
    with ThreadPoolExecutor(max_workers=workers) as ex:
        futures = {ex.submit(scan_port, host, p, timeout): p for p in ports}
        for fut in as_completed(futures):
            port, is_open = fut.result()
            if is_open:
                print(f"Porta {port} aberta")
                open_ports.append(port)

    if not open_ports:
        print("Nenhuma porta aberta encontrada no intervalo especificado.")
    else:
        print("Portas abertas:", ", ".join(map(str, sorted(open_ports))))

if __name__ == "__main__":
    main()
