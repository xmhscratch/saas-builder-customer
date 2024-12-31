#!/bin/bash

GO_DIR=$GOPATH/src/localdomain/customer;
APP_DIR=/home/web/repos/customer;
DATA_DIR=/home/web/data;

mkdir -p $GO_DIR;

if [ ! \( -e "${GO_DIR}/customer" \) ] ; then
    echo "creating symlink ${GO_DIR}/customer ..."
    ln -s /home/web/repos/customer/ $GO_DIR;
fi

if [ ! \( -e "${DATA_DIR}" \) ] ; then
    echo "creating symlink ${DATA_DIR} ..."
    ln -s $DATA_DIR /home/web/repos/customer/;
fi
