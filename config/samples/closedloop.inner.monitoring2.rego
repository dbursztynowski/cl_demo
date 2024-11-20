package closedloop.inner.monitoring

import rego.v1

default cpu := "ok"

cpu := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some cpu_idx
	  data.value[cpu_idx][metric_idx] == "cpu"
	data.value[cpu_idx][threshold_idx] == "inferior"
	to_number(input.cpu) < to_number(data.value[cpu_idx][value_idx])
	result := "low"
}

cpu := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some cpu_idx
	  data.value[cpu_idx][metric_idx] == "cpu"
	data.value[cpu_idx][threshold_idx] == "superior"
	to_number(input.cpu) > to_number(data.value[cpu_idx][value_idx])
	result := "high"
}

cpu := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some cpu_idx
	  data.value[cpu_idx][metric_idx] == "cpu"
	data.value[cpu_idx][threshold_idx] == "uniform"
	to_number(input.cpu) == to_number(data.value[cpu_idx][value_idx])
	result := "ok"
}

default memory := "ok"

memory := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some memory_idx
	  data.value[memory_idx][metric_idx] == "memory"
	data.value[memory_idx][threshold_idx] == "inferior"
	to_number(input.memory) < to_number(data.value[memory_idx][value_idx])
	result := "low"
}

memory := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some memory_idx
	  data.value[memory_idx][metric_idx] == "memory"
	data.value[memory_idx][threshold_idx] == "superior"
	to_number(input.memory) > to_number(data.value[memory_idx][value_idx])
	result := "high"
}

memory := result if {
	some metric_idx
    data.header[metric_idx] == "metric"
	some threshold_idx
    data.header[threshold_idx] == "threshold"
	some value_idx
    data.header[value_idx] == "value"
	some memory_idx
	  data.value[memory_idx][metric_idx] == "memory"
	data.value[memory_idx][threshold_idx] == "uniform"
	to_number(input.memory) == to_number(data.value[memory_idx][value_idx])
	result := "ok"
}
