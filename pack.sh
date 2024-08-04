# 创建文件夹
mkdir -p productshopapp/backend/web/
mkdir -p productshopapp/fonted/web/

# 复制本地资源到目标文件夹
cp -r ./assets productshopapp/
cp -r ./backend/web/views productshopapp/backend/web/
cp -r ./fonted/web/views productshopapp/fonted/web/

# 压缩目标文件夹
tar -czvf productshopapp.tar.gz productshopapp

# 输出提示信息
echo "pack finish"