
```bash
docker build -t workday-local-redis -f Dockerfile-redis .

```

```bash
# get ec2 instance
aws ec2 describe-instances --filters "Name=instance-state-name,Values=running" --query 'Reservations[*].Instances[*].InstanceId' --output text

# get Public IP by EC2 instance
# aws ec2 describe-instances --instance-ids i-032cfefe2a6898fc7 --query 'Reservations[*].Instances[*].PublicIpAddress' --output text
aws ec2 describe-instances --instance-ids [instanceId] --query 'Reservations[*].Instances[*].PublicIpAddress' --output text

# ssh login
# ssh -i workdayKeyPeir.pem ec2-user@175.41.241.237
ssh -i [/path/to/your-key.pem] ec2-user@[your-ec2-instance-public-ip]

# deploy
eb deploy
```


```bash

# eb-engine.logにアクセスする方法
# Elastic Beanstalk環境のEC2インスタンスにSSHで接続する
sudo less /var/log/eb-engine.log

```

HTMLテンプレートURl
https://get.foundation/templates-previews-sites-f6-xy-grid/news-magazine.html

ECRへPushまでの道のり
```bash
# 未ログインの場合
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/l0m5q1g1

# workday appの場合
pwd
> /workday

# docker build -t nybeyond/workday:0.0.1 -f Dockerfile-node .
docker build -t nybeyond/workday:[version] -f Dockerfile-node .
> 略

docker images
> REPOSITORY                  TAG        IMAGE ID       CREATED        SIZE
> nybeyond/workday            0.0.1      1214e5db6574   19 hours ago   317MB


# docker tag 1214e5db6574 public.ecr.aws/l0m5q1g1/nybeyond/workday:0.0.1
docker tag [workday_docker_IMAGE_ID] public.ecr.aws/l0m5q1g1/nybeyond/workday:[version]
> (empty)

docker images
> REPOSITORY                                 TAG        IMAGE ID       CREATED        SIZE
> workday-app                                latest     4a665bec1ae3   19 hours ago   317MB
> nybeyond/workday                           0.0.1      1214e5db6574   19 hours ago   317MB

# docker push public.ecr.aws/l0m5q1g1/nybeyond/workday:0.0.1
docker push public.ecr.aws/l0m5q1g1/nybeyond/workday:[version]

```
[version] : 適宜入れて下さい。基本は(MAJOR.MINOR.PATCH) で管理する（[参考サイト](https://learn.microsoft.com/ja-jp/dotnet/core/versions/#semantic-versioning)
[workday_docker_IMAGE_ID]: docker Build完了したときいのImage Id