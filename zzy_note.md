docker-compose -f docker-compose.dev.yml -p bluebell_dev down -v
用以上命令-v完全删除volume之后重建：
docker-compose -f docker-compose.dev.yml -p bluebell_dev up -d
之后，记得把migrations部分，seeds部分都要重新输入，放进新建的volume里面
