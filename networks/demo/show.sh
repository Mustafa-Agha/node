#!/bin/bash

# ./show.sh -l ADA_CE --from alice

chain_id=$CHAIN_ID

while true ; do
    case "$1" in
        -l|--list-pair )
            pair=$2
            shift 2
        ;;
		*)
            break
        ;;
    esac
done;

./cecli dex show -l $pair
