#!/usr/bin/python
# Hash Crack in Python3
import sys
from passlib.hash import phpass

def main(argv):
    ver = 1.0
    print("DevSec 360 Hash Crack v%.1f" %(ver))

    wordlist = "/usr/share/john/password.lst"
    hashfile = "/home/attacker/Desktop/passwordMd5.txt"
   
    print("Creating hash file %s from wordlist %s" %(hashfile,wordlist))

    with open(wordlist, "r") as file:
        passwords = file.readlines()

        with open(hashfile, "w") as arquivo_saida:

            for password in passwords:
                password = password.strip()
                hashed_password = phpass.hash(password)
                arquivo_saida.write(password + " - " + hashed_password + "\n")

    print("Finish crack hash")

if __name__ == "__main__":
    main(sys.argv[1:])