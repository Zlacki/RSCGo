#!/bin/bash

# Set the name of the binary file to create, global var
BINARY_NAME='GoPK'


# Clean the old binaries.  Might be unnecessary, not sure.
# Unsure if `find` args are POSIX compliant or not, maybe look into it later and if not fix it.
function cleanBinaries() {
    echo -n "Cleaning old $BINARY_NAME binaries..."
    find ./bin/ -executable -xtype f -delete
    echo 'done'
}

# Build new binaries.  Takes 2 arguments.
#
# arg1(arch): the target platforms architecture
# arg2(os): the target platforms operating system
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

if [[ $# == 0 ]]; then
	# clean house
	cleanBinaries

	# build house
	for os in linux openbsd freebsd netbsd darwin solaris windows; do
	    if [[ $os != 'solaris' ]]; then
	        buildBinaries 386 $os
	    fi
	    buildBinaries amd64 $os
	    echo
	done

	echo "Compilation for all of the platforms has completed.  Goodbye."
elif [[ $1 = 'clean' ]]; then
	cleanBinaries
elif [[ $1 = 'make' ]]; then
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