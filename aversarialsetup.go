package main

//
// import "math/rand"
//
// const (
// 	dx         = 16
// 	dy         = 16
// 	iterations = 700
// )
//
// // var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
// var alphabet = []string{"a"}
//
// // GenerateLetters returns data for white-background images with letters
// func GenerateLetters() [][]byte {
// 	var imgs = make([][]byte, len(alphabet))
// 	for i, char := range alphabet {
// 		data := UniformWhiteData(dx, dy)
// 		image := CreateGreyImage(dx, dy, data)
// 		AddLabel(image, 5, 11, char)
// 		imgs[i] = image.Pix
// 	}
// 	return imgs
// }
//
// func gen(inputs [][]byte) {
// 	ad := NewAdversarialNetwork([]int{1, 9, 50, dx * dy, 50, 9, 1}, 3)
//
// 	// train identifier
// 	for i := 0; i < iterations; i++ {
// 		for _, data := range inputs {
// 			ad.TrainReal(data)
// 		}
// 	}
//
// 	for i := 0; i < iterations; i++ {
// 		seed := rand.Float64()
// 		ad.TrainGenerated(seed)
// 		if i%10 == 0 {
// 			iterationData := ad.Generate(rand.Float64())
// 			iterationImage := CreateGreyImage(dx, dy, iterationData)
// 			SaveImage(iterationImage, "/Users/oskanberg/Desktop/result-"+string(i)+".png")
// 		}
// 	}
//
// 	finalData := ad.Generate(rand.Float64())
// 	finalImage := CreateGreyImage(dx, dy, finalData)
// 	SaveImage(finalImage, "/Users/oskanberg/Desktop/result.png")
// }
//
// func main() {
// 	letters := GenerateLetters()
// 	gen(letters)
// }
