      package closedloop.inner.monitoring
      import rego.v1
      default cpu := "ok"
      cpu := result if {
        some cpu_idx
        data.metric[cpu_idx] == "cpu"
        data.treshold.kind[cpu_idx] == "inferior"
        to_number(input.cpu) < to_number(data.treshold.value[cpu_idx])
        result := "low"
      }
      cpu := result if {
        some cpu_idx
        data.metric[cpu_idx] == "cpu"
        data.treshold.kind[cpu_idx] == "superior"
        to_number(input.cpu) > to_number(data.treshold.value[cpu_idx])
        result := "high"
      }
      cpu := result if {
        some cpu_idx
        data.metric[cpu_idx] == "cpu"
        data.treshold.kind[cpu_idx] == "uniform"
        to_number(input.cpu) == to_number(data.treshold.value[cpu_idx])
        result := "ok"
      }
      default memory := "ok"
      memory := result if {
        some memory_idx
        data.metric[memory_idx] == "memory"
        data.treshold.kind[memory_idx] == "inferior"
        to_number(input.memory) < to_number(data.treshold.value[memory_idx])
        result := "low"
      }
      memory := result if {
        some memory_idx
        data.metric[memory_idx] == "memory"
        data.treshold.kind[memory_idx] == "superior"
        to_number(input.memory) > to_number(data.tresholdvalue[memory_idx])
        result := "high"
      }
      memory := result if {
        some memory_idx
        data.metric[memory_idx] == "memory"
        data.treshold.kind[memory_idx] == "uniform"
        to_number(input.memory) == to_number(data.treshold.value[memory_idx])
        result := "ok"
      }
