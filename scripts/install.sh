#!/usr/bin/env bash

echo "Installing Harmony Validator Spammer + hmy"
curl -LO https://harmony.one/hmycli && mv hmycli hmy && chmod u+x hmy
curl -LOs http://tools.harmony.one.s3.amazonaws.com/release/linux-x86_64/stress && chmod u+x stress
curl -LOs https://raw.githubusercontent.com/SebastianJ/harmony-stress/master/config.yml
curl -LOs https://raw.githubusercontent.com/SebastianJ/harmony-stress/master/staking.yml
mkdir -p data && cd data && touch data.txt && curl -LOs https://gist.githubusercontent.com/SebastianJ/f0da9066f8f636df2a5f96eb0c4b07c8/raw/30b6fb182ad26d96186435fc636d8c66194e73ce/receivers.txt && cd ..
echo "Harmony Validator Spammer is now ready to use!"
echo "Invoke it using ./stress - see ./stress --help for all available options"
