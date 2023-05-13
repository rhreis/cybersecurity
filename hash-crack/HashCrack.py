#!/usr/bin/python
# Hash Crack in Python3
import sys
import getopt
import crypt

def main(argv):
    ver = 1.0
    print("DevSec 360 Hash Crack v%.1f" %(ver))

    hash = None
    salt = None
    wordlist = "mywordlist"

    print("argv ", argv)

    try:
        opts, args = getopt.getopt(argv, "h:s:w:", ["ifile=", "ofile="])
    except getopt.GetoptError:
        print('[*] ./HashCrack.py -h <hash> -s <salt> -w <wordlist>')
        sys.exit(1)
        
    for opt, arg in opts:
        if opt in ("-h", "--hash"):
            hash = arg
        elif opt in ("-s", "--salt"):
            salt = arg
        elif opt in ("-w", "--wordlist"):
            wordlist = arg

    print("hash", hash)
    print("salt", salt)
    print("wordlist", wordlist)
    #print("Try to crack \n\thash %s \n\tsalt %s \n\twordlist %s" %(hash,salt,wordlist))

    with open(wordlist, "r") as file:
        for row in file:
            row = row.strip()
            hashed_password = crypt.crypt(row, salt)
            if hash == hashed_password:
                print("Password founded! \n\t%s - %s" %(row, hash))

    print("Finish crack hash")

if __name__ == "__main__":
    main(sys.argv[1:])