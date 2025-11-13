This is a simple Load-Generator which mocks an mqtt-client that publishes Apple Health Data. 
This Load-Generator can be configured to ssh into a remote server to send requests.

Testing all tests
```shell
just test
```


Example config-file:
```YAML
hosts:
  local: {}
  server1:
    ip: 10.0.0.10
    username: user
    keyFile: ~/.ssh/<file>
```

Example Single Payload (single publish):
```JSON
{
  "type": "HKQuantityTypeIdentifierHeartRate",
  "sourceName": "ESP32-Wecker",
  "sourceVersion": "9.0",
  "unit": "count/min",
  "creationDate": "2022-09-17 16:05:23",
  "startDate": "2022-09-17 16:01:03",
  "endDate": "2022-09-17 16:01:03",
  "value": "93",
  "metadata": {
    "HKMetadataKeyHeartRateMotionContext": "0"
  },
  "device": "Apple-Health-Deivice-1"
}
```
