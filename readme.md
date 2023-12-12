
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