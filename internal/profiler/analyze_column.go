package profiler

import (
	"strconv"
	"strings"
)

type Column struct {
	Name   string
	Values []string
}

type ColumnResult struct {
	Name       string
	MainType   string
	Filled     float64
	TypeCounts map[string]int
	Stats map[string]string
}

func AnalyzeColumn(column Column) (result ColumnResult) {
	if len(column.Values) == 0 {
		return
	}
	result.Name = column.Name
	tipo := InferType(strings.TrimSpace(column.Values[0]))
	tipos := map[string]int{
		tipo: 0,
	}
	var stats []float64
	empty := 0.0
	for _, v := range column.Values {
		trimv := strings.TrimSpace(v)
		if trimv != "" {
			key := InferType(trimv)
			tipos[key] = tipos[key] + 1
			if key == "float" || key == "int" {
				number, _ := strconv.ParseFloat(trimv, 64)
				stats = append(stats, number)
			}
		} else {
			empty++
		}
	}
	quantidade := tipos[tipo]
	for k, v := range tipos {
		if v > quantidade {
			tipo = k
		}
	}
	
	result.TypeCounts = tipos
	result.MainType = tipo 
	if result.MainType == "float" || result.MainType == "int" {
		result.Stats = StatsCalc(stats)
	}
	result.Filled = ((float64(len(column.Values)) - empty) / float64(len(column.Values)))
	return
}
