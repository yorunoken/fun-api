ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build: clean
	@mkdir lib/
	@git clone https://github.com/KotRikD/rosu-pp-go && mv rosu-pp-go/rosu-pp-ffi ./ && rm -rf rosu-pp-go
	@cd rosu-pp-ffi/ && cargo build --release && cargo test
	@cp rosu-pp-ffi/target/release/librosu_pp_ffi.so lib/
	@cp rosu-pp-ffi/bindings/rosu_pp_ffi.h lib/

.PHONY: run
run: 
	@rm -rf fun_yorunoken_com
	@go build -ldflags="-r $(ROOT_DIR)lib" -o fun_yorunoken_com
	@./fun_yorunoken_com

.PHONY: clean
clean:
	rm -rf lib/ rosu-pp-ffi/ fun_yorunoken_com