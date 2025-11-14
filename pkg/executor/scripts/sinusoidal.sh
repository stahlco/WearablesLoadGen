#!/bin/zsh

# Usage: ./sinusoidal.sh --base=2000 --amp=500 --t=100

BASE=0
AMP=0
T=0

for ARG in "$@"; do
  case $ARG in
    --base=*)
        BASE="${ARG#*=}"
        ;;
    --amp=*)
        AMP="${ARG#*=}"
        ;;
    --t=*)
        T="${ARG#*=}"
        ;;
  esac
done

# Calculates: base + amp *(2*pi/900 * t)
VALUE=$(echo "scale=6; $BASE + $AMP * s(2*3.141592653589793*$T/900)" | bc -l)

echo $VALUE