#!/bin/bash

# Set the name of the binary file to create, global var
BINARY_NAME='gorsc'


# Clean the old binaries.  Might be unnecessary, not sure.
function cleanBinaries() {
    echo -n "Cleaning old $BINARY_NAME binaries..."
    find ./bin/ -executable -xtype f -delete
    echo 'done'
}

# Build new binaries.  Takes 3 arguments.
#
# arg1(arch): the target platforms architecture
# arg2(os): the target platforms operating system
# arg3(osName): human-readable version of the OS name, for printing messages
function buildBinaries() {
    local arch=$1
    local os=$2
    local shortArch='64'
    local binSuffix=''

    if [[ $arch = 'amd64' ]]; then
        shortArch='64'
    else
        shortArch='32'
    fi

    if [[ $os = 'windows' ]]; then
        binSuffix='.exe'
    fi

    echo -n "Building $BINARY_NAME for $os/$shortArch (dest='./bin/target/$os/$shortArch/$BINARY_NAME$binSuffix')..."
    GOARCH=$arch GOOS=$os go build -o "./bin/target/$os/$shortArch/$BINARY_NAME$binSuffix"
    echo 'done'
}

# clean house
if [[ $1 = 'clean' ]]; then
    cleanBinaries
else
    # build house
    for os in linux openbsd freebsd netbsd darwin solaris windows; do
        if [[ $os != 'solaris' ]]; then
            buildBinaries 386 $os
        fi
        buildBinaries amd64 $os
        echo
    done

    echo "Compilation for all of the platforms has completed.  Goodbye."
fi
