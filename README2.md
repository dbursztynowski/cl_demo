# blahblahblah



Figure - setup with external applications

'''
apiVersion: closedlooppooc.closedloop.io/v1
kind: ClosedLoop
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"closedlooppooc.closedloop.io/v1","kind":"ClosedLoop","metadata":{"annotations":{},"labels":{"app.kubernetes.io/created-by":"closedloop","app.kubernetes.io/instance":"closedloop-sample","app.kubernetes.io/managed-by":"kustomize","app.kubernetes.io/name":"closedloop","app.kubernetes.io/part-of":"closedloop"},"name":"closedloop-v2","namespace":"default"},"spec":{"decision":{"config":"{}","kind":"Decision","message":"{}","policy":{"data":{"body":"{\n}\n","name":"closedloop_5g/inner/decision"},"description":"decision","engine":{"api":{"data":"/v1/data/","policy":"/v1/policies/"},"kind":"opa","url":"http://192.168.49.2:32633"},"input":{"schema":"{\n  \"type\": \"object\",\n  \"properties\": {\n    \"open5gs_amf_metric\": {\n      \"type\": \"number\"\n    }\n  }\n}\n","value":"{\"open5gs_amf_metric\":\"#spec.message.open5gs_amf_metric\"}"},"kind":"priority","name":"policy/closedloop_5g/inner/decision","result":{"schema":"{\n  \"type\": \"object\",\n  \"properties\": {\n    \"cpu\": {\n      \"type\": \"string\"\n    }\n  }\n}\n","value":"{\"cpu\":\"cr:execution#spec.message.cpu\"}"},"rule":{"body":"package policy.closedloop_5g.inner.decision\nimport rego.v1\n\ndefault cpu := \"\"\n\ncpu :=  \"100m\" if {\n    to_number(input.open5gs_amf_metric) \u003e= 0\n        to_number(input.open5gs_amf_metric) \u003c 4\n    }\ncpu :=  \"150m\" if {\n    to_number(input.open5gs_amf_metric) \u003e= 4\n        to_number(input.open5gs_amf_metric) \u003c 8\n    }\ncpu :=  \"200m\" if {\n    to_number(input.open5gs_amf_metric) \u003e= 8\n        to_number(input.open5gs_amf_metric) \u003c 10\n    }\ncpu :=  \"250m\" if {\n    to_number(input.open5gs_amf_metric) \u003e= 10\n        to_number(input.open5gs_amf_metric) \u003c= 12\n    }\n","name":"policy.closedloop_5g.inner.decision"}}},"execution":{"config":"{\n  \"function\": {\n    \"name\": \"Podpatch\",\n    \"parameter\": \"cpu\"\n  }\n}\n","kind":"Execution","message":"{}"},"message":"{}","monitoring":{"config":"{requestedpod: true}","kind":"Monitoringv2","message":"{}","policy":{"data":{"body":"{\n}\n","name":"closedloop_5g/inner/monitoring"},"description":"monitoring","engine":{"api":{"data":"/v1/data/","policy":"/v1/policies/"},"kind":"opa","url":"http://192.168.49.2:32633"},"input":{"schema":"{\n  \"type\": \"object\",\n  \"properties\": {\n    \"open5gs_amf_metric\": {\n      \"type\": \"number\"\n    }\n  }\n}\n","value":"{\"open5gs_amf_metric\":\"#spec.message.data.result.0.value.1\"}"},"kind":"threshold","name":"policy/closedloop_5g/inner/monitoring","result":{"schema":"{\n  \"type\": \"object\",\n  \"properties\": {\n    \"open5gs_amf_metric\": {\n      \"type\": \"number\"\n    }\n  }\n}\n","value":"{\"open5gs_amf_metric\":\"cr:decision#spec.message.open5gs_amf_metric\"}"},"rule":{"body":"package policy.closedloop_5g.inner.monitoring\nimport rego.v1\ndefault open5gs_amf_metric := \"\"\nopen5gs_amf_metric := input.open5gs_amf_metric\n","name":"policy.closedloop_5g.inner.monitoring"}}}}}
  creationTimestamp: "2024-11-15T11:02:44Z"
  generation: 1
  labels:
    app.kubernetes.io/created-by: closedloop
    app.kubernetes.io/instance: closedloop-sample
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: closedloop
    app.kubernetes.io/part-of: closedloop
  name: closedloop-v2
  namespace: default
  resourceVersion: "44870323"
  uid: 502aec92-205a-472a-9714-e97819d2a7fc
spec:
  decision:
    config: '{}'
    kind: Decision
    message: '{}'
    policy:
      data:
        body: |
          {
          }
        name: closedloop_5g/inner/decision
      description: decision
      engine:
        api:
          data: /v1/data/
          policy: /v1/policies/
        kind: opa
        url: http://192.168.49.2:32633
      input:
        schema: |
          {
            "type": "object",
            "properties": {
              "open5gs_amf_metric": {
                "type": "number"
              }
            }
          }
        value: '{"open5gs_amf_metric":"#spec.message.open5gs_amf_metric"}'
      kind: priority
      name: policy/closedloop_5g/inner/decision
      result:
        schema: |
          {
            "type": "object",
            "properties": {
              "cpu": {
                "type": "string"
              }
            }
          }
        value: '{"cpu":"cr:execution#spec.message.cpu"}'
      rule:
        body: |
          package policy.closedloop_5g.inner.decision
          import rego.v1

          default cpu := ""

          cpu :=  "100m" if {
              to_number(input.open5gs_amf_metric) >= 0
                  to_number(input.open5gs_amf_metric) < 4
              }
          cpu :=  "150m" if {
              to_number(input.open5gs_amf_metric) >= 4
                  to_number(input.open5gs_amf_metric) < 8
              }
          cpu :=  "200m" if {
              to_number(input.open5gs_amf_metric) >= 8
                  to_number(input.open5gs_amf_metric) < 10
              }
          cpu :=  "250m" if {
              to_number(input.open5gs_amf_metric) >= 10
                  to_number(input.open5gs_amf_metric) <= 12
              }
        name: policy.closedloop_5g.inner.decision
  execution:
    config: |
      {
        "function": {
          "name": "Podpatch",
          "parameter": "cpu"
        }
      }
    kind: Execution
    message: '{}'
  message: '{}'
  monitoring:
    config: '{requestedpod: true}'
    kind: Monitoringv2
    message: '{}'
    policy:
      data:
        body: |
          {
          }
        name: closedloop_5g/inner/monitoring
      description: monitoring
      engine:
        api:
          data: /v1/data/
          policy: /v1/policies/
        kind: opa
        url: http://192.168.49.2:32633
      input:
        schema: |
          {
            "type": "object",
            "properties": {
              "open5gs_amf_metric": {
                "type": "number"
              }
            }
          }
        value: '{"open5gs_amf_metric":"#spec.message.data.result.0.value.1"}'
      kind: threshold
      name: policy/closedloop_5g/inner/monitoring
      result:
        schema: |
          {
            "type": "object",
            "properties": {
              "open5gs_amf_metric": {
                "type": "number"
              }
            }
          }
        value: '{"open5gs_amf_metric":"cr:decision#spec.message.open5gs_amf_metric"}'
      rule:
        body: |
          package policy.closedloop_5g.inner.monitoring
          import rego.v1
          default open5gs_amf_metric := ""
          open5gs_amf_metric := input.open5gs_amf_metric
        name: policy.closedloop_5g.inner.monitoring
status:
  increaserank: start
  increasetime: 2024-11-15 11:02:44.459552452 +0000 UTC m=+251.178921685
  name: closedloop-v2
'''

aaa
