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
  equal_distribution:
    formula: "min + rand*(max-min)"
    min: 800
    max: 1000
executor:
  name: "simulation"
  execution-pattern: "mixed"
  duration: 100
  steps:
    - distribution: "equal_distribution"
      duration: 100
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