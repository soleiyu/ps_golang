default:
	go run main.go > hoge
	gnuplot plot.txt
	eog res.png
