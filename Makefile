MAIN_VER=$(shell awk -F "v" 'NR==1{print $$2;exit}' CHANGELOG.md)
GIT_CNT=$(shell git rev-list --count HEAD)
VERSION=${MAIN_VER}.${GIT_CNT}
BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT_MSG=$(shell git log --pretty=format:"%B" -2)
REMOTE=$(shell git remote -v | grep "beeme" | awk '{print $$1;exit}')
REV=$(shell git rev-parse --short HEAD)
DATE=$(shell date "+%Y-%m-%d %H:%M:%S")
DEPLOY=$(shell awk '$$1=="Deploy"{print $$2;exit}' conf/config.toml)

all:
	go build -i -v -ldflags "-X 'main.version=version: ${VERSION}, git_version: ${GIT_CNT}(${REV}) date: ${DATE}'" -o beeme-${BRANCH}-v${VERSION}

linux:
	GOOS=linux go build -i -v -ldflags "-X 'main.version=version: ${VERSION}, git_version: ${GIT_CNT}(${REV}) date: ${DATE}'" -o beeme-${BRANCH}-v${VERSION}

debug:
	go build -i -v -race -gcflags "-N -l" -o beeme-debug

check_branch_master:

ifneq (${BRANCH}, master)
	echo "branch is not master"
	exit 1
endif

release: check_branch_master
	git pull ${REMOTE} ${BRANCH}
	$(eval NEW_MAIN_VER=$(shell awk -F "v" 'NR==1{print $$2;exit}' CHANGELOG.md))
	$(eval NEW_GIT_CNT=$(shell git rev-list --count HEAD))
	$(eval NEW_VERSION=${NEW_MAIN_VER}.${NEW_GIT_CNT})
	git tag -a v${NEW_MAIN_VER} -m "rc v${NEW_VERSION}"
	git push ${REMOTE} v${NEW_MAIN_VER}

run:
	rm -f routers/commentsRouter*
	bee run -gendoc=true

build:
	sh build.sh

.PHONY: all linux check_branch_master upload release run build
