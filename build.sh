#!/bin/bash

PNAME=wafer-session-server
VERSION=1.0.0
SCRATCH_DIR=$PNAME-$VERSION

CURRDIR=`dirname "$0"`
BASEDIR=`cd "$CURRDIR"; pwd`
export GOPATH=$BASEDIR
export GOBIN=$BASEDIR/bin

#GOROOT=/usr/local/go-1.9.2
GOROOT=/usr/local/go

$GOROOT/bin/go clean $PNAME
$GOROOT/bin/go install -gcflags "-N -l" $PNAME

rm -rf target
mkdir -p target/$SCRATCH_DIR

cp -r ./bin/$PNAME src/$PNAME/conf src/$PNAME/control.sh target/$SCRATCH_DIR

cd target

tar czf $SCRATCH_DIR.tar.gz $SCRATCH_DIR
#fpm -s dir -t rpm -n $PNAME -v $VERSION --iteration 1.el6 --epoch=`date +%s` --rpm-defattrfile=0755 --prefix=/usr/local/domob/prog.d $SCRATCH_DIR

rm -rf $SCRATCH_DIR
