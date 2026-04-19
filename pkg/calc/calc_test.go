package calc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Тест для функции Divide (Обычные тесты) [cite: 939]
func TestDivide(t *testing.T) {
	// Успешный кейс
	res, err := Divide(10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, res)

	// Ошибка: деление на ноль
	_, err = Divide(10, 0)
	assert.Error(t, err)
	assert.Equal(t, "division by zero", err.Error())
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"Both positive numbers", 10, 5, 5},       
		{"Positive minus zero", 10, 0, 10},        
		{"Negative minus positive", -5, 5, -10},   
		{"Both negative", -10, -5, -5},            
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { 
			got := Subtract(tt.a, tt.b) 
			assert.Equal(t, tt.want, got) 
		})
	}
}