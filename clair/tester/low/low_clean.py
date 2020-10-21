import re
 
INPUT_FILE_NAME = 'low_score.txt'
OUTPUT_FILE_NAME = 'low_score_clean.txt'
 
def clean_text(text):
    cleaned_text = re.sub('' , '', text)
    cleaned_text = re.sub('[\{\}\[\]\/?.,;:|\)*~`!^\-_+<>@\#$%&\\\=\(\'\"]','', cleaned_text)
    return cleaned_text
 
def main():
    read_file = open(INPUT_FILE_NAME, 'r')
    write_file = open(OUTPUT_FILE_NAME, 'w')
    text = read_file.read()
    text = clean_text(text)
    write_file.write(text)
    read_file.close()
    write_file.close() 
 
if __name__ == "__main__":
    main()
