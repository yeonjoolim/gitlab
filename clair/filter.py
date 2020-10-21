import re
import string

with open("result.txt", 'r') as infile:
    data = infile.readlines()
with open("result.txt", 'w') as outfile:
    for i in data:
        if not i.startswith("Image"):
            outfile.write(i)
