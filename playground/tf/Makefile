.PHONY: all clean init plan console test

all: init test plan
clean:   ; rm -rf .terraform*
init:    ; terraform $@
plan:    ; terraform $@
console: ; terraform $@

UNITS          = unit-(a|b|c)
STAGED_UNITS   = (staging|live|sandbox|unstable)\.$(UNITS)
COUNT_UNITS    = $(shell cat projects.auto.tfvars | grep -E '"$(UNITS)"' | wc -l)
COUNT_PROJECTS = $(shell echo local.projects_flattened | terraform console | grep -E '"$(STAGED_UNITS)"' | wc -l)
test:
	# comparing defined number of tfvars units:      $(UNITS)
	# with terraform-derived number of staged units: $(STAGED_UNITS)
	test $(COUNT_UNITS) -eq $(COUNT_PROJECTS)
