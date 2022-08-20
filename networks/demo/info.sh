#!/bin/bash

while true ; do
    case "$1" in
        -s|--symbol )
            symbol=$2
            shift 2
        ;;
		*)
            break
        ;;
    esac
done;

./cecli token info -s $symbol | sed 's/bnc//g'
