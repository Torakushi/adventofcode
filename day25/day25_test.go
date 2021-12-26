package day25

import "testing"

func BenchmarkCucumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readDatas()
		run := true
		for run {
			run = processEastBunchInParallel()
			run = processSouthBunchInParallel() || run
		}
	}
}
