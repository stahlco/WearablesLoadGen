This is a simple Load-Generator which mocks an mqtt-client that publishes Apple Health Data. 
This Load-Generator can be configured to ssh into a remote server to send requests.

Testing all tests
```shell
just test
```

---
### Examples

This is an example configuration for the executor, which controls the generation patterns
```yaml
distributions:
  sinusoidal:
    formula: "./scripts/sinusoidal.sh"
    base: 1000
    amp: 100
executor:
  name: "simulation"
  execution-pattern: "mixed"
  duration: 100
  steps:
    - distribution: "sinusoidal"
      duration: 100
```

To add new load-pattern please use a script in following format:

```shell
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
```


This is an example for the measurement-types which represent the blueprint for a measurement:
```yaml
measurement-types:
  heart-rate:
    type: "HKQuantityTypeIdentifierHeartRate"
    source-name: "ESP32-Wecker"
    source-version: "9.0"
    min: 100
    max: 200
    unit: "count/min"
```