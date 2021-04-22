CURR_DIR := $(shell cd `pwd` && pwd)
OUTPUT := grid_updater

$(OUTPUT):
	export GO111MODULE=auto
	export GOPROXY=https://goproxy.cn
	GO111MODULE=auto go get -v
	CGO_ENABLED=0 GOOS=linux GO111MODULE=auto go build -ldflags "-X git.oa.com/data_warehouse/util.APP_BUILD_DATE=`date +'%Y-%m-%d_%H:%M:%S'` -X git.oa.com/data_warehouse/util.APP_COMMIT_ID=`git rev-parse HEAD`" -a -installsuffix cgo -o $(OUTPUT) .

build:
	docker run -it --rm -v $(CURR_DIR):/go/src/git.oa.com/data_warehouse/grid_updater/ golang bash -c "cd /go/src/git.oa.com/data_warehouse/grid_updater; make clean; make"

clean:
	rm -f $(OUTPUT)
