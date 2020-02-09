CURR_DIR := $(shell cd `pwd` && pwd)
OUTPUT := grid_updater

$(OUTPUT):
	go get -v
	CGO_ENABLED=0 go build -ldflags "-X git.oa.com/data_warehouse/util.APP_BUILD_DATE=`date +'%Y-%m-%d_%H:%M:%S'` -X git.oa.com/data_warehouse/util.APP_COMMIT_ID=`git rev-parse HEAD`" -a -installsuffix cgo -o $(OUTPUT) .

build:
	docker run -it --rm -v $(CURR_DIR):/go/src/git.oa.com/data_warehouse/grid_updater/ golang bash -c "cd /go/src/git.oa.com/data_warehouse/grid_updater; make clean; make"

clean:
	rm -f $(OUTPUT)
