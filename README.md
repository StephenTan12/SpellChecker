# SpellChecker
A simple spellchecker using a bloom filter

To build:
Run the command bash ./start.sh

To build the bloom filter file by providing a path to a file words:
./tmp/spellcheck -build WORDS_FILEPATH -output OUTPUT_FILEPATH

To spell check a list of words:
./tmp/spellcheck -read BLOOM_FILTER_FILEPATH ...LIST_OF_WORDS