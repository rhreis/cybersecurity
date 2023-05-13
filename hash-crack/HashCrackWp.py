#!/usr/bin/python
# Hash Crack in Python3
import sys
from passlib.hash import phpass

def main(argv):
    ver = 1.0
    print("DevSec 360 Hash Crack v%.1f" %(ver))

    wordlist = "/usr/share/john/password.lst"
    hashedPass = "hashedPass.txt"
   
    print("Generate hash\n\twordlist %s" %(wordlist))
    
    with open(wordlist, "r") as passwords:
        with open(hashedPass, "w") as hashedFile:

            lines = passwords.readlines()
            total_lines = len(lines)
            percentage = 0

            print("Total Lines", total_lines)
            count = 1
            for index, password in enumerate(lines):
                current_percentage = int(((index + 1) / total_lines) * 100)
                if current_percentage > percentage:
                    percentage = current_percentage
                    sys.stdout.write("Progress: %d%%   \r" % (current_percentage) )
                    sys.stdout.flush()
                
                password = password.rstrip()
                hashed_password = phpass.hash(password)

                hashedFile.write("Hash " + hashed_password + " Pass " + password + "\n")

    print("Finish crack hash")

if __name__ == "__main__":
    main(sys.argv[1:])