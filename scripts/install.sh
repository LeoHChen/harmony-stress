#!/usr/bin/env bash

echo "Installing Harmony Validator Spammer + hmy"
curl -LO https://harmony.one/hmycli && mv hmycli hmy && chmod u+x hmy
curl -LOs http://tools.harmony.one.s3.amazonaws.com/release/linux-x86_64/stress && chmod u+x stress
curl -LOs https://raw.githubusercontent.com/SebastianJ/harmony-stress/master/config.yml
curl -LOs https://raw.githubusercontent.com/SebastianJ/harmony-stress/master/staking.yml
echo "Harmony Validator Spammer is now ready to use!"
echo "Invoke it using ./stress - see ./stress --help for all available options"
