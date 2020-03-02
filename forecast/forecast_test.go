package forecast

import (
	"context"
	"fmt"

	// "math/rand"
	"testing"
	"time"

	dataframe "github.com/rocketlaunchr/dataframe-go"
	utime "github.com/rocketlaunchr/dataframe-go/utils/utime"
)

func TestEtsSeries(t *testing.T) {
	ctx := context.Background()

	// data := dataframe.NewSeriesFloat64("Complete Data", nil, 445.43, 345.2, 565.56, 433.34, 585.23, 593.32, 641.43, 654.35, 234.65, 567.45, 645.45, 445.34, 564.65, 598.76, 676.54, 654.56, 564.76, 456.76, 656.57, 765.45, 755.43, 745.2, 665.56, 633.34, 585.23, 693.32, 741.43, 654.35, 734.65, 667.45, 545.45, 645.34, 754.65, 798.76, 776.54, 654.56, 664.76, 856.76, 776.57, 825.45, 815.43, 845.2, 765.56, 733.34, 785.23, 893.32, 841.43, 754.35, 524.65, 567.45, 715.45, 845.34, 864.65, 898.76, 876.54, 854.56, 864.76, 856.76, 726.57, 700.31, 815.43, 805.2, 855.56, 733.34, 785.23, 893.32, 641.43, 554.35, 734.63, 834.89)
	// m := 5
	alpha := 0.4

	data := dataframe.NewSeriesFloat64("unit increase", nil, 1.05, 2.5, 3, 4, 5, 6, 7.25, 8, 9.36, 10.04)
	m := 10

	// fmt.Println(data.Table())

	etsModel := NewEtsModel()

	cfg := ExponentialSmootheningConfig{
		Alpha: alpha,
	}

	etsModel.Configure(&cfg)

	etsModel.Load(data, &dataframe.Range{End: &[]int{5}[0]})

	// spew.Dump(etsModel)

	etsPredict, err := etsModel.Predict(ctx, m)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println(etsPredict.(*dataframe.SeriesFloat64).String())

	testData := data.Copy(dataframe.Range{Start: &[]int{6}[0]})
	errVal, err := etsModel.Validate(ctx, testData, nil, RMSE)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println("Root Mean Squared Error", errVal)
}

func TestEtsDF(t *testing.T) {

	ctx := context.Background()

	// prepare test  data
	size := 10
	opts := utime.NewSeriesTimeOptions{
		Size: &size,
	}
	now := time.Date(2020, 2, 13, 22, 25, 28, 0, time.UTC)
	timeRec, err := utime.NewSeriesTime(ctx, "Time Received", "2W1D", now, false, opts)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}

	qty := dataframe.NewSeriesFloat64("unit increase", nil, 1.05, 2.5, 3, 4, 5, 6, 7.25, 8, 9.36, 10.04)

	dataF := dataframe.NewDataFrame(timeRec, qty)

	fmt.Println(dataF.Table())

	// create model
	etsDfModel := NewEtsModel()

	alpha := 0.4

	cfg := ExponentialSmootheningConfig{
		Alpha: alpha,
	}

	etsDfModel.Configure(&cfg)

	etsDfModel.Load(dataF, &dataframe.Range{End: &[]int{5}[0]})

	//	spew.Dump(etsDfModel)
	m := 10
	etsDfPredict, err := etsDfModel.Predict(ctx, m)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println(etsDfPredict.(*dataframe.DataFrame).String())

	testData := dataF.Copy(dataframe.Range{Start: &[]int{6}[0]})
	errVal, err := etsDfModel.Validate(ctx, testData, nil, RMSE)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println("Root Mean Squared Error", errVal)
}

func TestHwSeries(t *testing.T) {
	ctx := context.Background()

	// 48 + 24 = 72 data pts + extra 12
	data := dataframe.NewSeriesFloat64("simple data", nil, 30, 21, 29, 31, 40, 48, 53, 47, 37, 39, 31, 29, 17, 9, 20, 24, 27, 35, 41, 38,
		27, 31, 27, 26, 21, 13, 21, 18, 33, 35, 40, 36, 22, 24, 21, 20, 17, 14, 17, 19,
		26, 29, 40, 31, 20, 24, 18, 26, 17, 9, 17, 21, 28, 32, 46, 33, 23, 28, 22, 27,
		18, 8, 17, 21, 31, 34, 44, 38, 31, 30, 26, 32, 45, 34, 30, 27, 25, 22, 28, 33, 42, 32, 40, 52)

	period := 12
	h := 24

	// fmt.Println(data.Table())
	alpha := 0.45
	beta := 0.03
	gamma := 0.73

	cfg := HoltWintersConfig{
		Alpha:  alpha,
		Beta:   beta,
		Gamma:  gamma,
		Period: period,
	}

	// create model
	hwModel := NewHwModel()

	hwModel.Configure(&cfg)

	hwModel.Load(data, &dataframe.Range{End: &[]int{71}[0]})

	// spew.Dump(hwModel)

	hwPredict, err := hwModel.Predict(ctx, h)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println(hwPredict.(*dataframe.SeriesFloat64).String())

	testData := data.Copy(dataframe.Range{Start: &[]int{72}[0]})
	errVal, err := hwModel.Validate(ctx, testData, nil, RMSE)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println("Root Mean Squared Error", errVal)
}

func TestHwDF(t *testing.T) {

	ctx := context.Background()

	// prepare test  data
	size := 84
	opts := utime.NewSeriesTimeOptions{
		Size: &size,
	}
	now := time.Date(2020, 2, 13, 22, 25, 28, 0, time.UTC)
	timeRec, err := utime.NewSeriesTime(ctx, "Time Received", "2W5D", now, false, opts)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}

	data := dataframe.NewSeriesFloat64("simple data", nil, 30, 21, 29, 31, 40, 48, 53, 47, 37, 39, 31, 29, 17, 9, 20, 24, 27, 35, 41, 38,
		27, 31, 27, 26, 21, 13, 21, 18, 33, 35, 40, 36, 22, 24, 21, 20, 17, 14, 17, 19,
		26, 29, 40, 31, 20, 24, 18, 26, 17, 9, 17, 21, 28, 32, 46, 33, 23, 28, 22, 27,
		18, 8, 17, 21, 31, 34, 44, 38, 31, 30, 26, 32, 45, 34, 30, 27, 25, 22, 28, 33, 42, 32, 40, 52)

	period := 12

	dataF := dataframe.NewDataFrame(timeRec, data)

	fmt.Println(dataF.Table())

	alpha := 0.45
	beta := 0.03
	gamma := 0.73

	h := 24

	cfg := HoltWintersConfig{
		Alpha:  alpha,
		Beta:   beta,
		Gamma:  gamma,
		Period: period,
	}

	// create model
	hwDfModel := NewHwModel()

	hwDfModel.Configure(&cfg)

	hwDfModel.Load(dataF, &dataframe.Range{End: &[]int{71}[0]})

	// spew.Dump(hwDfModel)

	hwDfPredict, err := hwDfModel.Predict(ctx, h)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println(hwDfPredict.(*dataframe.DataFrame).String())

	testData := dataF.Copy(dataframe.Range{Start: &[]int{72}[0]})
	errVal, err := hwDfModel.Validate(ctx, testData, nil, RMSE)
	if err != nil {
		t.Errorf("error encountered: %s", err)
	}
	fmt.Println("Root Mean Squared Error", errVal)
}
