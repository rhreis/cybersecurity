#!/usr/bin/python
# Hash Crack in Python3
import sys
import crypt

def main(argv):
    ver = 1.0
    print("DevSec 360 Hash Crack v%.1f" %(ver))

    hash = "$6$D/KSS.6J$mLl72m7xOpG8d1B5AKE79wa2oO37sTVBbCIWpjtWJntciPhWMWG61N/O2hKoNjLBb/lq59Fj.6UJvAJPOycjN."
    salt = "$6$D/KSS.6J$"
    wordlist = "/usr/share/wordlists/rock-mini.txt"
   
    print("Try to crack \n\thash %s \n\tsalt %s \n\twordlist %s" %(hash,salt,wordlist))

    with open(wordlist, "r") as file:
        for row in file:
            row = row.strip()
            hashed_password = crypt.crypt(row, salt)
            if hash == hashed_password:
                print("Password founded! \n\t%s - %s" %(row, hash))

    print("Finish crack hash")

if __name__ == "__main__":
    main(sys.argv[1:])