
.PHONY: all clean run

all: run

clean:
	rm face.jpg face.raw

face.raw:
	python3 generate.py

run: face.raw
	python3 transform.py
