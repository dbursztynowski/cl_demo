# blahblahblah



Figure - setup with external applications

'''

apiVersion: closedlooppooc.closedloop.io/v1
kind: ClosedLoop
metadata:
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
