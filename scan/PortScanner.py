#!/usr/bin/env python3
# Port Scanner written in Python3 - versão corrigida
import sys
import subprocess
import ipaddress

ver = 1.0
print(f"DevSec 360 Port Scanner v{ver:.1f}")

def usage():
    print("Uso: python3 scanner.py <IP> <PORT>")
    print("Exemplo: python3 scanner.py 192.168.1.10 80")

# valida número de argumentos
if len(sys.argv) < 3:
    print("Parâmetros insuficientes.")
    usage()
    sys.exit(1)

ip_arg = sys.argv[1]
port_arg = sys.argv[2]

# valida IP
try:
    # aceita IPv4 e IPv6
    ipaddress.ip_address(ip_arg)
except ValueError:
    print(f"IP inválido: {ip_arg}")
    usage()
    sys.exit(1)

# valida porta
try:
    port = int(port_arg)
    if not (1 <= port <= 65535):
        raise ValueError("fora do intervalo")
except ValueError:
    print(f"Porta inválida: {port_arg}. Deve ser um inteiro entre 1 e 65535.")
    usage()
    sys.exit(1)

print(f"Scanning Host: {ip_arg} port {port}")

# executa netstat (requer privilégios para alguns detalhes) — troque/expanda conforme necessário
try:
    completed = subprocess.run(["netstat", "-nlpt"], check=False, text=True, capture_output=True)
    # imprime saída (stdout + stderr)
    if completed.stdout:
        print(completed.stdout)
    if completed.stderr:
        print("netstat stderr:", completed.stderr, file=sys.stderr)
except FileNotFoundError:
    print("Comando 'netstat' não encontrado no sistema.")
except Exception as e:
    print("Erro ao executar netstat:", e)
