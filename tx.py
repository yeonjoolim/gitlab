import re
import string

allNums = []
total = 0
with open(r"1.txt", "r+") as f:
   data = f.readlines()
   for line in data:
      allNums += line.strip().split(" ")
   for num in allNums:
       total += float(num)
   if total > 200:
      print "Delete your docker image"
   else:
      print "It's a authorized docker image"
   print total

