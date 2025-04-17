package data

import (
	"math"
)

// BollingerBand 表示布林带指标的数据结构
type BollingerBand struct {
	Timestamp int     // 时间戳
	Middle    float64 // 中轨(移动平均线)
	Upper     float64 // 上轨(中轨+标准差乘数)
	Lower     float64 // 下轨(中轨-标准差乘数)
}

// CalculateBollingerBands 计算布林带指标
// 参数:
//
//	data: 原始数据点数组
//	period: 计算周期
//	multiplier: 标准差乘数(通常为2)
//
// 返回值:
//
//	布林带指标数组
func CalculateBollingerBands(data []Dot, period int, multiplier float64) []BollingerBand {
	var result []BollingerBand

	if len(data) < period {
		return result // 数据不足无法计算
	}

	// 滑动窗口计算布林带
	for i := period - 1; i < len(data); i++ {
		// 计算移动平均值(中轨)
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += float64(data[j].Count)
		}
		mean := sum / float64(period)

		// 计算标准差
		var variance float64
		for j := i - period + 1; j <= i; j++ {
			diff := float64(data[j].Count) - mean
			variance += diff * diff
		}
		stdDev := math.Sqrt(variance / float64(period))

		// 生成布林带数据点
		band := BollingerBand{
			Timestamp: data[i].Timestamp,
			Middle:    mean,
			Upper:     mean + multiplier*stdDev,
			Lower:     mean - multiplier*stdDev,
		}

		result = append(result, band)
	}

	return result
}
