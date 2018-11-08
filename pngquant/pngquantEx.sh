echo "**************PNGQUANT START"
rm -rf "./pngquantEx"
mkdir "./pngquantEx"

# $1 获取参数
pngquantFile () {
	f=$1
	dir=${f%/*} #获取目录
	if [ ! -d "./pngquantEx/${dir#*.}" ]; then
		mkdir "./pngquantEx/${dir#*.}"
	fi
	pngquant $f --output ./pngquantEx/${f#*/}
  	echo "pngquantFile $f"
}

# 将find搜索的结果给while循环用。
find . -name "*.png" | while read file; do pngquantFile "$file"; done

echo "**************PNGQUANT END"
# help url: 
# https://stackoverflow.com/questions/4321456/find-exec-a-shell-function-in-linux
# https://blog.csdn.net/ljianhui/article/details/43128465
# https://pngquant.org/