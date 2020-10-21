import re
import string

data = [] 
data = open("result.txt", "r").readlines()

number = re.findall("\d", data[0])
blank = list(map(float, number))

number0 = re.findall("\d", data[1])
black1 = list(map(float, number))

number1 = re.findall("\d", data[2])
unknown = list(map(float, number1))
unknown_count = unknown[0] * 0
#print(unknown_count)

number2 = re.findall("\d+", data[3])
negligible = list(map(float, number2))
negligible_count = negligible[0] * 0.5
#print(negligible_count)

number3 = re.findall("\d+", data[4])
low = list(map(float, number3))
low_count = low[0] * 2
#print(low_count)

number4 = re.findall("\d+", data[5])
medium = list(map(float, number4))
medium_count = medium[0] * 5.45
#print(medium_count)

number5 = re.findall("\d+", data[6])
high = list(map(float, number5))
high_count = high[0] * 7.95
#print(high_count)

number6 = re.findall("\d+", data[7])
critical = list(map(float, number6))
critical_count = critical[0] * 9.5
#print(critical_count)

#number7 = re.findall("\d+", data[8])
#defcon = list(map(float, number7))
#defcon_count = defcon[0] * 10

severity_value = unknown_count + negligible_count + low_count + medium_count + high_count + critical_count

if severity_value > 150:
	print "Delete your a docker image"
else:
	print "It is a authorized docker image"

print(severity_value)

