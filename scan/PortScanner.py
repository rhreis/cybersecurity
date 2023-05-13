#!/usr/bin/python
# Port Scanner written in Python3
import sys
import os

ver = 1.0
print("DevSec 360 Port Scanner v%.1f" %(ver))

ip = sys.argv[1]
port = int(sys.argv[2])

print("Scanning Host: %s port %d" %(ip,port))

os.system("netstat -nlpt")

