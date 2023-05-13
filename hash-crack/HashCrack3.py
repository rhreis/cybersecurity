#!/usr/bin/python
# Hash Crack in Python3
import sys
import crypt

def main(argv):
    ver = 1.0
    print("DevSec 360 Hash Crack v%.1f" %(ver))

    salt = "$6$D/KSS.6J$"
    wordlist = "/usr/share/wordlists/rockyou.txt"
    hashfile = "/home/attacker/Desktop/rockyouhashed.txt"
   
    print("Creating hash file %s from wordlist %s" %(hashfile,wordlist))

    with open(wordlist, "r") as file:

        with open(hashfile, "w") as arquivo_saida:

            for row in file:
                row = row.strip()
                hashed_password = crypt.crypt(row, salt)
                arquivo_saida.write(row + " - " + hashed_password + "\n")

    print("Finish crack hash")

if __name__ == "__main__":
    main(sys.argv[1:])