package closedloop.outer 

import rego.v1

default action_cpu := "0"

default action_memory := "0"

default priority := "none"

action_cpu := input.time if input.action == "cpu"

action_memory := input.time if input.action == "memory"

cpu_add := array.concat(data.cpu, [action_cpu]) if action_cpu > "0"

cpu_add := data.cpu if action_cpu == "0"

memory_add := array.concat(data.memory, [action_memory]) if action_memory > "0"

memory_add := data.memory if action_memory == "0"

cpu := array.slice(cpu_add, count(cpu_add) - 10, count(cpu_add))

memory := array.slice(memory_add, count(memory_add) - 10, count(memory_add))

priority := "cpu" if {
	full("metric")
	count(cpu) < count(memory)
}

priority := "cpu" if {
	full("metric")
	count(cpu) == count(memory)
	cpu[0] < memory[0]
}

priority := "memory" if {
	full("metric")
	count(memory) < count(cpu)
}

priority := "memory" if {
	full("metric")
	count(cpu) == count(memory)
	memory[0] < cpu[0]
}

full(name) if {
	name == "metric"
	count(cpu) == 10
} else if {
	name == "metric"
	count(memory) == 10
}

