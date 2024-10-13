# SpellChecker

## Description

A simle spellchecker using a bloomfilter.

## Usage

### Run
1) Build files
```
bash start.sh
```
2) Build the bloom filter file by providing a path to a file words:
```
./tmp/spellcheck -build WORDS_FILEPATH -output OUTPUT_FILEPATH
```
3) Spellcheck a list of words
```
./tmp/spellcheck -read BLOOM_FILTER_FILEPATH ...LIST_OF_WORDS
```

Example Spellcheck
```
./tmp/spellcheck -read BLOOM_FILTER_FILEPATH test hello valid code
```
