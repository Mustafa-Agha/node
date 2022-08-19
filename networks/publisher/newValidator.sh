#!/usr/bin/env bash

home=$HOME
src="${home}/go/src/github.com/Mustafa-Agha/node"
executable="${home}/go/src/github.com/Mustafa-Agha/node/build/tntchaind"
cli="${home}/go/src/github.com/Mustafa-Agha/node/build/tntcli"

## clean history data
#rm -r ${home}/.tntchaind_val2
#
## init a witness node
#${executable} init --name val2 --home ${home}/.tntchaind_val2 > ~/init2.json

# config witness node
cp ${home}/.tntchaind/config/genesis.json ${home}/.tntchaind_val2/config/

sed -i -e "s/26/30/g" ${home}/.tntchaind_val2/config/config.toml
sed -i -e "s/6060/10060/g" ${home}/.tntchaind_val2/config/config.toml

# get validator id
validator_pid=$(ps aux | grep "tntchaind start$" | awk '{print $2}')
validatorStatus=$(${cli} status)
validatorId=$(echo ${validatorStatus} | grep -o "\"id\":\"[a-zA-Z0-9]*\"" | sed "s/\"//g" | sed "s/id://g")
#echo ${validatorId}

# set witness peer to validator and start witness
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"${validatorId}@127.0.0.1:26656\"/g" ${home}/.tntchaind_val2/config/config.toml
sed -i -e "s/index_all_tags = false/index_all_tags = true/g" ${home}/.tntchaind_val2/config/config.toml
${executable} start --home ${home}/.tntchaind_val2 > ${home}/.tntchaind_val2/log.txt 2>&1 &
validator_pid=$!
echo ${validator_pid}
