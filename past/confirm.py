import re
import string


frequency = {}
document_text = open('result.txt', 'r')
text_string = document_text.read().lower()

match_pattern = re.findall(r'\b[unknown]{7,10}\b', text_string)
match_pattern1 = re.findall(r'\b[negligible]{10}\b', text_string)
match_pattern2 = re.findall(r'\b[low]{3,10}\b', text_string)
match_pattern3 = re.findall(r'\b[medium]{6,10}\b', text_string)
match_pattern4 = re.findall(r'\b[high]{4,10}\b', text_string)
match_pattern5 = re.findall(r'\b[critical]{7,10}\b', text_string)
match_pattern6 = re.findall(r'\b[defcon]{6,10}\b', text_string)

for word in match_pattern:
    count = frequency.get(word,0)
    frequency[word] = count + 1 * 0.0 * 1
for word1 in match_pattern1:
    count1 = frequency.get(word1,0)
    frequency[word1] = count1 + 1 * 0.5 * 1
for word2 in match_pattern2:
    count2 = frequency.get(word2,0)
    frequency[word2] = count2 + 1 * 2 * 1
for word3 in match_pattern3:
     count3 = frequency.get(word3,0)
     frequency[word3] = count3 + 1 * 5.45 * 1
for word4 in match_pattern4:
     count4 = frequency.get(word4,0)
     frequency[word4] = count4 + 1 * 7.95 * 1
for word5 in match_pattern5:
    count5 = frequency.get(word5,0)
    frequency[word5] = count5 + 1 * 9.5 * 1
for word6 in match_pattern6:
    count6 = frequency.get(word6,0)
    frequency[word6] = count6 + 1 *10 * 1



frequency_list = frequency.keys()


for words in frequency_list:
    print frequency[words]

