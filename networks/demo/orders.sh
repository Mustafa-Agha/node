#!/bin/bash

# ./orders.sh --list-pair BTC_CE --side 1 --price 1 --quantity 100 --from alice --tif 1

chain_id=$CHAIN_ID
id="$(cat /proc/sys/kernel/random/uuid)"

while true ; do
    case "$1" in
        -l|--list-pair )
            pair=$2
            shift 2
        ;;
        --side )
            side=$2
            shift 2
        ;;
		--price )
			price=$2
			shift 2
		;;
		--quantity )
			quantity=$2
			shift 2
		;;
		--from )
			from=$2
			shift 2
		;;
   		--tif )
			tif=$2
			shift 2
		;;
		*)
            break
        ;;
    esac
done;


for ((i=1; i<=10; i ++))
do
	id="$(cat /proc/sys/kernel/random/uuid)"
	expect ./order.exp $id $pair $side $i $quantity $from $chain_id $tif
	sleep 5
done


#expect ./order.exp $id $pair $side $price $quantity $from $chain_id $tif > /dev/null

echo "Order sent success."