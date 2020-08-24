package models

import (
	"image"
	"sync"

	"gocv.io/x/gocv"
)

// MLModel is a wrapper for gocv.Net.
type MLModel struct {
	Model            gocv.Net
	InputLayerName   string
	OutputLayerNames []string
	mutex            *sync.Mutex
}

// NewMLModel returns a new MLModel.
func NewMLModel(filepath string, input string, outputs []string) (*MLModel, error) {
	model := &MLModel{}
	model.Model = gocv.ReadNetFromTensorflow(filepath)
	model.InputLayerName = input
	model.OutputLayerNames = outputs
	model.mutex = &sync.Mutex{}
	return model, nil
}

// Close terminates MLModel.Model.
func (model *MLModel) Close() error {
	return model.Model.Close()
}

func (model *MLModel) predict(img *gocv.Mat) ([][]float32, error) {
	defer img.Close()

	blob := gocv.BlobFromImage(*img, 1.0, image.Pt(48, 48), gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()

	model.mutex.Lock()
	model.Model.SetInput(blob, model.InputLayerName)

	outputs := model.Model.ForwardLayers(model.OutputLayerNames)
	model.mutex.Unlock()

	result := make([][]float32, len(outputs))
	for idx, output := range outputs {
		defer output.Close()

		outputVal, err := output.DataPtrFloat32()
		if err != nil {
			return nil, err
		}

		result[idx] = outputVal
	}

	return result, nil
}

// Predict returns output of model.
func (model *MLModel) Predict(img *gocv.Mat) ([][]float32, error) {
	return model.predict(img)
}

// PredictWithImageFile calls predict with static image file.
func (model *MLModel) PredictWithImageFile(filepath string) ([][]float32, error) {
	img := gocv.IMRead(filepath, gocv.IMReadGrayScale)
	return model.predict(&img)
}
