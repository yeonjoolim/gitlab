import string
import re

f = open("severity.txt", "r")

lines = []

for paragraph in f:
    lines = string.split(paragraph, "\n")
    for each_line in lines:
        if each_line in lines:
           if each_line.find("Low:") > 0:
              numbers = re.findall("\d+", each_line)
              print numbers             
           else:
              pass
f.close()
