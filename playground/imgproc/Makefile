
.PHONY: all clean run

all: run

clean:
	rm face.jpg face.raw

face.raw:
	python generate.py

run: face.raw
	python transform.py
