#!/usr/bin/env bash

########################### SETUP #########################
home=$HOME
src="${home}/go/src/github.com/Mustafa-Agha/node"
deamonhome="${home}/.cechaind"
witnesshome="${home}/.cechaind_publisher"
clihome="${home}/.cecli"
chain_id='test-chain-n4b735'
echo $src
echo $deamonhome
echo $witnesshome
echo $clihome

key_seed_path="${home}"
executable="${src}/build/cechaind"
clipath="${src}/build/cecli"
cli="${clipath} --home ${clihome}"
scripthome="${src}/networks/publisher"
############################ END ##########################

# clean history data
rm -r ${witnesshome}
mkdir -p ${witnesshome}/config

# config witness node
cp ${deamonhome}/config/genesis.json ${witnesshome}/config/
cp ${deamonhome}/config/config.toml ${witnesshome}/config/
cp ${deamonhome}/config/app.toml ${witnesshome}/config/

sed -i -e "s/26/29/g" ${witnesshome}/config/config.toml
sed -i -e "s/6060/9060/g" ${witnesshome}/config/config.toml
sed -i -e "s/logToConsole = true/logToConsole = false/g" ${witnesshome}/config/app.toml

# get validator id
validator_pid=$(ps aux | grep "cechaind start$" | awk '{print $2}')
validatorStatus=$(${cli} status)
validator_id=$(echo ${validatorStatus} | grep -o "\"id\":\"[a-zA-Z0-9]*\"" | sed "s/\"//g" | sed "s/id://g")

# set witness peer to validator and start witness
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"${validator_id}@127.0.0.1:26656\"/g" ${witnesshome}/config/config.toml
sed -i -e "s/prometheus = false/prometheus = true/g" ${witnesshome}/config/config.toml
sed -i -e "s/publishOrderUpdates = false/publishOrderUpdates = true/g" ${witnesshome}/config/app.toml
sed -i -e "s/publishAccountBalance = false/publishAccountBalance = true/g" ${witnesshome}/config/app.toml
sed -i -e "s/publishOrderBook = false/publishOrderBook = true/g" ${witnesshome}/config/app.toml
sed -i -e "s/publishBlockFee = false/publishBlockFee = true/g" ${witnesshome}/config/app.toml
sed -i -e "s/accountBalanceTopic = \"accounts\"/accountBalanceTopic = \"test\"/g" ${witnesshome}/config/app.toml
sed -i -e "s/orderBookTopic = \"orders\"/orderBookTopic = \"test\"/g" ${witnesshome}/config/app.toml
sed -i -e "s/orderUpdatesTopic = \"orders\"/orderUpdatesTopic = \"test\"/g" ${witnesshome}/config/app.toml
sed -i -e "s/transferTopic = \"transfers\"/transferTopic = \"test\"/g" ${witnesshome}/config/app.toml
sed -i -e "s/blockFeeTopic = \"accounts\"/blockFeeTopic = \"test\"/g" ${witnesshome}/config/app.toml
sed -i -e "s/publishKafka = false/publishKafka = true/g" ${witnesshome}/config/app.toml
sed -i -e "s/publishLocal = false/publishLocal = true/g" ${witnesshome}/config/app.toml

# turn on debug level log
sed -i -e "s/log_level = \"main:info,state:info,\*:error\"/log_level = \"debug\"/g" ${witnesshome}/config/config.toml
sed -i -e "s/index_all_tags = false/index_all_tags = true/g" ${witnesshome}/config/config.toml

${executable} start --home ${witnesshome} > ${witnesshome}/log.txt 2>&1 &
witness_pid=$!
echo ${witness_pid}
