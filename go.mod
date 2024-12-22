module facedetection

go 1.20

require (
	github.com/montanaflynn/stats v0.7.1
	github.com/photoprism/photoprism v0.0.0-20241220012234-664f48ec42b0
	github.com/tensorflow/tensorflow v2.18.0+incompatible
)

replace (
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v0.0.0-20190301222220-23c4e26f81f5
)