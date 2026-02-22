# ----------------------------------------------------------------
# Run from repo root: make -C cmd/eval/generated [all]
# Or from here: make GOFFL_ROOT=../.. all
# Uses go run from module root so pgpg is resolvable.
GOFFL_ROOT ?= ../..
ROOT := $(abspath $(GOFFL_ROOT))

# Grammar specs are "name|TypePrefix" where TypePrefix is used to
# build FooLexer / FooParser type names.
LEX_SPECS=pemdas_int|PEMDASInt pemdas_mod|PEMDASMod pemdas_f2poly|PEMDASF2Poly
PARSE_SPECS=pemdas_int|PEMDASInt pemdas_mod|PEMDASMod pemdas_f2poly|PEMDASF2Poly

LEXGENS=$(foreach spec,$(LEX_SPECS),go/pkg/lexers/$(firstword $(subst |, ,$(spec))).go)
PARSEGENS=$(foreach spec,$(PARSE_SPECS),go/pkg/parsers/$(firstword $(subst |, ,$(spec))).go)

all: dirs $(LEXGENS) $(PARSEGENS)

dirs:
	@mkdir -p jsons
	@mkdir -p go/pkg/lexers
	@mkdir -p go/pkg/parsers

# ----------------------------------------------------------------
define LEX_GO_RULE
go/pkg/lexers/$(1).go: jsons/$(1)-lex.json
	cd $(ROOT) && go run github.com/johnkerl/pgpg/cmd/lexgen-code -o $(ROOT)/cmd/eval/generated/go/pkg/lexers/$(1).go -type $(2)Lexer $(ROOT)/cmd/eval/generated/jsons/$(1)-lex.json
endef

define LEX_JSON_RULE
jsons/$(1)-lex.json: ../bnfs/$(1).bnf
	cd $(ROOT) && go run github.com/johnkerl/pgpg/cmd/lexgen-tables -o $(ROOT)/cmd/eval/generated/jsons/$(1)-lex.json $(ROOT)/cmd/eval/bnfs/$(1).bnf
endef

define PARSE_GO_RULE
go/pkg/parsers/$(1).go: jsons/$(1)-parse.json
	cd $(ROOT) && go run github.com/johnkerl/pgpg/cmd/parsegen-code -o $(ROOT)/cmd/eval/generated/go/pkg/parsers/$(1).go -package parsers -type $(2)Parser $(ROOT)/cmd/eval/generated/jsons/$(1)-parse.json
endef

define PARSE_JSON_RULE
jsons/$(1)-parse.json: ../bnfs/$(1).bnf
	cd $(ROOT) && go run github.com/johnkerl/pgpg/cmd/parsegen-tables -o $(ROOT)/cmd/eval/generated/jsons/$(1)-parse.json $(ROOT)/cmd/eval/bnfs/$(1).bnf
endef

$(foreach spec,$(LEX_SPECS),$(eval $(call LEX_GO_RULE,$(firstword $(subst |, ,$(spec))),$(word 2,$(subst |, ,$(spec))))))
$(foreach spec,$(LEX_SPECS),$(eval $(call LEX_JSON_RULE,$(firstword $(subst |, ,$(spec))))))
$(foreach spec,$(PARSE_SPECS),$(eval $(call PARSE_GO_RULE,$(firstword $(subst |, ,$(spec))),$(word 2,$(subst |, ,$(spec))))))
$(foreach spec,$(PARSE_SPECS),$(eval $(call PARSE_JSON_RULE,$(firstword $(subst |, ,$(spec))))))


# ----------------------------------------------------------------
clean:

genclean:
	rm -rf jsons
	rm -rf go/pkg/lexers
	rm -rf go/pkg/parsers

# ----------------------------------------------------------------
# Formatting
# go fmt ./... finds Python files which we want to ignore.
fmt:
	-go fmt ./go/pkg/...

# ----------------------------------------------------------------
# Needs first: go install honnef.co/go/tools/cmd/staticcheck@latest
# See also: https://staticcheck.io
staticcheck:
	staticcheck ./...

# ================================================================
# Go does its own dependency management, outside of make.
.PHONY: build check fmt staticcheck dev clean
