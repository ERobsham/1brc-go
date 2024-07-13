INPUT_CMD=./cmd/gobrc
OUTPUT_BIN=gobrc

TEST_DATA_FILE=measurements-100m.txt
RUN_DATA_FILE=measurements.txt

DATA_FILE=${TEST_DATA_FILE}

LDFLAGS="-X 'main.FLAG_DebugLogs=${ENABLE_DEBUG_LOGS}'\
-X 'main.FLAG_CPUProf=${ENABLE_CPU_PROF}'\
-X 'main.FLAG_Output=${ENABLE_OUTPUT}'\
-X 'main.FLAG_DataFile=${DATA_FILE}'"

run: ENABLE_DEBUG_LOGS=""
run: ENABLE_CPU_PROF=""
run: ENABLE_OUTPUT=1
run: FLAG_DataFile=${RUN_DATA_FILE}
run: build
	time ./${OUTPUT_BIN}

profile: ENABLE_DEBUG_LOGS=1
profile: ENABLE_CPU_PROF=1
profile: ENABLE_OUTPUT=""
profile: build
	time ./${OUTPUT_BIN}

debug: ENABLE_DEBUG_LOGS=1
debug: ENABLE_CPU_PROF=""
debug: ENABLE_OUTPUT=""
debug: build
	./${OUTPUT_BIN}

build:
	go build -ldflags=${LDFLAGS} -o=${OUTPUT_BIN} ${INPUT_CMD}

view-prof:
	go tool pprof -http :8888 CPUProf.out

.PHONY: run profile debug