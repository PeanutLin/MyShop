# 创建文件夹
mkdir -p productshop/backend/web/
mkdir -p productshop/fonted/web/

# 复制本地资源到目标文件夹
cp -r ./assets productshop/
cp -r ./backend/web/views productshop/backend/web/
cp -r ./fonted/web/views productshop/fonted/web/

# 压缩目标文件夹
tar -czvf productshop.tar.gz productshop

# 输出提示信息
echo "pack finish"